// ExecutionPage.tsx

import { useState, useRef, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,  } from "@/components/ui/select"
import { TopNav } from '@/components/layout/top-nav'

import { ThemeSwitch } from '@/components/theme-switch'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { type PingResult, type DNSResult, type CurlResult, type TCPResult, type ExecResult, type ExecCommand, type ExecOutput, type Preset, HistoricCommand } from "../data/types"
import { cn } from "@/lib/utils"

// -------------------- Commands --------------------
const EXEC_COMMANDS: ExecCommand[] = [
  {
    value: "ping",
    label: "Ping",
    description: "ICMP ping a host",
    inputs: [{ name: "host", placeholder: "example.com" }],
  },
  {
    value: "curl",
    label: "HTTP GET",
    description: "Make an HTTP GET request",
    inputs: [{ name: "url", placeholder: "https://example.com" }],
  },
  {
    value: "dns",
    label: "DNS Lookup",
    description: "Resolve host to IPs",
    inputs: [{ name: "host", placeholder: "example.com" }],
  },
  {
    value: "tcp",
    label: "TCP Port Check",
    description: "Check if TCP port is open",
    inputs: [
      { name: "host", placeholder: "example.com" },
      { name: "port", placeholder: "80", type: "number" },
    ],
  },
]



// -------------------- Output Renderer --------------------
function Exec_Output({ viewMode, command, output }: ExecOutput) {
  if (!output) {
    return (
      <div className="text-muted-foreground italic">
          No output yet. Run a command.
      </div>
    )
  }
  if (viewMode === "json") {
    return (
      <pre className="whitespace-pre-wrap">
        {JSON.stringify(output, null, 2)}
      </pre>
    )
  }
  if (typeof output === "string") return <pre>{output}</pre>

  switch (command) {
    case "ping":
      const ping = output as PingResult
      return (
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4 text-sm">
          <div>Host: {ping.host}</div>
          <div>Sent: {ping.packets_sent}</div>
          <div>Received: {ping.packets_recv}</div>
          <div>Loss: {ping.loss_percent}%</div>
          <div>Min RTT: {ping.min_rtt}</div>
          <div>Max RTT: {ping.max_rtt}</div>
          <div>Avg RTT: {ping.avg_rtt}</div>
        </div>
      )
    case "dns":
      const dns = output as DNSResult
      return (
        <div className="space-y-1">
          <div>Host: {dns.host}</div>
          <div>CNAME: {dns.cname}</div>
          <div>
            IPs: {Array.isArray(dns.ips) ? dns.ips.join(", ") : "—"}
          </div>
        </div>
      )
    case "curl":
      const curl = output as CurlResult
      return (
        <div className="space-y-2">
          <div>URL: {curl.url}</div>
          <div>Status: {curl.status_code}</div>
          <pre className="bg-black/10 p-2 rounded">{curl.body}</pre>
        </div>
      )
    case "tcp":
      const tcp = output as TCPResult
      return (
        <div>
          {tcp.host}:{tcp.port} — {tcp.open ? "Open ✅" : `Closed ❌ ${tcp.error || ""}`}
        </div>
      )
    default:
      return (
        <pre>{JSON.stringify(output, null, 2)}</pre>
      )

  }
}

// -------------------- Page Component --------------------
const validators: Record<
  string,
  (value: string) => string | null
> = {
  host: (v) =>
    v.trim().length === 0 ? "Host is required" : null,

  url: (v) => {
    try {
      new URL(v)
      return null
    } catch {
      return "Invalid URL"
    }
  },

  port: (v) => {
    const n = Number(v)
    if (!Number.isInteger(n)) return "Port must be a number"
    if (n < 1 || n > 65535) return "Port must be 1–65535"
    return null
  },
}
type ProgressState = "idle" | "requesting" | "processing" | "done" | "error"
const progressColorMap: Record<ProgressState, string> = {
  idle: "",
  requesting: "bg-blue-500/10 text-blue-600 border-blue-500/20",
  processing: "bg-yellow-500/10 text-yellow-600 border-yellow-500/20",
  done: "bg-green-500/10 text-green-600 border-green-500/20",
  error: "bg-red-500/10 text-red-600 border-red-500/20",
}

