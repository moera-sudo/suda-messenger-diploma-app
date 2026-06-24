package mocks

import "context"

type MockAccessChecker struct {
	CheckAccessFn    func(ctx context.Context, userID, entityType, entityID string) (bool, error)
	CheckAccessCalls int
	CloseCalled      bool
}

func NewMockAccessChecker() *MockAccessChecker {
	return &MockAccessChecker{}
}

func (m *MockAccessChecker) CheckAccess(ctx context.Context, userID, entityType, entityID string) (bool, error) {
	m.CheckAccessCalls++
	if m.CheckAccessFn != nil {
		return m.CheckAccessFn(ctx, userID, entityType, entityID)
	}
	return true, nil
}

func (m *MockAccessChecker) Close() error {
	m.CloseCalled = true
	return nil
}