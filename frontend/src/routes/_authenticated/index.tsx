import { createFileRoute } from '@tanstack/react-router'
import { HostInfoPage } from "@/features/host-info"

export const Route = createFileRoute('/_authenticated/')({
  component: HostInfoPage,
})
