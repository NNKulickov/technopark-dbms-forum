package tests

import "github.com/bozaro/tech-db-forum/generated/client"
import (
	"crypto/md5"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Perf struct {
	c    *client.Forum
	data *PerfData
}

type PHash [16]byte
type PVersion uint32

type PerfTest struct {
	Name   string
	Mode   PerfMode
	Weight PerfWeight
	FnPerf func(p *Perf)
}

type PerfMode int

const (
	ModeRead PerfMode = iota
	ModeWrite
)

type PerfWeight int

const (
	WeightRare   PerfWeight = 1
	WeightNormal            = 10
)

var (
	registeredPerfsWeight int32 = 0
	registeredPerfs       []PerfTest
)

func PerfRegister(test PerfTest) {
	registeredPerfs = append(registeredPerfs, test)
	registeredPerfsWeight += int32(test.Weight)
}

func (self *Perf) Validate(callback func(validator PerfValidator)) {
	callback(&PerfSession{})
}

func Hash(data string) PHash {
	return PHash(md5.Sum([]byte(data)))
}

func GetRandomPerfTest() *PerfTest {
	index := rand.Int31n(registeredPerfsWeight)
	for _, item := range registeredPerfs {
		index -= int32(item.Weight)
		if index < 0 {
			return &item
		}
	}
	panic("Invalid state")
}

func (self *Perf) Run() {
	log.Info("BEFORE")

	var done int32 = 0
	var counter int64 = 0
	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < /*parallel*/ 20; i++ {
		wg.Add(1)
		go func() {
			for {
				if atomic.LoadInt32(&done) != 0 {
					break
				}
				p := GetRandomPerfTest()
				p.FnPerf(self)
				atomic.AddInt64(&counter, 1)
			}
			wg.Done()
		}()
	}

	log.Info(atomic.LoadInt64(&counter))
	time.Sleep(time.Second * 10)
	log.Info(atomic.LoadInt64(&counter))
	time.Sleep(time.Second * 10)
	log.Info(atomic.LoadInt64(&counter))
	done = 1

	// wait for the workers to finish
	wg.Wait()
}
