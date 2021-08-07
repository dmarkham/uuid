// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestNewUUIDv6(t *testing.T) {
	SetRand(badRand{})
	clockSeq = 0
	// Fake time.Now for this test to return a monotonically advancing time; restore it at end.
	defer func(orig func() time.Time) { timeNow = orig }(timeNow)
	monTime := time.Date(2021, 01, 4, 1, 50, 0, 0, time.UTC)
	timeNow = func() time.Time {
		monTime = monTime.Add(1 * time.Second)
		return monTime
	}

	tests := []struct {
		name    string
		want    UUID
		wantErr bool
	}{
		{
			want: UUID{1, 235, 78, 47, 131, 98, 106, 128, 128, 1, 0, 1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUUIDv6()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUUIDV6() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//fmt.Println(got.Time(), tt.want.Time())
			if !reflect.DeepEqual(got.NodeID(), tt.want.NodeID()) {
				t.Errorf("NewUUIDv6() Node MisMatch = %v, wantErr %v", got.NodeID(), tt.want.NodeID())
				return
			}
			if got.Time() != tt.want.Time() {
				t.Errorf("NewUUIDv6() Time MisMatch = %v, wantErr %v", got.Time(), tt.want.Time())
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUUIDv6() = %v, want %v", got, tt.want)
			}

		})
	}
}
func TestVersion6(t *testing.T) {
	SetRand(badRand{})
	uuid1, err := NewUUIDv6()
	if err != nil {
		t.Fatalf("could not create UUID: %v", err)
	}
	uuid2, err := NewUUIDv6()
	if err != nil {
		t.Fatalf("could not create UUID: %v", err)
	}

	if uuid1 == uuid2 {
		t.Errorf("%s:duplicate uuid", uuid1)
	}
	if v := uuid1.Version(); v != 6 {
		t.Errorf("%s: version %s expected 1", uuid1, v)
	}
	if v := uuid2.Version(); v != 6 {
		t.Errorf("%s: version %s expected 1", uuid2, v)
	}
	n1 := uuid1.NodeID()
	n2 := uuid2.NodeID()
	if !bytes.Equal(n1, n2) {
		t.Errorf("Different nodes %x != %x", n1, n2)
	}
	t1 := uuid1.Time()
	t2 := uuid2.Time()

	q1 := uuid1.ClockSequence()
	q2 := uuid2.ClockSequence()

	switch {
	case t1 == t2 && q1 == q2:
		t.Error("time stopped")
	case t1 > t2 && q1 == q2:
		t.Error("time reversed")
	case t1 < t2 && q1 != q2:
		t.Error("clock sequence changed unexpectedly")
	}
}
