package object

type blob struct {
	object

	data []byte
	err  error
}

func NewBlob(id Id) Blob {
	b := new(blob)
	b.id, b.kind = id, BlobType
	b.initFunc = func() {
		b.data, b.err = catBlob(id)
	}
	return b
}

func (b *blob) Contents() ([]byte, error) {
	b.Init()
	return b.data, b.err
}
