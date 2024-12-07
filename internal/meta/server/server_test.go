package server

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type httpServerTestSuite struct {
	suite.Suite
}

func TestHTTPServer(t *testing.T) {
	suite.Run(t, new(httpServerTestSuite))
}

func (h *httpServerTestSuite) TestSuccessHTTPServerWithOptions() {
	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello`))
	})

	s := New(handler,
		WithPort(8081),
		WithReadTimeout(time.Second),
		WithWriteTimeout(time.Second),
	)

	go s.Start()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/", nil)
	h.NoError(err)

	res, err := http.DefaultClient.Do(req)
	h.NoError(err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	h.NoError(err)
	h.Equal(res.StatusCode, http.StatusOK)
	h.Equal("hello", string(body))

	h.NoError(s.Stop(context.Background()))
}

func (h *httpServerTestSuite) TestSuccessHTTPServerWithDefaultOptions() {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	server := New(handler)
	go server.Start()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	h.NoError(err)

	res, err := http.DefaultClient.Do(req)
	h.NoError(err)
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	h.NoError(err)

	h.Equal("hello", string(resBody))
	h.Equal(http.StatusOK, res.StatusCode)

	h.NoError(server.Stop(context.Background()))
}

func (h *httpServerTestSuite) TestInternalServerErrorWithPanicAndRecoverMiddleware() {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	buf := NewMockWriter()

	logger := log.New(buf, "", log.LstdFlags)

	s := New(handler, WithErrorLogger(logger))
	go s.Start()

	r, err := http.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	h.NoError(err)

	resp, err := http.DefaultClient.Do(r)
	h.NoError(err)
	defer resp.Body.Close()

	h.Equal(http.StatusInternalServerError, resp.StatusCode)
	h.NoError(s.Stop(context.Background()))

	h.Regexp("recovered from panic", buf.String())
}

type mockWriter struct {
	buf   *bytes.Buffer
	mutex *sync.Mutex
}

func NewMockWriter() *mockWriter {
	return &mockWriter{
		buf:   bytes.NewBuffer([]byte{}),
		mutex: &sync.Mutex{},
	}
}

func (m *mockWriter) Write(b []byte) (n int, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.buf.Write(b)
}

func (m *mockWriter) String() string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.buf.String()
}
