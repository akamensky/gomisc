package ds

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestMap_Exists(t *testing.T) {
	m := &Map[string, string]{}

	if m.Exists("world") {
		t.Errorf("unexpected test result")
	}

	m.Insert("hello", "world")

	if !m.Exists("hello") {
		t.Errorf("unexpected test result")
	}
	if m.Exists("world") {
		t.Errorf("unexpected test result")
	}
}

func TestMap_Insert(t *testing.T) {
	m := &Map[string, int]{}
	m.Insert("zero", 0)
	m.Insert("one", 1)

	if len(m.data) != 2 {
		t.Errorf("expected map size 2, but got %d", len(m.data))
	}

	if i, ok := m.data["zero"]; !ok || i != 0 {
		t.Errorf("unexpected test result")
	}
	if i, ok := m.data["one"]; !ok || i != 1 {
		t.Errorf("unexpected test result")
	}
	if _, ok := m.data["two"]; ok {
		t.Errorf("unexpected test result")
	}
}

func TestMap_Get(t *testing.T) {
	m := &Map[string, int]{}
	m.Insert("zero", 0)
	m.Insert("one", 1)

	if i := m.Get("zero"); i != 0 {
		t.Errorf("unexpected test result")
	}
	if i := m.Get("one"); i != 1 {
		t.Errorf("unexpected test result")
	}
}

func TestM_Get_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("unexpected test result")
		}
	}()
	m := &Map[string, int]{}

	_ = m.Get("two")
}

func TestMap_Delete(t *testing.T) {
	m := &Map[int, int]{}

	m.Delete(0)

	if _, ok := m.data[0]; ok {
		t.Errorf("unexpected test result")
	}
}

func TestMap_Delete2(t *testing.T) {
	m := &Map[int, int]{}

	m.Insert(0, 1)

	if i, ok := m.data[0]; !ok || i != 1 {
		t.Errorf("unexpected test result")
	}

	m.Delete(0)

	if _, ok := m.data[0]; ok {
		t.Errorf("unexpected test result")
	}
}

func TestMap_Range(t *testing.T) {
	m := &Map[int, int]{}

	m.Range(func(k int, v int) bool {
		return true
	})

	m.Insert(0, 1)

	m.Range(func(k int, v int) bool {
		if k != 0 || v != 1 {
			t.Errorf("unexpected test result")
			return false
		}
		return true
	})
}

func TestMap_Range2(t *testing.T) {
	m := &Map[int, int]{}
	m.Insert(1, 2)
	m.Insert(3, 4)

	counter := 0
	m.Range(func(k int, v int) bool {
		counter++
		return false
	})
	if counter != 1 {
		t.Errorf("unexpected test result")
	}
}

func TestMap_Keys(t *testing.T) {
	m := &Map[int, int]{}

	_ = m.Keys()

	m.Insert(1, 2)
	m.Insert(3, 4)

	keys := m.Keys()
	sort.Ints(keys)
	if !reflect.DeepEqual(keys, []int{1, 3}) {
		t.Errorf("unexpected test result")
		fmt.Println(keys)
	}
}

func TestMap_Values(t *testing.T) {
	m := &Map[int, int]{}

	_ = m.Values()

	m.Insert(1, 2)
	m.Insert(3, 4)

	values := m.Values()
	sort.Ints(values)
	if !reflect.DeepEqual(values, []int{2, 4}) {
		t.Errorf("unexpected test result")
	}
}

func TestMap_MarshalJSON(t *testing.T) {
	m := &Map[int, int]{}
	m.Insert(1, 2)

	b, err := m.MarshalJSON()
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(b) != `{"1":2}` {
		t.Errorf(string(b))
	}
}

func TestMap_Async_Access(t *testing.T) {
	m := &Map[string, int]{}

	wg := &sync.WaitGroup{}
	for c := 0; c < 4; c++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000000; i++ {
				m.Insert(fmt.Sprintf("%d", i), i)
			}
			wg.Done()
		}()
	}

	for c := 0; c < 4; c++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000000; i++ {
				if m.Exists(fmt.Sprintf("%d", i)) {
					v := m.Get(fmt.Sprintf("%d", i))
					if v != i {
						t.Errorf("unexpected test results")
					}
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
