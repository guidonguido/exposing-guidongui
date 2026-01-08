package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

const METADATA_PREFIX = "{{"
const METADATA_SUFFIX = "}}"
const METADATA_TITLE = "Title"
const METADATA_AUTHOR = "Author"
const METADATA_DATE = "Date"
const METADATA_SUMMARY = "Summary"

func (fr FileReader) ReadMetadata(id string, data *PostData) error {
	f, err := os.Open(fmt.Sprintf("web/posts/%s.md", id))
	if err != nil {
		log.Printf("Error opening post file: %v", err)
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	// OxA indicates newline
	ln, err := reader.ReadBytes(0xA)
	for err == nil {

		if bytes.HasPrefix(ln, []byte("#")) {
			// Reached content section
			break
		}

		// Metadata lines start with {{ (0x7B7B)
		if bytes.HasPrefix(ln, []byte{0x7B, 0x7B}) {
			parts := bytes.SplitN(

				bytes.TrimSuffix(
					bytes.TrimPrefix(ln, []byte(METADATA_PREFIX)),
					[]byte(METADATA_SUFFIX+"\n")),
				[]byte(":"),
				2) // Split into key and value
			k := string(bytes.TrimSpace(parts[0]))
			v := string(bytes.TrimSpace(parts[1]))

			switch k {
			case METADATA_TITLE:
				data.Title = v
			case METADATA_AUTHOR:
				data.Author = v
			case METADATA_DATE:
				data.Date = v
			case METADATA_SUMMARY:
				data.Summary = v
			}
		}
		ln, err = reader.ReadBytes(0xA)
	}

	return nil
}
