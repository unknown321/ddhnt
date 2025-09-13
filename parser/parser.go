package parser

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
)

type RawEntry struct {
	Data   [260]byte
	Size   int32
	Offset int32
}

type Entry struct {
	Size   int32 `json:"-"`
	Offset int32 `json:"-"`
	Path   string
}

type PakFile struct {
	EntryCount int32
	Entries    []Entry
}

func Parse(reader io.Reader) (PakFile, error) {
	pk := PakFile{}

	ec := make([]byte, 4)
	_, err := reader.Read(ec)
	if err != nil {
		return pk, err
	}

	if _, err = binary.Decode(ec, binary.LittleEndian, &pk.EntryCount); err != nil {
		return pk, err
	}

	var count int32
	b := make([]byte, 268)
	for {
		_, err = reader.Read(b)

		if err != nil {
			if err == io.EOF {
				break
			}
			return pk, fmt.Errorf("read: %w", err)
		}

		e := RawEntry{}
		_, err = binary.Decode(b, binary.LittleEndian, &e)
		if err != nil {
			return pk, fmt.Errorf("decode: %w", err)
		}

		rawStr := []byte{}
		for _, bb := range e.Data {
			if bb == 0xcc {
				break
			}

			if bb == 0 {
				break
			}

			rawStr = append(rawStr, bb)
		}

		ss := charmap.Windows1251.NewDecoder()
		pp, err := ss.String(string(rawStr))
		if err != nil {
			return pk, err
		}
		pk.Entries = append(pk.Entries, Entry{
			Size:   e.Size,
			Offset: e.Offset,
			Path:   pp,
		})

		count++

		if count == pk.EntryCount {
			break
		}
	}

	return pk, nil
}
