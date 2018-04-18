// Copyright (c) 2016-2017 Tigera, Inc. All rights reserved.

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

package pod

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// IFNameForKey returns a deterministic interface name for the given key
// This is a copy-and-modified version of libcalico-go VethNameForWorkload
func IFNameForKey(key string) string {
	// A SHA1 is always 20 bytes long, and so is sufficient for generating the
	// veth name and mac addr.
	h := sha1.New()
	h.Write([]byte(key))
	return fmt.Sprintf("mido%s", hex.EncodeToString(h.Sum(nil))[:11])
}
