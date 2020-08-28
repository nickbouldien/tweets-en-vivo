export interface Tweet {
  id: string;
  text: string;
  tag?: string;
  authorId?: string;
  authorName?: string;
  authorUsername?: string;
}

export interface Rule {
  id: string;
  tag?: string;
}

export interface User {
  created_at: string;
  id: string;
  name: string;
  username: string;
}

// TODO - reformat this on the server
export interface IncludesSection {
  users: User[];
}

export interface TweetResponse {
  data: Tweet;
  includes: IncludesSection;
  matching_rules: Rule[];
}

export interface MainState {
  error: Error | null;
  tweets: Tweet[];
  websocket: Websocket | null;
}
