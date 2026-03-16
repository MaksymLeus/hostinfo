export type HealthStatus = 'healthy' | 'warning' | 'critical'

export function mapHealthStatus(apiStatus?: string): HealthStatus {
  switch (apiStatus) {
    case 'ok':
      return 'healthy'
    case 'degraded':
      return 'warning'
    case 'error':
    case 'fail':
      return 'critical'
    default:
      return 'critical'
  }
}
