package snowflake_test

import (
	"testing"

	"github.com/qicro/pkg/snowflake"
)

func TestGenCode(t *testing.T) {
	snowflake.Init(0)
	id, err := snowflake.GetID()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("successfully generated code: %d", id)
}
