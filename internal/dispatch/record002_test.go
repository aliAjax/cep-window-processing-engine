package dispatch

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestProducerClosesOnFailure(t *testing.T) {
	boom := errors.New("sink unavailable")
	jobs := make(chan Job, 2)
	jobs <- Job{Partition: 0}
	jobs <- Job{Partition: 1}
	close(jobs)
	output := make(chan Result, 2)
	start := make(chan struct{})
	go func() {
		<-start
		(Producer{Process: func(context.Context, Job) ([]byte, error) { return nil, boom }}).Emit(context.Background(), jobs, output)
	}()
	close(start)
	select {
	case <-time.After(500 * time.Millisecond):
		t.Fatal("producer did not close result channel")
	case _, ok := <-output:
		if ok {
			select {
			case <-time.After(500 * time.Millisecond):
				t.Fatal("producer left result channel open after failure")
			case _, ok = <-output:
				if ok {
					t.Fatal("unexpected second result")
				}
			}
		}
	}
}

func TestCoordinatorWaitsForWorkers(t *testing.T) {
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseWorkers := func() { releaseOnce.Do(func() { close(release) }) }
	t.Cleanup(releaseWorkers)
	jobs := []Job{{Partition: 0}, {Partition: 1}}
	results := (Coordinator{Workers: 2}).Run(context.Background(), jobs, func(_ context.Context, job Job) ([]byte, error) {
		entered <- struct{}{}
		<-release
		return []byte{byte(job.Partition)}, nil
	})
	for i := 0; i < len(jobs); i++ {
		select {
		case <-entered:
		case <-time.After(500 * time.Millisecond):
			t.Fatal("worker did not enter processing function")
		}
	}
	select {
	case _, ok := <-results:
		if !ok {
			t.Fatal("coordinator closed results before workers completed")
		}
		t.Fatal("coordinator published a result before workers were released")
	case <-time.After(20 * time.Millisecond):
	}
	releaseWorkers()
	count := 0
	for range results {
		count++
	}
	if count != len(jobs) {
		t.Fatalf("received %d results, want %d", count, len(jobs))
	}
}

func TestConsumerTerminatesAfterError(t *testing.T) {
	boom := errors.New("partition failed")
	input := make(chan Result, 2)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for _, result := range []Result{{Partition: 0}, {Partition: 1, Err: boom}} {
		result := result
		wg.Add(1)
		go func() { defer wg.Done(); <-start; input <- result }()
	}
	close(start)
	wg.Wait()
	close(input)
	_, err := (Consumer{}).Drain(context.Background(), input)
	if !errors.Is(err, boom) {
		t.Fatalf("consumer lost dispatch cause: %v", err)
	}
}

func TestResultStreamPublishesFailure(t *testing.T) {
	stream := NewResultStream(2)
	start := make(chan struct{})
	done := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() { <-start; done <- stream.Close() }()
	}
	close(start)
	<-done
	<-done
	select {
	case _, ok := <-stream.Errors:
		if ok {
			t.Fatal("error channel remained open")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("error channel was not closed")
	}
}
