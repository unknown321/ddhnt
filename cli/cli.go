package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/unknown321/ddhnt/coder"
	"github.com/unknown321/ddhnt/packer"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Unpack(filename string) {
	inFile, err := os.Open(filename)
	if err != nil {
		slog.Error("cannot open", "error", err.Error())
		os.Exit(1)
	}
	defer inFile.Close()

	decoded, err := os.CreateTemp("", "ddhunt-*")
	if err != nil {
		slog.Error("cannot create", "error", err.Error())
		os.Exit(1)
	}
	defer decoded.Close()
	defer os.Remove(decoded.Name())

	if err = coder.Encode(inFile, decoded); err != nil {
		slog.Error("cannot encode", "error", err.Error())
		os.Exit(1)
	}

	_, err = decoded.Seek(0, io.SeekStart)
	if err != nil {
		slog.Error("cannot seek", "error", err.Error())
		os.Exit(1)
	}

	var local packer.Localized
	if local, err = packer.Unpack(decoded, strings.ReplaceAll(filename, ".", "_")); err != nil {
		slog.Error("cannot unpack", "error", err.Error())
		os.Exit(1)
	}

	jj, err := json.MarshalIndent(local.Entries, "", "  ")
	if err != nil {
		slog.Error("cannot unpack", "error", err.Error())
		os.Exit(1)
	}

	js, err := os.Create(filename + ".json")
	if err != nil {
		slog.Error("cannot create json", "error", err.Error())
		os.Exit(1)
	}
	defer js.Close()

	if _, err = js.Write(jj); err != nil {
		slog.Error("cannot write json", "error", err.Error())
		os.Exit(1)
	}
}

func Pack(filename string) {
	d, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("cannot read json", "error", err.Error())
		os.Exit(1)
	}

	local := packer.Localized{}
	if err = json.Unmarshal(d, &local.Entries); err != nil {
		slog.Error("cannot unmarshal json", "error", err.Error())
		os.Exit(1)
	}

	tt, err := os.CreateTemp("", "ddhunt-*")
	if err != nil {
		slog.Error("cannot create temp file", "error", err.Error())
		os.Exit(1)
	}
	defer os.Remove(tt.Name())

	if err = packer.Pack(local, tt); err != nil {
		slog.Error("cannot pack", "error", err.Error())
		os.Exit(1)
	}

	oo, err := os.Create(strings.TrimSuffix(filename, ".json"))
	if err != nil {
		slog.Error("cannot create output file", "error", err.Error())
		os.Exit(1)
	}

	tt.Seek(0, io.SeekStart)

	if err = coder.Encode(tt, oo); err != nil {
		slog.Error("cannot encode", "error", err.Error())
		os.Exit(1)
	}
	oo.Close()
}

func Run() {
	var inName string
	flag.StringVar(&inName, "input", "", "input file")
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("Pack/unpack Deadhunt pk5 archives.")
		fmt.Println()
		fmt.Printf("Unpack: %s Data.pk5\n", os.Args[0])
		fmt.Printf("Pack: %s Data.pk5.json\n", os.Args[0])

		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if inName == "" {
		if len(os.Args) < 2 {
			slog.Error("no filename provided")
			flag.Usage()
			os.Exit(1)
		}

		inName = os.Args[1]
	}

	switch filepath.Ext(inName) {
	case ".json":
		Pack(inName)
	case ".pk5":
		Unpack(inName)
	default:
		slog.Error("Unrecognized file extension")
		os.Exit(1)
	}
}
