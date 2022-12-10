package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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

	sampleData := map[int]*User{
		1: &User{"Alice", 24, map[int]struct{}{}},
		2: &User{"Bob", 24, map[int]struct{}{}},
		3: &User{"Alex", 30, map[int]struct{}{
			4: struct{}{},
			5: struct{}{},
			6: struct{}{},
			7: struct{}{},
			8: struct{}{},
		}},
		4: &User{"Mondela", 31, map[int]struct{}{
			3: struct{}{},
			5: struct{}{},
			6: struct{}{},
			7: struct{}{},
			8: struct{}{},
		}},
		5: &User{"Rob", 32, map[int]struct{}{
			3: struct{}{},
			4: struct{}{},
			6: struct{}{},
			7: struct{}{},
			8: struct{}{},
		}},
		6: &User{"John", 33, map[int]struct{}{
			3: struct{}{},
			4: struct{}{},
			5: struct{}{},
			7: struct{}{},
			8: struct{}{},
		}},
		7: &User{"Jack", 34, map[int]struct{}{
			3: struct{}{},
			4: struct{}{},
			5: struct{}{},
			6: struct{}{},
			8: struct{}{},
		}},
		8: &User{"Fray", 35, map[int]struct{}{
			3: struct{}{},
			4: struct{}{},
			5: struct{}{},
			6: struct{}{},
			7: struct{}{},
		}},
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
			s := httptest.NewServer(handler())
			defer s.Close()

			storage = sampleData
			nextId = 14

			client := &http.Client{}
			r := strings.NewReader(tc.request.payload)
			req, err := http.NewRequest(tc.request.method, s.URL+tc.request.path, r)
			if err != nil {
				t.Fatal(err)
			}
			if tc.request.payload != "" {
				req.Header.Set("Content-Type", "application/json; charset=utf-8")
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
