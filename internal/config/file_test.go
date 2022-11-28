package config

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidateBasicAuthExampleConfigFile(t *testing.T) {
	// Save old command-line arguments so that we can restore them later
	// https://stackoverflow.com/questions/33723300/how-to-test-the-passing-of-arguments-in-golang
	oldArgs := os.Args

	defer func() {
		t.Log("Restoring os.Args to original value")
		os.Args = oldArgs
	}()

	normalizedFullyQualifiedFilename := filepath.Join("../../", "contrib", "list-emails", "basic-auth", "accounts.example.ini")

	// Mock flags for application.
	// Rely on default values for log and report output directories.
	//
	// Note to self: Don't add/escape double-quotes here. The shell strips
	// them away and the application never sees them.
	os.Args = []string{
		"/usr/local/bin/list-emails",
		"--config-file", normalizedFullyQualifiedFilename,
		"--log-file-dir", t.TempDir(),
	}

	// Setup configuration by parsing user-provided flags
	cfg, err := New(AppType{ReporterIMAPMailbox: true})
	switch {
	case errors.Is(err, ErrVersionRequested):
		t.Log(Version())
		t.SkipNow()

	case errors.Is(err, ErrHelpRequested):
		t.Log(cfg.Help())
		t.SkipNow()

	case err != nil:
		t.Fatalf("Error initializing application: %v", err)
	}
	t.Log("Config initialization complete")

	t.Cleanup(func() {
		if err := cfg.LogFileHandle.Close(); err != nil {
			// Ignore "file already closed" errors
			if !errors.Is(err, os.ErrClosed) {
				t.Fatalf("Failed to close log file: %v", err)
			}
		}

		// Log directory is removed by the testing harness due to our use
		// of the (testing.T).TempDir method.
	})

	// Assert that the config file has two accounts listed.
	if want, got := 2, len(cfg.Accounts); want != got {
		t.Errorf("ERROR: \nwant %d accounts\ngot %d accounts", want, got)
	} else {
		t.Logf("OK: want %d accounts, got %d accounts", want, got)
	}

	// Manually craft slice of accounts to reflect expected config file
	// values.
	wantedAccounts := []MailAccount{
		{
			Name:     "email1",
			Server:   "imap.example.com",
			Port:     993,
			AuthType: AuthTypeBasic,
			Folders: []string{
				"Inbox",
				"Junk EMail",
				"Trash",
				"INBOX/Current Reporting",
			},
			Username: "email1@example.com",
			Password: "keepingOnKeepingON",
		},
		{
			Name:     "email2",
			Server:   "imap.example.com",
			Port:     993,
			AuthType: AuthTypeBasic,
			Folders: []string{
				"Inbox",
				"Junk EMail",
			},
			Username: "email2@example.com",
			Password: "tired",
		},
	}

	// A difference in entry order between config file and our "wanted" slice
	// shouldn't be problematic, but it's worth noting.
	if want, got := cfg.Accounts[0].Name, wantedAccounts[0].Name; want != got {
		t.Logf(
			"WARN: Mismatch between first account entries: want name %q, got name %q (order mismatch?)",
			want, got,
		)
		t.Log("WARN: Check order of wanted accounts collection and example config file")
	}

	trans := cmp.Transformer("Sort", func(in []MailAccount) []MailAccount {
		out := append([]MailAccount(nil), in...) // Copy input to avoid mutating it

		sort.Slice(out, func(i, j int) bool {
			return out[i].Name < out[j].Name
		})
		return out
	})

	// NOTE: We override or "transform" the order of entries in the compared
	// slices so that equality comparison works as expected.
	if d := cmp.Diff(wantedAccounts, cfg.Accounts, trans); d != "" {
		t.Errorf("(-want, +got)\n:%s", d)
	} else {
		t.Logf("OK: All accounts from %s match expected values", normalizedFullyQualifiedFilename)
	}

}

