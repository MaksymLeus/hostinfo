import { useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { TopNav } from '@/components/layout/top-nav'
import { ThemeSwitch } from '@/components/theme-switch'

import { useHostInfo } from './hooks/use-host-info'
import { SystemCard } from './components/system-card'
import { ProviderCard } from './components/cloud-card'
import { NetworkCard } from './components/network-card'
import { HostMetrics } from './components/hostmetrics-card'
import { EnvTable } from './components/env-table'

// import { Logo } from '@/assets/logo'


export function HostInfoPage() {
  const { data, isLoading, error, refetch } = useHostInfo()
  const [copied, setCopied] = useState(false)

  const handleCopy = () => {
    if (!data) return
    navigator.clipboard.writeText(JSON.stringify(data, null, 2))
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <>
      {/* ===== Header ===== */}
      <Header>
        <div className="flex w-full items-center justify-center gap-4">
          {/* <Logo className='me-2' /> */}
          <TopNav links={topNav} />
        </div>
          <ThemeSwitch />
      </Header>

      {/* ===== Main ===== */}
      <Main>
        {/* Page title */}
        <div className="mb-2 flex items-center justify-between space-y-2">
          <h1 className="text-2xl font-bold tracking-tight">
            Host Information
          </h1>
          <div className="flex items-center space-x-2">
            <Button onClick={() => refetch()}>
              Refresh
            </Button>
          </div>
        </div>

        {isLoading && <div>Loading host info…</div>}
        {error && <div>Failed to load host info</div>}
        {!data && !isLoading && <div>No data</div>}

        {data && (
          <Tabs defaultValue="overview" className="space-y-4">
            {/* Tabs header */}
            <div className="w-full overflow-x-auto pb-2">

              <TabsList>
                <TabsTrigger value="overview">Overview</TabsTrigger>
                <TabsTrigger value="network">Network</TabsTrigger>
                <TabsTrigger value="environment">Environment Variables</TabsTrigger>
                <TabsTrigger value="raw">Raw JSON</TabsTrigger>
              </TabsList>
            </div>

            {/* ===== Overview ===== */}
            <TabsContent value="overview" className="space-y-4">

              <div className="grid grid-cols-2 gap-4">
                <SystemCard data={data} />
                <ProviderCard data={data} />
              </div>
              {/* Hardware / Metrics */}
              <HostMetrics data={data} />
            </TabsContent>

            {/* ===== Network ===== */}
            <TabsContent value="network" className="space-y-4">
              <Card>
                <CardHeader>
                  <CardTitle>Network Interfaces</CardTitle>
                  <CardDescription>
                    IP and MAC addresses detected on the host
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <NetworkCard ips={data.ips} macs={data.macs} />
                </CardContent>
              </Card>
            </TabsContent>

            {/* ===== Environment ===== */}
            <TabsContent value="environment" className="space-y-4">
              <Card>
                <CardHeader>
                  <CardTitle>Environment Variables</CardTitle>
                  <CardDescription>
                    Sensitive values are hidden by default
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  {data.env && Object.keys(data.env).length > 0 ? (
                    <EnvTable env={data.env} />
                  ) : (
                    <p className="text-sm text-muted-foreground">
                      Environment variables are not available. To enable set env {' '}
                      <code>FF_ENVIRONMENT_VARIABLES</code>  to <strong>true</strong>.
                    </p>
                  )}
                </CardContent>
              </Card>
            </TabsContent>

            {/* ===== Raw JSON ===== */}
            <TabsContent value="raw">
              <Card>
                <CardHeader className="flex items-center justify-between">
                  <div>
                    <CardTitle>Raw API Response</CardTitle>
                    <CardDescription>GET /api/hostinfo</CardDescription>
                  </div>
                  <Button size="sm" onClick={handleCopy}>
                    {copied ? 'Copied!' : 'Copy'}
                  </Button>
                </CardHeader>
                <CardContent>
                  <pre className="max-h-[500px] overflow-auto rounded-md bg-muted p-4 text-sm">
                    {JSON.stringify(data, null, 2)}
                  </pre>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        )}
      </Main>
    </>
  )
}

const topNav = [
  {
    title: 'Overview',
    href: '',
    isActive: true,
    disabled: true,
  },
  {
    title: 'Network Utilities',
    href: 'network-utilities',
    isActive: false,
    disabled: false,
  },
]