package packer

import (
	"bytes"
	"encoding/binary"
	"github.com/unknown321/ddhnt/parser"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Localized struct {
	Entries []parser.Entry
}

func Unpack(reader io.ReadSeeker, outdir string) (Localized, error) {
	pk, err := parser.Parse(reader)
	if err != nil {
		return Localized{}, err
	}

	local := Localized{}

	for _, e := range pk.Entries {
		ll, err := filepath.Localize(strings.ReplaceAll(e.Path, `\`, `/`))
		if err != nil {
			return Localized{}, err
		}
		l := filepath.Join(outdir, ll)
		local.Entries = append(local.Entries, parser.Entry{Path: l})

		if err = os.MkdirAll(filepath.Dir(l), 0755); err != nil {
			return Localized{}, err
		}

		slog.Info("Unpacking file", "path", ll)

		if _, err = reader.Seek(int64(e.Offset), io.SeekStart); err != nil {
			return Localized{}, err
		}

		of, err := os.Create(l)
		if err != nil {
			return Localized{}, err
		}

		bb := make([]byte, e.Size)
		if _, err = reader.Read(bb); err != nil {
			of.Close()
			return Localized{}, err
		}

		if _, err = of.Write(bb); err != nil {
			of.Close()
			return Localized{}, err
		}
		of.Close()
	}

	return local, nil
}

func Pack(local Localized, writer io.Writer) error {
	pk := parser.PakFile{
		EntryCount: int32(len(local.Entries)),
		Entries:    []parser.Entry{},
	}

	prevOffset := 4 + 268*len(local.Entries)
	for _, v := range local.Entries {
		fif, err := os.Stat(v.Path)
		if err != nil {
			return err
		}

		newPath := filepath.Join(strings.Split(v.Path, "/")[1:]...)              // remove first directory
		newPath = strings.TrimPrefix(strings.ReplaceAll(newPath, `/`, `\`), `\`) // to windows paths

		e := parser.Entry{
			Size:   int32(fif.Size()),
			Offset: int32(prevOffset),
			Path:   newPath,
		}

		pk.Entries = append(pk.Entries, e)
		prevOffset += int(fif.Size())
	}

	var err error
	if err = binary.Write(writer, binary.LittleEndian, pk.EntryCount); err != nil {
		return err
	}

	for _, v := range pk.Entries {

		dd := charmap.Windows1251.NewEncoder()
		np, err := dd.Bytes([]byte(v.Path))
		if err != nil {
			return err
		}

		if _, err = writer.Write(np); err != nil {
			return err
		}

		if _, err = writer.Write([]byte{0x0}); err != nil {
			return err
		}

		pad := bytes.Repeat([]byte{0xCC}, 260-len([]byte(np))-1)
		if _, err = writer.Write(pad); err != nil {
			return err
		}

		if err = binary.Write(writer, binary.LittleEndian, v.Size); err != nil {
			return err
		}

		if err = binary.Write(writer, binary.LittleEndian, v.Offset); err != nil {
			return err
		}
	}

	for _, v := range local.Entries {
		r, err := os.Open(v.Path)
		if err != nil {
			return err
		}

		slog.Info("writing file", "name", v.Path)
		if _, err = io.Copy(writer, r); err != nil {
			r.Close()
			return err
		}

		r.Close()
	}

	return nil
}
