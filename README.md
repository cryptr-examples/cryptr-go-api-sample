# Cryptr with Go API

## 04 - Protect API Endpoints

### Cryptr JWT middleware

🛠️️ Import the following packages:

```go
import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"

	// Add this packages:
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"github.com/form3tech-oss/jwt-go"
)
```

🛠️️ Next, add these types:

```go
type JWTMiddleware struct {
	Options Options
}

type errorHandler func(w http.ResponseWriter, r *http.Request, err string)

type TokenExtractor func(r *http.Request) (string, error)

type Options struct {
	ValidationKeyGetter jwt.Keyfunc
	UserProperty        string
	ErrorHandler        errorHandler
	CredentialsOptional bool
	Extractor           TokenExtractor
	Debug               bool
	EnableAuthOnOptions bool
	SigningMethod       jwt.SigningMethod
}
```

🛠️️ Next, add the following functions:

```go
func OnError(w http.ResponseWriter, r *http.Request, err string) {
 http.Error(w, err, http.StatusUnauthorized)
}
 
func CryptrJwtVerifier(token *jwt.Token, cryptrConfig CryptrConfig) (interface{}, error) {
 // validate "exp"
 checkExp := token.Claims.(jwt.MapClaims).VerifyExpiresAt(time.Now().Unix(), true)
 if !checkExp {
   return token, errors.New("token expired")
 }
 // validate "iat"
 checkIat := token.Claims.(jwt.MapClaims).VerifyIssuedAt(time.Now().Unix(), true)
 if !checkIat {
   return token, errors.New("token issued at error")
 }
 // validate "iss"
 iss := fmt.Sprintf("%v/t/%v", cryptrConfig.CRYPTR_BASE_URL, cryptrConfig.TENANT_DOMAIN)
 checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
 if !checkIss {
   return token, errors.New("invalid issuer")
 }
 // validate "aud"
 aud := cryptrConfig.AUDIENCE
 checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
 if !checkAud {
   return token, errors.New("invalid audience")
 }
 
 // validate Signature
 cert, err := getPemCert(token, cryptrConfig)
 if err != nil {
   panic(err.Error())
 }
 
 result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
 return result, nil
}
 
func NewCryptrJwtMiddleware(cryptrConfig CryptrConfig) *JWTMiddleware {
 options := Options{
   ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
     return CryptrJwtVerifier(token, cryptrConfig)
   },
   SigningMethod: jwt.SigningMethodRS256,
   UserProperty:  "user",
   ErrorHandler:  OnError,
   Extractor:     FromAuthHeader,
 }
 
 return &JWTMiddleware{
   Options: options,
 }
}
 
func FromAuthHeader(r *http.Request) (string, error) {
 authHeader := r.Header.Get("Authorization")
 if authHeader == "" {
   return "", nil
 }
 
 authHeaderParts := strings.Fields(authHeader)
 if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
   return "", errors.New("authorization header format must be Bearer {token}")
 }
 return authHeaderParts[1], nil
}
 
func (m *JWTMiddleware) logf(format string, args ...interface{}) {
 if m.Options.Debug {
   log.Printf(format, args...)
 }
}
 
func (m *JWTMiddleware) HandlerWithNext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
 err := m.CheckJWT(w, r)
 
 if err == nil && next != nil {
   next(w, r)
 }
}
 
func (m *JWTMiddleware) Handler(h http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
   err := m.CheckJWT(w, r)
 
   if err != nil {
     return
   }
 
   h.ServeHTTP(w, r)
 })
}
 
func (m *JWTMiddleware) CheckJWT(w http.ResponseWriter, r *http.Request) error {
 if r.Method == "OPTIONS" && !m.Options.EnableAuthOnOptions {
   return nil
 }
 
 token, err := m.Options.Extractor(r)
 
 if err != nil {
   m.logf("Error extracting JWT: %v", err)
   m.Options.ErrorHandler(w, r, err.Error())
   return fmt.Errorf("error extracting token: %w", err)
 }
 
 if token == "" {
   if m.Options.CredentialsOptional {
     m.logf("  No credentials found (CredentialsOptional=true)")
     return nil
   }
 
   errorMsg := "required authorization token not found"
   m.Options.ErrorHandler(w, r, errorMsg)
   m.logf("  Error: No credentials found (CredentialsOptional=false)")
   return fmt.Errorf(errorMsg)
 }
 
 parsedToken, err := jwt.Parse(token, m.Options.ValidationKeyGetter)
 
 if err != nil {
   m.logf("Error parsing token: %v", err)
   m.Options.ErrorHandler(w, r, err.Error())
   return fmt.Errorf("error parsing token: %w", err)
 }
 
 if m.Options.SigningMethod != nil && m.Options.SigningMethod.Alg() != parsedToken.Header["alg"] {
   message := fmt.Sprintf("Expected %s signing method but token specified %s",
     m.Options.SigningMethod.Alg(),
     parsedToken.Header["alg"])
   m.logf("Error validating token algorithm: %s", message)
   m.Options.ErrorHandler(w, r, errors.New(message).Error())
   return fmt.Errorf("error validating token algorithm: %s", message)
 }
 
 if !parsedToken.Valid {
   m.logf("Token is invalid")
   m.Options.ErrorHandler(w, r, "The token isn't valid")
   return errors.New("token is invalid")
 }
 
 newRequest := r.WithContext(context.WithValue(r.Context(), m.Options.UserProperty, parsedToken))
 *r = *newRequest
 return nil
}
 
func getPemCert(token *jwt.Token, cryptrConfig CryptrConfig) (string, error) {
 cert := ""
 JwksUri := fmt.Sprintf("%v/t/%v/.well-known", cryptrConfig.CRYPTR_BASE_URL, cryptrConfig.TENANT_DOMAIN)
 resp, err := http.Get(JwksUri)
 
 if err != nil {
   return cert, err
 }
 defer resp.Body.Close()
 
 var jwks = Jwks{}
 err = json.NewDecoder(resp.Body).Decode(&jwks)
 
 if err != nil {
   return cert, err
 }
 
 for k := range jwks.Keys {
   if token.Header["kid"] == jwks.Keys[k].Kid {
     cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
   }
 }
 
 if cert == "" {
   err := errors.New("unable to find appropriate key")
   return cert, err
 }
 
 return cert, nil
}
```

