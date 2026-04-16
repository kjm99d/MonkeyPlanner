package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	mphttp "github.com/kjm99d/monkey-planner/backend/internal/http"
	"github.com/kjm99d/monkey-planner/backend/internal/service"
	"github.com/kjm99d/monkey-planner/backend/internal/storage"
	_ "github.com/kjm99d/monkey-planner/backend/internal/storage/sqlite"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	dsn := "sqlite://" + filepath.Join(t.TempDir(), "httptest.db")
	repo, err := storage.NewRepo(dsn)
	if err != nil {
		t.Fatalf("repo: %v", err)
	}
	svc := service.New(repo, func() time.Time { return time.Now().UTC() })
	srv := httptest.NewServer(mphttp.NewRouter(svc, nil))
	t.Cleanup(func() {
		srv.Close()
		_ = repo.Close()
	})
	return srv
}

func doJSON(t *testing.T, method, url string, body any) (*http.Response, []byte) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		buf, _ := json.Marshal(body)
		reader = bytes.NewReader(buf)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	return resp, out
}

func mustStatus(t *testing.T, resp *http.Response, want int, body []byte) {
	t.Helper()
	if resp.StatusCode != want {
		t.Fatalf("status: got %d want %d; body=%s", resp.StatusCode, want, body)
	}
}

