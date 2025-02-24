export interface Group {
  id: number;
  name: string;
  word_count: number;
}

export interface PaginationInfo {
  current_page: number;
  items_per_page: number;
  total_items: number;
  total_pages: number;
}

export interface GroupsResponse {
  items: Group[];
  pagination: PaginationInfo;
}

export interface WordParts {
  gender: string;
  plural: string;
  type: string;
}

export interface Word {
  id: number;
  italian: string;
  english: string;
  parts: WordParts;
  correct_count: number;
  wrong_count: number;
}

export interface GenerateWordsRequest {
  category: string;
}

export interface GenerateWordsResponse {
  words: Word[];
}

export interface ImportWordsResponse {
  imported_count: number;
}