package routes

import (
	"cart-order-service/config"
	cart "cart-order-service/handler/cart"
	order "cart-order-service/handler/order"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Routes struct {
	Router *mux.Router
	Cart   *cart.Handler
	Order  *order.Handler
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
	r.SetupCart()
	r.SetupOrder()
}

func (r *Routes) SetupBaseURL() {
	baseURL := viper.GetString("BASE_URL_PATH")
	if baseURL != "" && baseURL != "/" {
		r.Router.PathPrefix(baseURL).HandlerFunc(URLRewriter(r.Router, baseURL))
	}
}

func (r *Routes) SetupCart() {
	cartRoutes := r.Router.PathPrefix("/cart").Subrouter()
	cartRoutes.HandleFunc("/{user_id}", r.Cart.GetCartByUserID).Methods(http.MethodGet, http.MethodOptions)
	cartRoutes.HandleFunc("/update/{user_id}", r.Cart.UpdateCart).Methods(http.MethodPut, http.MethodOptions)
	cartRoutes.HandleFunc("/add", r.Cart.AddCart).Methods(http.MethodPost, http.MethodOptions)
	cartRoutes.HandleFunc("/delete/{user_id}", r.Cart.DeleteCart).Methods(http.MethodDelete, http.MethodOptions)
}

func (r *Routes) SetupOrder() {
	orderRoutes := r.Router.PathPrefix("/order").Subrouter()
	orderRoutes.HandleFunc("/create", r.Order.CreateOrder).Methods(http.MethodPost, http.MethodOptions)
	orderRoutes.HandleFunc("/callback", r.Order.CallbackPayment).Methods(http.MethodPost, http.MethodOptions)
	orderRoutes.HandleFunc("/status/{user_id}", r.Order.GetOrderStatus).Methods(http.MethodGet, http.MethodOptions)
	orderRoutes.HandleFunc("/status/update", r.Order.UpdateStatus).Methods(http.MethodPut, http.MethodOptions)
	orderRoutes.HandleFunc("/shipping/update", r.Order.UpdateShipping).Methods(http.MethodPut, http.MethodOptions)
}
