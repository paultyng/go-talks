package main

import (
	"context"
	"math/rand"
	"time"

	"cloud.google.com/go/trace"
)

const maxRandSeconds = 3

func init() {
	rand.Seed(time.Now().Unix())
}

func doTracedWork(ctx context.Context, name string, seconds int64, next func(context.Context)) {
	if parentSpan := trace.FromContext(ctx); parentSpan != nil {
		childSpan := parentSpan.NewChild(name)
		ctx = trace.NewContext(ctx, childSpan) // OMIT
		defer childSpan.Finish()               // HL
	}

	doWork(seconds)
	if next != nil {
		next(ctx)
	}
}

func doWork(seconds int64) {
	sleep := time.Duration(seconds) * time.Second
	r := time.Duration(rand.Int63n(maxRandSeconds * int64(time.Second)))
	time.Sleep(sleep + r)
}
