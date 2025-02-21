import { QueryClient } from '@tanstack/react-query';
import { RootRoute, Route, Router } from '@tanstack/react-router';
import { DashboardPage } from './pages/dashboard';
import { WordsPage } from './pages/words';
import { WordDetailsPage } from './pages/word-details';
import { StudyActivitiesPage } from './pages/study-activities';
import { GroupsPage } from './pages/groups';
import { StudySessionsPage } from './pages/study-sessions';
import { RootLayout } from './components/layout/root-layout';

export const queryClient = new QueryClient();

const rootRoute = new RootRoute({
  component: RootLayout
});

const dashboardRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/',
  component: DashboardPage,
});

const wordsRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/words',
  component: WordsPage,
});

const wordDetailsRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/words/$wordId',
  component: WordDetailsPage,
});

const studyActivitiesRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/study-activities',
  component: StudyActivitiesPage,
});

const groupsRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/groups',
  component: GroupsPage,
});

const studySessionsRoute = new Route({
  getParentRoute: () => rootRoute,
  path: '/study-sessions',
  component: StudySessionsPage,
});

const routeTree = rootRoute.addChildren([
  dashboardRoute,
  wordsRoute,
  wordDetailsRoute,
  studyActivitiesRoute,
  groupsRoute,
  studySessionsRoute,
]);

export const router = new Router({ routeTree });
