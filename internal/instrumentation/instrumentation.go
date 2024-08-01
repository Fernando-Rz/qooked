package instrumentation

type Instrumentation interface {
	InitializeInstrumentation() error
	Log(message string)
	LogError(message string)
	LogIncomingHttpRequest(requestUrl string, statusCode int, durationInMilliseconds int64)
	EmitMetric(metricName string, value int)
}
