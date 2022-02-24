package id

import (
	"crypto/rand"
	"io"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

var service *ulidService

type safeEntropy struct { //nolint govet-fieldalignment
	mux     sync.Mutex
	entropy io.Reader
}

type ulidService struct {
	t           time.Time
	safeEntropy safeEntropy
	emptyULID   ulid.ULID
}

func NewULID() {
	t := time.Unix(time.Now().Unix(), 0)

	entropy := ulid.Monotonic(rand.Reader, 0)
	service = &ulidService{
		t:           t,
		safeEntropy: safeEntropy{entropy: entropy},
		emptyULID:   ulid.ULID{},
	}
}

func New() (ulid.ULID, error) {
	service.safeEntropy.mux.Lock()
	id, err := ulid.New(ulid.Timestamp(service.t), service.safeEntropy.entropy)
	service.safeEntropy.mux.Unlock()

	return id, err
}

func GetEmpty() ulid.ULID {
	return service.emptyULID
}
