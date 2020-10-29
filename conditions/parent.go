package conditions

type headParent struct {
	p Condition
}

func (h headParent) HasParent() bool {
	return h.p != nil
}

func (h *headParent) SetCursor() {
	h.p.SetCursor()
}

func (h *headParent) Parent(c Condition) {
	h.p = c
}
