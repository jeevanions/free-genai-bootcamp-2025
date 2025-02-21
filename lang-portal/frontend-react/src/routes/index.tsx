import { createRootRoute, createRoute, createRouter } from "@tanstack/react-router"
import { RootLayout } from "@/components/layout/root-layout"
import { QueryClient } from "@tanstack/react-query"
import { DashboardPage } from "@/pages/dashboard"
import { StudyActivitiesPage } from "@/pages/study-activities"
import { WordsPage } from "@/pages/words"
import { WordDetailsPage } from "@/pages/word-details"
import { GroupsPage } from "@/pages/groups"
import { StudySessionsPage } from "@/pages/study-sessions"
import { SettingsPage } from "@/pages/settings"

export const queryClient = new QueryClient()

const rootRoute = createRootRoute({
  component: RootLayout,
})

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: DashboardPage,
})

const studyActivitiesRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "study-activities",
  component: StudyActivitiesPage,
})

const studyActivityRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "study-activities/$activityId",
  component: () => <div>Study Activity</div>,
})

const wordsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "words",
  component: WordsPage,
})

const groupsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "groups",
  component: GroupsPage,
})

const groupRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "groups/$groupId",
  component: () => <div>Group Details</div>,
})

const studySessionsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "study-sessions",
  component: StudySessionsPage,
})

const studySessionRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "study-sessions/$sessionId",
  component: () => <div>Study Session Details</div>,
})

const settingsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "settings",
  component: SettingsPage,
})

const wordDetailsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "words/$wordId",
  component: WordDetailsPage,
})

const routeTree = rootRoute.addChildren([
  indexRoute,
  studyActivitiesRoute,
  studyActivityRoute,
  wordsRoute,
  groupsRoute,
  groupRoute,
  studySessionsRoute,
  studySessionRoute,
  settingsRoute,
  wordDetailsRoute,
])

export const router = createRouter({
  routeTree,
  defaultPreload: "intent",
})