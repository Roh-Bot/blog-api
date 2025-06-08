package api

import "context"

type MockLogger struct {
}

func (m *MockLogger) Infoln(args ...any) {

}

func (m *MockLogger) Infof(template string, args ...any) {

}

func (m *MockLogger) Errorln(args ...any) {

}

func (m *MockLogger) Errorf(template string, args ...any) {

}

func (m *MockLogger) InfolnWithRequestId(ctx context.Context, args ...any) {

}

func (m *MockLogger) InfofWithRequestId(ctx context.Context, template string, args ...any) {

}

func (m *MockLogger) ErrorlnWithRequestId(ctx context.Context, args ...any) {

}

func (m *MockLogger) ErrorfWithRequestId(ctx context.Context, template string, args ...any) {

}
