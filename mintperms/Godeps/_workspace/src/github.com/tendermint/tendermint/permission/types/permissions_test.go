package types

import (
	"fmt"
	"testing"
)

func TestPermFlagStrings(t *testing.T) {
	aP := NewDefaultAccountPermissions()
	m, r, err := AccountPermissionsToStrings(aP)
	fmt.Println(m)
	fmt.Println(r)
	fmt.Println(err)

}
