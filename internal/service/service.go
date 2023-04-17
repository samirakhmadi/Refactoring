package service

import (
	"net/http"

	"refactoring/internal/repository/storage"

	"github.com/go-chi/chi/v5"
)

type Service struct {
	store  storage.Repository
	router *chi.Mux
}

func NewService(store storage.Repository, router *chi.Mux) Service {
	return Service{
		store:  store,
		router: router,
	}
}

func (s Service) InitRoutes() {
	// инициализуруем heartbeat
	s.initHeartBeat(s.router)

	// иницилизируем api/v1
	s.initApiV1Routes(s.router)
}

func (s Service) initHeartBeat(router chi.Router, handlers ...func(http.Handler) http.Handler) chi.Router {
	return router.Group(func(r chi.Router) {
		// если добавятся middleware
		r.Use(handlers...)

		r.Get("/", heartbeat())
	})
}

func (s Service) initApiV1Routes(router chi.Router, handlers ...func(http.Handler) http.Handler) chi.Router {
	return router.Group(func(apiV1Route chi.Router) {
		apiV1Route.Route("/api/v1", func(r chi.Router) {
			// если добавятся middleware
			r.Use(handlers...)

			r.Route("/users", func(usersRoute chi.Router) {
				usersRoute.Get("/", searchUser(s.store))
				usersRoute.Post("/", createUser(s.store))

				s.InitUserByIdRoutes(usersRoute)
			})
		})
	})
}

func (s Service) InitUserByIdRoutes(root chi.Router, handlers ...func(http.Handler) http.Handler) chi.Router {
	return root.Group(func(userByIdRoute chi.Router) {
		// если добавятся middleware
		userByIdRoute.Use(handlers...)

		userByIdRoute.Route("/{id}", func(r chi.Router) {
			r.Get("/", getUser(s.store))
			r.Patch("/", updateUser(s.store))
			r.Delete("/", deleteUser(s.store))
		})
	})
}
