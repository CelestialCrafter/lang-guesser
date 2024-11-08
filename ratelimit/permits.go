package ratelimit

import (
	"time"
)

const GetBlob = "getblob"
const GetTree = "gettree"
const GetRepo = "getrepo"
const Search = "search"

type permits struct {
	Max int
	Current int
	channel chan struct{}
}

func (p permits) Aquire() {
	<-p.channel
	p.Current--
}

type concurrentPermits struct {
	permits
}

func (p concurrentPermits) Release() {
	p.channel <- struct{}{}
	p.Current++
}

type endpointPermits struct {
	permits
}

func (p endpointPermits) Aquire() {
	GlobalPermits.Aquire()
	p.permits.Aquire()
}

func timeRefresh(refresh time.Duration) func(permits) {
	return func(pool permits) {
		for {
			amount := pool.Max - pool.Current
			for range amount {
				pool.channel <- struct{}{}
				pool.Current++
			}

			time.Sleep(refresh)
		}
	}
}

func NewPermitPool(max int, strategy func(permits)) permits {
	pool := permits{
		Max: max,
		Current: max,
		channel: make(chan struct{}, max),
	}

	for range max {
		pool.channel <- struct{}{}
	}

	go strategy(pool)

	return pool
}

func cp(pool permits) concurrentPermits {
	return concurrentPermits{permits: pool}
}

func ep(pool permits) endpointPermits {
	return endpointPermits{permits: pool}
}

var GlobalPermits = NewPermitPool(5000, timeRefresh(time.Hour))
var EndpointPermits = map[string]endpointPermits{
	GetBlob: ep(NewPermitPool(900, timeRefresh(time.Minute))),
	GetTree: ep(NewPermitPool(900, timeRefresh(time.Minute))),
	GetRepo: ep(NewPermitPool(900, timeRefresh(time.Minute))),
	Search: ep(NewPermitPool(30, timeRefresh(time.Minute))),
}
var ConcurrentPermits = cp(NewPermitPool(100, func(p permits) {}))

