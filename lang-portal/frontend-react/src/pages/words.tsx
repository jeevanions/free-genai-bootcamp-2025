import { useQuery } from "@tanstack/react-query"
import { Word, getWords, PaginatedResponse } from "@/services/word-service"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Link } from "@tanstack/react-router"
import { ArrowLeft, Search } from "lucide-react"
import { useState, useEffect } from "react"
import { Pagination } from "@/components/ui/pagination"
import { useNavigate } from "@tanstack/react-router"

export function WordsPage() {
  const [page, setPage] = useState(1)
  const [searchTerm, setSearchTerm] = useState("")
  const [debouncedSearch, setDebouncedSearch] = useState(searchTerm)

  // Reset page when search term changes
  useEffect(() => {
    setPage(1)
  }, [debouncedSearch])

  // Debounce search term
  useEffect(() => {
    const timer = setTimeout(() => {
      // Only update search if term is not empty and at least 2 chars
      if (searchTerm === "" || searchTerm.length >= 2) {
        setDebouncedSearch(searchTerm)
      }
    }, 300)
    return () => clearTimeout(timer)
  }, [searchTerm])
  const navigate = useNavigate({ from: "/words" })

  const { data, isLoading, error } = useQuery<PaginatedResponse<Word>>({
    queryKey: ["words", page, debouncedSearch],
    queryFn: () => getWords(page, debouncedSearch),
    placeholderData: 'keepPrevious' // Keep previous data while fetching new data
  })

  const words = data?.items ?? []
  const totalPages = data?.pagination?.total_pages ?? 1

  const filteredWords = words

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold">Vocabulary Words</h1>
          <p className="text-muted-foreground">
            Browse and search your vocabulary words
          </p>
        </div>
        <Button variant="outline" asChild>
          <Link to="/">
            <ArrowLeft className="mr-2 h-4 w-4" /> Back to Dashboard
          </Link>
        </Button>
      </div>

      <div className="relative">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <input
          type="search"
          placeholder="Search words..."
          className="w-full pl-10 pr-4 py-2 border rounded-md"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>

      {error ? (
        <div className="text-center py-12">
          <p className="text-lg text-red-600">
            Error loading words. Please try again.
          </p>
        </div>
      ) : isLoading ? (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {[...Array(9)].map((_, i) => (
            <Card key={i} className="animate-pulse">
              <CardHeader>
                <div className="h-5 bg-muted rounded w-1/2"></div>
              </CardHeader>
              <CardContent className="space-y-2">
                <div className="h-4 bg-muted rounded w-3/4"></div>
                <div className="h-4 bg-muted rounded w-1/2"></div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : filteredWords.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-lg text-muted-foreground">
            {searchTerm
              ? "No words found matching your search"
              : "No vocabulary words available"}
          </p>
        </div>
      ) : (
        <>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {filteredWords.map((word) => (
              <Link key={word.id} to={`/words/${word.id}`}>
                <Card className="hover:shadow-md transition-shadow hover:border-primary cursor-pointer">
                  <CardHeader>
                    <CardTitle className="text-xl font-bold text-primary">
                      {word.italian}
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <p className="text-lg">{word.english}</p>
                      {word.parts && word.parts.length > 0 && (
                        <p className="text-sm text-muted-foreground mt-1">
                          {word.parts.join(", ")}
                        </p>
                      )}
                    </div>
                    <div className="flex items-center gap-4 text-sm">
                      <div className="flex items-center gap-1">
                        <span className="text-green-600">✓ {word.correct_count}</span>
                      </div>
                      <div className="flex items-center gap-1">
                        <span className="text-red-600">✗ {word.wrong_count}</span>
                      </div>
                      <div className="flex-1 text-right">
                        <span className="text-muted-foreground">
                          Success rate:{" "}
                          {Math.round(
                            (word.correct_count /
                              (word.correct_count + word.wrong_count || 1)) *
                              100
                          )}
                          %
                        </span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>

          {!searchTerm && (
            <Pagination
              currentPage={page}
              totalPages={totalPages}
              onPageChange={(newPage) => {
                setPage(newPage);
                refetch();
              }}
              className="mt-8"
            />
          )}
        </>
      )}
    </div>
  )
}