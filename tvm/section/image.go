package section

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// raw section is big endian for human readability
// Header is followed by []byte that contains the image
type Header struct {
	Version   uint64
	ImageSize uint64
	Flags     uint64
	SectionId uint64
	Digest    [sha256.Size]byte
}

// used to encode/decode []byte images
type Image struct {
	image       []byte
	sectionSeen map[uint64]bool
}

func NewImage() *Image {
	i := Image{
		sectionSeen: make(map[uint64]bool),
	}
	return &i
}

func (i *Image) AddSection(s *Section, compress bool) error {
	if seen := i.sectionSeen[s.Id]; seen == true {
		return fmt.Errorf("section 0x%0x already added", s.Id)
	}

	switch s.Id {
	case CodeId:
	case VariableId:
	case ConstId:
	case OsId:
	default:
		return fmt.Errorf("invalid image section id 0x%0x", s.Id)
	}

	b, err := s.Raw(compress)
	if err != nil {
		return err
	}

	// append to image
	buf := bytes.NewBuffer(i.image)
	_, err = buf.Write(b)
	if err != nil {
		return err
	}
	i.image = buf.Bytes()

	// mark section as seen
	i.sectionSeen[s.Id] = true

	return nil
}

func (i *Image) GetImage() []byte {
	img := make([]byte, len(i.image))
	copy(img, i.image)
	return img
}
