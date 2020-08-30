export interface ITweet {
  authorId?: string;
  authorName?: string;
  authorUsername?: string;
  createdAt?: string;
  id: string;
  matchingRules?: IRule[];
  text: string;
  tag?: string;
}

export interface IRule {
  id: string;
  tag?: string;
}

export interface IUser {
  createdAt: string;
  id: string;
  name: string;
  username: string;
}

export interface IIncludesSection {
  users: IUser[];
}

export interface ITweetResponse {
  data: ITweet;
  includes: IIncludesSection;
  matching_rules: IRule[];
}

export interface IMainState {
  error: Error | null;
  tweets: ITweet[];
  websocket: Websocket | null;
}
