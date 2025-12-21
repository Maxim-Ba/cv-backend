package router

import (
	//...

	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	m "github.com/Maxim-Ba/cv-backend/internal/middleware"
	"github.com/Maxim-Ba/cv-backend/internal/view/components/pages"
)

type Router struct {
	R         *chi.Mux
}

func New() *Router {
	r := chi.NewRouter()
	logger := &m.StructuredLogger{Logger: slog.Default()}
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	router := &Router{
		R:         r,
	}


	fs := http.FileServer(http.Dir("internal/view/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Route("/admin", func(r chi.Router) {
		r.Get("/", router.adminDashboard)
		r.Get("/tag", router.adminTags)
		r.Get("/tech", router.adminTech)
		r.Get("/history", router.admiHistory)
		r.Get("/education", router.adminEducation)
		r.Get("/login", router.adminLogin)
		r.Post("/login", router.adminLoginPost) 
	})

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

	return router
}


func (rt *Router) adminDashboard(w http.ResponseWriter, r *http.Request) {
    user := "Администратор"
    component := pages.AdminPage(user)
    component.Render(r.Context(), w)
}
func (rt *Router) adminEducation(w http.ResponseWriter, r *http.Request) {
    user := "Администратор"
    component := pages.EducationPage(user)
    component.Render(r.Context(), w)
}
func (rt *Router) admiHistory(w http.ResponseWriter, r *http.Request) {
    user := "Администратор"
    component := pages.HistoryPage(user)
    component.Render(r.Context(), w)
}

func (rt *Router) adminTech(w http.ResponseWriter, r *http.Request) {
    user := "Администратор"
    component := pages.TechPage(user)
    component.Render(r.Context(), w)
}

func (rt *Router) adminTags(w http.ResponseWriter, r *http.Request) {
    user := "Администратор"
    component := pages.TagPage(user)
    component.Render(r.Context(), w)
}


func (rt *Router) adminLogin(w http.ResponseWriter, r *http.Request) {
    component := pages.Login("")
    component.Render(r.Context(), w)
}

func (rt *Router) adminLoginPost(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
		if err != nil {
			slog.Error(err.Error())
		}
    username := r.FormValue("username")
    password := r.FormValue("password")

    if username != "admin" || password != "admin" {
        component := pages.Login("Неверный логин или пароль")
        component.Render(r.Context(), w)
        return
    }
    http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
