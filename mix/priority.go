package mix

type event struct {
	when int64
	loc  int
}

type priority []event

func (p priority) Len() int {
	return len(p)
}

func (p priority) Less(i, j int) bool {
	if p[i].when < p[j].when {
		return true
	}
	if p[i].when > p[j].when {
		return false
	}
	return p[i].loc > p[j].loc
}

func (p priority) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *priority) Push(x interface{}) {
	*p = append(*p, x.(event))
}

func (p *priority) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[:n-1]
	return x
}
