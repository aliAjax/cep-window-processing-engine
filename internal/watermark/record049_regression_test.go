package watermark

import (
	"sync"
	"testing"
	"time"
)

func TestWatermarkBarrierConcurrentPublish(t *testing.T) {
	var barrier Barrier
	base := time.Unix(1_700_000_000, 0)
	start := make(chan struct{})
	var group sync.WaitGroup
	for i := 1; i <= 64; i++ {
		candidate := base.Add(time.Duration(i) * time.Second)
		group.Add(1)
		go func() {
			defer group.Done()
			<-start
			barrier.Publish(candidate)
			_ = barrier.Current()
		}()
	}
	close(start)
	group.Wait()
	if got, want := barrier.Current(), base.Add(64*time.Second); !got.Equal(want) {
		t.Fatalf("barrier = %v, want %v", got, want)
	}
}

func TestWatermarkTrackerSnapshotIsolation(t *testing.T) {
	tracker := NewTracker()
	base := time.Unix(1_700_000_000, 0)
	if err := tracker.Advance(1, base); err != nil {
		t.Fatal(err)
	}

	start := make(chan struct{})
	var group sync.WaitGroup
	for i := 0; i < 16; i++ {
		partition := i + 2
		group.Add(2)
		go func() {
			defer group.Done()
			<-start
			_ = tracker.Advance(partition, base.Add(time.Duration(partition)*time.Second))
		}()
		go func() {
			defer group.Done()
			<-start
			_ = tracker.Snapshot()
		}()
	}
	close(start)
	group.Wait()

	snapshot := tracker.Snapshot()
	snapshot[1] = base.Add(time.Hour)
	if got := tracker.Snapshot()[1]; !got.Equal(base) {
		t.Fatalf("snapshot mutation changed tracker: %v", got)
	}
}

func TestWatermarkCoordinatorUsesUpdatedSnapshot(t *testing.T) {
	base := time.Unix(1_700_000_000, 0)
	tracker := NewTracker()
	if err := tracker.Advance(0, base); err != nil {
		t.Fatal(err)
	}
	if err := tracker.Advance(1, base.Add(5*time.Second)); err != nil {
		t.Fatal(err)
	}
	coordinator := Coordinator{Tracker: tracker, Barrier: &Barrier{}}

	got, err := coordinator.Advance(0, base.Add(10*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if want := base.Add(5 * time.Second); !got.Equal(want) {
		t.Fatalf("global watermark = %v, want %v", got, want)
	}
}

func TestWatermarkPolicyAcceptsFirstAdvance(t *testing.T) {
	policy := Policy{MaximumStep: time.Minute}
	if !policy.Accept(time.Time{}, time.Unix(1_700_000_000, 0)) {
		t.Fatal("first watermark was rejected")
	}
}
