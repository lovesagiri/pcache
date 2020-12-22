package view

type ByteView struct {
	b []byte
}

func NewByteView(bytes []byte) ByteView {
	return ByteView{b: cloneBytes(bytes)}
}
func (b ByteView) Len() int {
	return len(b.b)
}

func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.b)
}

func (b ByteView) String() string {
	return string(b.b)
}

func cloneBytes(b []byte) []byte {
	newone := make([]byte, len(b))
	copy(newone, b)
	return newone
}
