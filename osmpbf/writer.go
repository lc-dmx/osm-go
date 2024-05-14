package osmpbf

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"io"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/lc-dmx/osm-go/osmpbf/entity"
	"github.com/lc-dmx/osm-go/osmpbf/model_pb"
	"google.golang.org/protobuf/proto"
)

const (
	TYPE_OSM_DATA   = "OSMData"
	TYPE_OSM_HEADER = "OSMHeader"
)

const (
	FLUSH_SIZE = 16 * 1024 * 1024
)

type Writer struct {
	ctx context.Context

	w io.Writer

	encoder *Encoder
}

func NewWriter(ctx context.Context, w io.Writer) (*Writer, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	writer := &Writer{
		ctx:     ctx,
		w:       w,
		encoder: NewEncoder(),
	}
	headerBlockBytes, _ := proto.Marshal(writer.encoder.encodeHeader())
	if err := writer.writeHeader(headerBlockBytes); err != nil {
		return nil, err
	}

	return writer, nil
}

func (w *Writer) writeHeader(headerBlockBytes []byte) error {
	var in bytes.Buffer
	zw, err := zlib.NewWriterLevel(&in, zlib.BestCompression)
	if err != nil {
		return err
	}
	if _, err = zw.Write(headerBlockBytes); err != nil {
		return err
	}
	if err = zw.Close(); err != nil {
		return err
	}

	blob := &model_pb.Blob{
		RawSize: thrift.Int32Ptr(int32(len(headerBlockBytes))),
		Data:    &model_pb.Blob_ZlibData{ZlibData: in.Bytes()},
	}
	blobBytes, _ := proto.Marshal(blob)
	if err = w.writeBlob(TYPE_OSM_HEADER, blobBytes); err != nil {
		return err
	}

	return nil
}

func (w *Writer) writeData(data []byte) error {
	var in bytes.Buffer
	zw, err := zlib.NewWriterLevel(&in, zlib.BestCompression)
	if err != nil {
		return err
	}
	if _, err = zw.Write(data); err != nil {
		return err
	}
	if err = zw.Close(); err != nil {
		return err
	}

	blob := &model_pb.Blob{
		RawSize: thrift.Int32Ptr(int32(len(data))),
		Data:    &model_pb.Blob_ZlibData{ZlibData: in.Bytes()},
	}
	blobBytes, _ := proto.Marshal(blob)
	if err = w.writeBlob(TYPE_OSM_DATA, blobBytes); err != nil {
		return err
	}

	return nil
}

func (w *Writer) writeBlob(blobType string, blobBytes []byte) error {
	blobHeader := &model_pb.BlobHeader{
		Type:     thrift.StringPtr(blobType),
		Datasize: thrift.Int32Ptr(int32(len(blobBytes))),
	}
	blobHeaderBytes, _ := proto.Marshal(blobHeader)

	blobHeaderSize := make([]byte, 4)
	binary.BigEndian.PutUint32(blobHeaderSize, uint32(len(blobHeaderBytes)))

	if _, err := w.w.Write(blobHeaderSize); err != nil {
		return err
	}
	if _, err := w.w.Write(blobHeaderBytes); err != nil {
		return err
	}
	if _, err := w.w.Write(blobBytes); err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteEntity(exp entity.Exporter) error {
	w.encoder.encodeEntity(exp)

	// The length of the BlobHeader should be less than 32 KiB and must be less than 64 KiB.
	// The uncompressed length of a Blob should be less than 16 MiB and must be less than 32 MiB.
	if w.encoder.estimateBlockSize() >= FLUSH_SIZE {
		if err := w.Flush(); err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) Flush() error {
	primitiveBlock := w.encoder.encodeData()
	if primitiveBlock != nil {
		primitiveBlockBytes, _ := proto.Marshal(primitiveBlock)
		if err := w.writeData(primitiveBlockBytes); err != nil {
			return err
		}
	}

	w.encoder = NewEncoder()

	return nil
}

func (w *Writer) Close() error {
	return w.Flush()
}
