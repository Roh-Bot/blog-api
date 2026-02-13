package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Roh-Bot/blog-api/internal/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	requestIdKey  = "request_id"
	maxLogPreview = 4096 // limit to 4 KB for safety
)

func (s *Server) validateAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader, ok := ctx.Get("Authorization").(string)
		if !ok {
			return s.unauthorized(ctx, nil, stringEmpty)
		}
		if len(authHeader) < len("Bearer ") || authHeader[:7] != "Bearer " {
			return s.unauthorized(ctx, nil, stringEmpty)
		}
		token := authHeader[7:]
		isValid, err := s.App.Auth.ValidateToken(token)

		if err != nil || !isValid {
			s.Logger.Error(ctx.Request().Context(), err.Error())
			if errors.Is(err, auth.ErrTokenExpired) {
				return s.unauthorized(ctx, err, auth.ErrTokenExpired.Error())
			}
			return s.unauthorized(ctx, err, stringEmpty)
		}

		return next(ctx)
	}
}

func (s *Server) httpLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		res := ctx.Response()

		// Step 1: Generate and inject request ID into context
		reqID := uuid.New().String()
		ctxUser := context.WithValue(req.Context(), requestIdKey, reqID)
		ctx.SetRequest(req.WithContext(ctxUser))

		start := time.Now()

		// Step 2: Capture request body (safely restore after reading)
		var reqBodyBytes []byte
		if req.Body != nil {
			reqBodyBytes, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes)) // Restore so handlers can read it again
		}

		// Step 3: Wrap response writer to capture output
		var resBody bytes.Buffer
		mw := io.MultiWriter(res.Writer, &resBody)
		res.Writer = &bodyDumpResponseWriter{ResponseWriter: res.Writer, mw: mw}

		// Step 4: Call next middleware/handler
		err := next(ctx)
		latency := time.Since(start)

		// Step 5: Build log structure
		logEntry := map[string]any{
			"request_id": reqID,
			"time":       time.Now().Format(time.RFC3339),
			"method":     req.Method,
			"url":        req.URL.String(),
			"remote_ip":  ctx.RealIP(),
			"latency_ms": latency.Milliseconds(),
			"request": map[string]any{
				"headers": req.Header,
				"body":    truncateString(string(reqBodyBytes), 4096),
			},
			"response": map[string]any{
				"status": res.Status,
				"body":   truncateString(resBody.String(), 4096),
			},
		}

		if err != nil {
			logEntry["error"] = err.Error()
			s.Logger.Error(ctxUser, toJSON(logEntry))
		} else {
			s.Logger.Info(ctxUser, toJSON(logEntry))
		}

		return err
	}
}

func (s *Server) httpLoggerStream(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		res := ctx.Response()

		// Step 1: Attach request ID to context
		reqID := uuid.New().String()
		ctxUser := context.WithValue(req.Context(), requestIdKey, reqID)
		ctx.SetRequest(req.WithContext(ctxUser))
		start := time.Now()

		// Step 2: Capture request body safely (stream)
		var reqPreview bytes.Buffer
		if req.Body != nil {
			// io.TeeReader copies the read stream into reqPreview while still passing it along
			req.Body = io.NopCloser(io.TeeReader(req.Body, &limitWriter{Writer: &reqPreview, N: maxLogPreview}))
		}

		// Step 3: Capture response stream safely
		resPreview := &bytes.Buffer{}
		res.Writer = &bodyDumpResponseWriter{
			ResponseWriter: res.Writer,
			mw:             io.MultiWriter(res.Writer, &limitWriter{Writer: resPreview, N: maxLogPreview}),
		}

		// Step 4: Execute handler
		err := next(ctx)
		latency := time.Since(start)

		// Step 5: Build log entry
		logEntry := map[string]any{
			"request_id": reqID,
			"time":       time.Now().Format(time.RFC3339),
			"method":     req.Method,
			"url":        req.URL.String(),
			"remote_ip":  ctx.RealIP(),
			"latency_ms": latency.Milliseconds(),
			"request": map[string]any{
				"headers": req.Header,
				"body":    reqPreview.String(),
			},
			"response": map[string]any{
				"status": res.Status,
				"body":   resPreview.String(),
			},
		}

		if err != nil {
			logEntry["error"] = err.Error()
			s.Logger.Error(ctxUser, toJSON(logEntry))
		} else {
			s.Logger.Info(ctxUser, toJSON(logEntry))
		}

		return err
	}
}

// --- Helpers ---

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
	mw io.Writer
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.mw.Write(b)
}

// truncateString prevents massive logs by trimming large bodies
func truncateString(s string, max int) string {
	if len(s) > max {
		return s[:max] + "...[truncated]"
	}
	return s
}

func toJSON(v any) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// limitWriter writes up to N bytes and then ignores the rest — prevents large logs
type limitWriter struct {
	io.Writer
	N int
}

func (l *limitWriter) Write(p []byte) (int, error) {
	if len(p) > l.N {
		p = p[:l.N]
	}
	n, err := l.Writer.Write(p)
	l.N -= n
	if l.N <= 0 {
		// add truncation marker once limit reached
		l.Writer.Write([]byte("...[truncated]"))
		l.N = 0
	}
	return n, err
}
