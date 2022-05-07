package gnuplot

import (
	"fmt"
	"strings"
)

func TerminalTypeFromFilename(filename string) (string, error) {
	switch {
	case strings.HasSuffix(strings.ToLower(filename), ".png"):
		return "png", nil
	case strings.HasSuffix(strings.ToLower(filename), ".svg"):
		return "svg", nil
	default:
		return "", fmt.Errorf("unable to select suitable terminal type for %q", filename)

	}
}
