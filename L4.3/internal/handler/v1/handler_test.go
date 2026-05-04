package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"L4.3/internal/errs"
	"L4.3/internal/models"
	serviceMock "L4.3/internal/service/mocks"
	loggerMock "L4.3/pkg/logger/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_CreateEvent_Success(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)
	testHandler := NewHandler(mockService, mockLogger)

	gin.SetMode(gin.TestMode)

	eventDate := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02")
	body, _ := json.Marshal(CreateRequestV1{UserID: 1, EventDate: eventDate, Text: "ok"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.EXPECT().CreateEvent(gomock.Any()).Return("event-id", nil)

	testHandler.CreateEvent(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	result := resp["result"].(map[string]any)
	assert.Equal(t, "event-id", result["event_id"])

}

func TestHandler_CreateEvent_ErrInvalidJSON(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{invalid json}")))
	c.Request.Header.Set("Content-Type", "application/json")

	testHandler.CreateEvent(c)

	assertErrorResponse(t, w, http.StatusBadRequest, errs.ErrInvalidJSON.Error())

}

func TestHandler_CreateEvent_InvalidDate(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	body := []byte(`{"user_id": 1, "date": "THE DAY OF ABOBA", "text": "ok"}`)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	testHandler.CreateEvent(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, errs.ErrInvalidDateFormat.Error(), resp["error"])

}

func TestHandler_CreateEvent_ErrService(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	eventDate := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02")
	body, _ := json.Marshal(CreateRequestV1{UserID: 1, EventDate: eventDate, Text: "ok"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.EXPECT().CreateEvent(gomock.Any()).Return("", errors.New("service error"))

	testHandler.CreateEvent(c)

	assertErrorResponse(t, w, http.StatusInternalServerError, errs.ErrInternal.Error())

}

func TestHandler_UpdateEvent_Success(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	body, _ := json.Marshal(UpdateRequestV1{UserID: 1, EventID: "id", Text: "ok"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.EXPECT().UpdateEvent(gomock.Any()).Return(nil)

	testHandler.UpdateEvent(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	result := resp["result"].(map[string]any)
	assert.Equal(t, true, result["event_updated"])

}

func TestHandler_UpdateEvent_ErrInvalidJSON(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{invalid json}")))
	c.Request.Header.Set("Content-Type", "application/json")

	testHandler.UpdateEvent(c)

	assertErrorResponse(t, w, http.StatusBadRequest, errs.ErrInvalidJSON.Error())

}

func TestHandler_UpdateEvent_ErrInvalidDate(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{
		"user_id": 1,
		"event_id": "abc123",
		"new_date": "NOT-A-DATE"
	}`

	c.Request, _ = http.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewReader([]byte(body)),
	)

	c.Request.Header.Set("Content-Type", "application/json")

	testHandler.UpdateEvent(c)

	assertErrorResponse(t, w, http.StatusBadRequest, errs.ErrInvalidDateFormat.Error())

}

func TestHandler_UpdateEvent_ErrService(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	body, _ := json.Marshal(UpdateRequestV1{UserID: 1, EventID: "id", Text: "ok"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.EXPECT().UpdateEvent(gomock.Any()).Return(errors.New("service error"))

	testHandler.UpdateEvent(c)

	assertErrorResponse(t, w, http.StatusInternalServerError, "internal server error")

}

func TestHandler_DeleteEvent_Success(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	body, _ := json.Marshal(DeleteRequestV1{UserID: 1, EventID: "id"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.EXPECT().DeleteEvent(gomock.Any()).Return(nil)

	testHandler.DeleteEvent(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	result := resp["result"].(map[string]any)
	assert.Equal(t, true, result["event_deleted"])

}

func TestHandler_DeleteEvent_ErrInvalidJSON(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{invalid json}")))
	c.Request.Header.Set("Content-Type", "application/json")

	testHandler.DeleteEvent(c)

	assertErrorResponse(t, w, http.StatusBadRequest, errs.ErrInvalidJSON.Error())

}

func TestHandler_DeleteEvent_ErrService(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	body, _ := json.Marshal(DeleteRequestV1{UserID: 1, EventID: "id"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.EXPECT().DeleteEvent(gomock.Any()).Return(errors.New("service error"))

	testHandler.DeleteEvent(c)

	assertErrorResponse(t, w, http.StatusInternalServerError, "internal server error")

}

func TestHandler_GetEvents_Success(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?user_id=1&date=2025-12-03", nil)

	mockService.EXPECT().GetEvents(gomock.Any(), models.Day).Return([]models.Event{
		{Meta: models.Meta{UserID: 1, EventDate: time.Now()}, Data: models.Data{Text: "ok"}},
	}, nil)

	testHandler.GetEventsDay(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	result := resp["result"].(map[string]any)
	events := result["events"].([]any)
	event := events[0].(map[string]any)
	assert.Equal(t, "ok", event["text"])

}

func TestHandler_GetEvents_Week(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?user_id=1&date=2025-12-03", nil)

	mockService.EXPECT().GetEvents(gomock.Any(), models.Week).Return([]models.Event{
		{Meta: models.Meta{UserID: 1, EventDate: time.Now()}, Data: models.Data{Text: "ok"}},
	}, nil)

	testHandler.GetEventsWeek(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	result := resp["result"].(map[string]any)
	events := result["events"].([]any)
	event := events[0].(map[string]any)
	assert.Equal(t, "ok", event["text"])

}

func TestHandler_GetEvents_Month(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?user_id=1&date=2025-12-03", nil)

	mockService.EXPECT().GetEvents(gomock.Any(), models.Month).Return([]models.Event{
		{Meta: models.Meta{UserID: 1, EventDate: time.Now()}, Data: models.Data{Text: "ok"}},
	}, nil)

	testHandler.GetEventsMonth(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	result := resp["result"].(map[string]any)
	events := result["events"].([]any)
	event := events[0].(map[string]any)
	assert.Equal(t, "ok", event["text"])

}

func TestHandler_GetEvents_ParseQueryError(t *testing.T) {

	gin.SetMode(gin.TestMode)

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	router := gin.New()
	testHandler := NewHandler(mockService, mockLogger)

	router.GET("/events", func(c *gin.Context) {
		testHandler.GetEventsDay(c)
	})

	req, _ := http.NewRequest("GET", "/events?user_id=abc&date=2025-12-01", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, errs.ErrInvalidUserID.Error(), resp["error"])

}

func TestHandler_GetEvents_ErrService(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := serviceMock.NewMockService(controller)
	mockLogger := loggerMock.NewMockLogger(controller)

	testHandler := NewHandler(mockService, mockLogger)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?user_id=1&date=2025-12-03", nil)

	mockService.EXPECT().GetEvents(gomock.Any(), models.Day).Return(nil, errors.New("service error"))

	testHandler.GetEventsDay(c)

	assertErrorResponse(t, w, http.StatusInternalServerError, errs.ErrInternal.Error())

}

func assertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, wantStatus int, wantMsg string) {
	t.Helper()
	assert.Equal(t, wantStatus, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, wantMsg, resp["error"])
}
