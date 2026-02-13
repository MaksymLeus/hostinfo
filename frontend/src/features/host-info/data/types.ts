export interface CPUInfo {
  cores: number
  model: string
  usagePercent: number
  load1: number
  load5: number
  load15: number
}

export interface MemoryInfo {
  totalMB: number
  usedMB: number
  usedPercent: number
}

export interface LoadInfo {
  load1: number
  load5: number
  load15: number
}

export interface CloudInfo {
  provider: string
  region?: string
  zone?: string
  instance?: string
  extra?: Record<string, any> | null
}

export interface KubernetesInfo {
  enabled: boolean
  podName?: string
  podNamespace?: string
  podIP?: string
  nodeName?: string
  serviceAccount?: string
  container?: string
}

export interface RuntimeInfo {
  environment: string
  details?: string
}

export interface HostInfo {
  hostname: string
  os: string
  distro?: string
  arch: string
  goVersion: string

  cloud?: CloudInfo
  runtime?: RuntimeInfo
  kubernetes?: KubernetesInfo

  ips: string[]
  macs: string[]

  startTime: string
  now: string
  uptime: string
  updatedAt: string

  env: Record<string, string>

  cpu: CPUInfo
  memory: MemoryInfo
  load: LoadInfo
}



export interface PingResult {
  host: string
  packets_sent: number
  packets_recv: number
  loss_percent: number
  min_rtt: string
  max_rtt: string
  avg_rtt: string
}

export interface DNSResult {
  host: string
  ips: string[]
  cname: string
}

export interface CurlResult {
  url: string
  status_code: number
  body: string
}

export interface TCPResult {
  host: string
  port: string
  open: boolean
  error?: string
}

// Union type for generic output
export type ExecResult = PingResult | DNSResult | CurlResult | TCPResult


export type ExecCommand = {
  value: "ping" | "curl" | "dns" | "tcp"
  label: string
  description?: string
  inputs: { name: string; placeholder: string; type?: string }[]
}

export type ExecOutput = {
  viewMode: "pretty" | "json"
  command: string
  output: ExecResult | string
}

export type Preset = {
  label: string
  command: ExecCommand["value"]
  inputs: Record<string, string>
}

export type HistoricCommand = {
  parameters: string[]
  command: ExecCommand["value"]
}
