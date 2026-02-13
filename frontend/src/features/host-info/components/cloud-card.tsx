import {
  Card,
  CardContent,
  CardHeader,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible"

import { InfraIcon } from "@/features/host-info/components/infra/infra-icon"
import type { HostInfo } from "@/features/host-info/data/types"
const providerStyles = {
  aws: {
    accent: "before:bg-orange-500",
    glow: "hover:ring-orange-400/40",
  },
  gcp: {
    accent: "before:bg-blue-500",
    glow: "hover:ring-blue-400/40",
  },
  azure: {
    accent: "before:bg-cyan-500",
    glow: "hover:ring-cyan-400/40",
  },
  local: {
    accent: "before:bg-muted",
    glow: "hover:ring-muted/30",
  },
} as const

export function ProviderCard({ data }: { data: HostInfo }) {
  const provider =
    (data.cloud?.provider as
      | "aws"
      | "gcp"
      | "azure"
      | "local") ?? "local"


  const runtime =
    (data.runtime?.environment as
      | "kubernetes"
      | "docker"
      | "bare-metal") ?? "bare-metal"

  const styles = providerStyles[provider]

  return (
    <Card
      className={`
        relative overflow-hidden transition
        before:absolute before:left-0 before:top-0 before:h-full before:w-1
        ${styles.accent}
        hover:shadow-lg hover:ring-1 ${styles.glow}
      `}
    >
      {/* Sentence header */}
      <CardHeader className="pb-4">
        <div className="flex items-start gap-3">
          <InfraIcon name={provider} />

          <div className="space-y-1">
            <div className="text-sm font-medium">
              This workload is running on{" "}
              <span className="capitalize">{provider.toUpperCase()}</span>{" "}
              using{" "}
              <span className="capitalize">{runtime}</span>
            </div>

            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="capitalize">
                {provider.toUpperCase()}
              </Badge>
              <Badge variant="outline" className="capitalize">
                {runtime}
              </Badge>
            </div>
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-6">
        {/* Provider details */}
        <Section title="Provider details">
          <Info label="Region" value={data.cloud?.region} />

          {provider === "aws" && (
            <Info
              label="Account ID"
              value={data.cloud?.extra?.accountId}
            />
          )}

          {provider === "gcp" && (
            <Info
              label="Project"
              value={data.cloud?.extra?.projectId}
            />
          )}

          {provider === "azure" && (
            <Info
              label="Subscription"
              value={data.cloud?.extra?.subscriptionId}
            />
          )}
          {provider === "local" && (
            <>
              <Info label="Notes" value="Running locally" />

            </>
          )}


        </Section>

        {runtime === "kubernetes" && data.kubernetes && (
          <>
            <Separator />
            <KubernetesCollapsible kubernetes={data.kubernetes} />
          </>
        )}
      </CardContent>
    </Card>
  )
}

function KubernetesCollapsible({
  kubernetes,
}: {
  kubernetes: HostInfo["kubernetes"]
}) {
  return (
    <Collapsible defaultOpen>
      <CollapsibleTrigger className="flex items-center gap-2 text-sm font-semibold hover:underline">
        <InfraIcon name="kubernetes" />
        Kubernetes details
      </CollapsibleTrigger>

      <CollapsibleContent className="mt-3">
        <div className="rounded-lg border bg-muted/30 p-4 grid grid-cols-2 gap-4 text-sm">
          <Info
            label="Namespace"
            value={kubernetes?.podNamespace}
          />
          <Info
            label="Node"
            value={kubernetes?.nodeName}
          />
          <Info
            label="Pod"
            value={kubernetes?.podName}
          />
        </div>
      </CollapsibleContent>
    </Collapsible>
  )
}

/* ---------- small UI helpers ---------- */

function Section({
  title,
  children,
}: {
  title: string
  children: React.ReactNode
}) {
  return (
    <div className="space-y-3">
      <div className="text-sm font-semibold text-muted-foreground">
        {title}
      </div>
      <div className="grid grid-cols-2 gap-4">
        {children}
      </div>
    </div>
  )
}

function Info({
  label,
  value,
}: {
  label: string
  value?: string
}) {
  if (!value) return null
  return (
    <div>
      <div className="text-xs text-muted-foreground">
        {label}
      </div>
      <div className="font-medium truncate">{value}</div>
    </div>
  )
}
