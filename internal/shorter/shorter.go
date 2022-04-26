package shorter

import (
	"crypto/md5"
	"fmt"
)

//MakeShortner func for  hash.
func MakeShortner(s string) string {
	h := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", h)
}
