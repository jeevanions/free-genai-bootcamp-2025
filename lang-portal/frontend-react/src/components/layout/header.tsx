import { NavigationMenu, NavigationMenuItem, NavigationMenuLink, NavigationMenuList } from "@/components/ui/navigation-menu"
import { cn } from "@/lib/utils"
import { Link } from "@tanstack/react-router"
import { GraduationCap } from "lucide-react"
import { Button } from "@/components/ui/button"

const navigationItems = [
  { name: "Dashboard", path: "/" },
  { name: "Study Activities", path: "/study-activities" },
  { name: "Words", path: "/words" },
  { name: "Groups", path: "/groups" },
  { name: "Study Sessions", path: "/study-sessions" },
  { name: "Settings", path: "/settings" },
]

export function Header() {
  return (
    <header className="border-b">
      <div className="container flex h-16 items-center">
        <Link to="/" className="flex items-center space-x-2">
          <GraduationCap className="h-6 w-6" />
          <span className="text-xl font-bold">Language Learning Portal</span>
        </Link>
        <NavigationMenu className="ml-auto">
          <NavigationMenuList>
            {navigationItems.map((item) => (
              <NavigationMenuItem key={item.path}>
                <Button variant="ghost" className="h-auto px-3 py-2" asChild>
                  <Link to={item.path}>{item.name}</Link>
                </Button>
              </NavigationMenuItem>
            ))}
          </NavigationMenuList>
        </NavigationMenu>
      </div>
    </header>
  )
}