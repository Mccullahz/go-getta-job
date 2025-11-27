# Go Getta Job
![GGJ Preview](https://github.com/Mccullahz/go-getta-job/blob/main/demo.gif)
- This project is a Go powered TUI application that finds nearby business websites and searches them for career or job listing pages. The goal is to automate localized job hunting by surfacing hiring pages often buried in small business websites. By searching only these local businesses, applicants can find job opportunities that may not be listed on larger job boards, reducing resume traffic and potentially aiding in landing a desired position.
- To keep up with the developement of this project, please visit [the devlog for this project](https://mccullahz.github.io/#/articles/job-scraper-cli).

# Features
- Sleek terminal interface using [Bubbletea](https://github.com/charmbracelet/bubbletea) + [Lipgloss](https://github.com/charmbracelet/lipgloss)
- Scraping + Geo locational via ZIP built with Go's Standard Libraries + [Overpass API](https://wiki.openstreetmap.org/wiki/Overpass_API) + [Zippopotam.us](http://www.zippopotam.us/)

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

## Using Docker
- To run the application with persistent data storage you will need to use Docker to setup the containerized Mongo database., you can do this by using the provided `docker-compose.yml` file.
- With Docker installed, you will need to setup your environment variables. You can do this by creating a `.env` file in the root directory of the project. Here is an example of the file content:

```bash
# .env

TZ=America/New_York

# mongo root user
MONGO_INITDB_ROOT_USERNAME=mongo_username
MONGO_INITDB_ROOT_PASSWORD=mongo_password
MONGO_DB_NAME=job_search_db

# mongo URI
MONGODB_URI=mongodb://mongo_username:mongo_password@mongodb:27017/job_search_db?authSource=admin

# enable database mode (set to "true" to use mongo instead of files)
USE_DATABASE=true

# express
ME_CONFIG_MONGODB_ADMINUSERNAME=express_admin_username
ME_CONFIG_MONGODB_ADMINPASSWORD=express_admin_password
ME_CONFIG_MONGODB_SERVER=mongodb
ME_CONFIG_MONGODB_PORT=27017
ME_CONFIG_BASICAUTH_USERNAME=express_user_username
ME_CONFIG_BASICAUTH_PASSWORD=express_user_password
```
- Once your `.env` file is set up, you can start the Docker containers using Docker Compose:
  ```bash
  docker compose up -d
  ```

- After the containers are running, you can execute the Go Getta Job application with semi-persistent database support:
  ```bash
  go run ./cmd/tui
  ```
