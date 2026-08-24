package auditview

type Report struct {
	Labels  map[string]string
	Entries []Event
}

type Exporter struct {
	report Report
}

func (e *Exporter) Replace(report Report) {
	e.report = cloneReport(report)
}

func (e *Exporter) Export() Report {
	return cloneReport(e.report)
}

func cloneReport(report Report) Report {
	out := Report{
		Labels:  make(map[string]string, len(report.Labels)),
		Entries: make([]Event, len(report.Entries)),
	}
	for key, value := range report.Labels {
		out.Labels[key] = value
	}
	for i, event := range report.Entries {
		out.Entries[i] = cloneEvent(event)
	}
	return out
}
