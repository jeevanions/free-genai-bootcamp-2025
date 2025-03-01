# Frontend Technical Spec

## Business Goal
A language learning school wants to build a prototype of learning portal which will act as three things:
Inventory of possible vocabulary that can be learned
Act as a record store, providing correct and wrong score on practice vocabulary
A unified launchpad to launch different learning apps

You have been tasked with creating the frontend API of the application.
The fractional CTO has made strong recommendation that you settle on a frontend stack that is being commonly adopted and optimized for AI prototyping services.

The frontend application should powered by your backend API.


## Technical Restrictions

- The technical stack should be the following:
- Typescript (statically typed javascript)
- Tailwind CSS (css framework)
- Vite.js (frontend tool)
- ShadCN (UI components)


## Pages

### Dashboard `/dashboard`

#### Purpose
The purpose of this page is to provide a summary of Italian language learning progress
and act as the default page when a user visits the web-app.

#### Components
- Last Study Session
    - shows last activity used
    - shows when last activity used
    - summarizes wrong vs correct from last activity
    - has a link to the group
- Study Progress
    - total words studied eg. 150/5000
        - across all study sessions show the total Italian words studied out of all possible words
    - display a mastery progress eg. 35%
- Quick Stats
    - success rate eg. 80%
    - total study sessions eg. 4
    - total active groups eg. 3
    - study streak eg. 4 days
- Start Studying Button
    - goes to study activities page

#### Needed API Endpoints
- GET /api/dashboard/last_study_session - Returns last activity details including name, timestamp, scores
- GET /api/dashboard/study_progress - Returns total words studied and available
- GET /api/dashboard/quick-stats - Returns success rate, session count, active groups
- GET /api/dashboard/streak - Returns current study streak
- GET /api/dashboard/mastery - Returns overall mastery percentage

### Study Activities Index `/study_activities`

#### Purpose
The purpose of this page is to show a collection
of Italian study activities with a thumbnail and its name, to either launch or view the study activity.

#### Components
- Study Activity Card
    - show a thumbnail of the study activity
    - the name of the study activity (e.g., "Basic Verb Conjugation", "Food Vocabulary")
    - a launch button to take us to the launch page
    - the view page to view more information about past study sessions

#### Needed API Endpoints
- GET /api/study_activities

### Study Activity Show `/study_activities/:id`

#### Purpose
The purpose of this page is to show the details of a study activity and its past study sessions.

#### Components
- Name of study activity
- Thumbnail of study activity
- Description of study activity
- Launch button
- Study Activities Paginated List
    - id
    - activity name
    - group name
    - start time
    - end time (inferred by the last word_review_item submitted)
    - number of review items 

#### Needed API Endpoints
- GET /api/activities - Returns list of activities with details
- GET /api/activities/categories - Returns available activity categories
- GET /api/activities/recommended - Returns personalized activity suggestions/:id
- GET /api/sessions/groups/:groupId

### Study Activities Launch `/study_activities/:id/launch`

#### Purpose
The purpose of this page is to launch an Italian study activity.

#### Components
- Name of study activity
- Launch form
    - select field for group (e.g., "Basic Verbs", "Food Vocabulary")
    - launch now button

#### Behaviour
After the form is submitted a new tab opens with the study activity based on its URL provided in the database.
The page will redirect to the study session show page after form submission.

#### Needed API Endpoints
- POST /api/sessions - Creates new study session
- GET /api/groups/compatible/:activityId - Returns groups compatible with activity

### Words Index `/words`

#### Purpose
The purpose of this page is to show all Italian words in our database.

#### Components
- Paginated Word List
    - Columns
        - Italian
        - English
        - Parts of Speech
        - Correct Count
        - Wrong Count
    - Pagination with 100 items per page
    - Clicking the Italian word will take us to the word show page

#### Needed API Endpoints
- GET /api/words

### Word Show `/words/:id`

#### Purpose
The purpose of this page is to show information about a specific Italian word.

#### Components
- Italian
- English
- Parts of Speech
- Parts (JSON data showing conjugations, gender, etc.)
- Study Statistics
    - Correct Count
    - Wrong Count
- Word Groups 
    - show a series of pills eg. tags
    - when group name is clicked it will take us to the group show page

#### Needed API Endpoints
- GET /api/words/:id
- GET /api/reviews/words/:wordId/stats - For word statistics

### Word Groups Index `/groups`

#### Purpose
The purpose of this page is to show a list of Italian word groups in our database.

#### Components
- Paginated Group List
    - Columns
        - Group Name
        - Word Count
        - Category (grammar/thematic/situational)
        - Difficulty Level
    - Clicking the group name will take us to the group show page

#### Needed API Endpoints
- GET /api/groups

### Group Show `/groups/:id`

#### Purpose
The purpose of this page is to show information about a specific group of Italian words.

#### Components
- Group Name
- Group Statistics 
    - Total Word Count
    - Difficulty Level
    - Category
- Words in Group (Paginated List of Words)
    - Should use the same component as the words index page
- Study Sessions (Paginated List of Study Sessions)
    - Should use the same component as the study sessions index page

#### Needed API Endpoints
- GET /api/groups/:id - Returns group details
- GET /api/words - Returns words filtered by group ID
- GET /api/sessions/groups/:groupId - Returns group study sessions
- GET /api/groups/:id/progress - Returns group learning progress
- POST /api/groups/:id/words - Adds words to group
- DELETE /api/groups/:id/words/:wordId - Removes word from group

### Study Sessions Index `/study_sessions`

#### Purpose
The purpose of this page is to show a list of Italian study sessions in our database.

#### Components
- Paginated Study Session List
    - Columns
        - Id
        - Activity Name
        - Group Name
        - Start Time
        - End Time
        - Number of Review Items
        - Success Rate
    - Clicking the study session id will take us to the study session show page

#### Needed API Endpoints
- GET /api/study_sessions

### Study Session Show `/study_sessions/:id`

#### Purpose
The purpose of this page is to show information about a specific study session.

#### Components
- Study Session Details
    - Activity Name
    - Group Name
    - Start Time
    - End Time
    - Number of Review Items
    - Success Rate
- Word Review Items (Paginated List of Words)
    - Should use the same component as the words index page
    - Shows correct/incorrect status for each word

#### Needed API Endpoints
- GET /api/study_sessions/:id
- GET /api/study_sessions/:id/reviews

### Settings Page `/settings`

#### Purpose
The purpose of this page is to make configurations to the Italian study portal.

#### Components
- Theme Selection eg. Light, Dark, System Default
- Reset History Button
   - this will delete all study sessions and word review items
- Full Reset Button
   - this will drop all tables and re-create with seed data
- Language Level Selection
   - filter content based on CEFR levels (A1-C2)

#### Needed API Endpoints
- POST /api/reset 