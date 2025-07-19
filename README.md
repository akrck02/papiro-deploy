# Papiro deploy github action
Deploy your markdown documents as a wiki website on github pages.

## How to install
Create a `publish-docs.yaml` file inside the `.github/workflows` directory with the following configuration:


```yaml
name: Publish documents
on:
  push:
    branches: [main]
permissions:
  contents: write
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: akrck02/papiro-deploy@main
        with:
          title: "title"              # your project title
          description: "description"  # your project description
          logo: "./logo.svg"          # your logo path (only svg for now)
          path: "./your-docs-path"    # the file with your markdown documents
          isObsidianProject: false    # if it is an obsidian project
          showFooter: true            # if the footer of the page must be shown
          showBreadcrumb: true        # if the breadcrumb of the page must be shown
          showStartPage: true         # if the start page must be shown
      - name: Deploy to github actions 🚀
        uses: JamesIves/github-pages-deploy-action@v4.3.0 # please checkout and give a star to this amazing action.
        with:
          branch: gh-pages            # The branch the action should deploy to.
          folder: .                   # The folder the action should deploy.
```
