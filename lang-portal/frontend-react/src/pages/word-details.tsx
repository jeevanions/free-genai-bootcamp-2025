import { useQuery } from "@tanstack/react-query";
import { getWordById } from "@/services/word-service";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link, useParams } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";

export function WordDetailsPage() {
  const { wordId } = useParams({ from: "/words/$wordId" });

  const { data: word, isLoading } = useQuery({
    queryKey: ["word", wordId],
    queryFn: () => getWordById(parseInt(wordId)),
  });

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold">Word Details</h1>
          <p className="text-muted-foreground">
            Detailed information about the word
          </p>
        </div>
        <Button variant="outline" asChild>
          <Link to="/words">
            <ArrowLeft className="mr-2 h-4 w-4" /> Back to Words
          </Link>
        </Button>
      </div>

      {isLoading ? (
        <Card className="animate-pulse">
          <CardHeader>
            <div className="h-8 bg-muted rounded w-1/3"></div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="h-6 bg-muted rounded w-1/2"></div>
            <div className="h-4 bg-muted rounded w-3/4"></div>
            <div className="h-4 bg-muted rounded w-2/3"></div>
          </CardContent>
        </Card>
      ) : word ? (
        <Card>
          <CardHeader>
            <CardTitle className="text-3xl font-bold text-primary">
              {word.italian}
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div>
              <h2 className="text-xl font-semibold mb-2">Translation</h2>
              <p className="text-lg">{word.english}</p>
            </div>

            {word.parts && (
              <div>
                <h2 className="text-xl font-semibold mb-2">Grammar</h2>
                <div className="space-y-2">
                  <p><span className="font-medium">Type:</span> {word.parts.type}</p>
                  {word.parts.gender && (
                    <p><span className="font-medium">Gender:</span> {word.parts.gender}</p>
                  )}
                  {word.parts.plural && (
                    <p><span className="font-medium">Plural:</span> {word.parts.plural}</p>
                  )}
                  {word.parts.conjugation && (
                    <p><span className="font-medium">Conjugation:</span> {word.parts.conjugation}</p>
                  )}
                  {word.parts.irregular !== undefined && (
                    <p><span className="font-medium">Irregular:</span> {word.parts.irregular ? "Yes" : "No"}</p>
                  )}
                  {word.parts.usage && (
                    <p><span className="font-medium">Usage:</span> {word.parts.usage.join(", ")}</p>
                  )}
                </div>
              </div>
            )}

            <div>
              <h2 className="text-xl font-semibold mb-2">Study Statistics</h2>
              <div className="grid grid-cols-2 gap-4">
                <div className="p-4 bg-green-50 rounded-lg">
                  <p className="text-sm text-muted-foreground">Correct Answers</p>
                  <p className="text-2xl font-bold text-green-600">{word.correct_count}</p>
                </div>
                <div className="p-4 bg-red-50 rounded-lg">
                  <p className="text-sm text-muted-foreground">Wrong Answers</p>
                  <p className="text-2xl font-bold text-red-600">{word.wrong_count}</p>
                </div>
              </div>
              <p className="mt-4 text-center text-muted-foreground">
                Success Rate:{" "}
                {Math.round(
                  (word.correct_count / (word.correct_count + word.wrong_count || 1)) * 100
                )}%
              </p>
            </div>
          </CardContent>
        </Card>
      ) : (
        <div className="text-center py-12">
          <p className="text-lg text-muted-foreground">Word not found</p>
        </div>
      )}
    </div>
  );
}
