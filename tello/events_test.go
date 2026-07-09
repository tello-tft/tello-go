package tello

import "testing"

func TestParseUserTurn(t *testing.T) {
	event := ParseEvent(map[string]any{
		"type":       "user.turn",
		"version":    "1.0",
		"call_id":    "c1",
		"turn_index": float64(2),
		"text":       "hey",
		"timestamp":  "t",
	})

	if event.Type != EventTypeUserTurn || event.TurnIndex != 2 || event.Text != "hey" || event.CallID != "c1" {
		t.Fatalf("unexpected event: %+v", event)
	}
}

func TestParseErrorFrame(t *testing.T) {
	event := ParseEvent(map[string]any{
		"type":       "error",
		"version":    "1.0",
		"code":       "call_rejected",
		"message":    "Call rejected",
		"request_id": "r1",
		"question":   "why?",
	})

	if event.Code != "call_rejected" || event.RequestID != "r1" || event.Question != "why?" {
		t.Fatalf("unexpected error event: %+v", event)
	}
}

func TestParseAgentsListed(t *testing.T) {
	event := ParseEvent(map[string]any{
		"type":      "agents.listed",
		"version":   "1.0",
		"requestId": "agents-1",
		"agents": []any{
			map[string]any{
				"agentId":   "agent-1",
				"name":      "예약 확인",
				"role":      "AI 상담원",
				"isDefault": true,
				"status":    "published",
			},
		},
	})

	if event.Type != EventTypeAgentsListed || event.RequestID != "agents-1" {
		t.Fatalf("unexpected agents event: %+v", event)
	}
	if len(event.Agents) != 1 {
		t.Fatalf("unexpected agents: %+v", event.Agents)
	}
	agent := event.Agents[0]
	if agent.AgentID != "agent-1" || agent.Name != "예약 확인" || agent.Role != "AI 상담원" || !agent.IsDefault || agent.Status != "published" {
		t.Fatalf("unexpected agent: %+v", agent)
	}
}

func TestParseCallSummaryAndSmsSent(t *testing.T) {
	summary := ParseEvent(map[string]any{
		"type":            "call.summary",
		"version":         "1.0",
		"requestId":       "summary-1",
		"callId":          "call-1",
		"status":          "completed",
		"durationSeconds": float64(42),
		"transcript":      "고객: 예약 확인",
		"summary":         "예약 확인 완료",
		"creditCharged":   float64(15),
	})
	if summary.Type != EventTypeCallSummary || summary.RequestID != "summary-1" || summary.CallID != "call-1" || summary.DurationSeconds != 42 || summary.CreditCharged != 15 {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	sms := ParseEvent(map[string]any{
		"type":           "sms.sent",
		"version":        "1.0",
		"requestId":      "sms-1",
		"smsId":          "77",
		"status":         "queued",
		"to":             "01012345678",
		"messagePreview": "예약 확인",
		"callId":         "call-1",
	})
	if sms.Type != EventTypeSmsSent || sms.RequestID != "sms-1" || sms.SmsID != "77" || sms.CallID != "call-1" {
		t.Fatalf("unexpected sms: %+v", sms)
	}
}

func TestTerminalDetection(t *testing.T) {
	if !IsTerminal(ParseEvent(map[string]any{
		"type":      "call.completed",
		"version":   "1.0",
		"call_id":   "c1",
		"status":    "completed",
		"timestamp": "t",
	})) {
		t.Fatal("completed should be terminal")
	}
	if !IsTerminal(ParseEvent(map[string]any{
		"type":            "call.status_changed",
		"version":         "1.0",
		"call_id":         "c1",
		"status":          "cancelled",
		"previous_status": "in_progress",
		"timestamp":       "t",
	})) {
		t.Fatal("cancelled should be terminal")
	}
}
