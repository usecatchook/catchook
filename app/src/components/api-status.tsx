import { Card } from "@/components/ui/card"
import { useHealth } from "@/hooks/useHealth"
import Link from "next/link"

export function ApiStatus() {
  const { health, isLoading, error } = useHealth()

  // Détermine le statut de l'API
  const getStatusInfo = () => {
    if (isLoading) {
      return {
        colorClass: 'bg-green-500/80 group-hover:bg-green-500',
        pingClass: 'bg-green-500/40',
        text: 'API OPERATIONAL',
        animate: 'animate-pulse'
      }
    }

    if (error) {
      return {
        colorClass: 'bg-red-500/80 group-hover:bg-red-500',
        pingClass: 'bg-red-500/40',
        text: 'API ERROR',
        animate: 'animate-pulse'
      }
    }

    if (!health) {
      return {
        colorClass: 'bg-gray-500/80 group-hover:bg-gray-500',
        pingClass: 'bg-gray-500/40',
        text: 'API UNKNOWN',
        animate: 'animate-pulse'
      }
    }

    if (health.status === 'ok' || health.status === 'healthy') {
      return {
        colorClass: 'bg-green-500/80 group-hover:bg-green-500',
        pingClass: 'bg-green-500/40',
        text: 'API OPERATIONAL',
        animate: 'animate-pulse'
      }
    }

    // Statut inconnu ou dégradé
    return {
      colorClass: 'bg-orange-500/80 group-hover:bg-orange-500',
      pingClass: 'bg-orange-500/40',
      text: 'API WARNING',
      animate: 'animate-pulse'
    }
  }

  const statusInfo = getStatusInfo()

  return (
    <Link href="/status" className="self-center">
      <Card className="group px-4 py-2  border-accent/20">
        <div className="flex items-center gap-3">
          <div className="relative">
            <div className={`w-2 h-2 rounded-full ${statusInfo.colorClass} ${statusInfo.animate}`} />
            <div className={`absolute inset-0 w-2 h-2 rounded-full ${statusInfo.pingClass} animate-ping`} />
          </div>
          <span className="text-xs font-medium tracking-wide text-muted-foreground/80 group-hover:text-accent-foreground transition-colors">
            {statusInfo.text}
          </span>
        </div>
      </Card>
    </Link>
  )
}