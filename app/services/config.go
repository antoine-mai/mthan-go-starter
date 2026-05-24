package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Config holds the application configuration settings.
type Config struct {
	AppName         string
	Port            string
	Env             string
	AdminActive     bool
	AdminPath       string
	AdminUser       string
	AdminPass       string
	Client          bool
	ClientBuildPath string
	StoragePath     string
	DbDriver        string
	DbUrl           string
}

// Default YAML configuration template fallback
const defaultYamlTemplate = `# Mthan App Configuration File
app:
  name: "mthan-app"

server:
  port: "8080"
  env: "development"

# Client Serving
# If set to true, the Go server will host the built React static files from client/build.
client: false

# Database Settings
database:
  driver: "sqlite"
  url: ""

# Admin Panel Settings
# The admin panel is activated only if active: true and all credential values are defined.
admin:
  active: true
  path: "/admin"
  user: "admin"
  pass: "admin@123"
`

// LoadConfig resolves the configuration path, handles auto-creation of the config.yaml file if missing,
// parses it, and returns the loaded configuration.
func LoadConfig() (*Config, error) {
	// 1. Detect standalone mode
	isStandalone := false
	for _, arg := range os.Args {
		if arg == "--standalone" || arg == "-standalone" {
			isStandalone = true
			break
		}
	}

	// 2. Resolve basic directories
	exePath, err := os.Executable()
	if err != nil {
		// Fallback to current working directory if executable path cannot be resolved
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}
		exePath = filepath.Join(cwd, "main")
	}
	exeDir := filepath.Dir(exePath)

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	// 3. Detect if running in development mode (checking if go.mod exists in CWD)
	isDev := false
	if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
		isDev = true
	}

	// 4. Detect App Name from root config.yaml template
	appName := "mthan-app" // fallback default
	var templateData []byte
	var readErr error

	// Try reading CWD template first, then exeDir template
	templateData, readErr = os.ReadFile(filepath.Join(cwd, "config.yaml"))
	if readErr != nil {
		templateData, readErr = os.ReadFile(filepath.Join(exeDir, "config.yaml"))
	}

	if readErr == nil {
		parsedTemplate, err := parseYAML(string(templateData))
		if err == nil && parsedTemplate["app_name"] != "" {
			appName = parsedTemplate["app_name"]
		}
	}

	// 5. Resolve target configuration file path and storage directory path
	var configFilePath string
	var storageDirPath string

	if isDev {
		// Local Dev Mode: directly use config.yaml in workspace root and storage/ in workspace root
		configFilePath = filepath.Join(cwd, "config.yaml")
		storageDirPath = filepath.Join(cwd, "storage")
	} else if isStandalone || runtime.GOOS != "linux" {
		// Standalone Mode
		configFilePath = filepath.Join(exeDir, "config", "config.yaml")
		storageDirPath = filepath.Join(exeDir, "storage")
	} else {
		// Linux User Mode: ~/.[appName]/config/config.yaml and ~/.[appName]/storage
		configFilePath = filepath.Join(homeDir, "."+appName, "config", "config.yaml")
		storageDirPath = filepath.Join(homeDir, "."+appName, "storage")
	}

	// 6. Auto-create parent directory of config and storage directory if they do not exist
	configDir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	if err := os.MkdirAll(storageDirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory %s: %w", storageDirPath, err)
	}

	// 7. Auto-create target config file if it does not exist by copying template or using fallback
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if readErr == nil {
			// Copy the template content
			if err := os.WriteFile(configFilePath, templateData, 0644); err != nil {
				return nil, fmt.Errorf("failed to copy template config file to %s: %w", configFilePath, err)
			}
			fmt.Printf("Copied configuration template to: %s\n", configFilePath)
		} else {
			// Fallback to embedded default YAML template if template file cannot be found
			if err := os.WriteFile(configFilePath, []byte(defaultYamlTemplate), 0644); err != nil {
				return nil, fmt.Errorf("failed to create default config file at %s: %w", configFilePath, err)
			}
			fmt.Printf("Created default configuration file from embedded template at: %s\n", configFilePath)
		}
	} else {
		fmt.Printf("Loading configuration from: %s\n", configFilePath)
	}

	fmt.Printf("Storage directory set to: %s\n", storageDirPath)

	// 8. Parse the target YAML file and populate environment variables
	if err := loadYAMLFile(configFilePath); err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// Resolve the client build directory path
	var clientBuildPath string
	if isStandalone {
		clientBuildPath = filepath.Join(exeDir, "client", "build")
	} else {
		clientBuildPath = filepath.Join(cwd, "app", "mods", "client", "build")
	}

	dbDriver := getEnv("DB_DRIVER", "sqlite")
	dbUrl := getEnv("DB_URL", "")
	// If SQLite is selected and URL is empty, automatically place db.sqlite in the resolved storage/data directory
	if dbUrl == "" && dbDriver == "sqlite" {
		dbUrl = filepath.Join(storageDirPath, "data", "db.sqlite")
	}

	// 9. Populate and return Config struct with environment variables (using defaults if necessary)
	cfg := &Config{
		AppName:         getEnv("APP_NAME", appName),
		Port:            getEnv("PORT", "8080"),
		Env:             getEnv("ENV", "development"),
		AdminActive:     getEnv("ADMIN_ACTIVE", "false") == "true",
		AdminPath:       getEnv("ADMIN_PATH", ""),
		AdminUser:       getEnv("ADMIN_USER", ""),
		AdminPass:       getEnv("ADMIN_PASS", ""),
		Client:          getEnv("CLIENT", "false") == "true",
		ClientBuildPath: clientBuildPath,
		StoragePath:     storageDirPath,
		DbDriver:        dbDriver,
		DbUrl:           dbUrl,
	}

	return cfg, nil
}

