// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
)

type wbuff interface {
	io.Writer
	//io.Seeker
	//io.WriterAt
	//Len() int
}

// WBuffer is a write-only ROOT buffer for streaming.
type WBuffer struct {
	w      wbuff
	err    error
	offset uint32
	refs   map[int64]interface{}
}

func NewWBufferFrom(w wbuff, refs map[int64]interface{}, offset uint32) *WBuffer {
	if refs == nil {
		refs = make(map[int64]interface{})
	}
	return &WBuffer{
		w:      w,
		refs:   refs,
		offset: offset,
	}
}

func NewWBuffer(data []byte, refs map[int64]interface{}, offset uint32) *WBuffer {
	if refs == nil {
		refs = make(map[int64]interface{})
	}
	return &WBuffer{
		w:      bytes.NewBuffer(data),
		refs:   refs,
		offset: offset,
	}
}

/*
func (w *WBuffer) Pos() int64 {
	//pos, _ := w.w.Seek(0, ioSeekCurrent)
	pos := int64(w.w.Len())
	return pos + int64(w.offset)
}

func (r *WBuffer) setPos(pos int64) error {
	pos -= int64(r.offset)
	got, err := r.w.Seek(pos, ioSeekStart)
	if err != nil {
		return err
	}
	if got != pos {
		return errorf("rootio: WBuffer too short (got=%v want=%v)", got, pos)
	}
	return nil
}
*/

func (w *WBuffer) Err() error {
	return w.err
}

func (w *WBuffer) write(p []byte) {
	if w.err != nil {
		return
	}
	_, w.err = w.w.Write(p)
}

func (w *WBuffer) WriteString(s string) {
	if w.err != nil {
		return
	}

	switch {
	case len(s) > 254: // large string
		w.WriteU8(255)
		w.WriteU32(uint32(len(s)))
	default:
		w.WriteU8(uint8(len(s)))
	}
	if s == "" {
		w.WriteU8(0)
	}
	w.write([]byte(s))
}

func (w *WBuffer) WriteI8(v int8) {
	if w.err != nil {
		return
	}
	var buf = [1]byte{byte(v)}
	_, w.err = w.w.Write(buf[:])
}

func (w *WBuffer) WriteI16(v int16) {
	if w.err != nil {
		return
	}

	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], uint16(v))
	_, w.err = w.w.Write(buf[:])
}

func (w *WBuffer) WriteI32(v int32) {
	if w.err != nil {
		return
	}

	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(v))
	_, w.err = w.w.Write(buf[:])
}

func (w *WBuffer) WriteI64(v int64) {
	if w.err != nil {
		return
	}

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(v))
	_, w.err = w.w.Write(buf[:])
}

func (w *WBuffer) WriteU8(v uint8) {
	if w.err != nil {
		return
	}
	var buf = [1]byte{v}
	_, w.err = w.w.Write(buf[:])
}

func (w *WBuffer) WriteU16(v uint16) {
	if w.err != nil {
		return
	}

	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], v)
	_, w.err = w.w.Write(buf[:])
}

func (w *WBuffer) WriteU32(v uint32) {
	if w.err != nil {
		return
	}

	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], v)
	_, w.err = w.w.Write(buf[:])
}

func (w *WBuffer) WriteU64(v uint64) {
	if w.err != nil {
		return
	}

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], v)
	_, w.err = w.w.Write(buf[:])
}

func (w *WBuffer) WriteF32(v float32) {
	if w.err != nil {
		return
	}

	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], math.Float32bits(v))
	_, w.err = w.w.Write(buf[:])
}

func (w *WBuffer) WriteF64(v float64) {
	if w.err != nil {
		return
	}

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(v))
	_, w.err = w.w.Write(buf[:])
}
