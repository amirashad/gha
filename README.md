# gha
GitHub Access management tool for protected branches

You can use this tool to add someone to some branch as a merger.

## Installation
Download gha cli tool for OSX and set as executable: 
<pre>curl -L https://github.com/amirashad/gha/releases/download/v0.0.2/gha_darwin_amd64 -o /usr/local/bin/gha
chmod +x /usr/local/bin/gha</pre>

## Setup
To create Github personal token please follow: https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line

Add to your environment path: GITHUB_TOKEN={token-created-with-github-ui}

## Example
To give someone merge access: 
<pre>gha --org Some-Org --repos "repo1 repo2" --branches "develop master" --users "someuser1 someuser2" --operation add</pre>
