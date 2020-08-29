import { ITweetResponse } from '@/store/main/state'

export const tweetResponse1: ITweetResponse = {
  data: {
    authorId: "618818200",
    createdAt: "2020-08-29T15:30:38.000Z",
    id: "1299731189654155264",
    text: "sample tweet"
  },
  includes: {
    users: [
      {
        id: "618818200",
        name: "test user",
        username: "testuser",
        createdAt: "2020-08-29T15:30:38.000Z"
      }
    ]
  },
  matching_rules: [
    {
      id: "1299375087406252034",
      tag: "#tag"
    }
  ]
}

export const tweetResponse2: ITweetResponse = {
  "data": {
    "authorId": "618818200",
    "createdAt": "2020-08-29T15:30:38.000Z",
    "id": "129973118965415526",
    "text": "other test tweet"
  },
  "includes": {
    "users": [
      {
        "id": "618818200",
        "name": "other user",
        "username": "otheruser",
        "createdAt": "2020-08-29T15:30:38.000Z"
      }
    ]
  },
  "matching_rules": [
    {
      "id": "1299375087406252035",
      "tag": "#test"
    }
  ]
}

export const tweetResponses: ITweetResponse[] = [tweetResponse1, tweetResponse2];
