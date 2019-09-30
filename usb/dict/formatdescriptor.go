package dict

//csdf or fdsc in ascii
const FormatDescriptorMagic = 0x66647363

type FormatDescriptor struct {
}

func NewFormatDescriptorFromBytes([]byte) (FormatDescriptor, error) {
	return FormatDescriptor{}, nil
}
