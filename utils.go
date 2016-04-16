package main

import "time"

// Delete every empty strings in a slice of strings
func DeleteEmptyString(slice []string) []string {
	var localSlice []string
	for _, str := range slice {
		if str != "" {
			localSlice = append(localSlice, str)
		}
	}
	return localSlice
}

// Sleep a number of seconds
func DroneSleep(seconds int) {
	log.Debugf("sleeping for %d seconds", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)
}
