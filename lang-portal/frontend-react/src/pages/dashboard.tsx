import { useQuery } from "@tanstack/react-query"
import { getDashboardLastStudySession, getDashboardQuickStats, getDashboardStudyProgress } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { Button } from "@/components/ui/button"
import { Link } from "@tanstack/react-router"
import { ArrowRight, BookOpen, Calendar, Trophy, Users } from "lucide-react"
import { formatDistanceToNow, parseISO } from "date-fns"

export function DashboardPage() {
  const { data: lastSession } = useQuery({
    queryKey: ["lastStudySession"],
    queryFn: getDashboardLastStudySession
  })

  const { data: studyProgress } = useQuery({
    queryKey: ["studyProgress"],
    queryFn: getDashboardStudyProgress
  })

  const { data: quickStats } = useQuery({
    queryKey: ["quickStats"],
    queryFn: getDashboardQuickStats
  })

  const progressPercentage = studyProgress 
    ? (studyProgress.total_words_studied / studyProgress.total_available_words) * 100
    : 0

  const formattedDate = lastSession?.created_at 
    ? formatDistanceToNow(parseISO(lastSession.created_at), { addSuffix: true })
    : null

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <Button asChild>
          <Link to="/study-activities">
            Start Studying <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </div>

      {lastSession && lastSession.group_id && (
        <Card>
          <CardHeader>
            <CardTitle>Last Study Session</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {formattedDate && (
                <p className="text-sm text-muted-foreground">
                  {formattedDate}
                </p>
              )}
              <p className="font-medium">{lastSession.group_name}</p>
              <Link
                to={`/groups/${lastSession.group_id}`}
                className="text-sm text-primary hover:underline"
              >
                View Group
              </Link>
            </div>
          </CardContent>
        </Card>
      )}

      {studyProgress && (
        <Card>
          <CardHeader>
            <CardTitle>Study Progress</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between text-sm">
              <span>Words Studied: {studyProgress.total_words_studied}/{studyProgress.total_available_words}</span>
              <span>{Math.round(progressPercentage)}% Complete</span>
            </div>
            <Progress value={progressPercentage} />
          </CardContent>
        </Card>
      )}

      {quickStats && (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Success Rate</CardTitle>
              <Trophy className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{quickStats.success_rate}%</div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Study Sessions</CardTitle>
              <BookOpen className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{quickStats.total_study_sessions}</div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Active Groups</CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{quickStats.total_active_groups}</div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Study Streak</CardTitle>
              <Calendar className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{quickStats.study_streak_days} days</div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  )
}