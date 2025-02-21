export interface PaginatedResponse<T> {
  items: T[];
  pagination: {
    current_page: number;
    total_pages: number;
    total_items: number;
    per_page: number;
  };
}

export interface Group {
  id: number;
  name: string;
  description?: string;
  word_count: number;
  created_at: string;
  updated_at: string;
}

export interface StudySession {
  id: number;
  activity_name: string;
  group_name: string;
  review_items_count: number;
  correct_count: number;
  wrong_count: number;
  start_time: string;
  end_time?: string;
  duration_seconds?: number;
}
