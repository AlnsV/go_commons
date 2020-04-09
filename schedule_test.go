package acemon

import (
	"testing"
)

func TestNewIndexSchedule(t *testing.T) {
	idx := NewIndexSchedule("http://3.18.103.89:5100/")
	err := idx.UpdateSchedule()
	if err != nil {
		t.Log(err)
	}
	errTwo := idx.UpdateSchedule()
	if errTwo != nil {
		t.Log(errTwo)
	}
	idxTwo := NewIndexSchedule("http://3.18.103.89:5200/")
	err = idxTwo.UpdateSchedule()
	if err == nil {
		t.Log(err)
	}
	errTwo = idx.UpdateSchedule()
	if errTwo == nil {
		t.Log(errTwo)
	}
}

func BenchmarkIndexSchedule_UpdateSchedule(b *testing.B) {
	idx := NewIndexSchedule("http://3.18.103.89:5100/")
	for i := 0; i < b.N; i++ {
		errTwo := idx.UpdateSchedule()
		if errTwo == nil {
			b.Log(errTwo)
		}
	}
}
