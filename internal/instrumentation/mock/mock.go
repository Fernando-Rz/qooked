package mock

import "log"

type MockInstrumentation struct{}

func NewMockInstrumentation() *MockInstrumentation {
	return &MockInstrumentation{}
}

func (instrumentation *MockInstrumentation) InitializeInstrumentation() error {
	// No initialization required for mock implementation
	return nil
}

func (instrumentation *MockInstrumentation) Log(message string) {
	log.Printf("Info: %s\n\n", message)
}

func (instrumentation *MockInstrumentation) LogError(message string) {
	log.Printf("Error: %s\n\n", message)
}

func (instrumentation *MockInstrumentation) LogIncomingHttpRequest(requestUrl string, statusCode int, durationInMilliseconds int64) {
	log.Printf("Incoming Http Request - URL: %s, Status: %d, Duration: %dms\n\n", requestUrl, statusCode, durationInMilliseconds)
}

func (instrumentation *MockInstrumentation) EmitMetric(metricName string, value int) {
	log.Printf("Metric - %s: %d\n\n", metricName, value)
}
