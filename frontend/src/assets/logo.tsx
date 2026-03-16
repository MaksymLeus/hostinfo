import { type ImgHTMLAttributes } from 'react'
import { cn } from '@/lib/utils'

export function Logo({ className, ...props }: ImgHTMLAttributes<HTMLImageElement>) {
  return (
    <img
      src="/images/logo.png"
      alt="Hostinfo Logo"
      className={cn('h-10 w-10 object-contain', className)}
      {...props}
    />
  )
}
