package audit

import (
	"bytes"
	"encoding/json"

	"github.com/skwair/harmony/permission"
)

// StringValues holds a pair of string values.
type StringValues struct {
	Old, New string
}

// IntValues holds a pair of integer values.
type IntValues struct {
	Old, New int
}

// BoolValues holds a pair of boolean values.
type BoolValues struct {
	Old, New bool
}

func stringValues(oldValue, newValue json.RawMessage) (old string, new string, err error) {
	if len(oldValue) != 0 {
		if err = json.Unmarshal(oldValue, &old); err != nil {
			return "", "", err
		}
	}

	if len(newValue) != 0 {
		if err = json.Unmarshal(newValue, &new); err != nil {
			return "", "", err
		}
	}

	return old, new, nil
}

func stringValue(val json.RawMessage) (string, error) {
	var s string

	if len(val) != 0 {
		if err := json.Unmarshal(val, &s); err != nil {
			return "", err
		}
	}

	return s, nil
}

func intValues(oldValue, newValue json.RawMessage) (old int, new int, err error) {
	if len(oldValue) != 0 {
		if err = json.Unmarshal(oldValue, &old); err != nil {
			return 0, 0, err
		}
	}

	if len(newValue) != 0 {
		if err = json.Unmarshal(newValue, &new); err != nil {
			return 0, 0, err
		}
	}

	return old, new, nil
}

func intValue(val json.RawMessage) (int, error) {
	var i int

	if len(val) != 0 {
		if err := json.Unmarshal(val, &i); err != nil {
			return 0, err
		}
	}

	return i, nil
}

func boolValues(oldValue, newValue json.RawMessage) (old bool, new bool, err error) {
	if len(oldValue) != 0 {
		if err = json.Unmarshal(oldValue, &old); err != nil {
			return false, false, err
		}
	}

	if len(newValue) != 0 {
		if err = json.Unmarshal(newValue, &new); err != nil {
			return false, false, err
		}
	}

	return old, new, nil
}

func boolValue(val json.RawMessage) (bool, error) {
	var b bool

	if len(val) != 0 {
		if err := json.Unmarshal(val, &b); err != nil {
			return false, err
		}
	}

	return b, nil
}

func permissionOverwritesValue(val json.RawMessage) ([]permission.Overwrite, error) {
	var perm []permission.Overwrite

	if len(val) != 0 && !bytes.Equal(val, []byte("[]")) {
		if err := json.Unmarshal(val, &perm); err != nil {
			return nil, err
		}
	}

	return perm, nil
}
