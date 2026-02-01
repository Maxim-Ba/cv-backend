package router

import (
	//...

	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"

	m "github.com/Maxim-Ba/cv-backend/internal/middleware"
	"github.com/Maxim-Ba/cv-backend/internal/services"
	"github.com/Maxim-Ba/cv-backend/internal/view/components/pages"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

type Router struct {
	R    *chi.Mux
	Deps *Dependencies
}

type Dependencies struct {
	TagService         *services.TagService
	TechService        *services.TechService
	EducationService   *services.EducationService
	WorkHistoryService *services.WorkHistoryService
}

func New(deps *Dependencies) *Router {
	r := chi.NewRouter()

csrfMiddleware := csrf.Protect(
		[]byte("32-byte-long-auth-key"), 
		csrf.Secure(false),              
		csrf.FieldName("csrf_token"),
		csrf.CookieName("csrf_token"),
	)

	logger := &m.StructuredLogger{Logger: slog.Default()}
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(csrfMiddleware)

	router := &Router{
		R:    r,
		Deps: deps,
	}

	h := createHandlers(deps)

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
			r.Get("/{tagID}", h.TagHandler.TagGet)
			r.Get("/", h.TagHandler.TagList)
			r.Post("/", h.TagHandler.TagCreate)
			r.Delete("/", h.TagHandler.TagDelete)
			r.Put("/", h.TagHandler.TagUpdate)
		})
		//
		r.Route("/tech", func(r chi.Router) {
			r.Get("/{techID}", h.TechHandler.TechGet)
			r.Get("/", h.TechHandler.TechList)
			r.Post("/", h.TechHandler.TechCreate)
			r.Delete("/", h.TechHandler.TechDelete)
			r.Put("/", h.TechHandler.TechUpdate)
		})
		//
		r.Route("/wh", func(r chi.Router) {
			r.Get("/{whID}", h.WorkHistoryHandler.WorkHistoryGet)
			r.Get("/", h.WorkHistoryHandler.WorkHistoryList)
			r.Post("/", h.WorkHistoryHandler.WorkHistoryCreate)
			r.Delete("/", h.WorkHistoryHandler.WorkHistoryDelete)
			r.Put("/", h.WorkHistoryHandler.WorkHistoryUpdate)
		})
		//
		r.Route("/edu", func(r chi.Router) {
			r.Get("/{eduID}", h.EducationHandler.EducationGet)
			r.Get("/", h.EducationHandler.EducationList)
			r.Post("/", h.EducationHandler.EducationCreate)
			r.Delete("/", h.EducationHandler.EducationDelete)
			r.Put("/", h.EducationHandler.EducationUpdate)
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

type handlers struct {
	TagHandler         *TagHandler
	TechHandler        *TechHandler
	EducationHandler   *EducationHandler
	WorkHistoryHandler *WorkHistoryHandler
}

func createHandlers(deps *Dependencies) *handlers {
	tagHandler := NewTagHandler(*deps.TagService)
	techHandler := NewTechHandler(deps.TechService)
	educationHandler := NewEducationHandler(deps.EducationService)
	workHistoryHandler := NewWorkHistoryHandler(deps.WorkHistoryService)

	return &handlers{
		TagHandler:         tagHandler,
		TechHandler:        techHandler,
		EducationHandler:   educationHandler,
		WorkHistoryHandler: workHistoryHandler,
	}
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
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	techResult, err := rt.Deps.TechService.List(pagebleRq)
	if err != nil {
		slog.Error(err.Error())
	}
	csrfToken := csrf.Token(r)
	
	editID := r.URL.Query().Get("edit")
	component := pages.TechPage(user, techResult, editID, csrfToken )
	component.Render(r.Context(), w)
}

func (rt *Router) adminTags(w http.ResponseWriter, r *http.Request) {
	user := "Администратор"
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	tagsResult, err := rt.Deps.TagService.List(pagebleRq)
	if err != nil {
		slog.Error(err.Error())
	}
	component := pages.TagPage(user, tagsResult)
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
