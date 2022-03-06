package common

//EventStreamReader defines a way to retrieve events from a stream such as a response body.
//Next returns the next event in the stream. Implementations should return events in chronological order.
type EventIterator interface {
	Next() (*Event, error)
}