// loadYAMLFile reads a config.yaml file, parses it, and maps config keys to system environment variables.
func loadYAMLFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	yamlData, err := parseYAML(strings.Join(lines, "\n"))
	if err != nil {
		return err
	}

	// Maps parsed YAML keys to system environment variables
	keyMap := map[string]string{
		"app_name":        "APP_NAME",
		"server_port":     "PORT",
		"server_env":      "ENV",
		"admin_active":    "ADMIN_ACTIVE",
		"admin_path":      "ADMIN_PATH",
		"admin_user":      "ADMIN_USER",
		"admin_pass":      "ADMIN_PASS",
		"client":          "CLIENT",
		"database_driver": "DB_DRIVER",
		"database_url":    "DB_URL",
	}

	for yamlKey, envKey := range keyMap {
		if val, exists := yamlData[yamlKey]; exists {
			if os.Getenv(envKey) == "" {
				os.Setenv(envKey, val)
			}
		}
	}

	return nil
}

// parseYAML executes simple line-by-line parsing of indented key-value yaml configurations.
func parseYAML(data string) (map[string]string, error) {
	result := make(map[string]string)
	lines := strings.Split(data, "\n")
	inApp := false
	inAdmin := false
	inServer := false
	inDatabase := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			continue
		}

		// Check block headers
		if strings.HasPrefix(trimmed, "app:") {
			inApp = true
			inAdmin = false
			inServer = false
			inDatabase = false
			continue
		}
		if strings.HasPrefix(trimmed, "admin:") {
			inAdmin = true
			inApp = false
			inServer = false
			inDatabase = false
			continue
		}
		if strings.HasPrefix(trimmed, "server:") {
			inServer = true
			inApp = false
			inAdmin = false
			inDatabase = false
			continue
		}
		if strings.HasPrefix(trimmed, "database:") {
			inDatabase = true
			inApp = false
			inAdmin = false
			inServer = false
			continue
		}

		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Strip quotes
		if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) ||
			(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
			val = val[1 : len(val)-1]
		}

		// If block is indented
		if inApp && (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")) {
			result["app_"+key] = val
		} else if inAdmin && (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")) {
			result["admin_"+key] = val
		} else if inServer && (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")) {
			result["server_"+key] = val
		} else if inDatabase && (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")) {
			result["database_"+key] = val
		} else {
			inApp = false
			inAdmin = false
			inServer = false
			inDatabase = false
			result[key] = val
		}
	}

	return result, nil
}

// getEnv retrieves environment variable or returns a default value.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
