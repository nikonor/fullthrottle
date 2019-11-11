package fullthrottle

import (
	"fmt"
	"testing"
	"time"
)

func doSomething() {
	fmt.Println("\tDo something at " + time.Now().Format(time.RFC3339Nano))
}

func prn(str string) {
	fmt.Printf("%s at %s\n", str, time.Now().Format(time.RFC3339Nano))
}

func Test_One(t *testing.T) {
	d := make(chan struct{})
	tr := Throttling{MaxCount: 2, DoneChan: d, BeginDelay: 2 * time.Second}
	prn("Start")
	defer func() {
		prn("Finish")
	}()
	for {
		doSomething()
		ok := <-tr.Throttling()
		if ok {
			return
		}
	}
}

func Test_Two(t *testing.T) {
	d := make(chan struct{})
	tr := Throttling{DoneChan: d}
	prn("Start")
	defer func() {
		prn("Finish")
	}()
	c := 0
	for {
		doSomething()
		ok := <-tr.Throttling()
		if ok {
			return
		}
		if c == 2 {
			close(d)
		}
		c++
	}
}

func Test_Three(t *testing.T) {
	d := make(chan struct{})
	tr := Throttling{MaxCount: 4, DoneChan: d, BeginDelay: time.Second, IncDelay: time.Second}
	prn("Start")
	defer func() {
		prn("Finish")
	}()
	for {
		doSomething()
		ok := <-tr.Throttling()
		if ok {
			return
		}
	}
}

func Test_Four(t *testing.T) {
	d := make(chan struct{})
	tr := Throttling{MaxCount: 5, DoneChan: d, BeginDelay: time.Second, MultDelay: 2}

	prn("Start")
	defer func() {
		prn("Finish")
	}()

	for {
		doSomething()
		ok := <-tr.Throttling()
		if ok {
			return
		}
	}
}

func Test_Five(t *testing.T) {
	d := make(chan struct{})
	tr := Throttling{MaxCount: 6, DoneChan: d, BeginDelay: time.Second, MultDelay: 2, MaxDelay: 7 * time.Second}

	prn("Start")
	defer func() {
		prn("Finish")
	}()

	for {
		doSomething()
		if <-tr.Throttling() {
			return
		}
	}
}

func nextDelay(o Throttling) time.Duration {
	if o.count%2 == 0 {
		return time.Second
	}
	return 2 * time.Second
}

func Test_Six(t *testing.T) {
	d := make(chan struct{})
	tr := Throttling{MaxCount: 6, DoneChan: d, BeginDelay: time.Second, MaxDelay: 7 * time.Second}
	tr.NextDelay = nextDelay

	prn("Start")
	defer func() {
		prn("Finish")
	}()

	for {
		doSomething()
		if <-tr.Throttling() {
			return
		}
	}
}

func Test_Seven(t *testing.T) {
	d := make(chan struct{})
	tr := Throttling{MaxCount: 6, DoneChan: d, BeginDelay: 2 * time.Second}

	prn("Start")
	defer func() {
		prn("Finish")
	}()

	for i := 0; i < 5; i++ {
		doSomething()
		<-tr.Throttling()
	}
}
