package tello

const ProtocolVersion = "1.0"

const (
	EventTypeAuthOK            = "auth.ok"
	EventTypeCallCreated       = "call.created"
	EventTypeUserTurn          = "user.turn"
	EventTypeAgentTurn         = "agent.turn"
	EventTypeAgentsListed      = "agents.listed"
	EventTypeCallSummary       = "call.summary"
	EventTypeSmsSent           = "sms.sent"
	EventTypeAnswerAccepted    = "answer.accepted"
	EventTypeDtmfAccepted      = "dtmf.accepted"
	EventTypeCallStatusChanged = "call.statusChanged"
	EventTypeCallCompleted     = "call.completed"
	EventTypeCallNoAnswer      = "call.noAnswer"
	EventTypeCallFailed        = "call.failed"
	EventTypeError             = "error"
	EventTypeDisconnected      = "disconnected"
)

type Event struct {
	Type            string
	Version         string
	SessionID       string
	CallID          string
	Timestamp       string
	Raw             map[string]any
	TurnIndex       int
	Text            string
	MessageID       string
	Digits          string
	Status          string
	PreviousStatus  string
	FailureReason   string
	Code            string
	Message         string
	RequestID       string
	Question        string
	Agents          []AgentInfo
	DurationSeconds int
	Transcript      string
	Summary         string
	CreditCharged   int
	SmsID           string
	To              string
	MessagePreview  string
}

type AgentInfo struct {
	AgentID   string
	Name      string
	Role      string
	IsDefault bool
	Status    string
}
