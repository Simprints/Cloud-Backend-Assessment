# Prerequisites

Install Golang. Using [GVM](https://github.com/moovweb/gvm) and the latest version of Golang is recommended

# Running

```bash
GCLOUD_PROJECT=your-gcp-project go run cmd/mywebapp/main.go
```

# Testing

```bash
GCLOUD_PROJECT=your-gcp-project go test ./... -count=1
```



