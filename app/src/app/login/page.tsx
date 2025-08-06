"use client"
import { ApiStatus } from "@/components/api-status"
import { LoginForm } from "@/components/login-form"
import { DotPattern } from "@/components/magicui/dot-pattern"
import { ThemeToggleDropdown } from "@/components/theme-toggle-dropdown"
import Image from "next/image"
import Link from "next/link"

export default function LoginPage() {

  return (
    <div className="relative flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10 overflow-hidden">
      
      {/* Subtle Pattern Overlay */}
      <DotPattern className="absolute inset-0 opacity-40 dark:opacity-20" />
      
      {/* Floating Elements */}
      <div className="absolute top-20 left-10 w-32 h-32 rounded-full bg-primary/5 blur-xl animate-pulse" />
      <div className="absolute bottom-20 right-10 w-40 h-40 rounded-full bg-accent/10 blur-xl animate-pulse delay-1000" />
      <div className="absolute top-1/2 left-1/4 w-24 h-24 rounded-full bg-muted/20 blur-lg animate-pulse delay-500" />
      
      {/* Theme toggle in top right corner */}
      <div className="absolute top-6 right-6 z-20">
        <ThemeToggleDropdown />
      </div>
      
      <div className="relative z-20 flex w-full max-w-sm flex-col gap-8">
        <Link
          href="/"
          className="flex items-center gap-3 self-center font-jersey font-medium text-3xl tracking-wider"
        >
          {/* Logo that changes with theme */}
          <Image 
            priority={true}
            src="/logo_black.svg" 
            alt="Catchook" 
            width={32} 
            height={32}
            className="dark:hidden drop-shadow-sm"
          />
          <Image 
            priority={true}
            src="/logo.svg" 
            alt="Catchook" 
            width={32} 
            height={32}
            className="hidden dark:block drop-shadow-sm"
          />
          <span className="bg-gradient-to-r from-foreground to-foreground/80 bg-clip-text text-transparent">
            CATCHOOK
          </span>
        </Link>
        
        <ApiStatus />
        
        <LoginForm />
      </div>
    </div>
  )
}
