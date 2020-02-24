package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

// Data structures to hold JSON data
// this needs to match the fields in the JSON feed
// see sample data
/*
  {
    "postId": 100,
    "id": 499,
    "name": "excepturi sunt cum a et rerum quo voluptatibus quia",
    "email": "Wilburn_Labadie@araceli.name",
    "body": "et necessitatibus tempora ipsum quaerat inventore est quasi quidem\nea repudiandae laborum omnis ab reprehenderit ut\nratione sit numquam culpa a rem\natque aut et"
  },
{
    "postId": 100,
    "id": 499,
    "name": "excepturi sunt cum a et rerum quo voluptatibus quia",
    "email": "Wilburn_Labadie@araceli.name",
    "body": "et necessitatibus tempora ipsum quaerat inventore est quasi quidem\nea repudiandae laborum omnis ab reprehenderit ut\nratione sit numquam culpa a rem\natque aut et"
  },
*/
type Comment struct {
	PostId int    `json:"postId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

// SLice of WordCount
type WordCounts []WordCount

/*
	Data Structure of Word Count
*/
type WordCount struct {
	word  string
	count int
}

// wordCount will be displayed in this format
func (p WordCount) String() string {
	return fmt.Sprintf("%3d   %s", p.count, p.word)
}

// Len is part of sort.Interface.
func (wc WordCounts) Len() int {
	return len(wc)
}

// Swap is part of sort.Interface.
func (wc WordCounts) Less(i, j int) bool {
	return wc[i].count < wc[j].count
}

// Less is part of sort.Interface. We use count as the value to sort by
func (wc WordCounts) Swap(i, j int) {
	wc[i], wc[j] = wc[j], wc[i]
}

func main() {
	url := flag.String("url", "https://jsonplaceholder.typicode.com/comments", "fetch from this url")

	numOfWords := flag.Int("n", 4, "specify the number of word")

	flag.Parse()

	fmt.Println("Starting the application...")

	// create New HTTP Request
	req, err := http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		log.Fatalln("error creating request: ", err)
	}

	// Set Accept Header to application json
	req.Header.Set("Accept", "application/json")

	//  sends an HTTP request and returns an HTTP response, following
	// policy (such as redirects, cookies, auth) as configured on the
	// client.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("error fetching:", err)
	}

	// confirm we received an OK status
	if res.StatusCode != http.StatusOK {
		log.Fatalln("Error Status not OK:", res.Status)
	}

	// Process Response by counting each word
	wordCounts, err := ProcessComment(res.Body, *numOfWords)
	if err != nil {
		log.Fatalln( err )
	}

	fmt.Println("Words in text sorted by frequency low to high:")
	// Display WordCounts
	for _, counter := range wordCounts {
		fmt.Printf("%v\n", counter)
	}

	fmt.Println("Terminating the application...")
}

func ProcessComment (body io.ReadCloser, upperLimit int) (WordCounts, error)  {
	var comments []Comment

	// defer response close
	defer body.Close()

	// read the JSON-encoded value and decode into array of comments
	if  err := json.NewDecoder(body).Decode(&comments); err != nil {
		return nil, fmt.Errorf("unable to decode json: %v", err)
	}

	// Return array of words and their counts
	wordCounts := Counter(comments)

	return wordCounts[:upperLimit], nil

}

func Counter(comments []Comment) WordCounts {
	// Initialize WordCounts
	var wCs = WordCounts{}
	wordC := make(map[string]int)
	for _, comment := range comments {
		// change text to all lower characters
		text := strings.ToLower(comment.Body)
		// splits string into consecutive characters
		words := strings.Fields(text)

		for i := range words {
			wordC[words[i]]++
		}
	}
	for key, count := range wordC {
		c := WordCount{word: key, count: count}
		wCs = append(wCs, c)
	}
	sort.Sort(wCs)
	return wCs
}
