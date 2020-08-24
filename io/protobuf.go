package io

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/klauspost/pgzip"

	"github.com/dreadl0ck/netcap/defaults"
	"github.com/dreadl0ck/netcap/delimited"
	"github.com/dreadl0ck/netcap/types"
)

// ProtoWriter is a structure that supports writing protobuf audit records to disk.
type ProtoWriter struct {
	bWriter *bufio.Writer
	gWriter *pgzip.Writer
	dWriter *delimited.Writer
	pWriter *DelimitedProtoWriter

	file *os.File
	mu   sync.Mutex
	wc   *WriterConfig
}

// NewProtoWriter initializes and configures a new ProtoWriter instance.
func NewProtoWriter(wc *WriterConfig) *ProtoWriter {
	w := &ProtoWriter{}
	w.wc = wc

	if wc.MemBufferSize <= 0 {
		wc.MemBufferSize = defaults.BufferSize
	}

	if wc.Compress {
		w.file = createFile(filepath.Join(wc.Out, wc.Name), ".ncap.gz")
	} else {
		w.file = createFile(filepath.Join(wc.Out, wc.Name), ".ncap")
	}

	// buffer data?
	if wc.Buffer {
		if wc.Compress {
			// experiment: pgzip -> file
			var errGzipWriter error
			w.gWriter, errGzipWriter = pgzip.NewWriterLevel(w.file, defaults.CompressionLevel)

			if errGzipWriter != nil {
				panic(errGzipWriter)
			}
			// experiment: buffer -> pgzip
			w.bWriter = bufio.NewWriterSize(w.gWriter, defaults.BufferSize)
			// experiment: delimited -> buffer
			w.dWriter = delimited.NewWriter(w.bWriter)
		} else {
			w.bWriter = bufio.NewWriterSize(w.file, defaults.BufferSize)
			w.dWriter = delimited.NewWriter(w.bWriter)
		}
	} else {
		if w.wc.Compress {
			var errGzipWriter error
			w.gWriter, errGzipWriter = pgzip.NewWriterLevel(w.file, defaults.CompressionLevel)
			if errGzipWriter != nil {
				panic(errGzipWriter)
			}
			w.dWriter = delimited.NewWriter(w.gWriter)
		} else {
			w.dWriter = delimited.NewWriter(w.file)
		}
	}

	w.pWriter = NewDelimitedProtoWriter(w.dWriter)

	if w.gWriter != nil {
		// To get any performance gains, you should at least be compressing more than 1 megabyte of data at the time.
		// You should at least have a block size of 100k and at least a number of blocks that match the number of cores
		// your would like to utilize, but about twice the number of blocks would be the best.
		if err := w.gWriter.SetConcurrency(defaults.CompressionBlockSize, runtime.GOMAXPROCS(0)*2); err != nil {
			log.Fatal("failed to configure compression package: ", err)
		}
	}

	return w
}

// WriteProto writes a protobuf message.
func (w *ProtoWriter) Write(msg proto.Message) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.pWriter.PutProto(msg)
}

// WriteHeader writes a netcap file header for protobuf encoded audit record files.
func (w *ProtoWriter) WriteHeader(t types.Type) error {
	return w.pWriter.PutProto(NewHeader(t, w.wc.Source, w.wc.Version, w.wc.IncludesPayloads, w.wc.StartTime))
}

// Close flushes and closes the writer and the associated file handles.
func (w *ProtoWriter) Close() (name string, size int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.wc.Buffer {
		flushWriters(w.bWriter)
	}

	if w.wc.Compress {
		closeGzipWriters(w.gWriter)
	}

	return closeFile(w.wc.Out, w.file, w.wc.Name)
}
