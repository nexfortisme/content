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

	POSTS     []Post   = []Post{}
	POST_TAGS []string = []string{}

	// Map from file path to ID to preserve IDs across regenerations
	pathToIDMap map[string]int64 = make(map[string]int64)
	maxID       int64            = 0
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

	// Load existing index.json to preserve IDs
	loadExistingIndex()

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

	// IDs are already preserved from existing index or assigned during processing
	// No need to reassign them after sorting

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

func loadExistingIndex() {
	indexPath := "./index.json"
	file, err := os.Open(indexPath)
	if err != nil {
		// File doesn't exist yet, that's okay - we'll start fresh
		fmt.Println("No existing index.json found, starting fresh")
		return
	}
	defer file.Close()

	var existingPosts []Post
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&existingPosts)
	if err != nil {
		fmt.Println("Error reading existing index.json:", err)
		// Continue anyway - we'll start fresh
		return
	}

	fmt.Println("Loading existing IDs from index.json")
	for _, post := range existingPosts {
		pathToIDMap[post.Path] = post.ID
		if post.ID > maxID {
			maxID = post.ID
		}
		fmt.Println("  Loaded ID:", post.ID, "for path:", post.Path)
	}
	fmt.Println("Loaded", len(pathToIDMap), "existing post IDs")
}

func normalizePath(path string) string {
	// Convert absolute path to relative path from ROOT
	relPath, err := filepath.Rel(ROOT, path)
	if err != nil {
		// If relative path fails, try to clean up the path
		return filepath.Clean(path)
	}
	// Normalize path separators to forward slashes (for consistency)
	return filepath.ToSlash(relPath)
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

	// Normalize path to match format in index.json (relative to ROOT)
	normalizedPath := normalizePath(path)

	// Use existing ID if available, otherwise assign a new one
	if existingID, exists := pathToIDMap[normalizedPath]; exists {
		post.ID = existingID
		fmt.Println("  Using existing ID:", existingID, "for path:", normalizedPath)
	} else {
		maxID++
		post.ID = maxID
		pathToIDMap[normalizedPath] = maxID
		fmt.Println("  Assigning new ID:", maxID, "for path:", normalizedPath)
	}

	post.Path = normalizedPath
	post.GithubPath = GITHUB_REPO_URL + normalizedPath
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
