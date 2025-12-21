package router

import (
	"log/slog"
	"net/http"
)


func TagGet(w http.ResponseWriter, r *http.Request)  {}
func TagList(w http.ResponseWriter, r *http.Request) {
	slog.Info("TagList4",)
	w.Write([]byte("TagList"))

}

func TagCreate(w http.ResponseWriter, r *http.Request) {}

func TagDelete(w http.ResponseWriter, r *http.Request) {}

func TagUpdate(w http.ResponseWriter, r *http.Request) {}
