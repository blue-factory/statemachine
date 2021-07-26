package statemachine

import (
	"bytes"
	"fmt"
)

func (sm *StateMachine) RenderGraphviz() string {
	b := bytes.NewBufferString("")
	b.WriteString("digraph {\n")
	b.WriteString("\trankdir=LR;\n")
	b.WriteString("\tsize=\"8\"\n")
	b.WriteString("\tnode [shape = circle];\n")

	for current, s := range sm.states {
		// TODO(ca): Add label value to state struct, eg. [label = "label"]
		for _, dest := range s.Destination {
			b.WriteString(fmt.Sprintf("\t%s -> %s;\n", current, dest))
		}
	}

	b.WriteString("}")

	return b.String()
}

