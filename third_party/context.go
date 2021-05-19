package third_party

import (
	"context"
	"github.com/ciazhar/project-layout-rabbitmq/third_party/logger"
)

const SessionKey = 1

type Session struct {
	CID    string
	Logger logger.Logger
}

func SessionCid(ctx context.Context) string {
	session, ok := ctx.Value(SessionKey).(*Session)

	// Handle if session middleware is not used
	if !ok {
		return ""
	}

	return session.CID
}

func SessionLogger(ctx context.Context) logger.Logger {
	session, ok := ctx.Value(SessionKey).(*Session)

	// Handle if session middleware is not used
	if !ok {
		return logger.Log
	}

	return session.Logger
}

func NewSessionCtx(cid string, log logger.Logger) context.Context {
	session := Session{
		cid,
		log,
	}
	return context.WithValue(context.Background(), SessionKey, &session)
}
