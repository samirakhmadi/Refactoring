package storage

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	"refactoring/internal/repository/model"
)

type Repository interface {
	SearchUser() (*model.UserList, error)
	CreateUser(string, string) (model.User, error)
	GetUser(string) (model.User, error)
	UpdateUser(string, string) (model.User, error)
	DeleteUser(string) error
}

type Store struct {
	mx        *sync.Mutex
	storePath string
}

func NewJsonStorage(storePath string) (Store, error) {
	// проверяем есть ли доступ к нашему хранилищу
	_, err := ioutil.ReadFile(storePath)
	if err != nil {
		return Store{}, fmt.Errorf("error cannot reach store: %w", err)
	}

	return Store{
		mx:        &sync.Mutex{},
		storePath: storePath,
	}, nil
}

func (s Store) SearchUser() (*model.UserList, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	userStore, err := getUserStore(s.storePath)
	return &userStore.List, err
}

func (s Store) CreateUser(displayName, email string) (model.User, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	userStore, err := getUserStore(s.storePath)
	if err != nil {
		return model.User{}, err
	}

	userStore.Increment++

	user := model.User{
		CreatedAt:   time.Now(),
		DisplayName: displayName,
		Email:       email,
	}

	id := strconv.Itoa(userStore.Increment)
	userStore.List[id] = user

	if err = writeStore(userStore, s.storePath); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s Store) GetUser(id string) (model.User, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	userStore, err := getUserStore(s.storePath)
	if err != nil {
		return model.User{}, err
	}

	if _, ok := userStore.List[id]; !ok {
		return model.User{}, fmt.Errorf("error user not found: %s", id)
	}

	return userStore.List[id], nil
}

func (s Store) UpdateUser(id, displayName string) (model.User, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	userStore, err := getUserStore(s.storePath)
	if err != nil {
		return model.User{}, err
	}

	if _, ok := userStore.List[id]; !ok {
		return model.User{}, fmt.Errorf("error user not found: %s", id)
	}

	user := userStore.List[id]
	user.DisplayName = displayName
	userStore.List[id] = user

	if err = writeStore(userStore, s.storePath); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s Store) DeleteUser(id string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	userStore, err := getUserStore(s.storePath)
	if err != nil {
		return err
	}

	if _, ok := userStore.List[id]; !ok {
		return fmt.Errorf("error user not found: %s", id)
	}

	delete(userStore.List, id)

	if err = writeStore(userStore, s.storePath); err != nil {
		return err
	}

	return nil
}

func getUserStore(storePath string) (*model.UserStore, error) {
	// тк мы при инициализации проверяем есть ли доступ к файлу то проверка на ошибку не обязательна
	f, _ := ioutil.ReadFile(storePath)

	userStore := model.UserStore{}

	err := json.Unmarshal(f, &userStore)
	if err != nil {
		return nil, fmt.Errorf("error cannot find user %w", err)
	}

	return &userStore, nil
}

func writeStore(userStore *model.UserStore, storePath string) error {
	userBytes, err := json.Marshal(&userStore)
	if err != nil {
		return fmt.Errorf("error marshal during create user: %w", err)
	}

	err = ioutil.WriteFile(storePath, userBytes, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("error write during create user: %w", err)
	}

	return nil
}
