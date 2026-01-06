package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Post struct {
	ID int64 `json:"id"`

	Title            string   `json:"title"`
	Description      string   `json:"description"`
	DescriptionImage string   `json:"descriptionImage"`
	Tags             []string `json:"tags"`

	Path       string `json:"path"`
	GithubPath string `json:"githubPath"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var (
	// GITHUB_REPO_URL = os.Getenv("GITHUB_REPO_URL")
	GITHUB_REPO_URL = "https://raw.githubusercontent.com/nexfortisme/content/refs/heads/main/"
	ROOT            = os.Getenv("CONTENT_ROOT")

	TITLE_PREFIX             = "title: "
	DESCRIPTION_PREFIX       = "description: "
	DESCRIPTION_IMAGE_PREFIX = "descriptionImage: "
	TAGS_PREFIX              = "tags: "

	POST_ID_COUNTER int64 = 1

	POSTS     []Post   = []Post{}
	POST_TAGS []string = []string{}
)

func main() {
	fmt.Println("Starting dir walk")

	if ROOT == "" {
		fmt.Println("CONTENT_ROOT is not set")
		os.Exit(1)
	}

	if GITHUB_REPO_URL == "" {
		fmt.Println("GITHUB_REPO_URL is not set")
		os.Exit(1)
	}

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

	// Sort posts by CreatedAt (newest first)
	sort.Slice(POSTS, func(i, j int) bool {
		return POSTS[i].CreatedAt.After(POSTS[j].CreatedAt)
	})

	// Update IDs to be sequential after sorting
	for i := range POSTS {
		POSTS[i].ID = int64(i + 1)
	}

	indexFile, err := os.Create("./index.json")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer indexFile.Close()

	postIndexFile, err := os.Create("./tag_index.json")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	err = json.NewEncoder(indexFile).Encode(POSTS)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("Done Writing Post Index File")

	err = json.NewEncoder(postIndexFile).Encode(POST_TAGS)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("Done Writing Tag Index File")
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
	post.GithubPath = GITHUB_REPO_URL + path
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
			tagList := strings.Split(tags, ",")
			// Trim whitespace from each tag
			for i, tag := range tagList {
				tagList[i] = strings.TrimSpace(tag)
			}
			updateGlobalTagList(tagList)
			post.Tags = tagList
		}
	}

	POST_ID_COUNTER++
	return post
}

func updateGlobalTagList(newTags []string) {
	// Add unique newTags to POST_TAGS
	for _, tag := range newTags {
		exists := false
		for _, existingTag := range POST_TAGS {
			if tag == existingTag {
				exists = true
				break
			}
		}
		if !exists && tag != "" {
			POST_TAGS = append(POST_TAGS, tag)
		}
	}
}
