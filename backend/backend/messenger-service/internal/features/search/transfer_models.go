package search

type GlobalSearchResponse struct {
	Users    []SearchResult `json:"users"`
	Chats    []SearchResult `json:"chats"`
	Messages []SearchResult `json:"messages"`
}

type ChatSearchResponse struct {
	Messages []SearchResult `json:"messages"`
	Total    int            `json:"total"`
}