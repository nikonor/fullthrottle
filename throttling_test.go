package fullthrottle

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func doSomething(last, diff int64) (bool, string) {
	n := time.Now().Unix()
	fmt.Println("\tDo something at " + time.Now().Format(time.RFC3339Nano))
	if last != 0 {
		if n-last != diff {
			return false, strconv.FormatInt(n, 10) + "-" + strconv.FormatInt(last, 10) + "!=" + strconv.FormatInt(diff, 10)
		}
	}
	return true, ""

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
	var last int64
	for {
		if ok, errStr := doSomething(last, 2); !ok {
			t.Error("Ошибка:" + errStr)
			return
		}
		last = time.Now().Unix()
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
	var last int64
	for {
		if ok, errStr := doSomething(last, 1); !ok {
			t.Error("Ошибка:" + errStr)
			return
		}
		last = time.Now().Unix()
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
	var (
		last int64
		diff int64
	)
	for {
		if ok, errStr := doSomething(last, diff); !ok {
			t.Error("Ошибка:" + errStr)
			return
		}
		diff += 1
		last = time.Now().Unix()
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

	var (
		last int64
	)
	diffs := []int64{1, 1, 2, 4, 8}
	count := 0
	for {
		if ok, errStr := doSomething(last, diffs[count]); !ok {
			t.Error("Ошибка:" + errStr)
			return
		}
		count++
		last = time.Now().Unix()
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

	var (
		last int64
	)
	diffs := []int64{1, 1, 2, 4, 7, 7}
	count := 0
	for {
		if ok, errStr := doSomething(last, diffs[count]); !ok {
			t.Error("Ошибка:" + errStr)
			return
		}
		count++
		last = time.Now().Unix()
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
	tr := Throttling{MaxCount: 8, DoneChan: d, BeginDelay: time.Second, MaxDelay: 7 * time.Second}
	tr.NextDelay = nextDelay

	prn("Start")
	defer func() {
		prn("Finish")
	}()

	var (
		last int64
	)
	diffs := []int64{2, 2, 2, 1, 2, 1, 2, 1}
	count := 0
	for {
		if ok, errStr := doSomething(last, diffs[count]); !ok {
			t.Error("Ошибка:" + errStr)
			return
		}
		count++
		last = time.Now().Unix()
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

	var last int64
	for i := 0; i < 5; i++ {
		if ok, errStr := doSomething(last, 2); !ok {
			t.Error("Ошибка:" + errStr)
			return
		}
		last = time.Now().Unix()
		<-tr.Throttling()
	}
}
