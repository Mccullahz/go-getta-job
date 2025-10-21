# Go Getta Job
- This project is a Go powered TUI application that finds nearby business websites and searches them for career or job listing pages. The goal is to automate localized job hunting by surfacing hiring pages often buried in small business websites. By searching only these local businesses, applicants can find job opportunities that may not be listed on larger job boards, reducing resume traffic and potentially aiding in landing a desired position.
- To keep up with the developement of this project, please visit [the devlog for this project](https://mccullahz.github.io/#/articles/job-scraper-cli).

# Features
- Sleek terminal interface using [Bubbletea()] + [Lipgloss()]
- Scraping + Geo locational via ZIP built with Go's Standard Libraries + [Overpass API]()]

# Usage
- This project is still in early development and not yet ready for public use, however, if you are so inclined, you can either download the binary from the releases tab or build the project yourself.

## Building with Go
- Ensure you have Go installed on your machine. You can download it from [the official Go website](https://golang.org/dl/).
- For full functionality, you will also need to have Docker installed. You can download it from [the official Docker website](https://www.docker.com/get-started).

- Clone the repository to your local machine:
  ```bash
  git clone https://github.com/Mccullahz/go-getta-job
  cd go-getta-job
  ```
- From here you can either run the project directly:
  ```bash
  go run ./cmd/tui
  ```
- Or build the binary for your operating system:
  ```bash
  go build -o go-getta-job ./cmd/tui
  ```
- After building, you can run the binary:
  ```bash
  ./go-getta-job
  ```
- Note: If you choose to run the project directly without building, ensure that you have all necessary dependencies installed. You can use Go modules to manage dependencies via:
  ```bash
  go mod tidy
  ```
