package main

func isDeletionRequestToBeBlocked(resourceLabels map[string]string, resourceDeletionLockLabelKey string, resourceDeletionLockLabelValue string) bool {
	if len(resourceLabels) > 0 && string(resourceLabels[resourceDeletionLockLabelKey]) != "" && string(resourceLabels[resourceDeletionLockLabelKey]) == resourceDeletionLockLabelValue {
		return true
	}
	return false
}
