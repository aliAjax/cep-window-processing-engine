package subscription

import "testing"

func TestRegistryFirstAdd(t *testing.T) {
	var registry Registry
	registry.Add("alerts", "subscriber-a")
	members := registry.Members("alerts")
	if len(members) != 1 || members[0] != "subscriber-a" {
		t.Fatalf("first subscriber not retained: %v", members)
	}
}

func TestSessionStoreFirstOpen(t *testing.T) {
	var store SessionStore
	store.Open(Session{ID: "session-a", GroupID: "alerts"})
	session, ok := store.Get("session-a")
	if !ok || session.GroupID != "alerts" {
		t.Fatalf("first session not retained: %#v, %v", session, ok)
	}
}

func TestFanoutFirstEnqueue(t *testing.T) {
	var fanout Fanout
	fanout.Enqueue("alerts", Notification{ID: "notice-a"})
	pending := fanout.Pending("alerts")
	if len(pending) != 1 || pending[0].ID != "notice-a" {
		t.Fatalf("first notification not queued: %#v", pending)
	}
}

func TestSnapshotterFirstCapture(t *testing.T) {
	var snapshotter Snapshotter
	snapshotter.Capture("alerts", Snapshot{Version: 1, Members: []string{"subscriber-a"}})
	snapshot, ok := snapshotter.Latest("alerts")
	if !ok || snapshot.Version != 1 || len(snapshot.Members) != 1 {
		t.Fatalf("first snapshot not retained: %#v, %v", snapshot, ok)
	}
}
