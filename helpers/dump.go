// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/20 21:39:17

package helpers

import (
	"encoding/gob"
	"os"
)

// Marshal marshals object to dumpFile.
func Marshal(object interface{}, dumpFile string) error {
	file, err := os.OpenFile(dumpFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	return encoder.Encode(object)
}

// Unmarshal parses the dumpFile and stores the result in the value pointed to by object.
func Unmarshal(object interface{}, dumpFile string) error {
	file, err := os.Open(dumpFile)
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(file)
	return decoder.Decode(object)
}
