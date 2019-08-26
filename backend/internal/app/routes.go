package app

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/r-cbb/cbbpoll/internal/errors"
	models "github.com/r-cbb/cbbpoll/pkg"
)

func (s *Server) Routes() {
	s.router = mux.NewRouter()
	s.router.HandleFunc("/", s.handlePing())
	s.router.HandleFunc("/ping", s.handlePing())
	s.router.HandleFunc("/teams", s.handleAddTeam()).Methods(http.MethodPost)
	s.router.HandleFunc("/teams", s.handleListTeams()).Methods(http.MethodGet)
	s.router.HandleFunc("/teams/{id:[0-9]+}", s.handleGetTeam()).Methods(http.MethodGet)

	s.router.HandleFunc("/sessions", s.handleNewSession()).Methods(http.MethodPost)
}

func (s *Server) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, struct{ Version string }{Version: s.version()}, http.StatusOK)
	}
}

func (s *Server) handleAddTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newTeam models.Team
		err := s.decode(w, r, &newTeam)
		if err != nil {
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		id, err := s.Db.AddTeam(newTeam)

		if errors.Kind(err) == errors.KindConcurrencyProblem {
			// Retry once
			id, err = s.Db.AddTeam(newTeam)
		}

		if err != nil {
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		s.respond(w, r, id, http.StatusOK)
		return
	}
}

func (s *Server) handleGetTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		intId, err := strconv.Atoi(id)
		if err != nil {
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		team, err := s.Db.GetTeam(int64(intId))
		if err != nil {
			if errors.Kind(err) == errors.KindNotFound {
				s.respond(w, r, nil, http.StatusNotFound)
				return
			}

			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		s.respond(w, r, team, http.StatusOK)
		return
	}
}

func (s *Server) handleListTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teams, err := s.Db.GetTeams()
		if err != nil {
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		s.respond(w, r, teams, http.StatusOK)
	}
}

func (s *Server) handleNewSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := struct {
			AccessToken string
		} {}
		err := s.decode(w, r, &authToken)
		if err != nil {
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		name, err := usernameFromRedditToken(authToken.AccessToken)
		if err != nil {
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		// todo: get user from database; create if doesn't exist

		out, err := createJWT(models.User{Nickname: name, IsAdmin: true})
		if err != nil {
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		payload := struct {
			Nickname string
			Token string
		}{
			Nickname: name,
			Token: out,
		}

		s.respond(w, r, payload, http.StatusOK)
	}
}