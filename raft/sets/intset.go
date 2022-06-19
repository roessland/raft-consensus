package sets

type IntSet map[int]struct{}

func NewIntSet() IntSet {
	return make(map[int]struct{})
}
