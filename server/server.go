package server

import (
	"go-sms-gateway-api/config"
	"go-sms-gateway-api/handlers"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	router chi.Router
	db     *gorm.DB
	cfg    *config.Config
}

func NewServer(db *gorm.DB, cfg *config.Config) *Server {
	s := &Server{
		router: chi.NewRouter(),
		db:     db,
		cfg:    cfg,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(CORS)

	messagesHandler := handlers.NewMessagesHandler(s.db)

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			//r.Use(customMiddleware.AuthMiddleware(*s.cfg))
			r.Route("/messages", func(r chi.Router) {
				r.Post("/publish", messagesHandler.PublishMessage)
				r.Get("/ws", messagesHandler.HandleWebSocketConnection)
			})
		})
	})

	// Swagger documentation
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL pointing to API definition
	))

	// Health check
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
