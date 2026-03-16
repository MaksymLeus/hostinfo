import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export function RuntimeCard({
  startTime,
  now,
}: {
  startTime: string
  now: string
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>⏱ Runtime</CardTitle>
      </CardHeader>
      <CardContent className="space-y-1 text-sm">
        <div>Started: {startTime}</div>
        <div>Now: {now}</div>
      </CardContent>
    </Card>
  )
}
