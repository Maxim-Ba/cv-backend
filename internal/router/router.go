package router

import (
	//...

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	m "github.com/Maxim-Ba/cv-backend/internal/middleware"
)



type Router struct {
	R *chi.Mux
}

func New() *Router {
	r := chi.NewRouter()
	logger := &m.StructuredLogger{Logger: slog.Default()}
    r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/tag", func(r chi.Router) {
			r.Get("/{tagID}", TagGet)
			r.Get("/", TagList)
			r.Post("/", TagCreate)
			r.Delete("/", TagDelete)
			r.Put("/", TagUpdate)
		})
		//
		r.Route("/tech", func(r chi.Router) {
			r.Get("/{techID}", TechGet)
			r.Get("/", TechList)
			r.Post("/", TechCreate)
			r.Delete("/", TechDelete)
			r.Put("/", TechUpdate)
		})
		//
		r.Route("/wh", func(r chi.Router) {
				r.Get("/{whID}", WorkHistoryGet)
			r.Get("/", WorkHistoryList)
			r.Post("/", WorkHistoryCreate)
			r.Delete("/", WorkHistoryDelete)
			r.Put("/", WorkHistoryUpdate)
		})
		//
		r.Route("/edu", func(r chi.Router) {
				r.Get("/{eduID}", RducationGet)
			r.Get("/", RducationList)
			r.Post("/", RducationCreate)
			r.Delete("/", RducationDelete)
			r.Put("/", RducationUpdate)
		})
		//
		r.Route("/fb", func(r chi.Router) {
				r.Get("/{fbID}", FeedBackGet)
			r.Get("/", FeedBackList)
			r.Post("/", FeedBackCreate)
		})

	})

	return &Router{
		R: r,
	}
}
