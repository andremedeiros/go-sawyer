package httpcache

import (
	"bytes"
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	setup.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("{}"))
	})

	req, err := setup.Client.NewRequest("test")
	req.Header.Set("Accept", "application/json")
	assert.Equal(t, nil, err)

	orig := req.Get()
	assert.Equal(t, false, orig.IsError())
	assert.Equal(t, false, orig.IsApiError())

	var buf bytes.Buffer
	err = Encode(orig, &buf)
	assert.Equal(t, nil, err)

	cached := Decode(&buf)

	assert.Equal(t, false, cached.IsError())
	assert.Equal(t, 200, cached.StatusCode)
	assert.Equal(t, "", cached.Header.Get("Accept"))
	assert.Equal(t, "application/json", cached.Header.Get("Content-Type"))
	assert.Equal(t, "application/json", cached.MediaType.String())
}

type SetupServer struct {
	Client *sawyer.Client
	Server *httptest.Server
	Mux    *http.ServeMux
}

func Setup(t *testing.T) *SetupServer {
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	client, err := sawyer.NewFromString(srv.URL+"?a=1&b=1", nil)
	assert.Equalf(t, nil, err, "Unable to parse %s", srv.URL)

	return &SetupServer{client, srv, mux}
}

func (s *SetupServer) Teardown() {
	s.Server.Close()
}
