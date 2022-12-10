package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var nextId int

var storage map[int]*User

type User struct {
	Name    string             `json:"name"`
	Age     int                `json:"age"`
	Friends map[int](struct{}) `json:"friends"`
}

func (u *User) addFriend(id int) error {
	if _, ok := storage[id]; !ok {
		return errors.New("no such user")
	}
	u.Friends[id] = struct{}{}
	return nil
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", handler()))
}

func handler() http.Handler {
	storage = make(map[int]*User)
	nextId = 0
	r := chi.NewRouter()

	r.Post("/create", create)
	r.Post("/make_friends", makeFriends)
	r.Delete("/user", deleteUser)
	r.Get("/user", showUser)
	r.Get("/friends/{id}", friendsList)
	r.Put("/{id}", updateAge)

	return r
}

func create(w http.ResponseWriter, r *http.Request) {
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user := data.User
	id := nextId
	storage[id] = user

	for f, _ := range user.Friends {
		storage[f].addFriend(id)
	}

	nextId++

	render.Status(r, http.StatusCreated)
	render.DefaultResponder(w, r, id)
}

func makeFriends(w http.ResponseWriter, r *http.Request) {
	data := &FriendRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	src, err := strconv.Atoi(data.Source)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	u, ok := storage[src]
	if !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("no such user")))
		return
	}
	tgt, err := strconv.Atoi(data.Target)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if err := u.addFriend(tgt); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// судя по формулировке "X и Y теперь друзья", дружба транзитивна
	storage[tgt].addFriend(src)

	fmt.Fprintf(w, "%s и %s теперь друзья\n", u.Name, storage[tgt].Name)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	data := &DeleteRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	id, err := strconv.Atoi(data.Target)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	u, ok := storage[id]
	if !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("no such user")))
		return
	}

	name := u.Name

	for f := range u.Friends {
		friend := storage[f]
		delete(friend.Friends, id)
	}
	delete(storage, id)
	fmt.Fprintln(w, name)
}

func showUser(w http.ResponseWriter, r *http.Request) {
	data := &DeleteRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	id, err := strconv.Atoi(data.Target)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	u, ok := storage[id]
	if !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("no such user")))
		return
	}

	render.Render(w, r, &UserResponse{(u)})
}

func friendsList(w http.ResponseWriter, r *http.Request) {
	u, err := validateId(r)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var friends []int
	for f := range u.Friends {
		friends = append(friends, f)
	}
	sort.Ints(friends)

	render.DefaultResponder(w, r, friends)
}

func updateAge(w http.ResponseWriter, r *http.Request) {
	u, err := validateId(r)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	data := &UpdateAgeRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	new, err := strconv.Atoi(data.NewAge)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	u.Age = new
	fmt.Fprintln(w, "возраст пользователя успешно обновлён")
}

type UserRequest struct {
	*User
	Age     string   `json:"age"`
	Friends []string `json:"friends,omitempty"`
}

func (u *UserRequest) Bind(r *http.Request) error {
	if u.User == nil || u.Friends == nil {
		return errors.New("missing required fields")
	}

	age, err := strconv.Atoi(u.Age)
	if err != nil {
		return err
	}
	u.User.Age = age

	u.User.Friends = make(map[int]struct{})
	for _, fr := range u.Friends {
		f, err := strconv.Atoi(fr)
		if err != nil {
			return err
		}
		if err := u.addFriend(f); err != nil {
			return err
		}
	}

	return nil
}

type UserResponse struct {
	*User
}

func (ur UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type FriendRequest struct {
	Source string `json:"source_id"`
	Target string `json:"target_id"`
}

func (f *FriendRequest) Bind(r *http.Request) error {
	if f.Source == f.Target {
		return errors.New("can't friend oneself")
	}
	return nil
}

type DeleteRequest struct {
	Target string `json:"target_id"`
}

func (f *DeleteRequest) Bind(r *http.Request) error {
	return nil
}

type UpdateAgeRequest struct {
	NewAge string `json:"new age"`
}

func (f *UpdateAgeRequest) Bind(r *http.Request) error {
	return nil
}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func validateId(r *http.Request) (*User, error) {
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	u, ok := storage[id]
	if !ok {
		return nil, errors.New("no such user")
	}

	return u, nil
}
