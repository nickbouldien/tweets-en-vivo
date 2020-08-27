export interface Tweet {
  id: string;
  text: string;
  tag?: string;
}

export interface MainState {
  error: Error | null;
  tweets: Tweet[];
  websocket: Websocket | null;
}
