package elastic

import "strings"

// Error handling utilities

// IsNotFoundError checks if an error is a document not found error
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "404")
}

// IsConflictError checks if an error is a version conflict error
func IsConflictError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "version_conflict")
}

// IsTimeoutError checks if an error is a timeout error
func IsTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline")
}

// IsConnectionError checks if an error is a connection error
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "network") ||
		strings.Contains(errStr, "dial")
}

// IsIndexNotFoundError checks if an error is an index not found error
func IsIndexNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "index_not_found_exception") ||
		strings.Contains(errStr, "no such index")
}

// IsDocumentExistsError checks if an error is a document already exists error
func IsDocumentExistsError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "version_conflict_engine_exception") ||
		strings.Contains(errStr, "document already exists")
}

// IsMappingError checks if an error is a mapping related error
func IsMappingError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "mapping") ||
		strings.Contains(errStr, "illegal_argument_exception")
}

// IsNetworkError checks if an error is a network-related error (enhanced version)
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "network") ||
		strings.Contains(errStr, "dial") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "no route to host") ||
		strings.Contains(errStr, "connection refused")
}
