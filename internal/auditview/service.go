package auditview

type View struct {
	Sessions map[string][]string
	Report   Report
}

type Service struct {
	views map[string]View
}

func NewService() *Service {
	return &Service{views: make(map[string]View)}
}

func (s *Service) Publish(tenantID string, view View) {
	s.views[tenantID] = cloneView(view)
}

func (s *Service) View(tenantID string) View {
	return cloneView(s.views[tenantID])
}

func cloneView(view View) View {
	sessions := make(map[string][]string, len(view.Sessions))
	for sessionID, eventIDs := range view.Sessions {
		sessions[sessionID] = append([]string(nil), eventIDs...)
	}
	return View{Sessions: sessions, Report: cloneReport(view.Report)}
}
