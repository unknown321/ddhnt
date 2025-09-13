package coder

import "io"

func Encode(reader io.Reader, output io.Writer) error {
	buffer := make([]byte, 4096)
	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		for i := 0; i < n; i++ {
			buffer[i] ^= 0xFF
		}

		_, err = output.Write(buffer[:n])
		if err != nil {
			return err
		}
	}

	return nil
}
