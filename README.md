# Go API wrapper for Bobuild apps

This simple module helps you interact with API of apps built with Bobuild. More information is available in the official documentation at https://docs.bobuild.com/advanced/api/

## Installation

```bash
go get github.com/bobuild/bobuild-client-go
```

## Usage

```go
import "github.com/bobuild/bobuild-client-go"

// Initialize the client
client := bobuild.NewClient("yourapp.bobuild.com", "your-api-token")

// Define your record type
type Customer struct {
    Name  string `json:"name"`,
    Email string `json:"email"`,
}

customers, err := bobuild.GetList[Customer](client, "/customers")
if err != nil {
    log.Fatal(err)
}

```
