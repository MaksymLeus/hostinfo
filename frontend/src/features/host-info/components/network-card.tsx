import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export function NetworkCard({
  ips,
  macs,
}: {
  ips: string[]
  macs: string[]
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>🌐 Network</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <div>
          <div className="font-medium">IPs</div>
          {/* {ips.map(ip => (
            <code key={ip} className="block">{ip}</code>
          ))} */}
          {ips.map((ip, idx) => (
            <code key={`ip-${ip}-${idx}`} className="block">{ip}</code>
          ))}
        </div>

        <div>
          <div className="font-medium">MACs</div>
          {/* {macs.map(mac => (
            <code key={mac} className="block">{mac}</code>
          ))} */}
          {macs.map((mac, idx) => (
            <code key={`mac-${mac}-${idx}`} className="block">{mac}</code>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
