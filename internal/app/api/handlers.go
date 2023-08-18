package api

import (
	"context"
	"log"
	"net/http"

	"broabroad/internal/app/database"
)

// Root handler
type rootHandler struct{}

func newRootHandler() *rootHandler {
	return &rootHandler{}
}

func (h rootHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	_, _ = rw.Write([]byte("root handler"))
}

type SeekRequestsCreator interface {
	CreateSeekRequest(ctx context.Context, sr database.SeekRequest) error
}

// Create seek request
type createSeekRequestHandler struct {
	seekReqCreator SeekRequestsCreator
}

func newCreateSeekRequestHandler(seekReqCreator SeekRequestsCreator) *createSeekRequestHandler {
	return &createSeekRequestHandler{
		seekReqCreator: seekReqCreator,
	}
}

func (h createSeekRequestHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sr := database.SeekRequest{
		ID:     1,
		Street: "test-street",
	}
	err := h.seekReqCreator.CreateSeekRequest(ctx, sr)
	if err != nil {
		log.Fatalf("create seek request: %v", err)
	}
	log.Print("seek request created!")

	_, _ = rw.Write([]byte("seek request created"))
}

// Get members
type getMembersHandler struct{}

func newGetMembersHandler() *getMembersHandler {
	return &getMembersHandler{}
}

func (h getMembersHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	_, _ = rw.Write([]byte("get members handler"))
}

// Add member
type addMemberHandler struct{}

func newAddMemberHandler() *addMemberHandler {
	return &addMemberHandler{}
}

func (h addMemberHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	_, _ = rw.Write([]byte("add member handler"))
}
