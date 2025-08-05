"use client"

import { CommandDialog, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command"
import { useCommandMenu } from "@/contexts/command-menu-context"
import { IconCalculator, IconCalendar, IconDashboard, IconSearch, IconUsers } from "@tabler/icons-react"
import { useRouter } from "next/navigation"
import { useEffect } from "react"

interface Command {
  id: string
  title: string
  icon: React.ComponentType<{ className?: string }>
  action: () => void
  group: string
}

export function CommandMenu() {
  const { isOpen, openCommandMenu, closeCommandMenu, isEnabled } = useCommandMenu()
  const router = useRouter()

  const commands: Command[] = [
    // Navigation commands
    {
      id: "dashboard",
      title: "Dashboard",
      icon: IconDashboard,
      action: () => router.push("/"),
      group: "Navigation"
    },
    {
      id: "users",
      title: "Users",
      icon: IconUsers,
      action: () => router.push("/users"),
      group: "Navigation"
    },
    // Tool commands
    {
      id: "calendar",
      title: "Calendar",
      icon: IconCalendar,
      action: () => console.log("Calendar action"),
      group: "Tools"
    },
    {
      id: "search-emoji",
      title: "Search Emoji",
      icon: IconSearch,
      action: () => console.log("Search emoji action"),
      group: "Tools"
    },
    {
      id: "calculator",
      title: "Calculator",
      icon: IconCalculator,
      action: () => console.log("Calculator action"),
      group: "Tools"
    }
  ]

  const handleSelect = (commandId: string) => {
    const command = commands.find(cmd => cmd.id === commandId)
    if (command) {
      command.action()
      closeCommandMenu()
    }
  }

  const groupedCommands = commands.reduce((acc, command) => {
    if (!acc[command.group]) {
      acc[command.group] = []
    }
    acc[command.group].push(command)
    return acc
  }, {} as Record<string, Command[]>)

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault()
        if (isEnabled) {
          openCommandMenu()
        }
      }
    }
    document.addEventListener("keydown", down)
    return () => document.removeEventListener("keydown", down)
  }, [openCommandMenu, isEnabled])

  if (!isEnabled) {
    return null
  }

  return (
    <CommandDialog open={isOpen} onOpenChange={closeCommandMenu}>
      <CommandInput placeholder="Type a command or search..." />
      <CommandList>
        <CommandEmpty>No results found.</CommandEmpty>
        {Object.entries(groupedCommands).map(([groupName, groupCommands]) => (
          <CommandGroup key={groupName} heading={groupName}>
            {groupCommands.map((command) => (
              <CommandItem
                key={command.id}
                onSelect={() => handleSelect(command.id)}
              >
                <command.icon className="mr-2 h-4 w-4" />
                {command.title}
              </CommandItem>
            ))}
          </CommandGroup>
        ))}
      </CommandList>
    </CommandDialog>
  )
}