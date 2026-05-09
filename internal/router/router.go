package router

import (
	//...

	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgtype"
	httpSwagger "github.com/swaggo/http-swagger"

	m "github.com/Maxim-Ba/cv-backend/internal/middleware"
	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	"github.com/Maxim-Ba/cv-backend/internal/services"
	"github.com/Maxim-Ba/cv-backend/internal/view/components/pages"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

type Router struct {
	R           *chi.Mux
	Deps        *Dependencies
	db          *sql.DB
	adminUser   string
	adminPass   string
	adminSecret string
}

type Dependencies struct {
	TagService         *services.TagService
	TechService        *services.TechService
	EducationService   *services.EducationService
	WorkHistoryService *services.WorkHistoryService
	PDFService         *services.PDFService
}

func New(deps *Dependencies, db *sql.DB, allowedOrigin, adminUser, adminPass, appSecret string) *Router {
	r := chi.NewRouter()

	csrfKey := sha256.Sum256([]byte(appSecret))
	trustedOrigin := strings.TrimPrefix(strings.TrimPrefix(allowedOrigin, "https://"), "http://")
	csrfMiddleware := csrf.Protect(
		csrfKey[:],
		csrf.Secure(false),
		csrf.FieldName("csrf_token"),
		csrf.CookieName("csrf_token"),
		csrf.TrustedOrigins([]string{trustedOrigin, "localhost:3333"}),
		csrf.SameSite(csrf.SameSiteLaxMode),
	)

	corsMiddleware := cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigin, "http://localhost:3333"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	logger := &m.StructuredLogger{Logger: slog.Default()}
	r.Use(corsMiddleware)
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(csrfMiddleware)

	router := &Router{
		R:           r,
		Deps:        deps,
		db:          db,
		adminUser:   adminUser,
		adminPass:   adminPass,
		adminSecret: appSecret,
	}

	h := createHandlers(deps)

	fs := http.FileServer(http.Dir("internal/view/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Route("/admin", func(r chi.Router) {
		r.Get("/login", router.adminLogin)
		r.Post("/login", router.adminLoginPost)
		r.Group(func(r chi.Router) {
			r.Use(m.RequireAuth(adminUser, router.adminSecret))
			r.Get("/", router.adminDashboard)
			r.Get("/logout", router.adminLogout)
			r.Get("/tag", router.adminTags)
			r.Post("/tag", router.adminTagPost)
			r.Get("/tech", router.adminTech)
			r.Post("/tech", router.adminTechPost)
			r.Get("/history", router.admiHistory)
			r.Post("/history", router.adminHistoryPost)
			r.Get("/education", router.adminEducation)
			r.Post("/education", router.adminEducationPost)
		})
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/healthz", router.healthCheck)

	r.Get("/api/download-cv", h.PDFHandler.DownloadCV)

	r.Route("/api", func(r chi.Router) {
		r.Use(cacheControlMiddleware)
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
	PDFHandler         *PDFHandler
}

func createHandlers(deps *Dependencies) *handlers {
	tagHandler := NewTagHandler(deps.TagService)
	techHandler := NewTechHandler(deps.TechService)
	educationHandler := NewEducationHandler(deps.EducationService)
	workHistoryHandler := NewWorkHistoryHandler(deps.WorkHistoryService)
	pdfHandler := newPDFHandler(deps.PDFService)

	return &handlers{
		TagHandler:         tagHandler,
		TechHandler:        techHandler,
		EducationHandler:   educationHandler,
		WorkHistoryHandler: workHistoryHandler,
		PDFHandler:         pdfHandler,
	}
}

func (rt *Router) adminDashboard(w http.ResponseWriter, r *http.Request) {
	user := "Администратор"
	component := pages.AdminPage(user)
	component.Render(r.Context(), w)
}

func (rt *Router) adminLogout(w http.ResponseWriter, r *http.Request) {
	m.ClearSessionCookie(w)
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=60")
		} else {
			w.Header().Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	})
}

func (rt *Router) healthCheck(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := rt.db.PingContext(ctx); err != nil {
		slog.Error("health check: db ping failed", "error", err)
		dbStatus = "unavailable"
	}
	status := "ok"
	httpStatus := http.StatusOK
	if dbStatus != "ok" {
		status = "degraded"
		httpStatus = http.StatusServiceUnavailable
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(map[string]string{ //nolint:errcheck
		"status":  status,
		"db":      dbStatus,
		"version": "1.0.0",
	})
}

func (rt *Router) adminLogin(w http.ResponseWriter, r *http.Request) {
	component := pages.Login("", csrf.Token(r))
	component.Render(r.Context(), w)
}

func (rt *Router) adminLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		slog.Error(err.Error())
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != rt.adminUser || password != rt.adminPass {
		component := pages.Login("Неверный логин или пароль", csrf.Token(r))
		component.Render(r.Context(), w)
		return
	}
	m.SetSessionCookie(w, rt.adminUser, rt.adminSecret)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// --- Tag ---

func (rt *Router) adminTags(w http.ResponseWriter, r *http.Request) {
	user := "Администратор"
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	tagsResult, err := rt.Deps.TagService.List(pagebleRq)
	if err != nil {
		slog.Error(err.Error())
	}
	csrfToken := csrf.Token(r)

	var editTag models.Tag
	editID := r.URL.Query().Get("edit")
	if r.URL.Query().Get("create") == "1" {
		editID = "create"
	} else if editID != "" {
		id, convErr := strconv.ParseInt(editID, 10, 64)
		if convErr == nil {
			editTag, _ = rt.Deps.TagService.Get(id)
		}
	}
	component := pages.TagPage(user, tagsResult, editTag, editID, csrfToken)
	component.Render(r.Context(), w)
}

func (rt *Router) adminTagPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	switch r.FormValue("_method") {
	case "DELETE":
		idStr := r.FormValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			if _, err = rt.Deps.TagService.Delete(id); err != nil {
				slog.Error(err.Error())
			}
		}
	case "PUT":
		idStr := r.FormValue("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)
		tag := models.Tag{
			ID:       id,
			Name:     r.FormValue("name"),
			HexColor: r.FormValue("hexColor"),
		}
		if _, err := rt.Deps.TagService.Update(tag); err != nil {
			slog.Error(err.Error())
		}
	default:
		tag := models.Tag{
			Name:     r.FormValue("name"),
			HexColor: r.FormValue("hexColor"),
		}
		if _, err := rt.Deps.TagService.Create(tag); err != nil {
			slog.Error(err.Error())
		}
	}
	http.Redirect(w, r, "/admin/tag", http.StatusSeeOther)
}

// --- Tech ---

func (rt *Router) adminTech(w http.ResponseWriter, r *http.Request) {
	user := "Администратор"
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	techResult, err := rt.Deps.TechService.List(pagebleRq)
	if err != nil {
		slog.Error(err.Error())
	}
	csrfToken := csrf.Token(r)

	var editTech models.Technology
	editID := r.URL.Query().Get("edit")
	if r.URL.Query().Get("create") == "1" {
		editID = "create"
	} else if editID != "" {
		id, convErr := strconv.ParseInt(editID, 10, 64)
		if convErr == nil {
			editTech, _ = rt.Deps.TechService.Get(id)
		}
	}
	component := pages.TechPage(user, techResult, editTech, editID, csrfToken)
	component.Render(r.Context(), w)
}

func (rt *Router) adminTechPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	switch r.FormValue("_method") {
	case "DELETE":
		idStr := r.FormValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			if _, err = rt.Deps.TechService.Delete(id); err != nil {
				slog.Error(err.Error())
			}
		}
	case "PUT":
		idStr := r.FormValue("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)
		tech := models.Technology{
			ID:          id,
			Title:       r.FormValue("title"),
			Description: pgtype.Text{String: r.FormValue("description"), Valid: true},
			LogoUrl:     pgtype.Text{String: r.FormValue("logoUrl"), Valid: true},
		}
		if _, err := rt.Deps.TechService.Update(tech); err != nil {
			slog.Error(err.Error())
		}
	default:
		tech := models.Technology{
			Title:       r.FormValue("title"),
			Description: pgtype.Text{String: r.FormValue("description"), Valid: true},
			LogoUrl:     pgtype.Text{String: r.FormValue("logoUrl"), Valid: true},
		}
		if _, err := rt.Deps.TechService.Create(tech); err != nil {
			slog.Error(err.Error())
		}
	}
	http.Redirect(w, r, "/admin/tech", http.StatusSeeOther)
}

// --- Education ---

func (rt *Router) adminEducation(w http.ResponseWriter, r *http.Request) {
	user := "Администратор"
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	eduResult, err := rt.Deps.EducationService.List(pagebleRq)
	if err != nil {
		slog.Error(err.Error())
	}
	csrfToken := csrf.Token(r)

	var editEdu models.Education
	editID := r.URL.Query().Get("edit")
	if r.URL.Query().Get("create") == "1" {
		editID = "create"
	} else if editID != "" {
		id, convErr := strconv.ParseInt(editID, 10, 64)
		if convErr == nil {
			editEdu, _ = rt.Deps.EducationService.Get(id)
		}
	}
	component := pages.EducationPage(user, eduResult, editEdu, editID, csrfToken)
	component.Render(r.Context(), w)
}

func (rt *Router) adminEducationPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	switch r.FormValue("_method") {
	case "DELETE":
		idStr := r.FormValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			if _, err = rt.Deps.EducationService.Delete(id); err != nil {
				slog.Error(err.Error())
			}
		}
	case "PUT":
		idStr := r.FormValue("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)
		yearStr := r.FormValue("year")
		year, _ := strconv.Atoi(yearStr)
		edu := models.Education{
			ID:           id,
			Name:         pgtype.Text{String: r.FormValue("name"), Valid: true},
			Year:         int32(year),
			Course:       r.FormValue("course"),
			Organization: r.FormValue("organization"),
		}
		if _, err := rt.Deps.EducationService.Update(edu); err != nil {
			slog.Error(err.Error())
		}
	default:
		yearStr := r.FormValue("year")
		year, _ := strconv.Atoi(yearStr)
		edu := models.Education{
			Name:         pgtype.Text{String: r.FormValue("name"), Valid: true},
			Year:         int32(year),
			Course:       r.FormValue("course"),
			Organization: r.FormValue("organization"),
		}
		if _, err := rt.Deps.EducationService.Create(edu); err != nil {
			slog.Error(err.Error())
		}
	}
	http.Redirect(w, r, "/admin/education", http.StatusSeeOther)
}

// --- WorkHistory ---

func (rt *Router) admiHistory(w http.ResponseWriter, r *http.Request) {
	user := "Администратор"
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	whResult, err := rt.Deps.WorkHistoryService.List(pagebleRq)
	if err != nil {
		slog.Error(err.Error())
	}
	csrfToken := csrf.Token(r)

	var editWH models.WorkHistory
	editID := r.URL.Query().Get("edit")
	if r.URL.Query().Get("create") == "1" {
		editID = "create"
	} else if editID != "" {
		id, convErr := strconv.ParseInt(editID, 10, 64)
		if convErr == nil {
			editWH, _ = rt.Deps.WorkHistoryService.Get(id)
		}
	}
	component := pages.HistoryPage(user, whResult, editWH, editID, csrfToken)
	component.Render(r.Context(), w)
}

func parseLines(val string) []string {
	var result []string
	for _, line := range strings.Split(val, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

func parseDate(val string) pgtype.Date {
	if val == "" {
		return pgtype.Date{Valid: false}
	}
	t, err := time.Parse("2006-01-02", val)
	if err != nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: t, Valid: true}
}

func (rt *Router) adminHistoryPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	switch r.FormValue("_method") {
	case "DELETE":
		idStr := r.FormValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			if _, err = rt.Deps.WorkHistoryService.Delete(id); err != nil {
				slog.Error(err.Error())
			}
		}
	case "PUT":
		idStr := r.FormValue("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)
		wh := models.WorkHistory{
			ID:          id,
			Name:        r.FormValue("name"),
			About:       r.FormValue("about"),
			LogoUrl:     []byte(r.FormValue("logoUrl")),
			PeriodStart: parseDate(r.FormValue("periodStart")),
			PeriodEnd:   parseDate(r.FormValue("periodEnd")),
			WhatIDid:    parseLines(r.FormValue("whatIDid")),
			Projects:    parseLines(r.FormValue("projects")),
		}
		if _, err := rt.Deps.WorkHistoryService.Update(wh); err != nil {
			slog.Error(err.Error())
		}
	default:
		wh := models.WorkHistory{
			Name:        r.FormValue("name"),
			About:       r.FormValue("about"),
			LogoUrl:     []byte(r.FormValue("logoUrl")),
			PeriodStart: parseDate(r.FormValue("periodStart")),
			PeriodEnd:   parseDate(r.FormValue("periodEnd")),
			WhatIDid:    parseLines(r.FormValue("whatIDid")),
			Projects:    parseLines(r.FormValue("projects")),
		}
		if _, err := rt.Deps.WorkHistoryService.Create(wh); err != nil {
			slog.Error(err.Error())
		}
	}
	http.Redirect(w, r, "/admin/history", http.StatusSeeOther)
}
