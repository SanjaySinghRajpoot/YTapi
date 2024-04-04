# YouTube Video Fetcher

YouTube Video Fetcher is a Go application that continuously fetches the latest videos from YouTube based on a predefined search query and stores the fetched video data in a PostgreSQL database. It also provides a GET API endpoint to retrieve the stored video data in a paginated response sorted in descending order of published datetime.

## Features

- Continuous fetching of latest videos from YouTube.
- Storage of video data including title, description, publishing datetime, thumbnails URLs, etc. in a PostgreSQL database.
- GET API endpoint to retrieve stored video data in a paginated response.

## Prerequisites

Before you begin, ensure you have the following installed on your machine:

- Docker Compose

## Setup

1. Clone the repository:

```bash
git clone https://github.com/SanjaySinghRajpoot/YTapi.git
cd YTapi
```

2. Update API keys:

   - Open `main.go` file.
   - Replace the placeholder API keys in the `keys` array with your actual YouTube API keys.
   - A few keyys are already added remove this string "/REMOVE" from last to get it working. 

3. Start the PostgreSQL database:

```bash
docker-compose up -d
```

4. In the root folder run `go run main.go`

## Usage

### Continuous Fetching

The application will start fetching videos from YouTube continuously in the background with an interval of 1 minute.

### GET API Endpoint

- **Endpoint**: `/videos`
- **Method**: GET
- **Query Parameters**:
  - `page` (optional): Page number for pagination.
  - `limit` (optional): Number of records per page.
- **Example**:

```bash
curl -X GET "http://localhost:8080/videos?page=1&limit=10"
```
