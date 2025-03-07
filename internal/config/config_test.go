package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
    "errors"
)

func TestRead(t *testing.T) {
	tempDir := t.TempDir() // Create a temporary home directory for testing
	tests := []struct {
		name          string
		setupConfig   func(path string) error
		expectError   bool
		expectedConfig *Config
	}{
		{
			name: "Missing Config File",
			setupConfig: func(path string) error {
				// Do nothing, file does not exist
				return nil
			},
			expectError:   true,
			expectedConfig: nil,
		},
		{
			name: "Valid Config File",
			setupConfig: func(path string) error {
				cfg := Config{DbUrl: "postgres://localhost:5432/db", CurrentUsername: "testuser"}
				data, err := json.Marshal(cfg)
				if err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(path, ".gatorconfig"), data, 0644)
			},
			expectError: false,
			expectedConfig: &Config{
				DbUrl:          "postgres://localhost:5432/db",
				CurrentUsername: "testuser",
			},
		},
		{
			name: "Invalid JSON Format",
			setupConfig: func(path string) error {
				return os.WriteFile(filepath.Join(path, ".gatorconfig"), []byte("invalid json"), 0644)
			},
			expectError:   true,
			expectedConfig: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the test environment
			if err := tt.setupConfig(tempDir); err != nil {
				t.Fatalf("failed to set up test: %v", err)
			}

			// Override the home directory
			oldHome := os.Getenv("HOME")
			defer os.Setenv("HOME", oldHome)
			os.Setenv("HOME", tempDir)

			// Execute the function
			cfg, err := Read()

			// Validate the results
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if cfg == nil || *cfg != *tt.expectedConfig {
					t.Errorf("expected config %+v, got %+v", tt.expectedConfig, cfg)
				}
			}
		})
	}
}

type mockWriter struct{
    writeFunc func() error
}

func (m mockWriter)write() error {
    if m.writeFunc != nil {
        return m.writeFunc()
    }
    return nil
}

func TestSetUser(t *testing.T) {
    tests := []struct{
        name string
        username string
        writeFunc func() error
        expectError bool
        want Config
    }{
        {
            // Overwrite write method to return an error
            name: "Expect error",
            username: "Ronny",
            writeFunc: func() error {
                return errors.New("hello I am an error")
            },
            expectError: true,
            want: Config{CurrentUsername: "Ronny"},
        },
        {
            // Overwrite write method to not return an error
            name: "Successfully set user",
            username: "Ragge",
            writeFunc: func() error {
                return nil
            },
            expectError: false,
            want: Config{CurrentUsername: "Ragge"},
        },
    } 

    for _, tt := range tests {
        mockConfig := mockWriter{writeFunc: tt.writeFunc}

        cfg := Config{
            configWriter: mockConfig,
        }

        err := cfg.SetUser(tt.username)
        if tt.expectError {
            if err == nil {
                t.Errorf("expected error but got none")
            }
        } else {
            if err != nil {
				t.Errorf("unexpected error: %v", err)
            }
        }

        if cfg.CurrentUsername != tt.want.CurrentUsername {
            t.Errorf("%s -> Config.SetUser(%s) got %v, want %v", tt.name, tt.username, cfg, tt.want)
        }
    }
}
