package api

// Discussion summary information

type Discussion struct {
    Title    string
    URL      string
    Author   string
    Comments int
}

// DiscussionDetail holds detail of a discussion

type DiscussionDetail struct {
    Title  string
    Author string
    Body   string
    URL    string
}

