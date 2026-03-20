package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds CLI configuration
type Config struct {
	Server  string
	Token   string
	Profile string
}

// Profile represents a named configuration profile
type Profile struct {
	Server string `yaml:"server"`
	Token  string `yaml:"token"`
}

// ConfigFile represents the config file structure
type ConfigFile struct {
	CurrentProfile string              `yaml:"current_profile"`
	Profiles       map[string]*Profile `yaml:"profiles"`
}

// Load reads configuration from file and environment
func Load(profileName string) (*Config, error) {
	// Check environment variables first (highest priority)
	token := viper.GetString("token")
	server := viper.GetString("server")

	// If env vars are not set, read from config file
	if token == "" || server == "" {
		profiles := viper.GetStringMap("profiles")
		currentProfile := viper.GetString("current_profile")

		if currentProfile == "" {
			currentProfile = "default"
		}

		// Override with specified profile
		if profileName != "" && profileName != "default" {
			currentProfile = profileName
		}

		if profileData, ok := profiles[currentProfile]; ok {
			pmap, ok := profileData.(map[string]interface{})
			if ok {
				if token == "" {
					if t, ok := pmap["token"].(string); ok {
						token = t
					}
				}
				if server == "" {
					if s, ok := pmap["server"].(string); ok {
						server = s
					}
				}
			}
		}
	}

	if server == "" {
		return nil, fmt.Errorf("no server configured. Run: certimate config set --server URL --token TOKEN")
	}

	if token == "" {
		return nil, fmt.Errorf("no token configured. Run: certimate config set --server URL --token TOKEN")
	}

	return &Config{
		Server:  server,
		Token:   token,
		Profile: profileName,
	}, nil
}

// GetConfigPath returns the default config file path
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + "/.config/certimate-cli/config.yaml", nil
}

// Save saves configuration to file
func Save(profileName, server, token string) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create directory if not exists
	dir := configPath[:len(configPath)-len("/config.yaml")]
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Read existing config or create new
	cfg := &ConfigFile{
		CurrentProfile: "default",
		Profiles:       make(map[string]*Profile),
	}

	// Try to read existing config
	viper.ReadInConfig()
	if currentProfile := viper.GetString("current_profile"); currentProfile != "" {
		cfg.CurrentProfile = currentProfile
	}

	// Copy existing profiles
	profiles := viper.GetStringMap("profiles")
	for name, data := range profiles {
		if pmap, ok := data.(map[string]interface{}); ok {
			p := &Profile{}
			if s, ok := pmap["server"].(string); ok {
				p.Server = s
			}
			if t, ok := pmap["token"].(string); ok {
				p.Token = t
			}
			cfg.Profiles[name] = p
		}
	}

	// Update or create profile
	cfg.Profiles[profileName] = &Profile{
		Server: server,
		Token:  token,
	}
	cfg.CurrentProfile = profileName

	// Write config
	viper.Set("current_profile", cfg.CurrentProfile)
	viper.Set("profiles", cfg.Profiles)

	return viper.WriteConfigAs(configPath)
}

// GetCurrentProfile returns the current profile name
func GetCurrentProfile() string {
	current := viper.GetString("current_profile")
	if current == "" {
		return "default"
	}
	return current
}

// GetAllProfiles returns all configured profiles
func GetAllProfiles() (map[string]*Profile, string, error) {
	profiles := make(map[string]*Profile)

	pmap := viper.GetStringMap("profiles")
	for name, data := range pmap {
		if pm, ok := data.(map[string]interface{}); ok {
			p := &Profile{}
			if s, ok := pm["server"].(string); ok {
				p.Server = s
			}
			if t, ok := pm["token"].(string); ok {
				p.Token = t
			}
			profiles[name] = p
		}
	}

	current := viper.GetString("current_profile")
	if current == "" {
		current = "default"
	}

	return profiles, current, nil
}
