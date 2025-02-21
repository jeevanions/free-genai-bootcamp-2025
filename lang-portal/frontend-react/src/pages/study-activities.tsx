import { useQuery } from "@tanstack/react-query"
import { getStudyActivities } from "@/lib/api"
import { StudyActivity, getActivityType } from "@/types/study-activity"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Link } from "@tanstack/react-router"
import { ArrowLeft, Play } from "lucide-react"

// Import activity images
import flashcardsImage from "@/assets/activity-images/flashcards.svg"
import quizImage from "@/assets/activity-images/quiz.svg"
import matchingImage from "@/assets/activity-images/matching.svg"

// Map activity types to their images
const activityImages: Record<string, string> = {
  flashcards: flashcardsImage,
  quiz: quizImage,
  matching: matchingImage,
}

export function StudyActivitiesPage() {
  const { data, isLoading } = useQuery({
    queryKey: ["studyActivities"],
    queryFn: getStudyActivities
  })

  const activities: StudyActivity[] = data?.items ?? []

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold">Study Activities</h1>
          <p className="text-muted-foreground">
            Choose an activity to start studying your vocabulary
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
                <div className="h-32 bg-muted rounded"></div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : activities.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-lg text-muted-foreground">No study activities available</p>
        </div>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {activities.map((activity) => (
            <Card key={activity.id} className="group overflow-hidden">
              <CardHeader>
                <CardTitle>{activity.name}</CardTitle>
                <CardDescription>{activity.description}</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="aspect-video relative rounded-md overflow-hidden">
                  <img
                    src={activityImages[getActivityType(activity)] || activity.thumbnail_url}
                    alt={activity.name}
                    className="object-cover w-full h-full transition-transform group-hover:scale-105"
                  />
                  <div className="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                    <Button size="lg" asChild>
                      <Link
                        to={`/study-activities/${activity.id}`}
                        className="gap-2"
                      >
                        <Play className="h-4 w-4" />
                        Start Activity
                      </Link>
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}