package commute

import (
	"context"
	"database/sql"
	"fmt"
	"signalboard/internal/db"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CommuteDb struct {
	db *db.SQLite
}

func NewCommuteDb(db *db.SQLite) *CommuteDb {
	return &CommuteDb{db: db}
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
	DistanceMeters  *int
	DurationSeconds *time.Duration
	RecordedAt      *time.Time
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

func (s *CommuteDb) GetLocations(ctx context.Context) ([]LocationRow, error) {
	query := `
	SELECT id, name, latitude, longitude
	FROM locations;
	`

	rows, err := s.db.Conn().QueryContext(ctx, query)
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

func (s *CommuteDb) Save(ctx context.Context, c Route) error {
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

	_, err := s.db.Conn().ExecContext(
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

func (s *CommuteDb) GetRouteRows(ctx context.Context) ([]RouteRow, error) {
	query := `
	SELECT id, origin_id, destination_id,
	       distance_meters, duration_seconds, recorded_at
	FROM routes;
	`

	rows, err := s.db.Conn().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []RouteRow

	for rows.Next() {
		var r RouteRow

		var distance sql.NullInt64
		var duration sql.NullInt64
		var recorded sql.NullTime

		if err := rows.Scan(
			&r.ID,
			&r.OriginID,
			&r.DestinationID,
			&distance,
			&duration,
			&recorded,
		); err != nil {
			return nil, err
		}

		if distance.Valid {
			d := int(distance.Int64)
			r.DistanceMeters = &d
		}

		if duration.Valid {
			d := time.Duration(duration.Int64) * time.Second
			r.DurationSeconds = &d
		}

		if recorded.Valid {
			r.RecordedAt = &recorded.Time
		}

		result = append(result, r)
	}

	return result, rows.Err()
}

func (s *CommuteDb) GetSchedules(ctx context.Context) ([]ScheduleRow, error) {
	query := `
		SELECT route_id, day_of_week, start_time, end_time
		FROM schedules
	`

	rows, err := s.db.Conn().QueryContext(ctx, query)
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

func ScheduleRowsToDomain(rows []ScheduleRow) (map[int]Schedule, error) {
	result := make(map[int]Schedule)

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
			schedule.Days = make(map[time.Weekday][]TimeRange)
		}

		schedule.Days[weekday] = append(
			schedule.Days[weekday],
			TimeRange{
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

func (s *CommuteDb) UpdateRouteMeasurements(
	ctx context.Context,
	rows []RouteMeasurementRow,
) error {

	tx, err := s.db.Conn().BeginTx(ctx, nil)
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
