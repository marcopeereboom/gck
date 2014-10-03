package section

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

const (
	Version = 1 // image version

	InvalidId  = 0
	StackId    = 1
	CodeId     = 2
	ConstId    = 3
	VariableId = 4
	OsId       = 5

	FExecute  = 1 << 0
	FWrite    = 1 << 1
	FRead     = 1 << 2
	FCompress = 1 << 3
)

var (
	Sections = map[uint64]string{
		InvalidId:  ".INVALID",
		StackId:    ".STACK",
		CodeId:     ".CODE",
		ConstId:    ".CONST",
		VariableId: ".VAR",
		OsId:       ".OS",
	}
)

type Section struct {
	Version uint64
	Name    string
	Id      uint64

	Read    bool
	Write   bool
	Execute bool

	Payload interface{} // end result of encode/decode
}

func NewCodeSection(image []uint64) *Section {
	cs := Section{
		Version: Version,
		Name:    Sections[CodeId],
		Id:      CodeId,
		Read:    true,
		Write:   false,
		Execute: true,
		Payload: image,
	}
	return &cs
}

func encodeCode(code []uint64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, code)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func compressImage(image []byte) ([]byte, error) {
	if len(image) == 0 {
		return nil, fmt.Errorf("empty image")
	}

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(image)
	w.Close()

	return b.Bytes(), nil
}

// this really is a bit silly but we are going for correct, not fast
func uncompressedSize(image []byte) (int, error) {
	b := bytes.NewReader(image)
	r, err := zlib.NewReader(b)
	if err != nil {
		return 0, err
	}
	result, err := ioutil.ReadAll(r)
	r.Close()

	return len(result), nil
}

func (s *Section) Raw(compress bool) ([]byte, error) {
	var (
		image []byte
		err   error
	)

	switch p := s.Payload.(type) {
	case []uint64:
		switch s.Id {
		case CodeId:
			if len(p) == 0 {
				return nil, fmt.Errorf("can't encode empty " +
					"code section")
			}
			image, err = encodeCode(p)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid type %T for section %v",
				p, Sections[s.Id])
		}

	case []*Variable:
		switch s.Id {
		case VariableId:
			image, err = encodeVariableElements(p)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid type %T for section %v",
				p, Sections[s.Id])
		}

	case []*Const:
		switch s.Id {
		case ConstId:
			image, err = encodeConstElements(p)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid type %T for section %v",
				p, Sections[s.Id])
		}

	case []*Os:
		switch s.Id {
		case OsId:
			image, err = encodeOsElements(p)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid type %T for section %v",
				p, Sections[s.Id])
		}

	default:
		return nil, fmt.Errorf("unknown section id 0x%x", s.Id)
	}

	h := Header{
		Version:   Version,
		ImageSize: uint64(len(image)), // assume uncompressed
		SectionId: s.Id,
		Digest:    sha256.Sum256(image),
	}

	// flags
	if s.Execute {
		h.Flags |= FExecute
	}
	if s.Write {
		h.Flags |= FWrite
	}
	if s.Read {
		h.Flags |= FRead
	}

	var (
		compressed []byte
		b          *bytes.Buffer
	)
	if compress {
		compressed, err = compressImage(image)
		if err != nil {
			return nil, err
		}
		h.Flags |= FCompress
		h.ImageSize = uint64(len(compressed))
	} else {
		compressed = image
	}

	// encode header as big endian
	b = new(bytes.Buffer)
	if err := binary.Write(b, binary.BigEndian, h); err != nil {
		return nil, err
	}

	// write image as is
	_, err = b.Write(compressed)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func SectionsFromImage(b []byte) ([]*Section, error) {
	var (
		n, at    int
		sections []*Section
	)

	for todo := len(b); todo > 0; {
		v, err := SectionFromImage(b[at:], &n)
		if err != nil {
			return nil, err
		}
		sections = append(sections, v)
		at += n
		todo -= n
	}

	return sections, nil
}

func SectionFromImage(image []byte, consumed *int) (*Section, error) {
	s := Section{}
	h := &Header{}

	b := bytes.NewReader(image)
	// decode header
	if err := binary.Read(b, binary.BigEndian, h); err != nil {
		return nil, err
	}

	if Version != h.Version {
		return nil, fmt.Errorf("invalid version expected %v got %v",
			Version, h.Version)
	}

	// fill out section
	s.Version = h.Version
	s.Id = h.SectionId

	switch s.Id {
	case CodeId:
		s.Name = Sections[CodeId]
	case VariableId:
		s.Name = Sections[VariableId]
	case ConstId:
		s.Name = Sections[ConstId]
	case OsId:
		s.Name = Sections[OsId]
	default:
		return nil, fmt.Errorf("invalid image segment id 0x%x", s.Id)
	}

	// flags
	if h.Flags&FExecute == FExecute {
		s.Execute = true
	}
	if h.Flags&FWrite == FWrite {
		s.Write = true
	}
	if h.Flags&FRead == FRead {
		s.Read = true
	}

	// read image
	var (
		size int
		seek int64
		ir   io.Reader
		err  error
	)
	if h.Flags&FCompress == FCompress {
		seek, err = b.Seek(0, 1)
		if err != nil {
			return nil, err
		}

		size, err = uncompressedSize(image[seek:])
		if err != nil {
			return nil, err
		}

		uc := bytes.NewReader(image[seek:])
		zr, err := zlib.NewReader(uc)
		if err != nil {
			return nil, err
		}
		ir = zr
		defer zr.Close()
	} else {
		size = int(h.ImageSize)
		ir = io.Reader(b)
	}

	// calculate digest
	blob := make([]byte, size)
	_, err = ir.Read(blob)
	if err != nil {
		return nil, err
	}
	blobR := bytes.NewReader(blob)

	// decode image from blob
	switch s.Id {
	case CodeId:
		s.Payload = make([]uint64, size/8)
		err := binary.Read(blobR, binary.BigEndian, s.Payload)
		if err != nil {
			return nil, err
		}
	case VariableId:
		vars, err := decodeVariableElements(blob)
		if err != nil {
			return nil, err
		}

		s.Payload = vars

	case ConstId:
		consts, err := decodeConstElements(blob)
		if err != nil {
			return nil, err
		}

		s.Payload = consts

	case OsId:
		osc, err := decodeOsElements(blob)
		if err != nil {
			return nil, err
		}

		s.Payload = osc
	default:
		// can't happen due to test above
		return nil, fmt.Errorf("invalid segment id 0x%x", s.Id)
	}

	dt := sha256.Sum256(blob)
	if !bytes.Equal(dt[:], h.Digest[:]) {
		return nil, fmt.Errorf("invalid image digest")
	}

	if consumed != nil {
		*consumed = int(uint64(seek) + h.ImageSize)
	}

	return &s, nil
}
