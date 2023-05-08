package job

const (
	StatusPending    = "pending"    // StatusPending is the status of a job that is waiting to be processed.
	StatusProcessing = "processing" // StatusProcessing is the status of a job that is currently being processed.
	StatusCompleted  = "completed"  // StatusCompleted is the status of a job that has been successfully processed.
	StatusFailed     = "failed"     // StatusFailed is the status of a job that has failed to be processed.
)
