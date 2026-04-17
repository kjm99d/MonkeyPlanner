package main

import (
	"encoding/json"
	"fmt"
)

// decodeArgs unmarshals JSON-RPC tool arguments into the requested struct type.
// Missing arguments are tolerated (returns the zero value); malformed JSON returns an error.
func decodeArgs[T any](raw json.RawMessage) (T, error) {
	var v T
	if len(raw) == 0 {
		return v, nil
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return v, fmt.Errorf("invalid arguments: %w", err)
	}
	return v, nil
}

// unmarshalAny decodes an HTTP response body into a loosely-typed value.
// Errors from the HTTP layer are surfaced by the caller; JSON errors here are
// swallowed to keep the legacy behaviour of returning whatever parsed.
func unmarshalAny(raw json.RawMessage) any {
	var v any
	_ = json.Unmarshal(raw, &v)
	return v
}

func unknownToolError(name string) error {
	return fmt.Errorf("unknown tool: %s", name)
}
