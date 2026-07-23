[English](README.md) | **한국어**

# tello-go

Tello `/sdk` 프로토콜용 Go WebSocket SDK. SDK가 대화의 두뇌를 맡습니다.
게이트웨이는 진행 중인 통화에서 상대방이 말한 턴을 실시간으로 넘겨주고,
핸들러가 만든 답변은 다시 통화로 전달됩니다.

> 저장소: `tello-go` · 모듈: `github.com/tello-tft/tello-go` · 패키지: `tello`
>
> 전송 계층은 WebSocket뿐입니다. REST나 webhook은 제공하지 않습니다.

## 1. 설치

```bash
go get github.com/tello-tft/tello-go
```

Go 1.22 이상이 필요합니다. 런타임 의존성은 `github.com/gorilla/websocket`
하나뿐입니다.

## 2. API 키

`NewClient("")`는 `TELLO_API_KEY`에서 키를 읽습니다. URL은 `TELLO_URL`을 먼저
보고, 없으면 기본값 `ws://localhost:3000/sdk`를 씁니다. `WithURL`로 직접 지정할
수도 있습니다.

키 인증은 `Connect`가 내부에서 끝냅니다. 소켓이 열리면 `auth` 프레임
(`{"event":"auth","data":{"token":"<apiKey>"}}`)을 보내고, 서버가 `auth.ok`로
응답한 뒤에야 반환합니다. `Authorization` 헤더나 query string 토큰은 쓰지 않고,
키가 로그나 오류 메시지에 남지도 않습니다. 인증에 실패하거나, 서버가 `4401`로
연결을 닫거나, `WithOpenTimeout` 안에 `auth.ok`가 오지 않으면 `Connect`가 오류를
반환합니다. 인증이 끝나기 전에는 어떤 명령도 나가지 않습니다.

## 3. 연결 + 통화 시작

```go
package main

import (
	"context"
	"log"

	"github.com/tello-tft/tello-go/tello"
)

func main() {
	ctx := context.Background()
	client, err := tello.NewClient("")
	if err != nil {
		log.Fatal(err)
	}

	client.On(tello.EventTypeUserTurn, func(ctx context.Context, event tello.Event) error {
		return client.Answer(ctx, "heard: "+event.Text, "", "")
	})
	// DTMF 전송: client.SendDtmf(ctx, "1234#", "", "")

	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.CreateCall(ctx, "+821012345678", "예약 확인", nil, ""); err != nil {
		log.Fatal(err)
	}
	if err := client.WaitClosed(ctx); err != nil {
		log.Fatal(err)
	}
}
```

핸들러는 `Connect` 전에 등록하세요. 그래야 프레임을 놓치지 않습니다.

## 4. 실시간 턴 이벤트 (pub/sub)

`client.On(eventType string, handler func(context.Context, tello.Event) error)`로
이벤트 타입별 구독을 등록합니다. 모든 이벤트는 `tello.Event` 구조체 하나로
전달되고, 채워지는 필드는 타입마다 다릅니다. 디코딩된 원본 프레임은 언제나
`event.Raw`에 들어 있습니다.

| 상수 | 값 | 채워지는 필드 |
| --- | --- | --- |
| `EventTypeCallCreated` | `call.created` | `CallID`, `SessionID` |
| `EventTypeUserTurn` | `user.turn` | `TurnIndex`, `Text` |
| `EventTypeAgentTurn` | `agent.turn` | `TurnIndex`, `Text` |
| `EventTypeAnswerAccepted` | `answer.accepted` | `RequestID`, `MessageID` |
| `EventTypeDtmfAccepted` | `dtmf.accepted` | `RequestID`, `MessageID`, `Digits` |
| `EventTypeCallSummary` | `call.summary` | `RequestID`, `Status`, `DurationSeconds`, `Transcript`, `Summary`, `CreditCharged` |
| `EventTypeCallStatusChanged` | `call.statusChanged` | `Status`, `PreviousStatus` |
| `EventTypeCallCompleted` | `call.completed` | `Status` |
| `EventTypeCallNoAnswer` | `call.noAnswer` | `Status`, `FailureReason` |
| `EventTypeCallFailed` | `call.failed` | `Status`, `FailureReason` |
| `EventTypeError` | `error` | `Code`, `Message`, `RequestID`, `Question` |
| `EventTypeDisconnected` | `disconnected` | SDK 자체 이벤트. WS가 닫힐 때 발생 |

