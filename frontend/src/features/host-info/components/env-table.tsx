import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Switch } from "@/components/ui/switch"
import { useState } from "react"

const SENSITIVE = /(TOKEN|SECRET|KEY|PASSWORD)/i

export function EnvTable({ env }: { env: Record<string, string> }) {
  const [showSensitive, setShowSensitive] = useState(false)

  return (
    <div className="space-y-3">
      <div className="flex items-center gap-2">
        <Switch checked={showSensitive} onCheckedChange={setShowSensitive} />
        <span>Show sensitive values</span>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Key</TableHead>
            <TableHead>Value</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {Object.entries(env).map(([k, v]) => {
            const masked = SENSITIVE.test(k) && !showSensitive
            return (
              <TableRow key={k}>
                <TableCell>{k}</TableCell>
                <TableCell>
                  <code>{masked ? "••••••••" : v}</code>
                </TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>
    </div>
  )
}
