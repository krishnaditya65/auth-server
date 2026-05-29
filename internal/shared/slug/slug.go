package slug

import (
	"strings"
)

func FromEmail(email string) string {
	local := strings.Split(email, "@")[0]

	local = strings.ToLower(local)

	local = strings.ReplaceAll(local, ".", "-")

	return local
}
