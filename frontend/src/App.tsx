import { useEffect, useState } from "react"

// const URL = "http://localhost:3333"

type Route = {
    id: number
    origin: string
    destination: string
    duration_minutes: number | null
    distance_km: number | null
    recorded_at: string | null
    active_now: boolean
}

type Status = {
    running: boolean
    tick_rate: string
    update_rate: string
    locations: number
    last_tick: string
}

function timeAgo(dateString: string): string {
    if (!dateString) return "never"

    const diff = Math.floor((Date.now() - new Date(dateString).getTime()) / 1000)

    if (diff < 60) return `${diff}s ago`
    if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
    return `${Math.floor(diff / 3600)}h ago`
}

export default function App() {
    const [routes, setRoutes] = useState<Route[]>([])
    const [status, setStatus] = useState<Status | null>(null)
    const [refreshing, setRefreshing] = useState(false)

    async function fetchRoutes() {
        const res = await fetch("http://localhost:3333/commute/routes/active")
        const data = await res.json()
        setRoutes(data || [])
    }

    async function fetchStatus() {
        const res = await fetch("http://localhost:3333/engine/status")
        const data = await res.json()
        setStatus(data)
    }

    async function refreshRoutes() {
        setRefreshing(true)

        try {
            await fetch("http://localhost:3333/commute/routes/refresh", {
                method: "POST",
            })

            // After refreshing, fetch updated data
            await fetchRoutes()
            await fetchStatus()
        } catch (err) {
            console.error("Failed to refresh routes", err)
        } finally {
            setRefreshing(false)
        }
    }

    useEffect(() => {
        fetchRoutes()
        fetchStatus()

        const interval = setInterval(() => {
            fetchRoutes()
            fetchStatus()
        }, 10000) // every 10s

        return () => clearInterval(interval)
    }, [])

    const isStale =
        status?.last_tick &&
        Date.now() - new Date(status.last_tick).getTime() > 15000

    return (
        <div style={{ padding: "20px", fontFamily: "sans-serif" }}>
            <h1>Active Routes</h1>

            {/* Refresh button */}
            <button onClick={refreshRoutes} disabled={refreshing}>
                {refreshing ? "Refreshing..." : "Refresh Routes"}
            </button>

            {/* Routes */}
            {routes.length === 0 ? (
                <p>No active routes</p>
            ) : (
                <ul>
                    {routes.map(route => (
                        <div
                            key={route.id}
                            style={{
                                border: "1px solid #ccc",
                                borderRadius: "8px",
                                padding: "12px",
                                marginBottom: "12px"
                            }}
                        >
                            <h3>
                                {route.origin} → {route.destination}
                            </h3>

                            <p>
                                ⏱ {route.duration_minutes !== null ? `${route.duration_minutes} min` : "N/A"}
                            </p>

                            <p>
                                📍 {route.distance_km !== null
                                    ? `${route.distance_km.toFixed(1)} km`
                                    : "N/A"}
                            </p>
                        </div>
                    ))}
                </ul>
            )}

            {/* Engine Status */}
            {status && (
                <div style={{ marginBottom: "1rem" }}>
                    <p>
                        Engine: {status.running ? "🟢 Running" : "🔴 Stopped"}
                    </p>
                    <p style={{ color: isStale ? "red" : "green" }}>
                        Last updated: {timeAgo(status.last_tick)}
                    </p>
                </div>
            )}

        </div>
    )
}
