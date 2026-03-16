import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Cpu, Server, HardDrive, Activity } from 'lucide-react'
import type { HostInfo } from '../data/types'

type MetricCardProps = {
  title: string
  value: string
  details?: React.ReactNode
  icon: React.ReactNode
  status?: 'healthy' | 'warning' | 'critical'
}

export function MetricCard({ title, value, details, icon, status }: MetricCardProps) {
  const statusColorMap = {
    healthy: 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20',
    warning: 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20',
    critical: 'bg-red-500/10 text-red-500 border-red-500/20',
  }

  return (
    <Card className="relative overflow-hidden border-l-2 border-primary">
      {status && (
        <div className="absolute right-3 top-3">
          <Badge className={statusColorMap[status]}>{status.toUpperCase()}</Badge>
        </div>
      )}

      <CardHeader className="pb-3">
        <CardTitle className="flex items-center gap-2">
          {icon}
          {title}
        </CardTitle>
      </CardHeader>

      <CardContent className="space-y-2">
        <div className="text-2xl font-bold">{value}</div>
        {details && <div className="text-xs text-muted-foreground">{details}</div>}
      </CardContent>
    </Card>
  )
}

export function HostMetrics({ data }: { data: HostInfo }) {
  const cpuUsage = data.cpu?.usagePercent ?? 0
  const memUsage = data.memory?.usedPercent ?? 0
  const cores = data.cpu?.cores ?? 1

  const getStatus = (value: number, thresholds: number[]) => {
    if (value <= thresholds[0]) return 'healthy'
    if (value <= thresholds[1]) return 'warning'
    return 'critical'
  }

  const loadStatus = (value: number) => getStatus(value, [cores, cores * 1.5])
  const loadColor = (status: string) => {
    switch (status) {
      case 'healthy':
        return 'bg-emerald-500'
      case 'warning':
        return 'bg-yellow-500'
      case 'critical':
        return 'bg-red-500'
      default:
        return 'bg-gray-400'
    }
  }

  const renderLoadBar = (label: string, value: number) => {
    const status = loadStatus(value)
    const widthPercent = Math.min((value / cores) * 100, 100)
    return (
      <div className="flex items-center gap-2">
        <span className="w-8 text-xs font-semibold text-muted-foreground">{label}</span>
        <div className="w-full h-2 bg-sky-900 rounded">
          <div
            className={`${loadColor(status)} h-2 rounded`}
            style={{ width: `${widthPercent}%` }}
          />
        </div>
        <span className="w-10 text-xs text-right font-mono">{value.toFixed(2)}</span>
      </div>
    )
  }

  return (
    <Card className="overflow-hidden">
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center gap-2">
          <Server className="h-5 w-5 text-muted-foreground" />
          Hardware
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {/* CPU */}
          <MetricCard
            title="CPU Usage"
            value={`${cpuUsage.toFixed(1)}%`}
            details={`${data.cpu?.cores} cores`}
            icon={<Cpu className="h-5 w-5 text-muted-foreground" />}
            status={getStatus(cpuUsage, [50, 80])}
          />

          {/* Memory */}
          <MetricCard
            title="Memory Usage"
            value={`${memUsage.toFixed(1)}%`}
            details={`${data.memory?.usedMB} / ${data.memory?.totalMB} MB`}
            icon={<HardDrive className="h-5 w-5 text-muted-foreground" />}
            status={getStatus(memUsage, [50, 80])}
          />

          {/* Load */}
          <MetricCard
            title="Load Usage"
            value={`${data.load?.load1?.toFixed(2) ?? 0}`}
            icon={<Activity className="h-5 w-5 text-muted-foreground" />}
            status={loadStatus(data.load?.load1 ?? 0)}
            details={
              <div className="space-y-1">
                {renderLoadBar('1m', data.load?.load1 ?? 0)}
                {renderLoadBar('5m', data.load?.load5 ?? 0)}
                {renderLoadBar('15m', data.load?.load15 ?? 0)}
              </div>
            }
          />
        </div>
      </CardContent>
    </Card>
  )
}
