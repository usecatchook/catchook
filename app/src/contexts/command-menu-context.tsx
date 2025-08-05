"use client"

import { usePathname } from "next/navigation"
import React, { createContext, useContext, useState } from "react"

interface CommandMenuContextType {
  isOpen: boolean
  openCommandMenu: () => void
  closeCommandMenu: () => void
  toggleCommandMenu: () => void
  isEnabled: boolean
}

const CommandMenuContext = createContext<CommandMenuContextType | undefined>(undefined)

export function CommandMenuProvider({ children }: { children: React.ReactNode }) {
  const [isOpen, setIsOpen] = useState(false)
  const pathname = usePathname()
  
  const isEnabled = !pathname.startsWith('/login') && !pathname.startsWith('/setup')

  const openCommandMenu = () => {
    if (isEnabled) {
      setIsOpen(true)
    }
  }
  
  const closeCommandMenu = () => setIsOpen(false)
  
  const toggleCommandMenu = () => {
    if (isEnabled) {
      setIsOpen(!isOpen)
    }
  }

  return (
    <CommandMenuContext.Provider
      value={{
        isOpen,
        openCommandMenu,
        closeCommandMenu,
        toggleCommandMenu,
        isEnabled,
      }}
    >
      {children}
    </CommandMenuContext.Provider>
  )
}

export function useCommandMenu() {
  const context = useContext(CommandMenuContext)
  if (context === undefined) {
    throw new Error("useCommandMenu must be used within a CommandMenuProvider")
  }
  return context
} 