`auth.ok`는 `Connect`가 내부에서 소비하며 밖으로 다시 emit 하지 않습니다.

## 5. 명령

```go
client.CreateCall(ctx, to, prompt string, metadata map[string]any, requestID string) error
client.Answer(ctx, text, messageID, requestID string) error
client.SendDtmf(ctx, digits, messageID, requestID string) error
client.Cancel(ctx) error
client.GetSummary(ctx, callID, requestID string) error
```

선택 항목인 문자열 인자에는 `""`를 넘기면 됩니다. `requestID`는 명령과 응답
프레임을 짝지어 주는 값이며, 멱등성 키가 아닙니다.

`client.WaitClosed(ctx)`는 통화가 종료 상태(`call.completed` / `call.noAnswer` /
`call.failed`, 또는 cancelled 상태)에 이르거나 연결이 닫히면 반환합니다. 무한정
기다리지 않으려면 `context.WithTimeout`으로 상한을 거세요.

## 6. 오류 처리

게이트웨이 오류 프레임은 오류 타입과 1:1로 대응됩니다:

| 게이트웨이 `code` | 오류 타입 |
| --- | --- |
| `unauthenticated` | `*AuthenticationError` (인증 핸드셰이크. 4401 종료 포함) |
| `toRequired` | `*ValidationError` |
| `callIdRequired` | `*ValidationError` |
| `dtmfDigitsRequired` | `*ValidationError` |
| `dtmfDigitsInvalid` | `*ValidationError` |
| `callAlreadyActive` | `*CallAlreadyActiveError` |
| `noActiveCall` | `*NoActiveCallError` |
| `callRejected` | `*CallRejectedError` (`.Question` 포함) |
| `callNotFound` | `*TelloServerError` |
| `callNotCompleted` | `*TelloServerError` |
| `internalError` | `*TelloServerError` |

명령 단위 오류는 소켓을 닫지 않고 `EventTypeError` 구독자에게도 전달됩니다.
실패한 `CreateCall`(예: `toRequired`, `callRejected`)이 멈춘 채 남지 않도록,
`WaitClosed`가 그 오류를 그대로 반환합니다:

- 인증 실패(`unauthenticated` 프레임, 4401 종료, `auth.ok` 타임아웃) → `Connect`가 `*AuthenticationError` 반환
- 통화 시작 거부 → 위 표의 대응 오류
- 통화 도중 연결 끊김 → `*ConnectionClosedError`
- 다른 연결에 세션을 빼앗김(4429 종료) → `*SessionReplacedError`

이 오류 타입들은 공통 래퍼를 공유하지 않습니다. 그래서 `errors.As` 대상도
타입마다 따로 잡아야 합니다:

```go
var rejected *tello.CallRejectedError
if errors.As(err, &rejected) {
	log.Printf("call rejected: %s (%s)", rejected.Message, rejected.Question)
}
```

WS 수준 ping heartbeat는 게이트웨이가 주도하고, pong은 `gorilla/websocket`이
알아서 보냅니다. 재연결이나 세션 재개 프로토콜은 없습니다. 비정상 종료가 나면
재연결이 필요한 상황으로 보고 통화를 처음부터 다시 시작하세요.

## 7. 예제

바로 실행해 볼 수 있는 프로그램이 [`examples/`](examples/README.ko.md)에
있습니다:

```bash
go run ./examples/basic-call      # 연결, 통화 1건, 각 턴에 응답
go run ./examples/agent-callback  # 전체 수명주기, 이력, cancel, 타입별 오류
go run ./examples/call-summary    # 게이트로 막아 둔 라이브 시나리오 + call.summary
```

세 예제 모두 실제 통화를 겁니다. 먼저
[`examples/README.ko.md`](examples/README.ko.md)를 읽으세요.

## 8. 버전 호환성

`tello-go 0.1.x`는 Tello WS 프로토콜 `1.0`을 구현합니다
(`tello.ProtocolVersion`).
