package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/citadel/citadel"
	"github.com/gorilla/mux"
)

func containers(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	containers, err := manager.ListContainers(true)
	if err != nil {
		logger.WithField("error", err).Error("list containers")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	var out []*citadel.Container
	if id != "" {
		for _, c := range containers {
			if c.ID == id {
				out = []*citadel.Container{c}
				break
			}
		}
	} else {
		out = containers
	}

	if len(out) == 0 && id != "" {
		http.Error(w, http.StatusText(404), http.StatusNotFound)

		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		logger.WithField("error", err).Error("encode json")

		return
	}
}

func containersStart(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		logger.WithField("error", err).Error("parse form")
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	var (
		image *citadel.Image
		err   error
		pull  bool
	)

	if rp := r.FormValue("pull"); rp != "" {
		if pull, err = strconv.ParseBool(r.FormValue("pull")); err != nil {
			logger.WithField("error", err).Error("parse bool")
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}
	}

	if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
		logger.WithField("error", err).Error("decode json")
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	container, err := manager.Start(image, pull)
	if err != nil {
		logger.WithField("error", err).Error("start container")
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(container); err != nil {
		logger.WithField("error", err).Error("encode json")

		return
	}
}

func containersRemove(w http.ResponseWriter, r *http.Request) {
	var container *citadel.Container
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&container); err != nil {
		logger.WithField("error", err).Error("decode json")
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}
