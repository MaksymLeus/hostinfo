import { useEffect, useState } from 'react'

export type HealthStatus = 'healthy' | 'warning' | 'critical'
const API_BASE = import.meta.env.VITE_API_BASE_URL

export function useHealthStatus(
  url: string = `${API_BASE}/healthz`,
  pollInterval = 60_000
) {
  const [status, setStatus] = useState<HealthStatus>('critical')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true

    const fetchHealth = async () => {
      try {
        const res = await fetch(url)
        if (!res.ok) throw new Error('Health check failed')

        const json = await res.json()
        // console.log("🚀 Host info status:", json) 

        if (!mounted) return

        setStatus(
          String(json?.status).toLowerCase() === 'ok'
            ? 'healthy'
            : 'critical'
        )
      } catch {
        if (mounted) {
          setStatus('critical')
        }
      } finally {
        if (mounted) {
          setLoading(false)
        }
      }
    }

    fetchHealth()
    const interval = setInterval(fetchHealth, pollInterval)

    return () => {
      mounted = false
      clearInterval(interval)
    }
  }, [url, pollInterval])

  return { status, loading }
}
