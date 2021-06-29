package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"unicode"

	"github.com/joho/godotenv"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	query      = flag.String("query", "F0QXB5pw2qE", "")
	maxResults = flag.Int64("max-results", 25, "")
	jaComment  int64
)

func newClient() *http.Client {
	client := &http.Client{
		Transport: &transport.APIKey{Key: os.Getenv("YOUTUBE_API_KEY")},
	}
	return client
}

func newYoutubeService(client *http.Client) *youtube.Service {
	service, err := youtube.New(client)
	if err != nil {
		log.Fatal("Unable to create YouTube service: %v", err)
	}
	return service
}

func getComments() []*youtube.CommentThread {
	service := newYoutubeService(newClient())
	call := service.CommentThreads.List([]string{"id", "snippet"}).VideoId(*query).Order("time").MaxResults(*maxResults)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("%v", err)
	}
	return response.Items
}

func displayComments(comments []*youtube.CommentThread) {
	for _, comment := range comments {
		fmt.Println(comment.Snippet.TopLevelComment.Snippet.TextOriginal)
		fmt.Println(comment.Snippet.TopLevelComment.Snippet.LikeCount)
		fmt.Println("-------------")
	}
}

func checkJaorNot(comments []*youtube.CommentThread) {
	for _, comment := range comments {
		for _, r := range comment.Snippet.TopLevelComment.Snippet.TextOriginal {
			if unicode.In(r, unicode.Hiragana) || unicode.In(r, unicode.Katakana) {
				// fmt.Println(comment.Snippet.TopLevelComment.Snippet.TextOriginal)
				// fmt.Println(jaComment)
				jaComment++
				break
			}
		}
	}
	fmt.Println(jaComment)
}

func ratioOfJaCommnets() {
	ratio := (float64(jaComment) / float64(*maxResults)) * 100
	fmt.Println("-------------")
	fmt.Println(ratio, "% commnets were written in Japanese")
	fmt.Println("-------------")
}

func main() {
	flag.Parse()
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("%v", err)
	}
	comments := getComments()
	// displayComments(comments)
	checkJaorNot(comments)
	ratioOfJaCommnets()
}
