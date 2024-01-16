package identifier

import (
	"strings"

	"github.com/google/uuid"
)

func NewUUID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
