// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wycheproof

import (
	"bytes"
	"io"
	"testing"

	"golang.org/x/crypto/hkdf"
)

func TestHkdf(t *testing.T) {

	// HkdfTestVector
	type HkdfTestVector struct {

		// A brief description of the test case
		Comment string `json:"comment,omitempty"`

		// A list of flags
		Flags []string `json:"flags,omitempty"`

		// the key (input key material)
		Ikm string `json:"ikm,omitempty"`

		// additional information used in the key derivation
		Info string `json:"info,omitempty"`

		// the generated bytes (output key material)
		Okm string `json:"okm,omitempty"`

		// Test result
		Result string `json:"result,omitempty"`

		// the salt for the key derivation
		Salt string `json:"salt,omitempty"`

		// the size of the output in bytes
		Size int `json:"size,omitempty"`

		// Identifier of the test case
		TcId int `json:"tcId,omitempty"`
	}

	// Notes a description of the labels used in the test vectors
	type Notes struct {
	}

	// HkdfTestGroup
	type HkdfTestGroup struct {

		// the size of the ikm in bits
		KeySize int               `json:"keySize,omitempty"`
		Tests   []*HkdfTestVector `json:"tests,omitempty"`
		Type    any               `json:"type,omitempty"`
	}

	// Root
	type Root struct {

		// the primitive tested in the test file
		Algorithm string `json:"algorithm,omitempty"`

		// the version of the test vectors.
		GeneratorVersion string `json:"generatorVersion,omitempty"`

		// additional documentation
		Header []string `json:"header,omitempty"`

		// a description of the labels used in the test vectors
		Notes *Notes `json:"notes,omitempty"`

		// the number of test vectors in this test
		NumberOfTests int              `json:"numberOfTests,omitempty"`
		Schema        any              `json:"schema,omitempty"`
		TestGroups    []*HkdfTestGroup `json:"testGroups,omitempty"`
	}

	fileHashAlgorithms := map[string]string{
		"hkdf_sha1_test.json":   "SHA-1",
		"hkdf_sha256_test.json": "SHA-256",
		"hkdf_sha384_test.json": "SHA-384",
		"hkdf_sha512_test.json": "SHA-512",
	}

	for f := range fileHashAlgorithms {
		var root Root
		readTestVector(t, f, &root)
		for _, tg := range root.TestGroups {
			for _, tv := range tg.Tests {
				h := parseHash(fileHashAlgorithms[f]).New
				hkdf := hkdf.New(h, decodeHex(tv.Ikm), decodeHex(tv.Salt), decodeHex(tv.Info))
				key := make([]byte, tv.Size)
				wantPass := shouldPass(tv.Result, tv.Flags, nil)
				_, err := io.ReadFull(hkdf, key)
				if (err == nil) != wantPass {
					t.Errorf("tcid: %d, type: %s, comment: %q, wanted success: %t, got: %v", tv.TcId, tv.Result, tv.Comment, wantPass, err)
				}
				if err != nil {
					continue // don't validate output text if reading failed
				}
				if got, want := key, decodeHex(tv.Okm); !bytes.Equal(got, want) {
					t.Errorf("tcid: %d, type: %s, comment: %q, output bytes don't match", tv.TcId, tv.Result, tv.Comment)
				}
			}
		}
	}
}
