import { Moon, Sun } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useTheme } from "@/providers/theme-provider"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

export function ThemeToggle() {
  const { theme, setTheme } = useTheme()

  return (
    <Select value={theme} onValueChange={setTheme}>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Select theme" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="light">
          <div className="flex items-center">
            <Sun className="h-4 w-4 mr-2" />
            Light
          </div>
        </SelectItem>
        <SelectItem value="dark">
          <div className="flex items-center">
            <Moon className="h-4 w-4 mr-2" />
            Dark
          </div>
        </SelectItem>
        <SelectItem value="system">
          <div className="flex items-center">
            <Sun className="h-4 w-4 mr-2" />
            System
          </div>
        </SelectItem>
      </SelectContent>
    </Select>
  )
}
