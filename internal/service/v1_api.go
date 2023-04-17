package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"refactoring/internal/repository/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func heartbeat() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writeCors(writer)

		now := time.Now().String()

		writeResponse(writer, []byte(now), http.StatusOK)
	}
}

func searchUser(s storage.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writeCors(writer)

		usersList, err := s.SearchUser()
		if err != nil {
			_ = render.Render(writer, request, errInvalidRequest(err))
			log.Println(err)
			return
		}

		resp, err := json.Marshal(usersList)
		if err != nil {
			_ = render.Render(writer, request, errInternalError(err))
			log.Println(err)
			return
		}

		writeResponse(writer, resp, http.StatusOK)

		return
	}
}

func createUser(s storage.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writeCors(writer)

		req := CreateUserRequest{}

		if err := render.Bind(request, &req); err != nil {
			_ = render.Render(writer, request, errInvalidRequest(err))
			log.Println(err, req)
			return
		}

		user, err := s.CreateUser(req.DisplayName, req.Email)
		if err != nil {
			_ = render.Render(writer, request, errInternalError(err))
			log.Println(err, req)
			return
		}

		render.Status(request, http.StatusCreated)
		render.JSON(writer, request, user)

		return
	}
}

func updateUser(s storage.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writeCors(writer)

		req := UpdateUserRequest{}

		if err := render.Bind(request, &req); err != nil {
			_ = render.Render(writer, request, errInvalidRequest(err))
			log.Println(err, req)
			return
		}

		id := chi.URLParam(request, "id")

		user, err := s.UpdateUser(id, req.DisplayName)
		if err != nil {
			_ = render.Render(writer, request, errInvalidRequest(err))
			log.Println(err)
			return
		}

		render.Status(request, http.StatusCreated)
		render.JSON(writer, request, user)

		return
	}
}

func getUser(s storage.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writeCors(writer)

		id := chi.URLParam(request, "id")

		user, err := s.GetUser(id)
		if err != nil {
			_ = render.Render(writer, request, errInvalidRequest(err))
			log.Println(err)
			return
		}

		render.Status(request, http.StatusOK)
		render.JSON(writer, request, user)

		return
	}
}

func deleteUser(s storage.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writeCors(writer)

		id := chi.URLParam(request, "id")

		if err := s.DeleteUser(id); err != nil {
			_ = render.Render(writer, request, errInvalidRequest(err))
			log.Println(err)
			return
		}

		render.Status(request, http.StatusOK)
		render.JSON(writer, request, Status{
			id:      id,
			success: "ok",
		})

		return
	}
}
