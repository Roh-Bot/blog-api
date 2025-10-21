package logger

import "context"

type MockLogger struct {
}

func (m *MockLogger) Info(ctx context.Context, msg string, fields map[string]any) {

}
func (m *MockLogger) Error(ctx context.Context, msg string, fields map[string]any) {

}
func (m *MockLogger) Warn(ctx context.Context, msg string, fields map[string]any) {

}
func (m *MockLogger) Debug(ctx context.Context, msg string, fields map[string]any) {

}
func (m *MockLogger) Flush() error {
	return nil
}
