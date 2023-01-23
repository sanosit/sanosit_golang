package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const cType = "application/json; charset=utf-8"

func TestHandler(t *testing.T) {
	type request struct {
		method  string
		path    string
		payload string
	}

	type want struct {
		status int
		body   string
	}

	tests := map[string]struct {
		request request
		want    want
	}{
		"create": {
			request{
				"POST",
				"/create",
				`{"name":"some name","age":"24","friends":[]}`,
			},
			want{201, `14`},
		},
		"make_friends": {
			request{
				"POST",
				"/make_friends",
				`{"source_id":"1","target_id":"2"}`,
			},
			want{200, "Alice и Bob теперь друзья"},
		},
		"delete": {
			request{
				"DELETE",
				"/user",
				`{"target_id":"1"}`,
			},
			want{200, "Alice"},
		},
		"friends": {
			request{
				"GET",
				"/friends/8",
				"",
			},
			want{200, `[3,4,5,6,7]`},
		},
		"update_age": {
			request{
				"PUT",
				"/7",
				`{"new age":"28"}`,
			},
			want{200, "возраст пользователя успешно обновлён"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			initStorage(t)
			s := httptest.NewServer(handler())
			defer s.Close()

			client := &http.Client{}
			r := strings.NewReader(tc.request.payload)
			req, err := http.NewRequest(tc.request.method, s.URL+tc.request.path, r)
			if err != nil {
				t.Fatal(err)
			}
			if tc.request.payload != "" {
				req.Header.Set("Content-Type", cType)
			}
			res, err := client.Do(req)

			sc := res.StatusCode
			b, _ := io.ReadAll(res.Body)
			res.Body.Close()
			if sc != tc.want.status {
				t.Fatalf("unexpected status, want %v, got %v\n%s\n", tc.want.status, sc, b)
			}
			if string(b) != tc.want.body+"\n" {
				t.Fatalf("unexpected response body, want:\n%s\ngot:\n%s\n", tc.want.body, b)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	initStorage(t)
	s := httptest.NewServer(handler())
	defer s.Close()

	if _, err := http.Post(s.URL+"/make_friends", cType, strings.NewReader(`{"source_id":"2","target_id":"3"}`)); err != nil {
		t.Fatal(err)
	}
	st, err := readStorage()
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := st.Data[2].Friends[3]; !ok {
		t.Fatalf("Chandler is not Bob's friend!")
	}
	if _, ok := st.Data[3].Friends[2]; !ok {
		t.Fatalf("Bob is not Chandler's friend!")
	}
}

func TestMakeFriends(t *testing.T) {
	initStorage(t)
	s := httptest.NewServer(handler())
	defer s.Close()

	if _, err := http.Post(s.URL+"/create", cType, strings.NewReader(`{"name":"John Doe","age":"42","friends":[]}`)); err != nil {
		t.Fatal(err)
	}
	st, err := readStorage()
	if err != nil {
		t.Fatal(err)
	}

	if st.Data[14].Age != 42 {
		t.Fatalf("Expected user's age not in storage.")
	}
}

func TestDelete(t *testing.T) {
	initStorage(t)
	s := httptest.NewServer(handler())
	defer s.Close()

	client := &http.Client{}
	r := strings.NewReader(`{"target_id":"5"}`)
	req, err := http.NewRequest("DELETE", s.URL+"/user", r)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", cType)
	client.Do(req)

	st, err := readStorage()
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := st.Data[5]; ok {
		t.Fatalf("Rachel is still here.")
	}
	if _, ok := st.Data[6].Friends[5]; ok {
		t.Fatalf("Rachel is gone but she's still Ross's friend.")
	}
}

func TestUpdateAge(t *testing.T) {
	initStorage(t)
	s := httptest.NewServer(handler())
	defer s.Close()

	client := &http.Client{}
	r := strings.NewReader(`{"new age":"42"}`)
	req, err := http.NewRequest("PUT", s.URL+"/2", r)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", cType)
	client.Do(req)

	st, err := readStorage()
	if err != nil {
		t.Fatal(err)
	}

	if st.Data[2].Age != 42 {
		t.Fatalf("Bob didn't age.")
	}
}
func TestReadSave(t *testing.T) {
	initStorage(t)
}

func initStorage(t *testing.T) {
	t.Helper()
	filename = filepath.Join("../testdata", "storage.json")
	st, err := readStorage()
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.CreateTemp("", "storage")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	filename = f.Name()
	if err := st.save(); err != nil {
		t.Fatal(err)
	}
}
