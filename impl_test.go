// Copyright (c) 2019 Oasis Labs Inc. <info@oasislabs.com>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package deoxysii

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/oasislabs/deoxysii/internal/api"
	"github.com/oasislabs/deoxysii/internal/ct64"
	"github.com/oasislabs/deoxysii/internal/hardware"
	"github.com/oasislabs/deoxysii/internal/vartime"
)

var testImpls = []api.Impl{
	ct64.Impl,
	vartime.Impl,
}

func BenchmarkDeoxysII(b *testing.B) {
	oldImpl := impl
	defer func() {
		impl = oldImpl
	}()

	for _, testImpl := range testImpls {
		impl = testImpl
		doBenchmarkDeoxysII(b)
	}
}

func doBenchmarkDeoxysII(b *testing.B) {
	benchSizes := []int{8, 32, 64, 576, 1536, 4096, 1024768}

	for _, sz := range benchSizes {
		bn := "DeoxysII_" + impl.Name() + "_"
		sn := fmt.Sprintf("_%d", sz)
		b.Run(bn+"Encrypt"+sn, func(b *testing.B) { doBenchmarkAEADEncrypt(b, sz) })
		b.Run(bn+"Decrypt"+sn, func(b *testing.B) { doBenchmarkAEADDecrypt(b, sz) })
	}
}

func doBenchmarkAEADEncrypt(b *testing.B, sz int) {
	b.StopTimer()
	b.SetBytes(int64(sz))

	nonce, key := make([]byte, NonceSize), make([]byte, KeySize)
	m, c := make([]byte, sz), make([]byte, 0, sz+TagSize)
	_, _ = rand.Read(nonce)
	_, _ = rand.Read(key)
	_, _ = rand.Read(m)
	aead, _ := New(key)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c = c[:0]

		c = aead.Seal(c, nonce, m, nil)
		if len(c) != sz+TagSize {
			b.Fatalf("Seal failed")
		}
	}
}

func doBenchmarkAEADDecrypt(b *testing.B, sz int) {
	b.StopTimer()
	b.SetBytes(int64(sz))

	nonce, key := make([]byte, NonceSize), make([]byte, KeySize)
	m, c, d := make([]byte, sz), make([]byte, 0, sz+TagSize), make([]byte, 0, sz)
	_, _ = rand.Read(nonce)
	_, _ = rand.Read(key)
	_, _ = rand.Read(m)
	aead, _ := New(key)

	c = aead.Seal(c, nonce, m, nil)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		d = d[:0]

		var err error
		d, err = aead.Open(d, nonce, c, nil)
		if err != nil {
			b.Fatalf("Open failed")
		}
	}
	b.StopTimer()

	if !bytes.Equal(m, d) {
		b.Fatalf("Open output mismatch")
	}
}

func init() {
	if hardware.Impl != nil {
		testImpls = append([]api.Impl{hardware.Impl}, testImpls...)
	}
}
