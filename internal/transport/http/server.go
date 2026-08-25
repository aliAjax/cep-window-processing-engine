package httptransport

import (
	"encoding/json"
	"github.com/example/cep-window-engine/internal/checkpoint"
	"github.com/example/cep-window-engine/internal/domain"
	"github.com/example/cep-window-engine/internal/pattern/compiler"
	"github.com/example/cep-window-engine/internal/pattern/parser"
	"github.com/example/cep-window-engine/internal/runtime"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Store      *domain.Store
	Engine     *runtime.Engine
	Checkpoint *checkpoint.Manager
}

func New(s *domain.Store, e *runtime.Engine, c *checkpoint.Manager) *Server {
	return &Server{Store: s, Engine: e, Checkpoint: c}
}
func (s *Server) Routes() http.Handler {
	m := http.NewServeMux()
	m.Handle("/console/", consoleHandler())
	m.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"status":"ok"}`)) })
	m.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"ready":true}`)) })
	m.HandleFunc("/api/v1/streams", s.streams)
	m.HandleFunc("/api/v1/schemas", s.schemas)
	m.HandleFunc("/api/v1/patterns", s.patterns)
	m.HandleFunc("/api/v1/patterns/validate", s.validate)
	m.HandleFunc("/api/v1/patterns/publish", s.publish)
	m.HandleFunc("/api/v1/events", s.events)
	m.HandleFunc("/api/v1/matches", s.matches)
	m.HandleFunc("/api/v1/matches/", s.explain)
	m.HandleFunc("/api/v1/runtimes/checkpoint", s.save)
	m.HandleFunc("/api/v1/runtimes/replay", s.replay)
	m.HandleFunc("/api/v1/runtimes/pause", func(w http.ResponseWriter, r *http.Request) {
		s.Engine.Pause()
		json.NewEncoder(w).Encode(map[string]any{"paused": true})
	})
	m.HandleFunc("/api/v1/runtimes/resume", func(w http.ResponseWriter, r *http.Request) {
		s.Engine.Resume()
		json.NewEncoder(w).Encode(map[string]any{"paused": false})
	})
	return m
}
func decode(r *http.Request, v any) error { return json.NewDecoder(r.Body).Decode(v) }
func write(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
func (s *Server) streams(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var v domain.Stream
		if e := decode(r, &v); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		if v.ID == "" {
			v.ID = "stream-" + time.Now().Format("150405.000000")
		}
		v.CreatedAt = time.Now()
		if e := s.Store.AddStream(v); e != nil {
			http.Error(w, e.Error(), 409)
			return
		}
		write(w, v)
		return
	}
	write(w, s.Store.View().Streams)
}
func (s *Server) schemas(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var v domain.EventSchema
		if e := decode(r, &v); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		if v.ID == "" {
			v.ID = "schema-" + time.Now().Format("150405.000000")
		}
		if e := s.Store.AddSchema(v); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		write(w, v)
		return
	}
	write(w, s.Store.View().Schemas)
}
func (s *Server) patterns(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var v domain.PatternDefinition
		if e := decode(r, &v); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		if v.ID == "" {
			v.ID = "pattern-" + time.Now().Format("150405.000000")
		}
		n, e := parser.Parse(v.Expression)
		if e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		v.Compiled, e = compiler.Compiler{MaxComplexity: 100}.Compile(n)
		if e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		if e = s.Store.AddPattern(v); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		write(w, v)
		return
	}
	write(w, s.Store.View().Patterns)
}
func (s *Server) validate(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Expression string `json:"expression"`
	}
	if e := decode(r, &in); e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	n, e := parser.Parse(in.Expression)
	if e != nil {
		write(w, map[string]any{"valid": false, "error": e.Error()})
		return
	}
	p, e := compiler.Compiler{MaxComplexity: 100}.Compile(n)
	write(w, map[string]any{"valid": e == nil, "plan": p, "error": errString(e)})
}
func errString(e error) string {
	if e == nil {
		return ""
	}
	return e.Error()
}
func (s *Server) publish(w http.ResponseWriter, r *http.Request) { s.patterns(w, r) }
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", 405)
		return
	}
	var e domain.Event
	if er := decode(r, &e); er != nil {
		http.Error(w, er.Error(), 400)
		return
	}
	if e.EventTime.IsZero() {
		e.EventTime = time.Now()
	}
	if e.ID == "" {
		e.ID = e.Fingerprint()
	}
	m, er := s.Engine.Ingest(e)
	if er != nil {
		http.Error(w, er.Error(), 429)
		return
	}
	write(w, map[string]any{"accepted": true, "matches": m})
}
func (s *Server) matches(w http.ResponseWriter, r *http.Request) {
	write(w, s.Store.View().Matches)
}
func (s *Server) explain(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/matches/")
	if id == "" {
		http.Error(w, "match id required", http.StatusBadRequest)
		return
	}
	m, ok := s.Store.GetMatch(id)
	if !ok {
		http.Error(w, "match not found", http.StatusNotFound)
		return
	}
	write(w, map[string]any{"match_id": m.ID, "pattern_id": m.PatternID, "evidence": m.Evidence, "explain": m.Explain})
}
func (s *Server) save(w http.ResponseWriter, r *http.Request) {
	id, e := s.Checkpoint.Save(s.Store)
	if e != nil {
		http.Error(w, e.Error(), 500)
		return
	}
	write(w, map[string]string{"checkpoint_id": id})
}
func (s *Server) replay(w http.ResponseWriter, r *http.Request) {
	var in struct {
		StreamID string `json:"stream_id"`
	}
	decode(r, &in)
	n, e := s.Engine.Replay(strings.TrimSpace(in.StreamID))
	write(w, map[string]any{"replayed": n, "error": errString(e)})
}
