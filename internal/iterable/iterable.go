package iterable

type Iterable interface {
	/* Return the next key in set. */
	Iterate() <-chan interface{}
}
