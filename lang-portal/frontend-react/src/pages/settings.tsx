import { useMutation } from "@tanstack/react-query"
import { resetHistory, fullReset } from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Link } from "@tanstack/react-router"
import { AlertCircle, ArrowLeft, RotateCcw, Trash2 } from "lucide-react"
import { useState } from "react"

export function SettingsPage() {
  const [showConfirmReset, setShowConfirmReset] = useState(false)
  const [showConfirmFullReset, setShowConfirmFullReset] = useState(false)

  const resetHistoryMutation = useMutation({
    mutationFn: resetHistory,
    onSuccess: () => {
      setShowConfirmReset(false)
    },
  })

  const fullResetMutation = useMutation({
    mutationFn: fullReset,
    onSuccess: () => {
      setShowConfirmFullReset(false)
    },
  })

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold">Settings</h1>
          <p className="text-muted-foreground">
            Manage your application settings and data
          </p>
        </div>
        <Button variant="outline" asChild>
          <Link to="/">
            <ArrowLeft className="mr-2 h-4 w-4" /> Back to Dashboard
          </Link>
        </Button>
      </div>

      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Reset Study History</CardTitle>
            <CardDescription>
              Clear all your study session history and progress data. This action cannot be undone.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {showConfirmReset ? (
              <div className="space-y-4">
                <div className="flex items-start gap-3 p-4 bg-destructive/10 rounded-lg text-destructive">
                  <AlertCircle className="h-5 w-5 mt-0.5" />
                  <div className="space-y-1">
                    <p className="font-medium">Are you sure you want to reset your study history?</p>
                    <p className="text-sm text-destructive/80">
                      This will delete all your study sessions and progress data. Your word groups will remain unchanged.
                    </p>
                  </div>
                </div>
                <div className="flex gap-3">
                  <Button
                    variant="destructive"
                    onClick={() => resetHistoryMutation.mutate()}
                    disabled={resetHistoryMutation.isPending}
                  >
                    {resetHistoryMutation.isPending ? (
                      "Resetting..."
                    ) : (
                      <>
                        <RotateCcw className="mr-2 h-4 w-4" />
                        Confirm Reset
                      </>
                    )}
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => setShowConfirmReset(false)}
                    disabled={resetHistoryMutation.isPending}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            ) : (
              <Button
                variant="destructive"
                onClick={() => setShowConfirmReset(true)}
                className="w-full sm:w-auto"
              >
                <RotateCcw className="mr-2 h-4 w-4" />
                Reset Study History
              </Button>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Full Reset</CardTitle>
            <CardDescription>
              Reset the application to its initial state. This will delete all data including word groups.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {showConfirmFullReset ? (
              <div className="space-y-4">
                <div className="flex items-start gap-3 p-4 bg-destructive/10 rounded-lg text-destructive">
                  <AlertCircle className="h-5 w-5 mt-0.5" />
                  <div className="space-y-1">
                    <p className="font-medium">Are you sure you want to perform a full reset?</p>
                    <p className="text-sm text-destructive/80">
                      This will delete ALL data including word groups, study sessions, and progress data.
                      This action cannot be undone.
                    </p>
                  </div>
                </div>
                <div className="flex gap-3">
                  <Button
                    variant="destructive"
                    onClick={() => fullResetMutation.mutate()}
                    disabled={fullResetMutation.isPending}
                  >
                    {fullResetMutation.isPending ? (
                      "Resetting..."
                    ) : (
                      <>
                        <Trash2 className="mr-2 h-4 w-4" />
                        Confirm Full Reset
                      </>
                    )}
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => setShowConfirmFullReset(false)}
                    disabled={fullResetMutation.isPending}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            ) : (
              <Button
                variant="destructive"
                onClick={() => setShowConfirmFullReset(true)}
                className="w-full sm:w-auto"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Full Reset
              </Button>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}