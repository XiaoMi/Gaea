// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

const (
	c1 = 0xcc9e2d51
	c2 = 0x1b873593
)

// MurmurHash copy from guava Murmur3_32HashFunction
type MurmurHash struct {
	seed int
}

// NewMurmurHash constructor of MurmurHash
func NewMurmurHash(seed int) *MurmurHash {
	return &MurmurHash{
		seed: seed,
	}
}

// HashUnencodedChars check if has unencoded chars
func (m *MurmurHash) HashUnencodedChars(inputStr string) int {
	input := []rune(inputStr)
	h1 := int32(m.seed)

	for i := 1; i < len(input); i += 2 {
		k1 := int32(input[i-1] | (input[i] << 16))
		k1 = mixK1(k1)
		h1 = mixH1(h1, k1)
	}

	if len(input)&1 == 1 {
		k1 := int32(input[len(input)-1])
		k1 = mixK1(k1)
		h1 ^= k1
	}

	return int(fmix(h1, int32(2*len(input))))
}

func mixK1(k1 int32) int32 {
	k1_64 := int64(k1)
	k1_64 *= c1
	k1 = int32(k1_64)
	k1 = rotateLeft(k1, 15)
	k1_64 = int64(k1)
	k1_64 *= c2
	return int32(k1_64)
}

func mixH1(h1, k1 int32) int32 {
	h1 ^= k1
	h1 = rotateLeft(h1, 13)
	h1_64 := int64(h1)
	h1_64 = h1_64*5 + 0xe6546b64
	return int32(h1_64)
}

func fmix(h1, length int32) int32 {
	uh1 := uint32(h1)
	ulength := uint32(length)
	uh1 ^= ulength
	uh1 ^= uh1 >> 16
	uh1 *= 0x85ebca6b
	uh1 ^= uh1 >> 13
	uh1 *= 0xc2b2ae35
	uh1 ^= uh1 >> 16
	return int32(uh1)
}

func rotateLeft(i, distance int32) int32 {
	ui := uint32(i)
	udistance := uint32(distance)
	a1 := ui << udistance
	b1 := ui >> (32 - udistance)
	return int32(a1 | b1)
}
