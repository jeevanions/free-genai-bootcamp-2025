import { Outlet } from "@tanstack/react-router"
import { Header } from "./header"

export function RootLayout() {
  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="container py-8">
        <Outlet />
      </main>
    </div>
  )
}