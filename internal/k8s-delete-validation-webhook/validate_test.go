package main

import (
	"testing"
)

func Test_DeleteLockLabelKeyValueSet(t *testing.T) {
	resourceLabels := map[string]string{
		"app":         "myApp",
		"lockEnabled": "yes",
	}

	resourceDeletionLockLabelKey := "lockEnabled"
	resourceDeletionLockLabelValue := "yes"

	expected := true

	actual := isDeletionRequestToBeBlocked(resourceLabels, resourceDeletionLockLabelKey, resourceDeletionLockLabelValue)

	if actual != expected {
		t.Errorf("Expected deletionRequestToBeBlocked to be %t but instead got %t", expected, actual)
	}

}

func Test_NoResourceLabels(t *testing.T) {
	resourceLabels := map[string]string{}

	resourceDeletionLockLabelKey := "deleteLock"
	resourceDeletionLockLabelValue := "enabled"

	expected := false

	actual := isDeletionRequestToBeBlocked(resourceLabels, resourceDeletionLockLabelKey, resourceDeletionLockLabelValue)

	if actual != expected {
		t.Errorf("Expected deletionRequestToBeBlocked to be %t but instead got %t", expected, actual)
	}

}

func Test_DeleteLockLabelKeyValueNotSet(t *testing.T) {
	resourceLabels := map[string]string{
		"app":         "myApp",
		"lockEnabled": "yes",
	}

	resourceDeletionLockLabelKey := ""
	resourceDeletionLockLabelValue := ""

	expected := false

	actual := isDeletionRequestToBeBlocked(resourceLabels, resourceDeletionLockLabelKey, resourceDeletionLockLabelValue)

	if actual != expected {
		t.Errorf("Expected deletionRequestToBeBlocked to be %t but instead got %t", expected, actual)
	}

}

func Test_DeleteLockLabelValueMismatch(t *testing.T) {
	resourceLabels := map[string]string{
		"app":        "myApp",
		"deleteLock": "yes",
	}

	resourceDeletionLockLabelKey := "deleteLock"
	resourceDeletionLockLabelValue := "enabled"

	expected := false

	actual := isDeletionRequestToBeBlocked(resourceLabels, resourceDeletionLockLabelKey, resourceDeletionLockLabelValue)

	if actual != expected {
		t.Errorf("Expected deletionRequestToBeBlocked to be %t but instead got %t", expected, actual)
	}

}
