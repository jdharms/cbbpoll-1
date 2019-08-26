package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"github.com/r-cbb/cbbpoll/internal/errors"
	"github.com/r-cbb/cbbpoll/pkg"
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
		var newTeam pkg.Team
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

		req, err := http.NewRequest(http.MethodGet, "https://oauth.reddit.com/api/v1/me", nil)
		if err != nil {
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken.AccessToken))
		req.Header.Set("User-Agent", "cbbpoll_backend/0.1.0")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println(req.Header.Get("Authorization"))
			fmt.Println(resp.Status)
			s.respond(w, r, nil, http.StatusUnauthorized)
			return
		}

		content, err := ioutil.ReadAll(resp.Body)
		data := make(map[string]interface{})
		err = json.Unmarshal(content, &data)
		if err != nil {
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		name, ok := data["name"].(string)
		if !ok {
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		// todo: get user from database; create if doesn't exist

		// jwt stuff -- probably move to another function

		var claims jwt.MapClaims = make(map[string]interface{})
		claims["name"] = name
		claims["admin"] = true

		alg := jwt.GetSigningMethod("RS256")
		keytext, err := ioutil.ReadFile("jwtRS256.key")
		if err != nil {
			fmt.Println("couldn't read from secret file")
		}

		key, err := jwt.ParseRSAPrivateKeyFromPEM(keytext)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("couldn't parse key from byte")
		}

		token := jwt.NewWithClaims(alg, claims)

		out, err := token.SignedString(key)
		if err != nil {
			fmt.Printf("Error signing token: %v", err)
			return
		}

		s.respond(w, r, struct {
			Nickname string
			Token string
		}{
			Nickname: name,
			Token: out,
		}, http.StatusOK)
	}
}