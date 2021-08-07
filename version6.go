// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"encoding/binary"
)

// NewUUIDv6 returns a Version 6 UUID based on the current clock
// sequence, and the current time.  NodeID is sudorandom bytes
// or SetNodeInterface then it will be set automatically.
// If clock sequence has not been set by
// SetClockSequence then it will be set automatically.  If GetTime fails to
// return the current NewUUIDv6 returns nil and an error.
//
func NewUUIDv6() (UUID, error) {
	var uuid UUID
	now, seq, err := GetTime()
	if err != nil {
		return uuid, err
	}

	timeLow := uint16(now & 0x0fff)
	timeLow |= 0x6000 // Version 6
	timeMid := uint16((now >> 12) & 0xffff)
	timeHi := uint32((now >> 32) & 0xffffffff)
	binary.BigEndian.PutUint32(uuid[0:], timeHi)
	binary.BigEndian.PutUint16(uuid[4:], timeMid)
	binary.BigEndian.PutUint16(uuid[6:], timeLow)
	binary.BigEndian.PutUint16(uuid[8:], seq)
	nodeID := [6]byte{}
	randomBits(nodeID[:])
	copy(uuid[10:], nodeID[:])

	return uuid, nil
}
