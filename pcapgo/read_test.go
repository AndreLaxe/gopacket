// Copyright 2014 Damjan Cvetko. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.
package pcapgo

import (
	"bytes"
	"testing"
	"time"
)

// test header read
func TestCreatePcapReader(t *testing.T) {
	test := []byte{
		0xd4, 0xc3, 0xb2, 0xa1, 0x02, 0x00, 0x04, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
	}
	buf := bytes.NewBuffer(test)
	_, err := NewReader(buf)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// test big endian file read
func TestCreatePcapReaderBigEndian(t *testing.T) {
	test := []byte{
		0xa1, 0xb2, 0xc3, 0xd4, 0x00, 0x02, 0x00, 0x04,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x01,
	}
	buf := bytes.NewBuffer(test)
	_, err := NewReader(buf)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// test opening invalid data
func TestCreatePcapReaderFail(t *testing.T) {
	test := []byte{
		0xd0, 0xc3, 0xb2, 0xa1, 0x02, 0x00, 0x04, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
	}
	buf := bytes.NewBuffer(test)
	_, err := NewReader(buf)
	if err == nil {
		t.Error("Should fail but did not")
		t.FailNow()
	}
}

func TestPacket(t *testing.T) {
	test := []byte{
		0xd4, 0xc3, 0xb2, 0xa1, 0x02, 0x00, 0x04, 0x00, // magic, maj, min
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // tz, sigfigs
		0xff, 0xff, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, // snaplen, linkType
		0x5A, 0xCC, 0x1A, 0x54, 0x01, 0x00, 0x00, 0x00, // sec, usec
		0x04, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, // cap len, full len
		0x01, 0x02, 0x03, 0x04, // data
	}

	buf := bytes.NewBuffer(test)
	r, err := NewReader(buf)
	if err != nil {
		t.Errorf("Failed to get new reader object: %v", err)
	}

	data, ci, err := r.ReadPacketData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if !ci.Timestamp.Equal(time.Date(2014, 9, 18, 12, 13, 14, 1000, time.UTC)) {
		t.Error("Invalid time read")
		t.FailNow()
	}
	if ci.CaptureLength != 4 || ci.Length != 8 {
		t.Error("Invalid CapLen or Len")
	}
	want := []byte{1, 2, 3, 4}
	if !bytes.Equal(data, want) {
		t.Errorf("buf mismatch:\nwant: %+v\ngot:  %+v", want, data)
	}
}

func TestPacketNano(t *testing.T) {
	test := []byte{
		0x4d, 0x3c, 0xb2, 0xa1, 0x02, 0x00, 0x04, 0x00, // magic, maj, min
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // tz, sigfigs
		0xff, 0xff, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, // snaplen, linkType
		0x5A, 0xCC, 0x1A, 0x54, 0x01, 0x00, 0x00, 0x00, // sec, usec
		0x04, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, // cap len, full len
		0x01, 0x02, 0x03, 0x04, // data
	}

	buf := bytes.NewBuffer(test)
	r, err := NewReader(buf)
	if err != nil {
		t.Errorf("Failed to get new reader object: %v", err)
	}

	data, ci, err := r.ReadPacketData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if !ci.Timestamp.Equal(time.Date(2014, 9, 18, 12, 13, 14, 1, time.UTC)) {
		t.Error("Invalid time read")
		t.FailNow()
	}
	if ci.CaptureLength != 4 || ci.Length != 8 {
		t.Error("Invalid CapLen or Len")
	}
	want := []byte{1, 2, 3, 4}
	if !bytes.Equal(data, want) {
		t.Errorf("buf mismatch:\nwant: %+v\ngot:  %+v", want, data)
	}
}

func TestGzipPacket(t *testing.T) {
	test := []byte{
		0x1f, 0x8b, 0x08, 0x08, 0x92, 0x4d, 0x81, 0x57,
		0x00, 0x03, 0x74, 0x65, 0x73, 0x74, 0x00, 0xbb,
		0x72, 0x78, 0xd3, 0x42, 0x26, 0x06, 0x16, 0x06,
		0x18, 0xf8, 0xff, 0x9f, 0x81, 0x81, 0x11, 0x48,
		0x47, 0x9d, 0x91, 0x0a, 0x01, 0xd1, 0x20, 0x19,
		0x0e, 0x20, 0x66, 0x64, 0x62, 0x66, 0x01, 0x00,
		0xe4, 0x76, 0x9b, 0x75, 0x2c, 0x00, 0x00, 0x00,
	}

	buf := bytes.NewBuffer(test)
	r, err := NewReader(buf)
	if err != nil {
		t.Error("Unexpected error returned:", err)
		t.FailNow()
	}

	data, ci, err := r.ReadPacketData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if !ci.Timestamp.Equal(time.Date(2014, 9, 18, 12, 13, 14, 1000, time.UTC)) {
		t.Error("Invalid time read")
		t.FailNow()
	}
	if ci.CaptureLength != 4 || ci.Length != 8 {
		t.Error("Invalid CapLen or Len")
	}
	want := []byte{1, 2, 3, 4}
	if !bytes.Equal(data, want) {
		t.Errorf("buf mismatch:\nwant: %+v\ngot:  %+v", want, data)
	}
}

func TestTruncatedGzipPacket(t *testing.T) {
	test := []byte{
		0x1f, 0x8b, 0x08,
	}

	buf := bytes.NewBuffer(test)

	if _, err := NewReader(buf); err == nil {
		t.Error("Should fail but did not")
		t.FailNow()
	}
}

func TestPacketBufferReuse(t *testing.T) {
	test := []byte{
		0xd4, 0xc3, 0xb2, 0xa1, 0x02, 0x00, 0x04, 0x00, // magic, maj, min
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // tz, sigfigs
		0xff, 0xff, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, // snaplen, linkType
		0x5A, 0xCC, 0x1A, 0x54, 0x01, 0x00, 0x00, 0x00, // sec, usec
		0x04, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, // cap len, full len
		0x01, 0x02, 0x03, 0x04, // data
		0x5A, 0xCC, 0x1A, 0x54, 0x01, 0x00, 0x00, 0x00, // sec, usec
		0x04, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, // cap len, full len
		0x01, 0x02, 0x03, 0x04, // data
	}

	buf := bytes.NewBuffer(test)
	r, err := NewReader(buf)
	if err != nil {
		t.Errorf("Failed to get new reader object: %v", err)
	}

	data1, _, err := r.ReadPacketData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if want := []byte{1, 2, 3, 4}; !bytes.Equal(data1, want) {
		t.Errorf("buf mismatch:\nwant: %+v\ngot:  %+v", want, data1)
	}
	data2, _, err := r.ReadPacketData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for i := range data1 {
		data1[i] = 0xff // modify data1 after getting data2, make sure we don't overlap buffers.
	}
	if want := []byte{1, 2, 3, 4}; !bytes.Equal(data2, want) {
		t.Errorf("buf mismatch:\nwant: %+v\ngot:  %+v", want, data2)
	}
}

func TestPacketZeroCopy(t *testing.T) {
	test := []byte{
		0xd4, 0xc3, 0xb2, 0xa1, 0x02, 0x00, 0x04, 0x00, // magic, maj, min
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // tz, sigfigs
		0xff, 0xff, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, // snaplen, linkType
		0x5A, 0xCC, 0x1A, 0x54, 0x01, 0x00, 0x00, 0x00, // sec, usec
		0x04, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, // cap len, full len
		0x01, 0x02, 0x03, 0x04, // data
		0x5A, 0xCC, 0x1A, 0x54, 0x01, 0x00, 0x00, 0x00, // sec, usec
		0x04, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, // cap len, full len
		0x05, 0x06, 0x07, 0x08, // data
	}

	buf := bytes.NewBuffer(test)
	r, err := NewReader(buf)
	if err != nil {
		t.Errorf("Failed to get new reader object: %v", err)
	}

	data1, _, err := r.ZeroCopyReadPacketData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if want := []byte{1, 2, 3, 4}; !bytes.Equal(data1, want) {
		t.Errorf("buf mismatch:\nwant: %+v\ngot:  %+v", want, data1)
	}
	data2, _, err := r.ZeroCopyReadPacketData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if want := []byte{5, 6, 7, 8}; !bytes.Equal(data2, want) {
		t.Errorf("buf mismatch:\nwant: %+v\ngot:  %+v", want, data2)
	}

	if &data1[0] != &data2[0] {
		t.Error("different buffers returned by subsequent ZeroCopyReadPacketData calls")
	}
}
