package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v25/github"
	"github.com/jessevdk/go-flags"
	"golang.org/x/oauth2"
)

var appVersion = "v0.0.3"

var opts struct {
	Version       bool   `long:"version" description:"Show version"`
	OrgName       string `short:"o" long:"org" description:"Organization name"`
	RepoNames     string `short:"r" long:"repos" description:"Repository names. Can be multiple"`
	BranchNames   string `short:"b" long:"branches" default:"develop master" description:"Protected branch names. Can be multiple"`
	OperationName string `short:"p" long:"operation" default:"add" description:"Operation name [add, remove]"`
	UserNames     string `short:"u" long:"users" description:"User's login to change settings. Can be multiple"`
}

func main() {
	flags.Parse(&opts)
	if opts.Version {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present. Please, add GITHUB_TOKEN environment variable")
	}

	if len(opts.OrgName) == 0 {
		log.Fatal("org must be described!")
	}
	if len(opts.RepoNames) == 0 {
		log.Fatal("repos must be described!")
	}
	if len(opts.BranchNames) == 0 {
		log.Fatal("branches must be described!")
	}
	if len(opts.OperationName) == 0 {
		log.Fatal("operation must be described!")
	}
	if len(opts.UserNames) == 0 {
		log.Fatal("users must be described!")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_repos := strings.Split(opts.RepoNames, " ")
	_branches := strings.Split(opts.BranchNames, " ")
	_users := strings.Split(opts.UserNames, " ")
	for _, repoName := range _repos {
		fmt.Println(repoName)

		branches, _, err := client.Repositories.ListBranches(ctx, opts.OrgName, repoName, &github.ListOptions{PerPage: 100})
		if err == nil {
			for _, branch := range branches {

				if contains(_branches, *branch.Name) {
					fmt.Println("\t", *branch.Name)

					protection, _, err := client.Repositories.GetBranchProtection(ctx, opts.OrgName, repoName, *branch.Name)
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

						if opts.OperationName == "add" {
							for _, u := range _users {
								preq.Restrictions.Users = append(preq.Restrictions.Users, u)
							}
						} else if opts.OperationName == "remove" {
							fmt.Println("Operation not supported yet: remove")
							return
						} else {
							fmt.Println("Operation not supported: ", opts.OperationName)
							return
						}
						fmt.Println("\t\tpreq.Restrictions.Users ", preq.Restrictions.Users)

						_, _, err := client.Repositories.UpdateBranchProtection(ctx, opts.OrgName, repoName, *branch.Name, preq)
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
