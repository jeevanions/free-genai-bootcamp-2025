export interface StudyActivity {
    id: number;
    name: string;
    description: string;
    thumbnail_url?: string;
}

export type ActivityType = 'flashcards' | 'quiz' | 'matching';

export const getActivityType = (activity: StudyActivity): ActivityType => {
    // Determine activity type based on activity name or other properties
    const name = activity.name.toLowerCase();
    if (name.includes('flash') || name.includes('card')) return 'flashcards';
    if (name.includes('quiz')) return 'quiz';
    if (name.includes('match')) return 'matching';
    return 'flashcards'; // default type
};
