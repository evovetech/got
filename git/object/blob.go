package object

type blob struct {
	Object

	data []byte
	err  error
}

func NewBlob(id Id) Blob {
	b := &blob{Object: New(id, BlobType)}
	b.SetInitFunc(func() {
		b.data, b.err = catBlob(id)
	})
	return b
}

func (b *blob) Contents() ([]byte, error) {
	b.init()
	return b.data, b.err
}

func (b *blob) init() *blob {
	b.Object.Init()
	return b
}
