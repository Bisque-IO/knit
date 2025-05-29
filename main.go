package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func main() {
	var conf Config
	conf.Name = "knit"
	conf.Version = "0.1.0"
	conf.Backend = Bolt
	conf.InitialData = new(data)

	// Since we are not holding onto much data we can used the built-in JSON
	// snapshot system. You just need to make sure all the important fields in
	// the data are exportable (capitalized) to JSON.
	conf.UseJSONSnapshots = true

	conf.ConnOpened = func(addr string) (context interface{}, accept bool) {
		return nil, true
	}

	conf.ConnClosed = func(context interface{}, addr string) {

	}

	conf.StateChange = func(state State) {
		fmt.Printf("state changed: %d\n", int(state))
	}

	conf.AddReadCommand("SERVER", cmdSERVER)
	conf.AddWriteCommand("ADDANDGET", cmdADDANDGET)
	conf.AddWriteCommand("GETANDADD", cmdGETANDADD)

	Main(conf)
}

type data struct {
	Nodes    map[string]*Node
	Maps     map[string]*MultiMap
	Counters map[string]*Counter
}

func cmdSERVER(m Machine, args []string) (interface{}, error) {
	data := m.Data().(*data)
	_ = data

	return map[string]string{
		"version": "1.0.1",
	}, nil
}

func cmdJOIN(m Machine, args []string) (interface{}, error) {
	data := m.Data().(*data)
	_ = data

	// Return the new ticket to caller
	return "PONG", nil
}

func cmdLEAVE(m Machine, args []string) (interface{}, error) {
	// The the current data from the machine
	data := m.Data().(*data)
	_ = data

	// Return the new ticket to caller
	return "PONG", nil
}

func cmdPING(m Machine, args []string) (interface{}, error) {
	// The the current data from the machine
	data := m.Data().(*data)
	_ = data

	// Return the new ticket to caller
	return "PONG", nil
}

func cmdADDANDGET(m Machine, args []string) (interface{}, error) {
	if len(args) < 3 {
		return nil, errors.New("not enough arguments")
	}

	by, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return nil, err
	}

	name := args[1]

	data := m.Data().(*data)
	_ = data

	if data.Counters == nil {
		data.Counters = make(map[string]*Counter)
	}

	counter := data.Counters[name]
	if counter == nil {
		counter = new(Counter)
		data.Counters[name] = counter
	}

	// Return the new ticket to caller
	return counter.AddAndGet(by), nil
}

func cmdGETANDADD(m Machine, args []string) (interface{}, error) {
	if len(args) < 3 {
		return nil, errors.New("not enough arguments")
	}

	by, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return nil, err
	}

	name := args[1]

	// The the current data from the machine
	data := m.Data().(*data)
	_ = data

	if data.Counters == nil {
		data.Counters = make(map[string]*Counter)
	}

	counter := data.Counters[name]
	if counter == nil {
		counter = new(Counter)
		data.Counters[name] = counter
	}

	// Return the new ticket to caller
	return counter.GetAndAdd(by), nil
}

type Node struct {
	id       string
	joinedAt time.Time
	lastSeen time.Time
	locks    map[string]*Lock
}

type Lock struct {
	Name         string
	HeldBy       *Node
	ExpiresAt    time.Time
	LastAccessed time.Time
}

func (self *Lock) Lock() bool {
	return false
}

type Counter struct {
	value int64
}

func (c *Counter) GetAndIncrement() int64 {
	value := c.value
	c.value++
	return value
}

func (c *Counter) IncrementAndGet() int64 {
	c.value++
	return c.value
}

func (c *Counter) GetAndDecrement() int64 {
	value := c.value
	c.value--
	return value
}

func (c *Counter) DecrementAndGet() int64 {
	c.value--
	return c.value
}

func (c *Counter) GetAndAdd(by int64) int64 {
	value := c.value
	c.value = c.value + by
	return value
}

func (c *Counter) AddAndGet(by int64) int64 {
	c.value = c.value + by
	return c.value
}

func (c *Counter) CompareAndSet(expected, value int64) bool {
	if c.value == expected {
		c.value = value
		return true
	}
	return false
}

type MultiMap struct {
}
