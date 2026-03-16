// provider-meta.ts
import type { HostInfo } from "../data/types"

export type ProviderKey = "aws" | "gcp" | "azure" | "local"

type ProviderField = {
  label: string
  value: (data: HostInfo) => string | undefined
  span?: 1 | 2
}

type ProviderMeta = {
  title: string
  accentClass: string
  glowClass: string
  sentence: (data: HostInfo) => string
  fields: ProviderField[]
}

export const providerMetaConfig: Record<ProviderKey, ProviderMeta> = {
  aws: {
    title: "Amazon Web Services",
    accentClass: "before:bg-orange-500",
    glowClass: "hover:ring-orange-400/40",
    sentence: () => "This workload is running on AWS",
    fields: [
      {
        label: "Region",
        value: d => d.cloud?.region,
      },
      {
        label: "Account ID",
        value: d => d.cloud?.extra?.accountId,
      },
    ],
  },

  gcp: {
    title: "Google Cloud Platform",
    accentClass: "before:bg-blue-500",
    glowClass: "hover:ring-blue-400/40",
    sentence: () => "This workload is running on GCP",
    fields: [
      {
        label: "Region",
        value: d => d.cloud?.region,
      },
      {
        label: "Project ID",
        value: d => d.cloud?.extra?.projectId,
      },
    ],
  },

  azure: {
    title: "Microsoft Azure",
    accentClass: "before:bg-cyan-500",
    glowClass: "hover:ring-cyan-400/40",
    sentence: () => "This workload is running on Azure",
    fields: [
      {
        label: "Region",
        value: d => d.cloud?.region,
      },
      {
        label: "Subscription ID",
        value: d => d.cloud?.extra?.subscriptionId,
      },
    ],
  },

  local: {
    title: "Local environment",
    accentClass: "before:bg-muted",
    glowClass: "hover:ring-muted/30",
    sentence: () => "This workload is running locally",
    fields: [
      {
        label: "Environment",
        value: () => "Local development",
      },
      {
        label: "Hostname",
        value: d => d.hostname,
      },
      {
        label: "OS",
        value: d => d.os,
      },
      {
        label: "Architecture",
        value: d => d.arch,
      },
    ],
  },
}
