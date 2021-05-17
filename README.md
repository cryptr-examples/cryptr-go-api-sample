# Cryptr with Go API

## 03 - Add your Cryptr credentials

🛠️️ Define cryptr config structure type:

```go
type CryptrConfig struct {
	AUDIENCE        string
	CRYPTR_BASE_URL string
	TENANT_DOMAIN   string
}
```

🛠️️ Instantiate project config in the main function with the variables that you get when creating your application at the end of Cryptr Onboarding or on your Cryptr application. Don't forget to replace `YOUR_DOMAIN`:

```go
func main() {
  // Instantiate project config:
	cryptrConfig := CryptrConfig{
		"http://localhost:8081",
		"https://auth.cryptr.eu",
		"YOUR_DOMAIN",
	}

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

Note: __If you are from the EU, you must add `https://auth.cryptr.eu/` and if you are from the US, you must add `https://auth.cryptr.us/`__

[Next](https://github.com/cryptr-examples/cryptr-go-api-sample/tree/04-protect-api-endpoints)
