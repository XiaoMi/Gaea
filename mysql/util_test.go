// Copyright 2016 The kingshard Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

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

package mysql

import (
	"encoding/hex"
	"testing"

	"github.com/XiaoMi/Gaea/util/hack"
)

func TestCalcPassword(t *testing.T) {
	/*
		// **** JDBC ****
		seed:
			@jx=d_3z42;sS$YrS)p|
		hex:
			406a783d645f337a34323b73532459725329707c
		pass:
			kingshard
		scramble:
			fbc71db5ac3d7b51048d1a1d88c1677f34bcca11
	*/
	test, _ := RandomBuf(20)
	hexTest := hex.EncodeToString(test)
	t.Logf("rnd seed: %s, %s", hack.String(test), hexTest)

	seed := hack.Slice("@jx=d_3z42;sS$YrS)p|")
	hexSeed := hex.EncodeToString(seed)

	t.Logf("seed: %s equal %s, pass: %v", "406a783d645f337a34323b73532459725329707c", hexSeed, "406a783d645f337a34323b73532459725329707c" == hexSeed)
	scramble := CalcPassword(seed, hack.Slice("kingshard"))

	hexScramble := hex.EncodeToString(scramble)
	t.Logf("scramble: %s equal %s, pass: %v", "fbc71db5ac3d7b51048d1a1d88c1677f34bcca11", hexScramble, "fbc71db5ac3d7b51048d1a1d88c1677f34bcca11" == hexScramble)
}

func TestCalcPassword2(t *testing.T) {
	seed := hack.Slice("D?Y.SZC@v,P${ GdT0e{")
	hexSeed := hex.EncodeToString(seed)
	t.Logf("seed: %s equal %s, pass: %v", "443f592e535a4340762c50247b2047645430657b", hexSeed, "443f592e535a4340762c50247b2047645430657b" == hexSeed)
	scramble := CalcCachingSha2Password(seed, "123456")

	hexScramble := hex.EncodeToString(scramble)
	t.Logf("scramble: %s equal %s, pass: %v", "896d208c92429f5d1d5cb67f1ca3a639d7abdf335b05f58894e7f11d90608ca4", hexScramble, "896d208c92429f5d1d5cb67f1ca3a639d7abdf335b05f58894e7f11d90608ca4" == hexScramble)
}
