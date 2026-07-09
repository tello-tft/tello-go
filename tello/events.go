package tello

func ParseEvent(frame map[string]any) Event {
	eventType := stringValue(frame["type"])
	event := Event{
		Type:      eventType,
		Version:   stringValue(frame["version"]),
		CallID:    stringValue(frame["call_id"]),
		Timestamp: stringValue(frame["timestamp"]),
		Raw:       frame,
	}

	switch eventType {
	case EventTypeAgentsListed:
		event.RequestID = stringValue(frame["requestId"])
		event.Agents = agentInfos(frame["agents"])
	case EventTypeCallSummary:
		event.RequestID = stringValue(frame["requestId"])
		event.CallID = stringValue(frame["callId"])
		event.Status = stringValue(frame["status"])
		event.DurationSeconds = intValue(frame["durationSeconds"])
		event.Transcript = stringValue(frame["transcript"])
		event.Summary = stringValue(frame["summary"])
		event.CreditCharged = intValue(frame["creditCharged"])
	case EventTypeSmsSent:
		event.RequestID = stringValue(frame["requestId"])
		event.SmsID = stringValue(frame["smsId"])
		event.Status = stringValue(frame["status"])
		event.To = stringValue(frame["to"])
		event.MessagePreview = stringValue(frame["messagePreview"])
		event.CallID = stringValue(frame["callId"])
	case EventTypeUserTurn, EventTypeAgentTurn:
		event.TurnIndex = intValue(frame["turn_index"])
		event.Text = stringValue(frame["text"])
	case EventTypeCallStatusChanged:
		event.Status = stringValue(frame["status"])
		event.PreviousStatus = stringValue(frame["previous_status"])
	case EventTypeCallCompleted, EventTypeCallNoAnswer, EventTypeCallFailed:
		event.Status = stringValue(frame["status"])
		event.FailureReason = stringValue(frame["failure_reason"])
	case EventTypeError:
		event.Code = stringValue(frame["code"])
		event.Message = stringValue(frame["message"])
		event.RequestID = stringValue(frame["request_id"])
		event.Question = stringValue(frame["question"])
	}

	return event
}

func IsTerminal(event Event) bool {
	switch event.Type {
	case EventTypeCallCompleted, EventTypeCallNoAnswer, EventTypeCallFailed:
		return true
	case EventTypeCallStatusChanged:
		return event.Status == "cancelled"
	default:
		return false
	}
}

func agentInfos(value any) []AgentInfo {
	rows, ok := value.([]any)
	if !ok {
		return nil
	}
	agents := make([]AgentInfo, 0, len(rows))
	for _, row := range rows {
		agent, ok := row.(map[string]any)
		if !ok {
			continue
		}
		agents = append(agents, AgentInfo{
			AgentID:   stringValue(agent["agentId"]),
			Name:      stringValue(agent["name"]),
			Role:      stringValue(agent["role"]),
			IsDefault: boolValue(agent["isDefault"]),
			Status:    stringValue(agent["status"]),
		})
	}
	return agents
}

func stringValue(value any) string {
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}

func boolValue(value any) bool {
	boolean, ok := value.(bool)
	return ok && boolean
}

func intValue(value any) int {
	switch number := value.(type) {
	case int:
		return number
	case int64:
		return int(number)
	case float64:
		return int(number)
	default:
		return 0
	}
}
