package statemachine

import (
	"bytes"
	"fmt"
)

func (sm *StateMachine) RenderMermaid() string {
	b := bytes.NewBufferString("")
	b.WriteString("stateDiagram-v2\n")

	for current, s := range sm.states {
		for _, dest := range s.Destination {
			c := current
			if current == PristineState {
				c = "[*]"
			}

			d := dest
			if dest == EventAbort {
				d = "[*]"
			}

			b.WriteString(fmt.Sprintf("\t%s --> %s\n", c, d))
		}
	}

	return b.String()
}
