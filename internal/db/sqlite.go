package db

import (
	"commuteboard/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	conn *sql.DB
}

type ScheduleRow struct {
	RouteID   int
	DayOfWeek int
	StartTime string
	EndTime   string
}

type RouteRow struct {
	ID              int
	OriginID        int
	DestinationID   int
	DistanceMeters  int
	DurationSeconds int
	RecordedAt      time.Time
}

type LocationRow struct {
	ID        int
	Name      string
	Latitude  float64
	Longitude float64
}

type RouteMeasurementRow struct {
	ID              int
	DistanceMeters  int
	DurationSeconds int
	RecordedAt      time.Time
}

func NewSQLite(path string) (*SQLite, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Enable WAL
	if _, err := conn.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		return nil, err
	}

	// Optional but recommended
	if _, err := conn.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		return nil, err
	}

	s := &SQLite{conn: conn}

	if err := s.initSchema(); err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)

	return s, nil
}

func (s *SQLite) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS routes (
		id TEXT PRIMARY KEY,
		origin_id TEXT NOT NULL,
		destination_id TEXT NOT NULL,
		distance_meters INTEGER NOT NULL,
		duration_seconds INTEGER NOT NULL,
		recorded_at DATETIME NOT NULL
	);
	`
	_, err := s.conn.Exec(schema)
	return err
}

func (s *SQLite) Save(ctx context.Context, c domain.Route) error {
	query := `
	INSERT INTO commutes (
		id,
		origin_id,
		destination_id,
		distance_meters,
		duration_seconds,
		recorded_at
	)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		distance_meters=excluded.distance_meters,
		duration_seconds=excluded.duration_seconds,
		recorded_at=excluded.recorded_at;
	`

	_, err := s.conn.ExecContext(
		ctx,
		query,
		c.ID,
		c.Origin,
		c.Destination,
		c.DistanceMeters,
		int(c.DurationSeconds.Seconds()),
		c.RecordedAt,
	)

	return err
}

func (s *SQLite) GetLocations(ctx context.Context) ([]LocationRow, error) {
	query := `
	SELECT id, name, latitude, longitude
	FROM locations;
	`

	rows, err := s.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []LocationRow

	for rows.Next() {
		var r LocationRow

		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Latitude,
			&r.Longitude,
		); err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, rows.Err()
}

func (s *SQLite) GetRouteRows(ctx context.Context) ([]RouteRow, error) {
	query := `
	SELECT id, origin_id, destination_id,
	       distance_meters, duration_seconds, recorded_at
	FROM routes;
	`

	rows, err := s.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []RouteRow

	for rows.Next() {
		var r RouteRow

		if err := rows.Scan(
			&r.ID,
			&r.OriginID,
			&r.DestinationID,
			&r.DistanceMeters,
			&r.DurationSeconds,
			&r.RecordedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, rows.Err()
}

func (db *SQLite) GetSchedules(ctx context.Context) ([]ScheduleRow, error) {
	rows, err := db.conn.QueryContext(ctx, `
		SELECT route_id, day_of_week, start_time, end_time
		FROM schedules
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []ScheduleRow

	for rows.Next() {
		var s ScheduleRow

		if err := rows.Scan(
			&s.RouteID,
			&s.DayOfWeek,
			&s.StartTime,
			&s.EndTime,
		); err != nil {
			return nil, err
		}

		schedules = append(schedules, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func ScheduleRowsToDomain(rows []ScheduleRow) (map[int]domain.Schedule, error) {
	result := make(map[int]domain.Schedule)

	for _, row := range rows {
		if row.DayOfWeek < 0 || row.DayOfWeek > 6 {
			return nil, fmt.Errorf("invalid day_of_week: %d", row.DayOfWeek)
		}

		start, err := parseHHMMToDuration(row.StartTime)
		if err != nil {
			return nil, err
		}

		end, err := parseHHMMToDuration(row.EndTime)
		if err != nil {
			return nil, err
		}

		weekday := time.Weekday(row.DayOfWeek)

		schedule := result[row.RouteID]

		if schedule.Days == nil {
			schedule.Days = make(map[time.Weekday][]domain.TimeRange)
		}

		schedule.Days[weekday] = append(
			schedule.Days[weekday],
			domain.TimeRange{
				Start: start,
				End:   end,
			},
		)

		result[row.RouteID] = schedule
	}

	return result, nil
}

func parseHHMMToDuration(value string) (time.Duration, error) {
	if len(value) != 5 {
		return 0, fmt.Errorf("invalid time format %q: expected HH:MM", value)
	}

	t, err := time.Parse("15:04", value)
	if err != nil {
		return 0, fmt.Errorf("invalid time format %q: %w", value, err)
	}

	return time.Duration(t.Hour())*time.Hour +
			time.Duration(t.Minute())*time.Minute,
		nil
}

func (s *SQLite) UpdateRouteMeasurements(
	ctx context.Context,
	rows []RouteMeasurementRow,
) error {

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE routes
		SET distance_meters = ?,
		    duration_seconds = ?,
		    recorded_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range rows {
		_, err := stmt.ExecContext(
			ctx,
			r.DistanceMeters,
			r.DurationSeconds,
			r.RecordedAt,
			r.ID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
