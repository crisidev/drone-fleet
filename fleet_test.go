package main

import "testing"

func TestSplitUnitPathOk(t *testing.T) {
	fleet := Fleet{}
	testUnitPath := "/var/test/myunit.myservice.org.service"
	testUnitName := "myunit.myservice.org.service"
	unitName, unitPathSplit := fleet.SplitUnitPath(testUnitPath)
	if unitName != testUnitName {
		t.Fatalf("expected unitName == %s, got %s", testUnitName, unitName)
	}
	if len(unitPathSplit) != 3 {
		t.Fatalf("expected unitPath lenght == 3, got %d", len(unitPathSplit))
	}
}
func TestSplitUnitPathFail(t *testing.T) {
	fleet := Fleet{}
	testUnitPath := "/var"
	unitName, unitPathSplit := fleet.SplitUnitPath(testUnitPath)
	if unitName != "" {
		t.Fatalf("expected empty unitName , got %s", unitName)
	}
	if len(unitPathSplit) != 0 {
		t.Fatalf("expected unitPath lenght == 0, got %d", len(unitPathSplit))
	}
}

func TestHandleScalableUnitOk(t *testing.T) {
	fleet := Fleet{}
	testUnitName := "myunit.myservice.org@.service"
	scaledUnitName := "myunit.myservice.org@0.service"
	unitName, _ := fleet.HandleScalableUnit(0, testUnitName)
	if unitName != scaledUnitName {
		t.Fatalf("expected unitName == %s, got %s", scaledUnitName, unitName)
	}
}

func TestHandleScalableUnitError(t *testing.T) {
	fleet := Fleet{}
	testUnitName := "myunit.myservice.org.service"
	_, err := fleet.HandleScalableUnit(0, testUnitName)
	if err == nil {
		t.Fatalf("expected err =! nil, got %s", err)
	}
}
