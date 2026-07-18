package analytics

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

// --- Mocks ---

type mockRepo struct {
	summary        Summary
	summaryErr     error
	volume         []VolumePoint
	volumeErr      error
	dist           []StatusDistribution
	distErr        error
	resolution     []ResolutionMetrics
	resolutionErr  error
}

func (m *mockRepo) GetSummary(ctx context.Context, tenantID uuid.UUID) (Summary, error) {
	return m.summary, m.summaryErr
}

func (m *mockRepo) GetVolume(ctx context.Context, tenantID uuid.UUID) ([]VolumePoint, error) {
	return m.volume, m.volumeErr
}

func (m *mockRepo) GetStatusDistribution(ctx context.Context, tenantID uuid.UUID) ([]StatusDistribution, error) {
	return m.dist, m.distErr
}

func (m *mockRepo) GetResolutionTime(ctx context.Context, tenantID uuid.UUID) ([]ResolutionMetrics, error) {
	return m.resolution, m.resolutionErr
}

// --- Tests ---

func TestService_GetSummary_Delegates(t *testing.T) {
	want := Summary{TotalTickets: 10, OpenTickets: 3}
	repo := &mockRepo{summary: want}
	svc := NewService(repo)

	got, err := svc.GetSummary(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestService_GetSummary_Error(t *testing.T) {
	repo := &mockRepo{summaryErr: errors.New("db")}
	svc := NewService(repo)
	if _, err := svc.GetSummary(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

func TestService_GetVolume_Delegates(t *testing.T) {
	want := []VolumePoint{{Date: "2025-01-01", Count: 5}}
	repo := &mockRepo{volume: want}
	svc := NewService(repo)

	got, err := svc.GetVolume(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Count != 5 {
		t.Errorf("got %v", got)
	}
}

func TestService_GetVolume_Error(t *testing.T) {
	repo := &mockRepo{volumeErr: errors.New("db")}
	svc := NewService(repo)
	if _, err := svc.GetVolume(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

func TestService_GetStatusDistribution_Delegates(t *testing.T) {
	want := []StatusDistribution{{Status: "open", Count: 3}}
	repo := &mockRepo{dist: want}
	svc := NewService(repo)

	got, err := svc.GetStatusDistribution(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Status != "open" {
		t.Errorf("got %v", got)
	}
}

func TestService_GetStatusDistribution_Error(t *testing.T) {
	repo := &mockRepo{distErr: errors.New("db")}
	svc := NewService(repo)
	if _, err := svc.GetStatusDistribution(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

func TestService_GetResolutionTime_Delegates(t *testing.T) {
	want := []ResolutionMetrics{{Priority: "high", AvgHours: 4.5}}
	repo := &mockRepo{resolution: want}
	svc := NewService(repo)

	got, err := svc.GetResolutionTime(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].AvgHours != 4.5 {
		t.Errorf("got %v", got)
	}
}

func TestService_GetResolutionTime_Error(t *testing.T) {
	repo := &mockRepo{resolutionErr: errors.New("db")}
	svc := NewService(repo)
	if _, err := svc.GetResolutionTime(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

// --- GetCharts tests ---

func TestService_GetCharts_Success(t *testing.T) {
	repo := &mockRepo{
		summary:    Summary{TotalTickets: 10, SLABreached: 3},
		volume:     []VolumePoint{{Date: "2025-01-01", Count: 5}},
		dist:       []StatusDistribution{{Status: "open", Count: 3}},
		resolution: []ResolutionMetrics{{Priority: "high", AvgHours: 4.5}},
	}
	svc := NewService(repo)

	data, err := svc.GetCharts(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data.Volume) != 1 {
		t.Errorf("volume = %d", len(data.Volume))
	}
	if len(data.StatusDistribution) != 1 {
		t.Errorf("dist = %d", len(data.StatusDistribution))
	}
	if len(data.ResolutionTime) != 1 {
		t.Errorf("resolution = %d", len(data.ResolutionTime))
	}
	if data.SLACompliance.Total != 10 || data.SLACompliance.Breached != 3 || data.SLACompliance.Compliant != 7 {
		t.Errorf("sla = %+v", data.SLACompliance)
	}
}

func TestService_GetCharts_VolumeError(t *testing.T) {
	repo := &mockRepo{volumeErr: errors.New("db")}
	svc := NewService(repo)
	_, err := svc.GetCharts(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestService_GetCharts_DistError(t *testing.T) {
	repo := &mockRepo{distErr: errors.New("db")}
	svc := NewService(repo)
	_, err := svc.GetCharts(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestService_GetCharts_ResolutionError(t *testing.T) {
	repo := &mockRepo{resolutionErr: errors.New("db")}
	svc := NewService(repo)
	_, err := svc.GetCharts(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestService_GetCharts_SummaryError(t *testing.T) {
	repo := &mockRepo{summaryErr: errors.New("db")}
	svc := NewService(repo)
	_, err := svc.GetCharts(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error")
	}
}
