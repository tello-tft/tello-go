package tello

import (
	"encoding/json"
	"testing"
)

func TestAuthFrameUsesEnvelopeAndOmitsEmptyRequestID(t *testing.T) {
	assertJSONEqual(t, map[string]any{
		"event": "auth",
		"data":  map[string]any{"token": "tello_live_x", "requestId": "auth-1"},
	}, AuthFrame("tello_live_x", "auth-1"))
	assertJSONEqual(t, map[string]any{
		"event": "auth",
		"data":  map[string]any{"token": "tello_live_x"},
	}, AuthFrame("tello_live_x", ""))
}

func TestCreateCallFrameUsesEnvelopeAndCamelCase(t *testing.T) {
	frame := CreateCallFrame("+821012345678", "hi", map[string]any{"src": "test"}, "r1")

	want := map[string]any{
		"event": "createCall",
		"data": map[string]any{
			"to":        "+821012345678",
			"prompt":    "hi",
			"metadata":  map[string]any{"src": "test"},
			"requestId": "r1",
		},
	}
	assertJSONEqual(t, want, frame)
}

func TestCreateCallFrameOmitsOptionalFields(t *testing.T) {
	assertJSONEqual(t, map[string]any{
		"event": "createCall",
		"data":  map[string]any{"to": "+821012345678", "prompt": ""},
	}, CreateCallFrame("+821012345678", "", nil, ""))
}

func TestCreateCallFrameNeverIncludesAgentID(t *testing.T) {
	frames := []CommandFrame{
		CreateCallFrame("+821012345678", "hi", nil, ""),
		CreateCallFrame("+821012345678", "hi", map[string]any{"src": "test"}, "r1"),
	}
	for _, frame := range frames {
		data, ok := frame["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data payload, got %v", frame["data"])
		}
		if _, exists := data["agentId"]; exists {
			t.Fatalf("createCall data must not contain agentId key, got %v", data)
		}
	}
}

func TestAnswerAndCancelFrames(t *testing.T) {
	assertJSONEqual(t, map[string]any{
		"event": "answer",
		"data":  map[string]any{"text": "yo", "messageId": "m1"},
	}, AnswerFrame("yo", "m1", ""))
	assertJSONEqual(t, map[string]any{"event": "cancel", "data": map[string]any{}}, CancelFrame())
}

func TestSendDtmfFrame(t *testing.T) {
	assertJSONEqual(t, map[string]any{
		"event": "sendDtmf",
		"data":  map[string]any{"digits": "1234#", "messageId": "m1", "requestId": "r1"},
	}, SendDtmfFrame("1234#", "m1", "r1"))
	assertJSONEqual(t, map[string]any{
		"event": "sendDtmf",
		"data":  map[string]any{"digits": "1234#"},
	}, SendDtmfFrame("1234#", "", ""))
}

func TestSummaryAndSmsFrames(t *testing.T) {
	assertJSONEqual(t, map[string]any{
		"event": "getSummary",
		"data":  map[string]any{"callId": "call-1", "requestId": "summary-1"},
	}, GetSummaryFrame("call-1", "summary-1"))
	assertJSONEqual(t, map[string]any{
		"event": "sendSms",
		"data": map[string]any{
			"to":        "01012345678",
			"message":   "예약 확인",
			"requestId": "sms-1",
		},
	}, SendSmsFrame("01012345678", "예약 확인", "sms-1"))
}

func assertJSONEqual(t *testing.T, want, got any) {
	t.Helper()
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if string(wantJSON) != string(gotJSON) {
		t.Fatalf("want %s, got %s", wantJSON, gotJSON)
	}
}
