package payloadtype

import "io"

type Payload map[string]any

type ResultStorer interface {
	io.Closer
	Storer
	Open()
}

type Storer interface {
	Store(label string, payload Payload)
	StoreSysInfo(payload Payload)
}
