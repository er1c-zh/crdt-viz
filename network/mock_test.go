package network

import (
	"testing"
	"time"
)

func TestMock(t *testing.T) {
	cloud := NewCloud()

	f := func(addr string) Interface {
		i, err := cloud.NewClient(addr)
		if err != nil {
			t.Error(err)
			return nil
		}
		go func() {
			ch, _ := i.GetRecv()
			for {
				select {
				case msg := <-ch:
					t.Logf("[%s->%s] get %s", msg.From, msg.To, msg.Msg)
				}
			}
		}()
		return i
	}
	i1 := f("1")
	i2 := f("2")
	if i1 == nil || i2 == nil {
		return
	}

	err := i1.Enable()
	if err != nil {
		t.Error(err)
		return
	}
	err = i2.Enable()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("start!")

	t1 := time.NewTicker(1 * time.Second)
	t2 := time.NewTicker(300 * time.Millisecond)
	for {
		select {
		case <-t1.C:
			err = i1.Write("2", "hello")
			if err != nil {
				t.Error(err)
				return
			}
		case <-t2.C:
			err = i2.Write("1", "t2")
			if err != nil {
				t.Error(err)
				return
			}
		}
	}
}
