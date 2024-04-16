package commands

import (
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestHelp(t *testing.T) {
	blocks := compileHelp()

	blocksJSON, _ := json.Marshal(blocks)
	spew.Dump(blocksJSON)
}
