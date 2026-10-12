package main

import (
	"fmt"
)

func validateRecoveryRequest(name string, migrateEnabled bool, request recoveryRequest, modeExplicit, stepsExplicit, diagnose, unsafe bool, migrateFlag, forceVersionFlag, forceClearDirtyFlag, modeFlag, stepsFlag, unsafeFlag string) error {
	if !request.requested() {
		return nil
	}

	if request.forceVersion < -1 {
		return fmt.Errorf("%s requires %s to be >= 0", name, flagRef(forceVersionFlag))
	}
	if request.forceVersion == -1 || !request.forceClearDirty {
		return fmt.Errorf("%s requires %s together with %s", name, flagRef(forceVersionFlag), flagRef(forceClearDirtyFlag))
	}
	if !migrateEnabled {
		return fmt.Errorf("%s requires %s", name, flagRef(migrateFlag))
	}
	if diagnose {
		return fmt.Errorf("%s cannot be combined with %s", name, flagRef(flagDiagnoseMigrations))
	}
	if modeExplicit {
		return fmt.Errorf("%s cannot be combined with %s", name, flagRef(modeFlag))
	}
	if stepsExplicit {
		return fmt.Errorf("%s cannot be combined with %s", name, flagRef(stepsFlag))
	}
	if unsafe {
		return fmt.Errorf("%s cannot be combined with %s", name, flagRef(unsafeFlag))
	}

	return nil
}