type createdIssue struct {
	ID         string     `json:"id"`
	Status     string     `json:"status"`
	ApprovedAt *time.Time `json:"approvedAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

type createdBoard struct {
	ID string `json:"id"`
}

func TestFullFlow(t *testing.T) {
	srv := newTestServer(t)

	// 1) 헬스체크
	resp, body := doJSON(t, http.MethodGet, srv.URL+"/api/health", nil)
	mustStatus(t, resp, 200, body)

	// 2) 보드 생성
	resp, body = doJSON(t, http.MethodPost, srv.URL+"/api/boards",
		map[string]any{"name": "Backlog", "viewType": "kanban"})
	mustStatus(t, resp, 201, body)
	var b createdBoard
	_ = json.Unmarshal(body, &b)

	// 3) 이슈 생성 → Pending
	resp, body = doJSON(t, http.MethodPost, srv.URL+"/api/issues",
		map[string]any{"boardId": b.ID, "title": "첫 작업", "body": "내용"})
	mustStatus(t, resp, 201, body)
	var iss createdIssue
	_ = json.Unmarshal(body, &iss)
	if iss.Status != "Pending" || iss.ApprovedAt != nil {
		t.Fatalf("create: status=%s approvedAt=%v", iss.Status, iss.ApprovedAt)
	}

	// 4) PATCH 로 Approved 전이 시도 → 409 (use_approve_endpoint)
	resp, body = doJSON(t, http.MethodPatch, srv.URL+"/api/issues/"+iss.ID,
		map[string]any{"status": "Approved"})
	mustStatus(t, resp, 409, body)
	if !bytes.Contains(body, []byte("use POST")) {
		t.Fatalf("expected use-approve-endpoint hint, got %s", body)
	}

	// 5) Approve 전용 엔드포인트 호출 → Approved + approved_at
	resp, body = doJSON(t, http.MethodPost, srv.URL+"/api/issues/"+iss.ID+"/approve", nil)
	mustStatus(t, resp, 200, body)
	var afterApprove createdIssue
	_ = json.Unmarshal(body, &afterApprove)
	if afterApprove.Status != "Approved" || afterApprove.ApprovedAt == nil {
		t.Fatalf("after approve: status=%s approvedAt=%v", afterApprove.Status, afterApprove.ApprovedAt)
	}
	firstApproved := *afterApprove.ApprovedAt

	// 6) Approve 두 번째 호출 → 멱등(approved_at 불변)
	time.Sleep(15 * time.Millisecond)
	resp, body = doJSON(t, http.MethodPost, srv.URL+"/api/issues/"+iss.ID+"/approve", nil)
	mustStatus(t, resp, 200, body)
	var afterApprove2 createdIssue
	_ = json.Unmarshal(body, &afterApprove2)
	if !afterApprove2.ApprovedAt.Equal(firstApproved) {
		t.Fatalf("approve not idempotent: first=%v second=%v", firstApproved, *afterApprove2.ApprovedAt)
	}

	// 7) InProgress → Done 흐름
	resp, body = doJSON(t, http.MethodPatch, srv.URL+"/api/issues/"+iss.ID,
		map[string]any{"status": "InProgress"})
	mustStatus(t, resp, 200, body)
	resp, body = doJSON(t, http.MethodPatch, srv.URL+"/api/issues/"+iss.ID,
		map[string]any{"status": "Done"})
	mustStatus(t, resp, 200, body)
	var done createdIssue
	_ = json.Unmarshal(body, &done)
	if done.Status != "Done" || done.CompletedAt == nil {
		t.Fatalf("done: status=%s completedAt=%v", done.Status, done.CompletedAt)
	}

	// 8) 역행 차단
	resp, body = doJSON(t, http.MethodPatch, srv.URL+"/api/issues/"+iss.ID,
		map[string]any{"status": "InProgress"})
	mustStatus(t, resp, 400, body)

	// 9) 캘린더 월 집계에 반영됨
	now := time.Now().UTC()
	resp, body = doJSON(t, http.MethodGet,
		srv.URL+"/api/calendar?year="+itoa(now.Year())+"&month="+itoa(int(now.Month())), nil)
	mustStatus(t, resp, 200, body)
	if !bytes.Contains(body, []byte(`"created":`)) {
		t.Fatalf("expected created key in month stats: %s", body)
	}
}

func TestEmptyCalendarDayReturnsArrays(t *testing.T) {
	// 회귀: 빈 날짜에 대한 /api/calendar/day 응답에서 created/approved/completed가
	// null이 아닌 [] 로 직렬화되어야 한다 (프론트의 .length 접근 TypeError 방지).
	srv := newTestServer(t)
	resp, body := doJSON(t, http.MethodGet, srv.URL+"/api/calendar/day?date=2099-01-01", nil)
	mustStatus(t, resp, 200, body)
	for _, k := range []string{`"created":[]`, `"approved":[]`, `"completed":[]`} {
		if !bytes.Contains(body, []byte(k)) {
			t.Fatalf("expected %q in response, got %s", k, body)
		}
	}
	if bytes.Contains(body, []byte("null")) {
		t.Fatalf("response must not contain null for empty arrays: %s", body)
	}
}

func TestCycleBlocked(t *testing.T) {
	srv := newTestServer(t)

	resp, body := doJSON(t, http.MethodPost, srv.URL+"/api/boards",
		map[string]any{"name": "B", "viewType": "kanban"})
	mustStatus(t, resp, 201, body)
	var b createdBoard
	_ = json.Unmarshal(body, &b)

	mk := func(title string, parent *string) createdIssue {
		p := map[string]any{"boardId": b.ID, "title": title}
		if parent != nil {
			p["parentId"] = *parent
		}
		resp, body := doJSON(t, http.MethodPost, srv.URL+"/api/issues", p)
		mustStatus(t, resp, 201, body)
		var i createdIssue
		_ = json.Unmarshal(body, &i)
		return i
	}
	a := mk("A", nil)
	bi := mk("B", &a.ID)
	c := mk("C", &bi.ID)

	// A의 parent를 C로 → 순환 (A→C→B→A) → 400 cycle
	resp, body = doJSON(t, http.MethodPatch, srv.URL+"/api/issues/"+a.ID,
		map[string]any{"parentId": c.ID})
	mustStatus(t, resp, 400, body)
	if !bytes.Contains(body, []byte("cycle")) {
		t.Fatalf("expected cycle error code, got %s", body)
	}
}

func itoa(n int) string {
	b := make([]byte, 0, 4)
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}
