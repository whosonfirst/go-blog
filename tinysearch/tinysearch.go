// package tinysearch provides methods for generating tinysearch indexes from a collection of blog (Markdown) posts.
package tinysearch

// Record is a struct representing the data structure used to index context in tinysearch.
type Record struct {
	// Title is the title of the document being indexed.
	Title string `json:"title"`
	// URL is the URL of the document being indexed.
	URL string `json:"url"`
	// Body is the body of the document (the "content") being indexed.
	Body string `json:"body"`
}
