## DORA MATRIX BACKEND SERVER
```markdown
Dora Server is a lightweight deployment tracking API built with Go and MongoDB.  
It allows you to store, retrieve, and visualize deployment information such as commit data, build status, and deployment metrics.
```

## Table of Contents
```markdown

- [Features](#features)
- [File Structure](#file-structure)
- [Setup](#setup)
- [Environment Variables](#environment-variables)
- [Local Development](#local-development)
- [API Endpoints](#api-endpoints)
- [Example cURL Requests](#example-curl-requests)
- [Visualization](#visualization)
- [License](#license)

```


## Features
```markdown
- Store deployment details (commit hash, author, files changed, build info, etc.)
- Retrieve deployments with optional filtering
- Supports MongoDB for persistence
- Can be deployed locally or on Vercel
- Ready to integrate with Grafana for deployment metrics visualization

```

## File Structure
```markdown
Dora Server/
├─ go.mod
├─ go.sum
├─ pkg/                       # Core package with main handler
│     └─ handler.go
├─ local/                     # Local server entrypoint
│     └─ main.go
└─ api/                       # API endpoints for Vercel deployment
      └─ deployments.go

```

## Setup

### Environment Variables

Create a `.env` file in the root directory and define the following:

```

MONGODB\_URI=mongodb://localhost:27017
PORT=8080

````

You can use different values for local development or production.

### Local Development

1. Install dependencies:

```bash
go mod tidy
````

2. Run the server locally:

```bash
cd local
go run main.go
```

The server will start on `http://localhost:8080` (default port).

---

## API Endpoints

| Method | Endpoint       | Description                    |
| ------ | -------------- | ------------------------------ |
| GET    | `/`            | Welcome message                |
| GET    | `/deployments` | Retrieve all deployments       |
| POST   | `/deployments` | Create a new deployment record |

---

## Example cURL Requests

### Create a deployment

```bash
curl -X POST https://<your-server-url>/deployments \
-H "Content-Type: application/json" \
-d '{
  "commit_hash": "a538c7c9270f138e59b116688a3e8bcc16034d18",
  "commit_subject": "Merge pull request #4958",
  "commit_body": "Add deployment status",
  "commit_timestamp": "2025-08-16 14:45:43 +0600",
  "commit_author": "Aman",
  "commit_author_email": "aman@example.com",
  "release_version": "v1.0.0",
  "previous_commit": "02cc49e7d86f72aca1e8a5b95101b59e6bd8eb7f",
  "files_changed": "1",
  "lines_changed": "0",
  "jenkins_build_number": "3527",
  "jenkins_build_url": "https://jenkins.example.com/job/3527/",
  "deployment_status": "SUCCESS"
}'
```

### Retrieve all deployments

```bash
curl -X GET https://<your-server-url>/deployments
```

---

## Visualization

* You can connect Grafana to your MongoDB instance and create dashboards to visualize deployment metrics.
* Example aggregation query for Grafana:

```js
db.deployments.aggregate([
  { $sort: { inserted_at: -1 } },
  {
    $project: {
      _id: 0,
      commit_subject: 1,
      commit_body: 1,
      commit_timestamp: 1,
      release_version: 1,
      deployment_status: 1,
      inserted_at: 1
    }
  }
])
```

---

## License

MIT License © 2025

```
I can also make a **version with badges, setup for Vercel deployment, and detailed Docker instructions** if you want it to be more production-ready.  

Do you want me to make that enhanced version?
```
