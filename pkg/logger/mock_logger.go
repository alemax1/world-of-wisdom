package logger

import "go.uber.org/zap/zapcore"

type MockLogger struct{}

func NewMockLogger() *MockLogger { return &MockLogger{} }

func (MockLogger) Info(msg string, fields ...zapcore.Field)  {}
func (MockLogger) Error(msg string, fields ...zapcore.Field) {}
func (MockLogger) Warn(msg string, fields ...zapcore.Field)  {}
