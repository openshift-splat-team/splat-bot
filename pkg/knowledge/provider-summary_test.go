package knowledge

import (
	"fmt"
	"os"
	"testing"
)

func TestRssFeed(t *testing.T) {

	os.Setenv("OLLAMA_ENDPOINT", "http://192.168.0.77:11434")
	summary, err := getFeedSummary(20, "aws")
	fmt.Printf("%s, %v\n", summary, err)

	t.Error("booom")
}
