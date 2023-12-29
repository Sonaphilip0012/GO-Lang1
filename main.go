package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Define structs to unmarshal data from API endpoints
type Comment struct {
	PostID int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Body   string `json:"body"`
}

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// Combine data from comments, posts, and users
func combineData() ([]map[string]interface{}, error) {
	// Fetch data from endpoints
	comments, err := fetchData("https://jsonplaceholder.typicode.com/comments")
	if err != nil {
		return nil, err
	}

	posts, err := fetchData("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		return nil, err
	}

	users, err := fetchData("https://jsonplaceholder.typicode.com/users")
	if err != nil {
		return nil, err
	}

	// Combine data into desired format
	combinedData := make([]map[string]interface{}, 0)

	for _, comment := range comments {

		comment := &Comment{
			PostID: int(comment["postId"].(float64)),
			ID:     int(comment["id"].(float64)),
			Name:   comment["name"].(string),
			Body:   comment["body"].(string),
		}

		post := findPostByID(posts, comment.PostID)
		if post != nil {
			user := findUserByID(users, post.UserID)
			combined := map[string]interface{}{
				"postId":        post.ID,
				"postName":      post.Title,
				"commentsCount": len(getCommentsForPost(comments, post.ID)),
				"userName":      user.Username,
				"body":          comment.Body,
			}
			combinedData = append(combinedData, combined)
		}
	}

	return combinedData, nil
}

// Fetch data from an API endpoint and unmarshal into corresponding structs
func fetchData(url string) ([]map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// Find post by ID from the posts slice
func findPostByID(posts []map[string]interface{}, id int) *Post {
	for _, p := range posts {
		if p["id"].(float64) == float64(id) {
			return &Post{
				UserID: int(p["userId"].(float64)),
				ID:     int(p["id"].(float64)),
				Title:  p["title"].(string),
				Body:   p["body"].(string),
			}
		}
	}
	return nil
}

// Find user by ID from the users slice
func findUserByID(users []map[string]interface{}, id int) *User {
	for _, u := range users {
		if u["id"].(float64) == float64(id) {
			return &User{
				ID:       int(u["id"].(float64)),
				Username: u["username"].(string),
			}
		}
	}
	return nil
}

// Get comments for a particular post
func getCommentsForPost(comments []map[string]interface{}, postID int) []Comment {
	var postComments []Comment
	for _, c := range comments {
		if c["postId"].(float64) == float64(postID) {
			comment := Comment{
				PostID: int(c["postId"].(float64)),
				ID:     int(c["id"].(float64)),
				Name:   c["name"].(string),
				Body:   c["body"].(string),
			}
			postComments = append(postComments, comment)
		}
	}
	return postComments
}

// Handler function to serve the combined data
func combinedDataHandler(w http.ResponseWriter, r *http.Request) {
	combined, err := combineData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal combined data to JSON
	response, err := json.Marshal(combined)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func main() {
	// Define route
	http.HandleFunc("/combinedData", combinedDataHandler)

	// Start the server
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
