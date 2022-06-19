package sets

type IntSet map[int]struct{}

func NewIntSet(members ...int) IntSet {
	s := make(map[int]struct{})
	for _, member := range members {
		s[member] = struct{}{}
	}
	return s
}

func (s IntSet) Add(member int) {
	s[member] = struct{}{}
}

func (s IntSet) Len() int {
	return len(s)
}
