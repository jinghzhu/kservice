# kservice

Welcome to kservice, which is a simple api interface to run the command in [Kubernetes](https://kubernetes.io/). 


## How to contribute
One of the most effective ways to collaborate on GitHub is by using a forking/branching model as described in the [Pull Request](https://docs.github.com/en/free-pro-team@latest/github/collaborating-with-issues-and-pull-requests/proposing-changes-to-your-work-with-pull-requests):
1. [Fork](https://docs.github.com/en/free-pro-team@latest/github/collaborating-with-issues-and-pull-requests/about-forks) the main repository to your personal GitHub space.
2. Clone this new fork locally to your computer. Make sure you use the SSH URL, not the HTTPS URL. This will be your origin remote.
3. Add an upstream remote whose URL is the SSH URL of the main repository - `git remote add upstream {{url}}`, replacing `{url}}` with the main repo's URL.
4. Each time you begin doing work on a new story, check out the master branch by doing `git checkout master`. You will only be able to do this if you don't have any changes in your local codebase.
5. Pull in the latest changes from upstream's master branch - `git pull upstream master`.
6. Create a new [feature branch](https://docs.github.com/en/free-pro-team@latest/github/getting-started-with-github/github-glossary#feature-branch), named something relevant to the story being worked on - `git checkout -b {{branch-name}}`, replacing `{{branch-name}}` with the name of your branch.
7. Push your new branch to your origin remote - `git push -u origin {{branch-name}}`.
8. Add your commits and push to that branch - `git push origin {{branch-name}}`.
9. Issue a Pull Request in to the upstream repository when the work is done. 
10. Once the Pull Request is merged, delete the local and remote branch you worked on - `git branch -d {{branch-name}}` for local, `git push origin :{{branch-name}}` for remote. **Please avoid to reuse a branch after it has been merged.**