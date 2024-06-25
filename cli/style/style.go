package style

import (
	"fmt"

	"github.com/andrwui/hyprsources/constants/colors"
)

type Style struct {
	state string
}

func CreateStyle() *Style {
	return &Style{
		state: "\033[",
	}
}

func (s *Style) AddStyles(styles ...string) *Style {

	for i, style := range styles {
		if i > 0 {
			s.state += ";"
		}

		s.state += style

	}
	s.state += "m"
	return s
}

func (s *Style) Use(text string) string {
	return fmt.Sprintf(s.state + text + colors.TEXT_RESET)
}
