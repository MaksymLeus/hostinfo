import { createFileRoute } from '@tanstack/react-router'
import { ExecutionPage } from "@/features/host-info/components/ExecutionPage"

export const Route = createFileRoute('/network-utilities')({
  component: ExecutionPage,
})
