"use client"

import Image from "next/image"
import Link from "next/link"
import { ThemeToggleDropdown } from "./theme-toggle-dropdown"

interface NavbarProps {
  showLogo?: boolean
  showToggle?: boolean
  className?: string
}

export function Navbar({ 
  showLogo = true, 
  showToggle = true, 
  className = "" 
}: NavbarProps) {
  return (
    <nav className={`flex items-center justify-between p-4 ${className}`}>
      {showLogo && (
        <Link
          href="/"
          className="flex items-center gap-2 font-jersey font-medium text-xl"
        >
          {/* Logo that changes with theme */}
          <Image 
            src="/logo_black.svg" 
            alt="Catchook" 
            width={24} 
            height={24}
            className="dark:hidden"
          />
          <Image 
            src="/logo.svg" 
            alt="Catchook" 
            width={24} 
            height={24}
            className="hidden dark:block"
          />
          CATCHOOK
        </Link>
      )}

      {showToggle && (
        <div className="flex items-center gap-4">
          <ThemeToggleDropdown />
        </div>
      )}
    </nav>
  )
}