package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var filename string

type storage struct {
	Data   map[int]*User
	NextId int
}

func (st *storage) save() error {
	bs, err := json.MarshalIndent(st, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, bs, 0644)
	return err
}

func readStorage() (*storage, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	s := storage{}
	err = json.Unmarshal(bs, &s)
	return &s, err
}

type User struct {
	Name    string             `json:"name"`
	Age     int                `json:"age"`
	Friends map[int](struct{}) `json:"friends"`
}

func (u *User) addFriend(id int) error {
	st, err := readStorage()
	if err != nil {
		return err
	}
	if _, ok := st.Data[id]; !ok {
		return errors.New("no such user")
	}
	u.Friends[id] = struct{}{}
	return nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = fmt.Sprintf(":%s", port)

	filename = os.Getenv("FILENAME")
	if filename == "" {
		filename = "storage.json"
	}
	if _, err := os.Stat(filename); err != nil && errors.Is(err, os.ErrNotExist) {
		s := storage{}
		s.Data = make(map[int]*User)
		s.NextId = 0
		s.save()
	}

	log.Fatal(http.ListenAndServe(port, handler()))
}

func handler() http.Handler {
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
	st, err := readStorage()
	if err != nil {
		render.Render(w, r, ErrStorage(err))
		return
	}
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user := data.User
	id := st.NextId
	st.Data[id] = user

	for f, _ := range user.Friends {
		st.Data[f].addFriend(id)
	}

	st.NextId++
	if err := st.save(); err != nil {
		render.Render(w, r, ErrStorage(err))
		return
	}

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

	st, err := readStorage()
	if err != nil {
		render.Render(w, r, ErrStorage(err))
		return
	}

	u, ok := st.Data[src]
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
	st.Data[tgt].addFriend(src)

	if err := st.save(); err != nil {
		render.Render(w, r, ErrStorage(err))
		return
	}

	fmt.Fprintf(w, "%s и %s теперь друзья\n", u.Name, st.Data[tgt].Name)
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

	st, err := readStorage()
	if err != nil {
		render.Render(w, r, ErrStorage(err))
		return
	}

	u, ok := st.Data[id]
	if !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("no such user")))
		return
	}

	name := u.Name

	for f := range u.Friends {
		friend := st.Data[f]
		delete(friend.Friends, id)
	}
	delete(st.Data, id)

	if err := st.save(); err != nil {
		render.Render(w, r, ErrStorage(err))
		return
	}

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

	st, err := readStorage()
	if err != nil {
		render.Render(w, r, ErrStorage(err))
		return
	}

	u, ok := st.Data[id]
	if !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("no such user")))
		return
	}

	render.Render(w, r, &UserResponse{(u)})
}

func friendsList(w http.ResponseWriter, r *http.Request) {
	st, err := readStorage()
	if err != nil {
		render.Render(w, r, ErrStorage(err))
		return
	}
	u, err := validateId(st, r)
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
	st, err := readStorage()
	if err != nil {
		render.Render(w, r, ErrStorage(err))
		return
	}
	u, err := validateId(st, r)
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

	if err := st.save(); err != nil {
		render.Render(w, r, ErrStorage(err))
		return
	}

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

func ErrStorage(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Server storage error.",
		ErrorText:      err.Error(),
	}
}

func validateId(st *storage, r *http.Request) (*User, error) {
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	u, ok := st.Data[id]
	if !ok {
		return nil, errors.New("no such user")
	}

	return u, nil
}
