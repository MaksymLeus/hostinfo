import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Cpu, Clock, MonitorCog, Hexagon } from 'lucide-react'
import { GiSolarSystem } from "react-icons/gi";
import type { HostInfo } from '../data/types'
import { useHealthStatus } from '../hooks/use-health-status'

type SystemCardProps = {
  data: HostInfo
}

const statusColorMap = {
  healthy: 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20',
  warning: 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20',
  critical: 'bg-red-500/10 text-red-500 border-red-500/20',
}

export function SystemCard({ data }: SystemCardProps) {
  const { status, loading } = useHealthStatus()

  return (
    <Card className="relative overflow-hidden ">
      <div className="absolute right-3 top-3">
        <Badge className={statusColorMap[loading ? 'warning' : status]}>
          {(loading ? 'warning' : status).toUpperCase()}
        </Badge>
      </div>

      <CardHeader className="pb-3">
        <CardTitle className="flex items-center gap-2">
          <GiSolarSystem className="h-8 w-8 text-muted-foreground" />
          System
        </CardTitle>
      </CardHeader>

      <CardContent className="space-y-4 ">

        <div>
          <div className="text-xs text-muted-foreground">Hostname</div>
          <div className="text-xl font-bold">{data.hostname}</div>
        </div>

        <div className="grid grid-cols-2 gap-4 ">
          <div className="rounded-md border bg-muted/20 px-4 py-3 border-l-4 border-primary">
          
            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <MonitorCog className="h-3 w-3" />
              Operating system
            </div>
            <div className="mt-1 flex gap-2">
              <Badge variant="outline">{data.os}</Badge>
              {data.distro && (
                <Badge variant="outline">{data.distro}</Badge>
              )}
            </div>
          </div>

          <div className="rounded-md border bg-muted/20 px-4 py-3 text-right border-r-4 border-primary">

            <div className="flex items-center gap-1 justify-end text-xs text-muted-foreground">
              <Hexagon className="h-3 w-3" />
              GoVersion
            </div>
            <Badge variant="secondary" className="mt-1">
              {data.goVersion}
            </Badge>
          </div>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div className="rounded-md border bg-muted/20 px-4 py-3 border-l-4 border-primary">

            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <Cpu className="h-3 w-3" />
              Architecture
            </div>
            <div className="font-medium">{data.arch}</div>
          </div>

          <div className="rounded-md border bg-muted/20 px-4 py-3 text-right border-r-4 border-primary">

            <div className="flex items-center gap-1 justify-end text-xs text-muted-foreground">
              <Clock className="h-3 w-3" />
              Uptime
            </div>
            <div className="font-medium">{data.uptime}</div>
          </div>
        </div>
        <div className="flex justify-end text-xs text-muted-foreground pt-2">
          Last update:&nbsp;
          <span className="text-foreground">
            {new Date(data.updatedAt).toLocaleString()}
          </span>
        </div>

      </CardContent>
    </Card>
  )
}
