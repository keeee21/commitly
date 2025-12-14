package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	return db, mock
}

func TestHealthCheck(t *testing.T) {
	e := echo.New()
	db, _ := setupTestDB(t)

	SetupRoutes(e, db)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ok"`)
}

func TestRoutesRegistered(t *testing.T) {
	e := echo.New()
	db, _ := setupTestDB(t)

	SetupRoutes(e, db)

	routes := e.Routes()

	// 期待されるルートのリスト
	expectedRoutes := map[string][]string{
		"/health":              {http.MethodGet},
		"/api/auth/callback":   {http.MethodPost},
		"/api/auth/logout":     {http.MethodPost},
		"/api/me":              {http.MethodGet},
		"/api/rivals":          {http.MethodGet, http.MethodPost},
		"/api/rivals/:id":      {http.MethodDelete},
		"/api/dashboard/weekly": {http.MethodGet},
		"/api/dashboard/monthly": {http.MethodGet},
		"/api/notifications":    {http.MethodGet, http.MethodPost},
		"/api/notifications/:id": {http.MethodPut, http.MethodDelete},
	}

	// ルートが登録されていることを確認
	registeredRoutes := make(map[string][]string)
	for _, route := range routes {
		registeredRoutes[route.Path] = append(registeredRoutes[route.Path], route.Method)
	}

	for path, methods := range expectedRoutes {
		for _, method := range methods {
			found := false
			for _, registeredMethod := range registeredRoutes[path] {
				if registeredMethod == method {
					found = true
					break
				}
			}
			assert.True(t, found, "Route %s %s should be registered", method, path)
		}
	}
}

func TestAuthCallbackEndpoint(t *testing.T) {
	e := echo.New()
	db, _ := setupTestDB(t)

	SetupRoutes(e, db)

	// 不正なリクエストでエンドポイントが存在することを確認
	req := httptest.NewRequest(http.MethodPost, "/api/auth/callback", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	// Bad Requestが返ることを期待（エンドポイントは存在する）
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProtectedEndpointWithoutAuth(t *testing.T) {
	e := echo.New()
	db, _ := setupTestDB(t)

	SetupRoutes(e, db)

	// 認証なしでプロテクトされたエンドポイントにアクセス
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	// Unauthorizedが返ることを期待
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRivalsEndpointWithoutAuth(t *testing.T) {
	e := echo.New()
	db, _ := setupTestDB(t)

	SetupRoutes(e, db)

	req := httptest.NewRequest(http.MethodGet, "/api/rivals", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDashboardEndpointWithoutAuth(t *testing.T) {
	e := echo.New()
	db, _ := setupTestDB(t)

	SetupRoutes(e, db)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/weekly", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestNotificationsEndpointWithoutAuth(t *testing.T) {
	e := echo.New()
	db, _ := setupTestDB(t)

	SetupRoutes(e, db)

	req := httptest.NewRequest(http.MethodGet, "/api/notifications", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestNotFoundEndpoint(t *testing.T) {
	e := echo.New()
	db, _ := setupTestDB(t)

	SetupRoutes(e, db)

	// /api以下は認証が必要なため401が返る
	// 認証不要なルートの存在しないパスをテスト
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
