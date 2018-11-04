package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("You should specify message")
	}

	msg := os.Args[1]

	githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		log.Fatalln("Cannot find GITHUB_TOKEN")
	}

	prURL, ok := os.LookupEnv("CIRCLE_PULL_REQUEST")
	if !ok {
		log.Fatalln("Cannot find CIRCLE_PULL_REQUEST")
	}

	owner, repo, prNo := parse(prURL)

	oauthClient := oauth2.NewClient(nil, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	))

	client := github.NewClient(oauthClient)

	if hasCommented(client, msg, owner, repo, prNo) {
		log.Println("Already notified")
		return
	}

	createComment(client, msg, owner, repo, prNo)
}

func parse(prURL string) (string, string, int) {
	re := regexp.MustCompile("https://github.com/([a-zA-Z0-9_-]+)/([a-zA-Z0-9_-]+)/pull/([0-9]+)")
	submatched := re.FindSubmatch([]byte(prURL))

	if len(submatched) == 0 {
		log.Fatalln("Not found submatches")
	}

	prNo, err := strconv.Atoi(string(submatched[3]))
	if err != nil {
		log.Fatalf("Failed to atoi : %v\n", err)
	}

	return string(submatched[1]), string(submatched[2]), prNo
}

func hasCommented(client *github.Client, msg, owner, repo string, prNo int) bool {
	var currentPage = 1
	ctx := context.Background()
	for {
		opt := &github.IssueListCommentsOptions{
			ListOptions: github.ListOptions{PerPage: 100, Page: currentPage},
		}
		comments, resp, err := client.Issues.ListComments(ctx, owner, repo, prNo, opt)
		if err != nil {
			log.Fatalf("failed to list comments")
		}

		for _, comment := range comments {
			if strings.Contains(*comment.Body, msg) {
				return true
			}
		}

		if resp.NextPage == 0 {
			break
		}
	}

	return false
}

func createComment(client *github.Client, msg, owner, repo string, prNo int) {
	comment := &github.IssueComment{Body: &msg}

	ctx := context.Background()
	if _, _, err := client.Issues.CreateComment(ctx, owner, repo, prNo, comment); err != nil {
		log.Fatalf("failed to post comment : %v\n", err)
	}
}
