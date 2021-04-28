# Cryptr with Go API

## 01 - Configuration

🛠️️ First we create a new directory for our new go project:

```bash
mkdir cryptr-go-api-sample
cd cryptr-go-api-sample
```

🛠️️ Now that we're in our new project directory, we can run `go mod init`:

```bash
go mod init cryptr.com/sample
```

To tell Go modules what the name of our module is, we use go mod init, with the fully qualified path to our module. We have a new file, called go.mod, that includes our module and the Go version we used. When we add imports to our Go code later, they'll also be added to this file.

🛠️️ Next, create a file `main.go`:

```bash
touch main.go
```

🛠️️ Open `main.go` file and add `package main` inside it

🛠️️ Now copy paste this structure for the project inside the `main.go`:

```go
type Teacher struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Course struct {
	Id        int      `json:"id"`
	User_id   string   `json:"user_id"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	Img       string   `json:"img"`
	Desc      string   `json:"desc"`
	Date      string   `json:"date"`
	Timestamp string   `json:"timestamp"`
	Teacher   Teacher  `json:"teacher"`
}
```

[Next](https://github.com/cryptr-examples/cryptr-laravel-api-sample/tree/02-validate-access-tokens)