export function ExecutionPage() {
  const API_BASE = import.meta.env.VITE_API_BASE_URL
  
  const [selectedCommand, setSelectedCommand] = useState<ExecCommand>(EXEC_COMMANDS[0])
  const [inputs, setInputs] = useState<Record<string, string>>({})
  const [output, setOutput] = useState<ExecResult | string>("")
  const [history, setHistory] = useState<HistoricCommand[]>([])
  const [loading, setLoading] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [viewMode, setViewMode] = useState<ExecOutput["viewMode"]>("pretty")

  const abortRef = useRef<AbortController | null>(null)

  const [progress, setProgress] = useState<ProgressState>("idle")
  useEffect(() => {
    if (progress === "done") {
      const timer = setTimeout(() => {
        setProgress("idle")
      }, 4000) // 3–5 seconds is perfect

      return () => clearTimeout(timer)
    }
  }, [progress])


  useEffect(() => {
    abortRef.current?.abort()
    setOutput("")
    setErrors({})
    setProgress("idle")
    setLoading(false)
  }, [selectedCommand.value])

  // Handlers
  const hasEmptyInputs = selectedCommand.inputs.some(
    (i) => !inputs[i.name]?.trim()
  )
  const validateInputs = () => {
    const nextErrors: Record<string, string> = {}

    selectedCommand.inputs.forEach((input) => {
      const value = inputs[input.name] || ""
      const validate = validators[input.name]

      if (validate) {
        const error = validate(value)
        if (error) nextErrors[input.name] = error
      }
    })

    setErrors(nextErrors)
    return Object.keys(nextErrors).length === 0
  }

  const handleExecute = async () => {
    if (loading) return
    if (!validateInputs()) return

    const params = new URLSearchParams(inputs).toString()
    const url = `${API_BASE}/api/v1/${selectedCommand.value}?${params}`

    try {
      setLoading(true)
      setProgress("requesting")

      const controller = new AbortController()
      abortRef.current = controller

      const timeout = setTimeout(() => {
        controller.abort()
      }, 10_000)


      const res = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        signal: controller.signal,
      })
      clearTimeout(timeout)

      if (!res.ok) throw new Error(`Request failed ${res.status}`)

      setProgress("processing")
      const json = await res.json()
      setOutput(json)

      // update history (max 10)
      setHistory((h) => {
        const newEntry = {
          command: selectedCommand.value,
          parameters: Object.entries(inputs).map(([k, v]) => `--${k}=${v}`),
        }
        const exists = h.some(
          (e) =>
            e.command === newEntry.command &&
            JSON.stringify(e.parameters) === JSON.stringify(newEntry.parameters)
        )
        if (exists) return h

        return [newEntry, ...h].slice(0, 10)
      })



      setProgress("done")
    } catch (err: any) {
      if (err.name === "AbortError") {
        setOutput("Request timed out")
      } else {
        setOutput(err.message)
      }
      setProgress("error")
    } finally {
      abortRef.current = null
      setLoading(false)
    }
  }
  const PRESETS: Preset[] = [
    {
      label: "Ping Google DNS",
      command: "ping",
      inputs: { host: "8.8.8.8" },
    },
    {
      label: "DNS Cloudflare",
      command: "dns",
      inputs: { host: "1.1.1.1" },
    },
    {
      label: "HTTP Health",
      command: "curl",
      inputs: { url: "https://example.com/health" },
    },
  ]




  return (
    <>
      {/* ===== Header ===== */}
      <Header>
        <div className="flex w-full items-center justify-center gap-4">
          <TopNav links={topNav} />
        </div>
          <ThemeSwitch />
      </Header>
      
      {/* ===== Main ===== */}
      <Main>
        <div className="p-6 flex flex-col items-center space-y-6">
          {/* Command Selection */}
          <div className="w-full max-w-4xl text-left">
          <h2 className="text-xl font-semibold">Network Utilities</h2>
          <p className="text-sm text-muted-foreground">
            Run common network diagnostics like ping, DNS lookup, HTTP checks, and TCP port scans.
          </p>
          </div>

          <Card className="w-full max-w-4xl">
            <CardHeader>
              <CardTitle>Command</CardTitle>
            </CardHeader>

            <div className="flex gap-2 flex-wrap max-w-4xl px-6">
              {PRESETS.map(p => (
                <Badge

                  key={p.label}
                  variant="secondary"
                  className="
                    cursor-pointer
                    px-3 py-1
                    font-mono
                    hover:bg-blue-500/20
                    hover:text-blue-400
                    transition-colors
                  "
                  onClick={() => {
                    const cmd = EXEC_COMMANDS.find(c => c.value === p.command)
                    if (!cmd) return

                    abortRef.current?.abort()

                    setSelectedCommand(cmd)
                    setInputs(p.inputs)
                    setErrors({})
                    setOutput("")
                    setProgress("idle")
                  }}

                >
                  {p.label}
                </Badge>
              ))}
            </div>
            <CardContent className="space-y-6">
              {/* <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center"> */}
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-[1fr_2fr_auto] items-start">

                {/* Command selector */}
                {/* <div className="flex-1"> */}
                <div className="relative flex flex-col items-center gap-4 self-start">
                  <Select
                    value={selectedCommand.value}
                    onValueChange={(val) => {
                      const cmd = EXEC_COMMANDS.find((c) => c.value === val)
                      if (cmd) {
                        setSelectedCommand(cmd)
                        setInputs({})
                      }
                    }}
                  >
                    <SelectTrigger>
                      <SelectValue>{selectedCommand.label}</SelectValue>
                    </SelectTrigger>

                    <SelectContent>
                      <SelectGroup>
                        {EXEC_COMMANDS.map((cmd) => (
                          <SelectItem key={cmd.value} value={cmd.value}>
                            <div className="flex flex-col gap-1">
                              <span className="font-medium">{cmd.label}</span>
                              {cmd.description && (
                                <span className="text-xs text-muted-foreground">
                                  {cmd.description}
                                </span>
                              )}
                            </div>
                          </SelectItem>
                        ))}
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </div>


                {/* Inputs */}
                <div className="flex flex-wrap gap-2 flex-1 min-h-[76px]">

                  {selectedCommand.inputs.map((i) => (
                    <div key={i.name} className="flex flex-col min-w-[160px] flex-1">
                      <Input
                        className="
                          bg-black/30 
                          border border-white/10 
                          focus:border-blue-500 
                          focus:ring-1 focus:ring-blue-500
                          rounded-lg
                          font-mono
                          h-11
                        "
                        placeholder={i.placeholder}
                        type={i.type || "text"}
                        value={inputs[i.name] || ""}
                        onChange={(e) =>
                          setInputs((p) => ({ ...p, [i.name]: e.target.value }))
                        }
                        onKeyDown={(e) => {
                              if (e.key === "Enter" && !loading) handleExecute()
                            }}

                      />
                      <div className="h-4 mt-1">
                        <span
                          className={`
                            text-xs text-red-500
                            transition-opacity duration-200
                            ${errors[i.name] ? "opacity-100" : "opacity-0"}
                          `}
                        >
                          {errors[i.name] ?? "placeholder"}
                        </span>
                      </div>

                    </div>
                  ))}
                </div>

                {/* Execute button */}
                <div className="flex gap-2 self-start">
                  <Button
                    size="lg"
                    className="
                      px-10
                      font-mono
                      bg-gradient-to-r
                      from-primary to-primary/70
                      hover:from-primary/90 hover:to-primary/60
                      shadow-md
                      transition-all
                    "
                    onClick={handleExecute}
                    disabled={loading || hasEmptyInputs}
                  >
                    {loading && (
                      <svg
                        className="mr-2 h-4 w-4 animate-spin"
                        viewBox="0 0 24 24"
                        fill="none"
                      >
                        <circle
                          className="opacity-25"
                          cx="12"
                          cy="12"
                          r="10"
                          stroke="currentColor"
                          strokeWidth="4"
                        />
                        <path
                          className="opacity-75"
                          fill="currentColor"
                          d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
                        />
                      </svg>
                    )}
                    {loading ? "Executing…" : "Execute"}
                  </Button>

                  <Button
                    variant="outline"
                    disabled={!loading}
                    className="font-mono"
                    onClick={() => abortRef.current?.abort()}
                  >
                    Cancel
                  </Button>
                </div>
              </div>



              {/* History */}
              {history.length > 0 && (
                <div className="space-y-2">
                  <div className="text-xs uppercase tracking-wide text-muted-foreground">
                    Recent executions
                  </div>

                  <div className="flex flex-wrap gap-2">
                    {history.map((h) => {
                      const searchParams = Object.fromEntries(
                        h.parameters.map((p) => {
                          const [key, value] = p.replace(/^--/, "").split("=")
                          return [key, value]
                        }))

                      const isActive = JSON.stringify(searchParams) === JSON.stringify(inputs)

                      return (
                      <Badge
                        key={h.command + h.parameters.join(" ")}
                        variant="secondary"
                        className={cn(
                          "font-mono cursor-pointer transition-colors border px-3 py-1",
                          isActive
                            ? " text-primary border-primary"
                            : "border-border bg-muted/30 hover:bg-primary/20 hover:text-primary hover:border-primary"
                        )}
                        onClick={() => {
                          const cmd = EXEC_COMMANDS.find(c => c.value === h.command)
                          
                          if (!cmd) return

                          abortRef.current?.abort()

                          setSelectedCommand(cmd)
                          setInputs(searchParams)
                          setErrors({})
                          setOutput("")
                          setProgress("idle")
                        }}
                      >
                        {h.command} {h.parameters.join(" ")}
                      </Badge>
                      )
                    })}
                  </div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Output */}
          <Card className="w-full max-w-4xl rounded-xl shadow-md bg-card p-6 relative ">

            {/* <CardHeader> */}
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>Output</CardTitle>
            </CardHeader>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => setViewMode(v => v === "pretty" ? "json" : "pretty")}
            >
              {viewMode === "json" ? "Pretty view" : "Raw JSON"}
            </Button>

            {progress !== "idle" && (
              <div className="absolute right-3 top-3 z-10">

                <Badge
                  className={`${progressColorMap[progress]} transition-opacity duration-300`}
                >
                  {progress.toUpperCase()}
                </Badge>
              </div>
            )}
            <CardContent className="max-h-[500px] overflow-auto rounded-lg bg-black/20 p-4 font-mono text-sm">
              {output && (
                <div className="mb-2 text-xs text-muted-foreground">
                  $ {selectedCommand.value}{" "}
                  {Object.entries(inputs)
                    .map(([k, v]) => `--${k}=${v}`)
                    .join(" ")}
                </div>
              )}
              
              <Exec_Output viewMode={viewMode} command={selectedCommand.value} output={output} />
            </CardContent>
          </Card>

        </div>
      </Main>
    </>
  )
}

const topNav = [
  {
    title: 'Overview',
    href: '/',
    isActive: false,
    disabled: false,
  },
  {
    title: 'Network Utilities',
    href: 'network-utilities',
    isActive: true,
    disabled: true,
  },
]