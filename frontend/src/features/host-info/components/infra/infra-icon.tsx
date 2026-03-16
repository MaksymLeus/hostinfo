import {
  SiAmazonwebservices,
  SiGooglecloud,
  SiKubernetes,
  SiDocker,
  SiGnubash,
} from "react-icons/si"
import { VscAzure } from "react-icons/vsc"
import { FaGolang } from "react-icons/fa6";


type InfraKey =
  | "aws"
  | "gcp"
  | "azure"
  | "local"
  | "kubernetes"
  | "docker"
  | "bare-metal"

const iconMap: Record<
  InfraKey,
  React.ComponentType<{ className?: string }>
> = {
  aws: SiAmazonwebservices,
  gcp: SiGooglecloud,
  azure: VscAzure,
  kubernetes: SiKubernetes,
  docker: SiDocker,
  "bare-metal": FaGolang,
  local: SiGnubash,
}

const colorMap: Record<InfraKey, string> = {
  aws: "text-orange-500",
  gcp: "text-blue-500",
  azure: "text-indigo-500",
  kubernetes: "text-sky-500",
  docker: "text-cyan-500",
  "bare-metal": "text-muted-foreground",
  local: "text-muted-foreground",
}

export function InfraIcon({
  name,
  className = "h-8 w-8 text-muted-foreground",
}: {
  name: InfraKey
  className?: string
}) {
  const Icon = iconMap[name] ?? iconMap.local
  return <Icon className={`${className} ${colorMap[name]}`} />
}
