package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Post struct {
	ID int64 `json:"id"`

	Title            string   `json:"title"`
	Description      string   `json:"description"`
	DescriptionImage string   `json:"descriptionImage"`
	Tags             []string `json:"tags"`

	Path string `json:"path"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var (
	GITHUB_REPO_URL = os.Getenv("GITHUB_REPO_URL")
	ROOT            = os.Getenv("CONTENT_ROOT")

	TITLE_PREFIX             = "title: "
	DESCRIPTION_PREFIX       = "description: "
	DESCRIPTION_IMAGE_PREFIX = "descriptionImage: "
	TAGS_PREFIX              = "tags: "

	POST_ID_COUNTER int64 = 1

	POSTS []Post = []Post{}
)

func main() {
	fmt.Println("Starting dir walk")

	if ROOT == "" {
		fmt.Println("CONTENT_ROOT is not set")
		os.Exit(1)
	}

	// if GITHUB_REPO_URL == "" {
	// 	fmt.Println("GITHUB_REPO_URL is not set")
	// 	os.Exit(1)
	// }

	fmt.Println("Root: ", ROOT)

	err := filepath.WalkDir(ROOT, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// Handles permission errors, broken symlinks, etc.
			fmt.Println("Error: ", err)
			return err
		}

		if d.IsDir() {
			return nil // skip directories if desired
		}

		post := processFile(path)
		POSTS = append(POSTS, post)
		return nil
	})

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("Done Processing Files")

	indexFile, err := os.Create("../index.json")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer indexFile.Close()

	err = json.NewEncoder(indexFile).Encode(POSTS)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("Done Writing Index File")
}

func processFile(path string) Post {
	fmt.Println("Processing file: ", path)

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	post := Post{}
	post.ID = POST_ID_COUNTER
	post.Path = path
	post.CreatedAt = fileInfo.ModTime()
	// post.UpdatedAt = fileInfo.ModTime()

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, TITLE_PREFIX) {
			title := strings.TrimPrefix(line, TITLE_PREFIX)
			post.Title = title
		}
		if strings.HasPrefix(line, DESCRIPTION_PREFIX) {
			description := strings.TrimPrefix(line, DESCRIPTION_PREFIX)
			post.Description = description
		}
		if strings.HasPrefix(line, DESCRIPTION_IMAGE_PREFIX) {
			descriptionImage := strings.TrimPrefix(line, DESCRIPTION_IMAGE_PREFIX)
			post.DescriptionImage = descriptionImage
		}
		if strings.HasPrefix(line, TAGS_PREFIX) {
			tags := strings.TrimPrefix(line, TAGS_PREFIX)
			tags = strings.Trim(tags, "[")
			tags = strings.Trim(tags, "]")
			tags = strings.Trim(tags, "\"")
			tags = strings.Trim(tags, " ")
			post.Tags = strings.Split(tags, ",")
		}
	}

	POST_ID_COUNTER++
	return post
}
