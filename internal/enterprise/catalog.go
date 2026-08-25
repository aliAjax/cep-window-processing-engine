package enterprise

// The catalog package keeps operational metadata separate from the event evaluator.
// These services expose dependency-free ports for production adapters.
type Tenant struct {
	ID, Name, Status      string
	CPUQuota, MemoryQuota int64
}
type Quota struct {
	TenantID                                string
	MaxEventsPerSecond, MaxStates, MaxBytes int64
}
type SchemaRef struct{ ID, Version, Digest string }
type SinkRef struct {
	ID, Kind, Endpoint string
	Enabled            bool
}
type RuntimeRef struct {
	ID, StreamID, PatternID, Status string
	Generation                      int64
}
type SnapshotRef struct {
	ID, RuntimeID, Checksum, Offset string
	Size                            int64
}
type Service struct {
	tenants   map[string]Tenant
	quotas    map[string]Quota
	schemas   map[string]SchemaRef
	sinks     map[string]SinkRef
	runtimes  map[string]RuntimeRef
	snapshots map[string]SnapshotRef
}

func New() *Service {
	return &Service{tenants: map[string]Tenant{}, quotas: map[string]Quota{}, schemas: map[string]SchemaRef{}, sinks: map[string]SinkRef{}, runtimes: map[string]RuntimeRef{}, snapshots: map[string]SnapshotRef{}}
}
func (s *Service) PutTenant(v Tenant)                 { s.tenants[v.ID] = v }
func (s *Service) GetTenant(id string) (Tenant, bool) { v, ok := s.tenants[id]; return v, ok }
func (s *Service) DeleteTenant(id string)             { delete(s.tenants, id) }
func (s *Service) ListTenants() []Tenant {
	o := []Tenant{}
	for _, v := range s.tenants {
		o = append(o, v)
	}
	return o
}
func (s *Service) PutQuota(v Quota)                          { s.quotas[v.TenantID] = v }
func (s *Service) GetQuota(id string) (Quota, bool)          { v, ok := s.quotas[id]; return v, ok }
func (s *Service) PutSchema(v SchemaRef)                     { s.schemas[v.ID] = v }
func (s *Service) GetSchema(id string) (SchemaRef, bool)     { v, ok := s.schemas[id]; return v, ok }
func (s *Service) PutSink(v SinkRef)                         { s.sinks[v.ID] = v }
func (s *Service) GetSink(id string) (SinkRef, bool)         { v, ok := s.sinks[id]; return v, ok }
func (s *Service) PutRuntime(v RuntimeRef)                   { s.runtimes[v.ID] = v }
func (s *Service) GetRuntime(id string) (RuntimeRef, bool)   { v, ok := s.runtimes[id]; return v, ok }
func (s *Service) PutSnapshot(v SnapshotRef)                 { s.snapshots[v.ID] = v }
func (s *Service) GetSnapshot(id string) (SnapshotRef, bool) { v, ok := s.snapshots[id]; return v, ok }

// Validator performs bounded checks used by API handlers and workers.
type Validator struct {
	MaxNameLength int
	MaxFields     int
	MaxExpression int
}

func (v Validator) CheckTenant(t Tenant) error {
	if t.ID == "" {
		return err("tenant id required")
	}
	if len(t.Name) > v.MaxNameLength && v.MaxNameLength > 0 {
		return err("tenant name too long")
	}
	return nil
}
func (v Validator) CheckQuota(q Quota) error {
	if q.MaxEventsPerSecond < 1 {
		return err("event quota required")
	}
	if q.MaxStates < 1 {
		return err("state quota required")
	}
	return nil
}
func (v Validator) CheckSchema(s SchemaRef) error {
	if s.ID == "" || s.Version == "" {
		return err("schema identity required")
	}
	return nil
}
func (v Validator) CheckSink(s SinkRef) error {
	if s.ID == "" || s.Kind == "" {
		return err("sink identity required")
	}
	return nil
}
func (v Validator) CheckRuntime(r RuntimeRef) error {
	if r.ID == "" || r.StreamID == "" || r.PatternID == "" {
		return err("runtime binding required")
	}
	return nil
}

type validationError string

func (e validationError) Error() string { return string(e) }
func err(s string) error                { return validationError(s) }

// State migration descriptors document hot rule update behavior.
type Migration struct {
	FromVersion, ToVersion    string
	PreserveWindows, DrainOld bool
}

func (m Migration) Valid() bool {
	return m.FromVersion != "" && m.ToVersion != "" && m.FromVersion != m.ToVersion
}

type Drain struct {
	RuntimeID string
	Remaining int
	Completed bool
}

func (d *Drain) Step() {
	if d.Remaining > 0 {
		d.Remaining--
	}
	d.Completed = d.Remaining == 0
}

type Budget struct{ CPU, Memory, Events int64 }

func (b Budget) Exceeded(u Budget) bool {
	return u.CPU > b.CPU || u.Memory > b.Memory || u.Events > b.Events
}

// The following records are used by adapters for deterministic audit and replay.
type EventOffset struct {
	StreamID  string
	Partition int
	Offset    int64
}
type ReplayRequest struct {
	RuntimeID string
	From, To  EventOffset
	DryRun    bool
}
type ReplayResult struct {
	Accepted, Skipped, Matches int
	Errors                     []string
}
type ExplainRequest struct{ MatchID, Format string }
type ExplainResponse struct {
	MatchID, Pattern, Rendered string
	Events                     []string
}
type WatermarkPolicy struct {
	IdleTimeoutMillis    int64
	RequireAllPartitions bool
}
type LateEventPolicy struct {
	Mode             string
	MaxAgeMillis     int64
	EmitCompensation bool
}
type RetentionPolicy struct {
	MaxAgeMillis int64
	MaxBytes     int64
}
type DedupPolicy struct {
	TTLMillis int64
	Key       string
}
type BackpressurePolicy struct {
	MaxInFlight    int
	RejectWhenFull bool
}
type SecurityPolicy struct {
	RequireTenant   bool
	MaxPayloadBytes int
	AllowedTypes    []string
}
type Health struct {
	Status string
	Checks map[string]string
}

func Healthy() Health {
	return Health{Status: "ok", Checks: map[string]string{"runtime": "ok", "checkpoint": "ok"}}
}

// No-op hooks keep orchestration testable without forcing a storage technology.
type Hook interface {
	Before(name string) error
	After(name string, cause error)
}
type Hooks []Hook

func (h Hooks) Before(n string) error {
	for _, x := range h {
		if e := x.Before(n); e != nil {
			return e
		}
	}
	return nil
}
func (h Hooks) After(n string, e error) {
	for _, x := range h {
		x.After(n, e)
	}
}

// Operation is a bounded unit in worker queues.
type Operation struct {
	ID, Kind, State string
	Attempts        int
	Payload         []byte
}

func (o *Operation) Retry(max int) bool {
	if o.Attempts >= max {
		o.State = "dead"
		return false
	}
	o.Attempts++
	o.State = "queued"
	return true
}
func (o *Operation) Complete() { o.State = "completed" }
func (o *Operation) Fail()     { o.State = "failed" }
