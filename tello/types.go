package tello

const ProtocolVersion = "1.0"

const (
	EventTypeAuthOK            = "auth.ok"
	EventTypeUserTurn          = "user.turn"
	EventTypeAgentTurn         = "agent.turn"
	EventTypeAgentsListed      = "agents.listed"
	EventTypeCallSummary       = "call.summary"
	EventTypeSmsSent           = "sms.sent"
	EventTypeCallStatusChanged = "call.status_changed"
	EventTypeCallCompleted     = "call.completed"
	EventTypeCallNoAnswer      = "call.no_answer"
	EventTypeCallFailed        = "call.failed"
	EventTypeError             = "error"
	EventTypeDisconnected      = "disconnected"
)

type Event struct {
	Type            string
	Version         string
	CallID          string
	Timestamp       string
	Raw             map[string]any
	TurnIndex       int
	Text            string
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
