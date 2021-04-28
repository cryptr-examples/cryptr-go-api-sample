# Cryptr with Go API

## 01 - Validate access tokens

### Install dependencies

🛠️️ Run the following commands in your terminal:

```bash
go get -d github.com/form3tech-oss/jwt-go
go get -d github.com/codegangsta/negroni
go get -d github.com/gorilla/mux
```

Note:   
__- [form3tech-oss/jwt-go](https://github.com/form3tech-oss/jwt-go) to verify incoming JWTs__  
__- [codegangsta/negroni](https://github.com/urfave/negroni) for HTTP middleware__  
__- [gorilla/mux](https://github.com/gorilla/mux) to handle our routes__  

### Create simple rendering courses

🛠️️ Import `encoding/json`, `net/http`, and `github.com/gorilla/mux` in `main.go`:

```go
import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)
```

🛠️️ Add the courses function:

```go
func courses() []Course {
	t := Teacher{"Max", "https://images.unsplash.com/photo-1558531304-a4773b7e3a9c?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80"}
	cTags := []string{"colaborate", "git", "cli", "commit", "versionning"}
	c := Course{1, "eba25511-afce-4c8e-8cab-f82822434648", "learn git", cTags, "https://carlchenet.com/wp-content/uploads/2019/04/git-logo.png", "Learn how to create, manage, fork, and collaborate on a project. Git stays a major part of all companies projects. Learning git is learning how to make your project better everyday", "5 nov", "1604577600000", t}
	return []Course{c}
}
```

🛠️️ Add `func main()`:

```go
func main() {
	r := mux.NewRouter()
	r.Handle("/api/v1/courses", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse, err := json.Marshal(courses())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)

	}))

	http.ListenAndServe(":8000", r)
}
```

The main function will manage the routes, accept requests on `/api/v1/courses`, and will return the Courses (in `func courses()`).

🛠️️ Run the code with command `go run .` and open **insomnia** or **postman** to make a `GET` request that should end with `200`

### JWT authentication types

🛠️️Add authentication types in `main.go`:

```go
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}
```

[Next](https://github.com/cryptr-examples/cryptr-go-api-sample/tree/03-add-your-cryptr-credentials)
