package routes

import (
	"Base-Project/config"
	"Base-Project/middleware"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

type Routes struct {
	Router *mux.Router
}

func EnabledCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS, PUT, PATCH")
		header.Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func URLRewriter(router *mux.Router, baseURLPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = func(url string) string {
			if strings.Index(url, baseURLPath) == 0 {
				url = url[len(baseURLPath):]
			}
			return url
		}(r.URL.Path)

		router.ServeHTTP(w, r)
	}
}

func LoggerMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "notifications") {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			recorder := httptest.NewRecorder()
			next.ServeHTTP(recorder, r)

			for k, v := range recorder.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(recorder.Code)
			recorder.Body.WriteTo(w)

			responseTime := time.Since(start).Seconds()
			formattedResponseTime := fmt.Sprintf("%.9f", responseTime)
			formattedResponseTime = fmt.Sprintf("%sµs", formattedResponseTime)

			log.Printf("%s - [%s] - [%s] \"%s %s %s\" %d %s\n",
				r.RemoteAddr,
				time.Now().Format(time.RFC1123),
				formattedResponseTime,
				r.Method,
				r.URL.Path,
				r.Proto,
				recorder.Code,
				r.UserAgent(),
			)
		})
	}
}

func (r *Routes) Run(port string) {
	r.SetupRouter()

	log.Printf("[HTTP SRV] clients on localhost port :%s", port)
	srv := &http.Server{
		Handler:      r.Router,
		Addr:         "localhost:" + port,
		WriteTimeout: config.WriteTimeout() * time.Second,
		ReadTimeout:  config.ReadTimeout() * time.Second,
	}

	log.Panic(srv.ListenAndServe())
}

func (r *Routes) SetupRouter() {
	r.Router = mux.NewRouter()
	r.Router.Use(EnabledCors, LoggerMiddleware())

	r.SetupBaseURL()
	r.SetupUser()
}

func (r *Routes) SetupBaseURL() {
	baseURL := viper.GetString("BASE_URL_PATH")
	if baseURL != "" && baseURL != "/" {
		r.Router.PathPrefix(baseURL).HandlerFunc(URLRewriter(r.Router, baseURL))
	}
}

func (r *Routes) SetupUser() {
	userRoutes := r.Router.PathPrefix("/users").Subrouter()
	userRoutes.Use(middleware.Authentication)
}
