import { useQuery } from "@tanstack/react-query"
import type { HostInfo } from "../data/types"

const API_BASE = import.meta.env.VITE_API_BASE_URL

export function useHostInfo() {
  return useQuery<HostInfo>({
    queryKey: ['host-info'],
    queryFn: async () => {
      try {
        const res = await fetch(`${API_BASE}/api/v1/info`)

        if (!res.ok) {
          const text = await res.text()
          throw new Error(
            `Failed to load host info: ${res.status} ${res.statusText}${
              text ? ` - ${text}` : ''
            }`
          )
        }

        const json = await res.json()
        // console.log("🚀 Host info payload:", json) // optional debug
        return json
      } catch (err: any) {
        console.error("❌ Host info fetch error:", err)
        throw new Error(
          err.message || "Unknown error while fetching host info"
        )
      }
    },
    // retry: 2, // automatically retry twice on failure
    // staleTime: 15_000, // cache for 15 seconds
  })
}