**How the middleware proceeds:**    
- It checks the header, if there are errors the request is blocked  
- `RS256` verifies the tokens against the public key for your Cryptr account  
- It accepts the request and processes it, otherwise it returns `401`  

### Secure courses route

We can now instantiate the middleware and protect the course route

🛠️️ First, import `negroni` package:

```go
import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"

	// 1. Import negroni package
	"github.com/codegangsta/negroni"
)
```

🛠️ ️Next, instantiate the cryptr middleware with config:

```go
func main() {
	cryptrConfig := CryptrConfig{
		"http://localhost:8081",
		"https://cleeck-umbrella-staging-eu.onrender.com",
		"shark-academy",
	}
	// 2. Instantiate cryptr jwt middleware:
	jwtMiddleware := NewCryptrJwtMiddleware(cryptrConfig)
	r := mux.NewRouter()
	// ...
}
```

🛠️️ Then you can now protect the route:

```go
func main() {
	cryptrConfig := CryptrConfig{
		"http://localhost:8081",
		"https://cleeck-umbrella-staging-eu.onrender.com",
		"shark-academy",
	}
	jwtMiddleware := NewCryptrJwtMiddleware(cryptrConfig)
	r := mux.NewRouter()

	// 3. Secure the courses route:
	r.Handle("/api/v1/courses", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jsonResponse, err := json.Marshal(courses())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", cryptrConfig.AUDIENCE)
			w.Header().Set("access-control-allow-headers", "authorization,content-type,sentry-trace")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResponse)
		}))))

	http.ListenAndServe(":8000", r)
}
```

### Test with a Cryptr Vue app

Let’s try this on an application. For this purpose, we have an example app on Vue.

🛠 Run your code with `go run .`

🛠 Clone our `cryptr-vue-sample`:

```bash
git clone --branch 07-backend-courses-api https://github.com/cryptr-examples/cryptr-vue2-sample.git
```

🛠 Install the Vue project dependencies with `yarn`

🛠️️ Add `.env.local` file with your variables:

```javascript
VUE_APP_AUDIENCE=http://localhost:8080
VUE_APP_CLIENT_ID=YOUR_CLIENT_ID
VUE_APP_CRYPTR_BASE_URL=YOUR_BASE_URL
VUE_APP_DEFAULT_LOCALE=fr
VUE_APP_DEFAULT_REDIRECT_URI=http://localhost:8080
VUE_APP_TENANT_DOMAIN=YOUR_DOMAIN
VUE_APP_CRYPTR_TELEMETRY=FALSE
```

🛠️️ Run vue server with `yarn serve` and try to connect. Your Vue application redirects you to your sign form page, where you can sign in or sign up with an email.

Note: __You can log in with a sandbox email and we send you a magic link which should directly arrive in your personal inbox.__

Once you're connected, click on "Protected route". You can now view the list of the courses.

It’s done, congratulations if you made it to the end!

I hope this was helpful, and thanks for following this tutorial! 🙂
