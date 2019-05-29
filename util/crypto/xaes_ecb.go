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

package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length <= 0 {
		return origData, nil
	}

	unpadding := int(origData[length-1])
	if length < unpadding {
		err := fmt.Errorf("pkcs5UnPadding failed, invalid padding length")
		return nil, err
	}

	return origData[:(length - unpadding)], nil
}

type ecbEncrypter struct {
	b         cipher.Block
	blockSize int
}

// EncryptECB encrypt data in ecb mode
func EncryptECB(key string, data []byte) (result []byte, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	dataPadding := pkcs5Padding(data, block.BlockSize())
	result = make([]byte, len(dataPadding))

	mode := newECBEncrypter(block)
	err = mode.cryptBlocks(result, dataPadding)
	if err != nil {
		return
	}

	return
}

// DecryptECB decrypt data in ecb mode
func DecryptECB(key string, data []byte) (result []byte, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	mode := newECBDecrypter(block)
	err = mode.cryptBlocks(data, data)
	if err != nil {
		return
	}

	result, err = pkcs5UnPadding(data)
	if err != nil {
		return
	}

	return
}

func newECBEncrypter(b cipher.Block) (c *ecbEncrypter) {
	c = new(ecbEncrypter)
	c.b = b
	c.blockSize = b.BlockSize()
	return
}

func (e *ecbEncrypter) BlockSize() (size int) {
	size = e.blockSize
	return
}

func (e *ecbEncrypter) cryptBlocks(dst []byte, src []byte) (err error) {
	if len(src)%e.blockSize != 0 {
		err = fmt.Errorf("len(src)[%v] not full blocks, blocksize[%v]", len(src), e.blockSize)
		return
	}

	if len(dst) < len(src) {
		err = fmt.Errorf("output[%v] is smaller than input[%v]", len(dst), len(src))
		return
	}

	for len(src) > 0 {
		e.b.Encrypt(dst, src[:e.blockSize])
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}
	return
}

type ecbDecrypter struct {
	b         cipher.Block
	blockSize int
}

func newECBDecrypter(b cipher.Block) (c *ecbDecrypter) {
	c = new(ecbDecrypter)
	c.b = b
	c.blockSize = b.BlockSize()
	return
}

func (e *ecbDecrypter) BlockSize() (size int) {
	size = e.blockSize
	return
}

func (e *ecbDecrypter) cryptBlocks(dst []byte, src []byte) (err error) {
	if len(src)%e.blockSize != 0 {
		err = fmt.Errorf("len(src)[%v] not full blocks, blocksize[%v]", len(src), e.blockSize)
		return
	}
	if len(dst) < len(src) {
		err = fmt.Errorf("output[%v] is smaller than input[%v]", len(dst), len(src))
		return
	}

	for len(src) > 0 {
		e.b.Decrypt(dst, src[:e.blockSize]) // different from encrypter here.
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}

	return
}
