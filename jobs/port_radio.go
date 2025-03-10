package jobs

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	pa "github.com/gordonklaus/portaudio"
)

const rate = 44100
const seconds = 10 

// borked
func PlayFromFile(fname string) {
    path, err := filepath.Abs(fname)
    chk(err)

	f, err := os.Open(path)
	chk(err)
    defer f.Close()

    id, data, err := readChunk(f)
    chk(err)
    if id.String() != "FORM" {
        log.Println("Bad file format")
        return
    }

    _, err = data.Read(id[:])
    chk(err)
    if id.String() != "AIFF" {
        log.Println("Bad file format")
        return
    }

    var c commonChunk
    var audio io.Reader
	for {
		id, chunk, err := readChunk(data)
		if err == io.EOF {
			break
		}
		chk(err)
		switch id.String() {
		case "COMM":
			chk(binary.Read(chunk, binary.BigEndian, &c))
		case "SSND":
			chunk.Seek(8, 1) //ignore offset and block
			audio = chunk
		default:
			log.Printf("Ignoring unknown chunk '%s'\n", id)
		}
	}

	pa.Initialize()
	defer pa.Terminate()
    out := make([]int32, 8192)
    stream, err := pa.OpenDefaultStream(0, 1, rate, len(out), &out)
    chk(err)
    defer stream.Close()

    chk(stream.Start())
    defer stream.Stop()
    for remaining := int(c.NumSamples); remaining > 0; remaining -= len(out) {
		if len(out) > remaining {
			out = out[:remaining]
		}
		err := binary.Read(audio, binary.BigEndian, out)
		if err == io.EOF {
			break
		}
		chk(err)
		chk(stream.Write())
		select {
		case <-time.After(5 * time.Minute):
			return
		default:
		}
    }
}

func readChunk(r readerAtSeeker) (id ID, data *io.SectionReader, err error) {
	_, err = r.Read(id[:])
	if err != nil {
		return
	}
	var n int32
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return
	}
	off, _ := r.Seek(0, 1)
	data = io.NewSectionReader(r, off, int64(n))
	_, err = r.Seek(int64(n), 1)
	return
}

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

type ID [4]byte

func (id ID) String() string {
	return string(id[:])
}

type commonChunk struct {
	NumChans      int16
	NumSamples    int32
	BitsPerSample int16
	SampleRate    [10]byte
}
