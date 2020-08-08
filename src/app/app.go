package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"mobilewallet/core"
	"mobilewallet/transfer"
	"mobilewallet/user"
	// for driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "mobilewallet/docs"
)

// App level struct containing its dependencies
type App struct {
	AppCtx *core.AppContext
}

const dbDriver = "mysql"
const dbName = "mobilewallet"

var dbUsername = os.Getenv("DB_USERNAME")
var dbPassword = os.Getenv("DB_PASSWORD")

// Define our struct
type authenticationMiddleware struct {
	tokenUsers map[string]string
}

// Initialize it somewhere
func (amw *authenticationMiddleware) Populate() {
	amw.tokenUsers["11111111"] = "Rahul"
	amw.tokenUsers["22222222"] = "Mike"
	amw.tokenUsers["33333333"] = "Pierre"
}

// Middleware function, which will be called for each request
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if user, found := amw.tokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user: %s\n", user)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

// @title MobileWallet API
// @version 1.0
// @description This is the API for MobileWallet POC
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email rd4704@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
func (a *App) Initialize() {
	connectionString := fmt.Sprintf("%s:%s@tcp(db:3306)/%s", dbUsername, dbPassword, dbName)

	var err error
	a.AppCtx = &core.AppContext{}

	a.AppCtx.DB, err = sql.Open(dbDriver, connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.AppCtx.MainRouter = mux.NewRouter()
	a.AppCtx.APIRouter = a.AppCtx.MainRouter.PathPrefix("/api").Subrouter()

	amw := authenticationMiddleware{make(map[string]string)}
	amw.Populate()
	a.AppCtx.APIRouter.Use(amw.Middleware)

	a.initializeRoutes()
}

// Run app
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.AppCtx.MainRouter))
}

func (a *App) initializeRoutes() {
	user.New(a.AppCtx.DB, a.AppCtx.APIRouter)
	transfer.New(a.AppCtx.DB, a.AppCtx.APIRouter)
	a.AppCtx.MainRouter.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
}
