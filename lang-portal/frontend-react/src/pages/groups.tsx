import { useQuery } from "@tanstack/react-query"
import { getGroups } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Link } from "@tanstack/react-router"
import { ArrowLeft, Book } from "lucide-react"
import { useState } from "react"
import { Pagination } from "@/components/ui/pagination"

export function GroupsPage() {
  const [page, setPage] = useState(1)

  const { data, isLoading } = useQuery({
    queryKey: ["groups", page],
    queryFn: () => getGroups(page),
  })

  const groups = data?.items ?? []
  const totalPages = data?.pagination?.total_pages ?? 1

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold">Word Groups</h1>
          <p className="text-muted-foreground">
            Browse and study vocabulary word groups
          </p>
        </div>
        <Button variant="outline" asChild>
          <Link to="/">
            <ArrowLeft className="mr-2 h-4 w-4" /> Back to Dashboard
          </Link>
        </Button>
      </div>

      {isLoading ? (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {[...Array(6)].map((_, i) => (
            <Card key={i} className="animate-pulse">
              <CardHeader className="space-y-2">
                <div className="h-4 w-1/2 bg-muted rounded"></div>
                <div className="h-3 w-3/4 bg-muted rounded"></div>
              </CardHeader>
              <CardContent>
                <div className="h-20 bg-muted rounded"></div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : groups.length === 0 ? (
        <Card className="p-12 text-center">
          <div className="flex flex-col items-center gap-4">
            <Book className="h-12 w-12 text-muted-foreground" />
            <div className="space-y-2">
              <h3 className="text-xl font-semibold">No Word Groups Available</h3>
              <p className="text-muted-foreground">
                Please wait while your instructor adds vocabulary groups
              </p>
            </div>
          </div>
        </Card>
      ) : (
        <>
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {groups.map((group) => (
              <Link
                key={group.id}
                to={`/groups/${group.id}`}
                className="block group"
              >
                <Card className="h-full transition-all hover:shadow-md hover:border-primary/50">
                  <CardHeader>
                    <CardTitle className="flex items-center justify-between">
                      <span className="text-xl font-bold group-hover:text-primary transition-colors">
                        {group.name}
                      </span>
                      <span className="text-sm font-normal text-muted-foreground">
                        {group.word_count} words
                      </span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div className="flex items-center justify-between text-sm text-muted-foreground">
                        <span>View group details</span>
                        <ArrowLeft className="h-4 w-4 rotate-180 transition-transform group-hover:translate-x-1" />
                      </div>
                      <div className="h-1 w-full bg-muted rounded-full overflow-hidden">
                        <div
                          className="h-full bg-primary transition-all"
                          style={{
                            width: `${(group.word_count / 100) * 100}%`,
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