package web

import (
	"context"
	"log/slog"
	"net/http"
)

func HelloWebHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	name := r.FormValue("name")
	component := HelloPost(name)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.LogAttrs(context.Background(), slog.LevelError, "rendering HelloPost", slog.Any("err", err))
	}
}
