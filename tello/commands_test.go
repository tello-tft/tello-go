package tello

import (
	"encoding/json"
	"testing"
)

func TestCreateCallFrameUsesEnvelopeAndCamelCase(t *testing.T) {
	frame := CreateCallFrame("+821012345678", "agent-1", "hi", map[string]any{"src": "test"}, "r1")

	want := map[string]any{
		"event": "create_call",
		"data": map[string]any{
			"to":        "+821012345678",
			"agentId":   "agent-1",
			"prompt":    "hi",
			"metadata":  map[string]any{"src": "test"},
			"requestId": "r1",
		},
	}
	assertJSONEqual(t, want, frame)
}

func TestCreateCallFrameOmitsOptionalFields(t *testing.T) {
	assertJSONEqual(t, map[string]any{
		"event": "create_call",
		"data":  map[string]any{"to": "+821012345678", "agentId": "agent-1", "prompt": ""},
	}, CreateCallFrame("+821012345678", "agent-1", "", nil, ""))
}

func TestAnswerAndCancelFrames(t *testing.T) {
	assertJSONEqual(t, map[string]any{
		"event": "answer",
		"data":  map[string]any{"text": "yo", "messageId": "m1"},
	}, AnswerFrame("yo", "m1", ""))
	assertJSONEqual(t, map[string]any{"event": "cancel", "data": map[string]any{}}, CancelFrame())
}

func TestListAgentsFrame(t *testing.T) {
	assertJSONEqual(t, map[string]any{
		"event": "listAgents",
		"data":  map[string]any{"requestId": "agents-1"},
	}, ListAgentsFrame("agents-1"))
	assertJSONEqual(t, map[string]any{"event": "listAgents", "data": map[string]any{}}, ListAgentsFrame(""))
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
			"callId":    "call-1",
			"requestId": "sms-1",
		},
	}, SendSmsFrame("01012345678", "예약 확인", "call-1", "sms-1"))
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
