package app

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/libdns/libdns"
	"github.com/orvice/ddns/internal/config"
)

type fakeIPGetter struct {
	ip  string
	err error
}

func (f *fakeIPGetter) GetIP() (string, error) { return f.ip, f.err }

type fakeDNS struct {
	getRecords    []libdns.Record
	getErr        error
	appendCalls   [][]libdns.Record
	setCalls      [][]libdns.Record
	appendErr     error
	setErr        error
	lastZoneGet   string
	lastZoneApp   string
	lastZoneSet   string
}

func (f *fakeDNS) GetRecords(_ context.Context, zone string) ([]libdns.Record, error) {
	f.lastZoneGet = zone
	return f.getRecords, f.getErr
}

func (f *fakeDNS) AppendRecords(_ context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	f.lastZoneApp = zone
	f.appendCalls = append(f.appendCalls, recs)
	return recs, f.appendErr
}

func (f *fakeDNS) SetRecords(_ context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	f.lastZoneSet = zone
	f.setCalls = append(f.setCalls, recs)
	return recs, f.setErr
}

type fakeNotifier struct {
	messages []string
	err      error
}

func (f *fakeNotifier) Send(_ context.Context, s string) error {
	f.messages = append(f.messages, s)
	return f.err
}

func newTestApp(getter *fakeIPGetter, dns *fakeDNS, notif *fakeNotifier) *App {
	return New(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		&config.Config{Domain: "host.example.com"},
		dns,
		getter,
		notif,
	)
}

func TestZoneFromDomain(t *testing.T) {
	cases := []struct {
		in, name, zone string
	}{
		{"host.example.com", "host", "example.com"},
		{"example.com", "example", "com"},
		{"single", "", ""},
	}
	for _, c := range cases {
		n, z := zoneFromDomain(c.in)
		if n != c.name || z != c.zone {
			t.Errorf("zoneFromDomain(%q) = (%q,%q), want (%q,%q)", c.in, n, z, c.name, c.zone)
		}
	}
}

func TestUpdateIP_AppendWhenNotFound(t *testing.T) {
	ctx := context.Background()
	dns := &fakeDNS{getRecords: nil}
	notif := &fakeNotifier{}
	a := newTestApp(&fakeIPGetter{ip: "1.2.3.4"}, dns, notif)

	if err := a.UpdateIP(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dns.appendCalls) != 1 || len(dns.appendCalls[0]) != 1 {
		t.Fatalf("expected 1 append call, got %+v", dns.appendCalls)
	}
	rec := dns.appendCalls[0][0]
	if rec.Name != "host" || rec.Value != "1.2.3.4" || rec.Type != "A" {
		t.Errorf("unexpected record: %+v", rec)
	}
	if len(notif.messages) != 0 {
		t.Errorf("should not notify on first-time create, got %v", notif.messages)
	}
}

func TestUpdateIP_SkipWhenSame(t *testing.T) {
	ctx := context.Background()
	dns := &fakeDNS{getRecords: []libdns.Record{{Name: "host", Value: "1.2.3.4", Type: "A"}}}
	notif := &fakeNotifier{}
	a := newTestApp(&fakeIPGetter{ip: "1.2.3.4"}, dns, notif)

	if err := a.UpdateIP(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dns.setCalls) != 0 || len(dns.appendCalls) != 0 {
		t.Errorf("no updates expected, got set=%v append=%v", dns.setCalls, dns.appendCalls)
	}
	if len(notif.messages) != 0 {
		t.Errorf("no notification expected, got %v", notif.messages)
	}
}

func TestUpdateIP_SetWhenChanged(t *testing.T) {
	ctx := context.Background()
	dns := &fakeDNS{getRecords: []libdns.Record{{Name: "host", Value: "1.1.1.1", Type: "A"}}}
	notif := &fakeNotifier{}
	a := newTestApp(&fakeIPGetter{ip: "9.9.9.9"}, dns, notif)

	if err := a.UpdateIP(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dns.setCalls) != 1 || dns.setCalls[0][0].Value != "9.9.9.9" {
		t.Fatalf("expected set to 9.9.9.9, got %+v", dns.setCalls)
	}
	if len(notif.messages) != 1 {
		t.Fatalf("expected 1 notification, got %v", notif.messages)
	}
}

func TestUpdateIP_PropagatesGetIPError(t *testing.T) {
	a := newTestApp(&fakeIPGetter{err: errors.New("boom")}, &fakeDNS{}, &fakeNotifier{})
	if err := a.UpdateIP(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateIP_PropagatesGetRecordsError(t *testing.T) {
	a := newTestApp(&fakeIPGetter{ip: "1.2.3.4"}, &fakeDNS{getErr: errors.New("dns boom")}, &fakeNotifier{})
	if err := a.UpdateIP(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}