func TestValidateOAuth2ExampleConfigFile(t *testing.T) {
	// NOTE: Not sure if running in parallel would work well
	// with os.Args manipulation.
	//
	// t.Parallel()

	// Save old command-line arguments so that we can restore them later
	oldArgs := os.Args

	// Defer restoring original command-line arguments
	defer func() { os.Args = oldArgs }()

	normalizedFullyQualifiedFilename := filepath.Join("../../", "contrib", "list-emails", "oauth2", "accounts.example.ini")

	// Mock flags for application.
	// Rely on default values for log and report output directories.
	//
	// Note to self: Don't add/escape double-quotes here. The shell strips
	// them away and the application never sees them.
	os.Args = []string{
		"/usr/local/bin/list-emails",
		"--config-file", normalizedFullyQualifiedFilename,
		"--log-file-dir", t.TempDir(),
	}

	// Setup configuration by parsing user-provided flags
	cfg, err := New(AppType{ReporterIMAPMailbox: true})
	switch {
	case errors.Is(err, ErrVersionRequested):
		t.Log(Version())
		t.SkipNow()

	case errors.Is(err, ErrHelpRequested):
		t.Log(cfg.Help())
		t.SkipNow()

	case err != nil:
		t.Fatalf("Error initializing application: %v", err)
	}
	t.Log("Config initialization complete")

	t.Cleanup(func() {
		if err := cfg.LogFileHandle.Close(); err != nil {
			// Ignore "file already closed" errors
			if !errors.Is(err, os.ErrClosed) {
				t.Fatalf("Failed to close log file: %v", err)
			}
		}

		// Log directory is removed by the testing harness due to our use
		// of the (testing.T).TempDir method.
	})

	// Assert that the config file has two accounts listed.
	if want, got := 2, len(cfg.Accounts); want != got {
		t.Errorf("ERROR: \nwant %d accounts\ngot %d accounts", want, got)
	} else {
		t.Logf("OK: want %d accounts, got %d accounts", want, got)
	}

	// Manually craft slice of accounts to reflect expected config file
	// values.
	wantedAccounts := []MailAccount{
		{
			Name:     "email1",
			Server:   "outlook.office365.com",
			Port:     993,
			AuthType: AuthTypeOAuth2ClientCreds,
			Folders: []string{
				"Inbox",
				"Junk EMail",
				"Trash",
				"INBOX/Current Reporting",
			},
			OAuth2Settings: OAuth2ClientCredentialsFlow{
				SharedMailbox: "email1@example.com",
				ClientID:      "ZYDPLLBWSK3MVQJSIYHB1OR2JXCY0X2C5UJ2QAR2MAAIT5Q",
				ClientSecret:  "_djgA8heFo0WSIMom7U39WmGTQFHWkcD8x-A1o-4sro",
				Scopes:        []string{"https://outlook.office365.com/.default"},
				TokenURL:      "https://login.microsoftonline.com/6029c1d9-aa2f-4227-8f7c-0c23224a0fa9/oauth2/v2.0/token",
			},
		},
		{
			Name:     "email2",
			Server:   "outlook.office365.com",
			Port:     993,
			AuthType: AuthTypeOAuth2ClientCreds,
			Folders: []string{
				"Inbox",
				"Junk EMail",
			},
			OAuth2Settings: OAuth2ClientCredentialsFlow{
				SharedMailbox: "email2@example.com",
				ClientID:      "ZYDPLLBWSK3MVQJSIYHB1OR2JXCY0X2C5UJ2QAR2MAAIT5Q",
				ClientSecret:  "_djgA8heFo0WSIMom7U39WmGTQFHWkcD8x-A1o-4sro",
				Scopes:        []string{"https://outlook.office365.com/.default"},
				TokenURL:      "https://login.microsoftonline.com/6029c1d9-aa2f-4227-8f7c-0c23224a0fa9/oauth2/v2.0/token",
			},
		},
	}

	// A difference in entry order between config file and our "wanted" slice
	// shouldn't be problematic, but it's worth noting.
	if want, got := cfg.Accounts[0].Name, wantedAccounts[0].Name; want != got {
		t.Logf(
			"WARN: Mismatch between first account entries: want name %q, got name %q (order mismatch?)",
			want, got,
		)
		t.Log("WARN: Check order of wanted accounts collection and example config file")
	}

	trans := cmp.Transformer("Sort", func(in []MailAccount) []MailAccount {
		out := append([]MailAccount(nil), in...) // Copy input to avoid mutating it

		sort.Slice(out, func(i, j int) bool {
			return out[i].Name < out[j].Name
		})
		return out
	})

	// NOTE: We override or "transform" the order of entries in the compared
	// slices so that equality comparison works as expected.
	if d := cmp.Diff(wantedAccounts, cfg.Accounts, trans); d != "" {
		t.Errorf("(-want, +got)\n:%s", d)
	} else {
		t.Logf("OK: All accounts from %s match expected values", normalizedFullyQualifiedFilename)
	}
}

// TestParseExampleConfigFiles tests whether the config files can be loaded,
// parsed and pre/post validation of configuration passes. Assertion of
// specific values is not performed.
func TestParseExampleConfigFiles(t *testing.T) {
	// NOTE: Not sure if running in parallel would work well
	// with os.Args manipulation.
	//
	// t.Parallel()

	tests := map[string]struct {
		filePath string
	}{
		"basic-auth": {
			filePath: filepath.Join("../../", "contrib", "list-emails", "basic-auth", "accounts.example.ini"),
		},
		"oauth2": {
			filePath: filepath.Join("../../", "contrib", "list-emails", "oauth2", "accounts.example.ini"),
		},
	}

	for testName, testCase := range tests {

		t.Run(testName, func(t *testing.T) {

			// Save old command-line arguments so that we can restore them later
			// https://stackoverflow.com/questions/33723300/how-to-test-the-passing-of-arguments-in-golang
			oldArgs := os.Args

			defer func() {
				t.Log("Restoring os.Args to original value")
				os.Args = oldArgs
			}()

			// Mock flags for application.
			//
			// Note to self: Don't add/escape double-quotes here. The shell strips
			// them away and the application never sees them.
			os.Args = []string{
				"/usr/local/bin/list-emails",
				"--config-file", testCase.filePath,
				"--log-file-dir", t.TempDir(),
			}

			// Setup configuration by parsing user-provided flags
			cfg, err := New(AppType{ReporterIMAPMailbox: true})
			switch {
			case errors.Is(err, ErrVersionRequested):
				t.Log(Version())
				t.SkipNow()

			case errors.Is(err, ErrHelpRequested):
				t.Log(cfg.Help())
				t.SkipNow()

			case err != nil:
				t.Fatalf("Error initializing application: %v", err)
			}
			t.Log("Config initialization complete")

			t.Cleanup(func() {
				if err := cfg.LogFileHandle.Close(); err != nil {
					// Ignore "file already closed" errors
					if !errors.Is(err, os.ErrClosed) {
						t.Fatalf("Failed to close log file: %v", err)
					}
				}

				// Log directory is removed by the testing harness due to our use
				// of the (testing.T).TempDir method.
			})

			t.Logf("Config struct: +%v", cfg)

		})
	}

}
