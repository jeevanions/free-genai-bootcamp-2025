import { useQuery } from "@tanstack/react-query";
import { Group, getGroupById } from "@/services/group-service";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link, useParams } from "@tanstack/react-router";
import { ArrowLeft, Book, CheckCircle2, XCircle } from "lucide-react";

export function GroupDetailsPage() {
  const { groupId } = useParams({ from: "/groups/$groupId" });

  const { data: group, isLoading, error } = useQuery<Group>({
    queryKey: ["group", groupId],
    queryFn: () => getGroupById(parseInt(groupId)),
  });

  if (error) {
    return (
      <div className="space-y-8">
        <div className="flex items-center">
          <Button variant="outline" asChild>
            <Link to="/groups">
              <ArrowLeft className="mr-2 h-4 w-4" /> Back to Groups
            </Link>
          </Button>
        </div>
        <Card className="p-12 text-center">
          <div className="flex flex-col items-center gap-4">
            <Book className="h-12 w-12 text-muted-foreground" />
            <div className="space-y-2">
              <h3 className="text-xl font-semibold">Error Loading Group</h3>
              <p className="text-muted-foreground">
                Failed to load group details. Please try again.
              </p>
            </div>
          </div>
        </Card>
      </div>
    );
  }

  if (isLoading || !group) {
    return (
      <div className="space-y-8">
        <div className="flex items-center">
          <Button variant="outline" asChild>
            <Link to="/groups">
              <ArrowLeft className="mr-2 h-4 w-4" /> Back to Groups
            </Link>
          </Button>
        </div>
        <div className="grid gap-6">
          <Card className="animate-pulse">
            <CardHeader className="space-y-2">
              <div className="h-6 w-1/3 bg-muted rounded"></div>
              <div className="h-4 w-2/3 bg-muted rounded"></div>
            </CardHeader>
          </Card>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
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
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <Button variant="outline" asChild>
              <Link to="/groups">
                <ArrowLeft className="mr-2 h-4 w-4" /> Back to Groups
              </Link>
            </Button>
          </div>
          <h1 className="text-3xl font-bold mt-4">{group.name}</h1>
          {group.description && (
            <p className="text-muted-foreground">{group.description}</p>
          )}
        </div>
      </div>

      {!group.words || group.words.length === 0 ? (
        <Card className="p-12 text-center">
          <div className="flex flex-col items-center gap-4">
            <Book className="h-12 w-12 text-muted-foreground" />
            <div className="space-y-2">
              <h3 className="text-xl font-semibold">No Words Available</h3>
              <p className="text-muted-foreground">
                This group doesn't have any words yet
              </p>
            </div>
          </div>
        </Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {group.words.map((word) => (
          <Link key={word.id} to={`/words/${word.id}`}>
            <Card className="hover:shadow-md transition-shadow hover:border-primary cursor-pointer h-full">
              <CardHeader>
                <CardTitle className="text-xl font-bold text-primary">
                  {word.italian}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <p className="text-lg">{word.english}</p>
                <div className="flex items-center gap-4 text-sm">
                  <div className="flex items-center gap-1">
                    <CheckCircle2 className="h-4 w-4 text-green-600" />
                    <span>{word.correct_count}</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <XCircle className="h-4 w-4 text-red-600" />
                    <span>{word.wrong_count}</span>
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
      )}
    </div>
  );
}
