import { useQuery } from "@tanstack/react-query"
import { getStudySessions } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Link } from "@tanstack/react-router"
import { ArrowLeft, Clock, History } from "lucide-react"
import { useState } from "react"
import { Pagination } from "@/components/ui/pagination"
import { formatDistanceToNow, parseISO } from "date-fns"

export function StudySessionsPage() {
  const [page, setPage] = useState(1)

  const { data, isLoading } = useQuery({
    queryKey: ["studySessions", page],
    queryFn: () => getStudySessions(page),
  })

  const sessions = data?.items ?? []
  const totalPages = data?.pagination?.total_pages ?? 1

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold">Study Sessions</h1>
          <p className="text-muted-foreground">
            Track your learning progress and review past study sessions
          </p>
        </div>
        <Button variant="outline" asChild>
          <Link to="/">
            <ArrowLeft className="mr-2 h-4 w-4" /> Back to Dashboard
          </Link>
        </Button>
      </div>

      {isLoading ? (
        <div className="grid gap-6 md:grid-cols-2">
          {[...Array(6)].map((_, i) => (
            <Card key={i} className="animate-pulse">
              <CardHeader className="space-y-2">
                <div className="h-4 w-1/2 bg-muted rounded"></div>
                <div className="h-3 w-3/4 bg-muted rounded"></div>
              </CardHeader>
              <CardContent>
                <div className="h-24 bg-muted rounded"></div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : sessions.length === 0 ? (
        <Card className="p-12 text-center">
          <div className="flex flex-col items-center gap-4">
            <History className="h-12 w-12 text-muted-foreground" />
            <div className="space-y-2">
              <h3 className="text-xl font-semibold">No Study Sessions Yet</h3>
              <p className="text-muted-foreground">
                Start a study activity to begin tracking your progress
              </p>
            </div>
            <Button asChild className="mt-4">
              <Link to="/study-activities">
                Start Studying
              </Link>
            </Button>
          </div>
        </Card>
      ) : (
        <>
          <div className="grid gap-6 md:grid-cols-2">
            {sessions.map((session) => (
              <Link
                key={session.id}
                to={`/study-sessions/${session.id}`}
                className="block group"
              >
                <Card className="h-full transition-all hover:shadow-md hover:border-primary/50">
                  <CardHeader>
                    <CardTitle className="flex items-center justify-between">
                      <span className="text-xl font-bold group-hover:text-primary transition-colors">
                        {session.activity_name}
                      </span>
                      <span className="text-sm font-normal text-muted-foreground">
                        {session.review_items_count} words
                      </span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div className="flex items-center justify-between">
                        <div className="space-y-1">
                          <p className="text-sm font-medium">{session.group_name}</p>
                          <div className="flex items-center text-sm text-muted-foreground">
                            <Clock className="mr-1 h-3 w-3" />
                            {formatDistanceToNow(parseISO(session.start_time), { addSuffix: true })}
                          </div>
                        </div>
                        <div className="flex items-center text-sm text-muted-foreground">
                          <span>View details</span>
                          <ArrowLeft className="ml-1 h-4 w-4 rotate-180 transition-transform group-hover:translate-x-1" />
                        </div>
                      </div>
                      <div className="h-1 w-full bg-muted rounded-full overflow-hidden">
                        <div
                          className="h-full bg-primary transition-all"
                          style={{
                            width: `${(session.review_items_count / 20) * 100}%`,
                          }}
                        />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>

          <Pagination
            currentPage={page}
            totalPages={totalPages}
            onPageChange={setPage}
            className="mt-8"
          />
        </>
      )}
    </div>
  )
}