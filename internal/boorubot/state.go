package boorubot

import (
	"context"
	"encoding/json"
	"os"
)

// State data model
type State struct {
	LastPost uint64 `json:"last_post"`
}

// StateProvider persistence interface
type StateProvider interface {
	StateLoad(ctx context.Context) (*State, error)
	StateSave(ctx context.Context, state *State) error
}

// FileStateProvider jsonb encoded file-based implementation for StateProvider
type FileStateProvider struct {
	FileName string
}

// StateLoad implementation
func (fs *FileStateProvider) StateLoad(ctx context.Context) (*State, error) {
	var state State

	bs, err := os.ReadFile(fs.FileName)
	if os.IsNotExist(err) {
		return &state, nil
	}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, &state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

// StateSave implementation
func (fs *FileStateProvider) StateSave(ctx context.Context, state *State) error {
	bs, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(fs.FileName, bs, 0666)
}
