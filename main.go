package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v25/github"
	"golang.org/x/oauth2"
)

var (
	orgName       = flag.String("org", "", "Organization name")
	repoNames     = flag.String("repos", "", "Repository names. Can be multiple")
	branchNames   = flag.String("branches", "develop master", "Protected branch names. Can be multiple")
	operationName = flag.String("operation", "add", "Operation name [add, remove]")
	userNames     = flag.String("users", "", "User's login to change settings. Can be multiple")
)

func main() {
	flag.Parse()

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	if len(*orgName) == 0 {
		log.Fatal("org must be described!")
	}
	if len(*repoNames) == 0 {
		log.Fatal("repos must be described!")
	}
	if len(*branchNames) == 0 {
		log.Fatal("branches must be described!")
	}
	if len(*operationName) == 0 {
		log.Fatal("operation must be described!")
	}
	if len(*userNames) == 0 {
		log.Fatal("users must be described!")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_repos := strings.Split(*repoNames, " ")
	_branches := strings.Split(*branchNames, " ")
	_users := strings.Split(*userNames, " ")
	for _, repoName := range _repos {
		fmt.Println(repoName)

		branches, _, err := client.Repositories.ListBranches(ctx, *orgName, repoName, &github.ListOptions{PerPage: 100})
		if err == nil {
			for _, branch := range branches {

				if contains(_branches, *branch.Name) {
					fmt.Println("\t", *branch.Name)

					protection, _, err := client.Repositories.GetBranchProtection(ctx, *orgName, repoName, *branch.Name)
					if err == nil {
						preq := &github.ProtectionRequest{
							RequiredStatusChecks: protection.RequiredStatusChecks,
							EnforceAdmins:        protection.EnforceAdmins.Enabled,
							Restrictions: &github.BranchRestrictionsRequest{
								Teams: []string{},
								Users: []string{},
							},
						}

						if protection.RequiredPullRequestReviews != nil {
							preq.RequiredPullRequestReviews = &github.PullRequestReviewsEnforcementRequest{
								DismissStaleReviews:          protection.RequiredPullRequestReviews.DismissStaleReviews,
								RequireCodeOwnerReviews:      protection.RequiredPullRequestReviews.RequireCodeOwnerReviews,
								RequiredApprovingReviewCount: protection.RequiredPullRequestReviews.RequiredApprovingReviewCount,
							}
						}

						for _, v := range protection.Restrictions.Users {
							preq.Restrictions.Users = append(preq.Restrictions.Users, *v.Login)
						}

						if *operationName == "add" {
							for _, u := range _users {
								preq.Restrictions.Users = append(preq.Restrictions.Users, u)
							}
						} else if *operationName == "remove" {
							fmt.Println("Operation not supported yet: remove")
							return
						} else {
							fmt.Println("Operation not supported: ", *operationName)
							return
						}
						fmt.Println("\t\tpreq.Restrictions.Users ", preq.Restrictions.Users)

						_, _, err := client.Repositories.UpdateBranchProtection(ctx, *orgName, repoName, *branch.Name, preq)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println("\t\tDone!")
					}
				}
			}
		}
	}
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
