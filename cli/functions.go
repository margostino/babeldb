package cli

import (
	"fmt"
)

func createSource(params interface{}) {
	source := params.(string)
	fmt.Printf("source: %s\n", source)
}
