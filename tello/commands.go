package tello

import "encoding/json"

type CommandFrame map[string]any

func AuthFrame(apiKey, requestID string) CommandFrame {
	data := map[string]any{"token": apiKey}
	if requestID != "" {
		data["requestId"] = requestID
	}
	return CommandFrame{"event": "auth", "data": data}
}

func CreateCallFrame(to, prompt string, metadata map[string]any, requestID string) CommandFrame {
	data := map[string]any{
		"to":     to,
		"prompt": prompt,
	}
	if metadata != nil {
		data["metadata"] = metadata
	}
	if requestID != "" {
		data["requestId"] = requestID
	}
	return CommandFrame{"event": "createCall", "data": data}
}

func AnswerFrame(text, messageID, requestID string) CommandFrame {
	data := map[string]any{"text": text}
	if messageID != "" {
		data["messageId"] = messageID
	}
	if requestID != "" {
		data["requestId"] = requestID
	}
	return CommandFrame{"event": "answer", "data": data}
}

func SendDtmfFrame(digits, messageID, requestID string) CommandFrame {
	data := map[string]any{"digits": digits}
	if messageID != "" {
		data["messageId"] = messageID
	}
	if requestID != "" {
		data["requestId"] = requestID
	}
	return CommandFrame{"event": "sendDtmf", "data": data}
}

func CancelFrame() CommandFrame {
	return CommandFrame{"event": "cancel", "data": map[string]any{}}
}

func GetSummaryFrame(callID, requestID string) CommandFrame {
	data := map[string]any{"callId": callID}
	if requestID != "" {
		data["requestId"] = requestID
	}
	return CommandFrame{"event": "getSummary", "data": data}
}

func SendSmsFrame(to, message, requestID string) CommandFrame {
	data := map[string]any{"to": to, "message": message}
	if requestID != "" {
		data["requestId"] = requestID
	}
	return CommandFrame{"event": "sendSms", "data": data}
}

func Encode(frame CommandFrame) ([]byte, error) {
	return json.Marshal(frame)
}
