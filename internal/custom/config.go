package custom

import "os"

// GetEnv retrieves the value of an environment variable or returns a default.
//
// Parameters:
//   - key: the name of the environment variable to look up. This should match
//          the variable name set in your OS, Docker container, or Kubernetes pod.
//          Example: "HOSTINFO_PORT"
//   - def: the default value to return if the environment variable is not set
//          or is empty. This allows your program to have sensible defaults
//          without requiring all environment variables to be defined.
//
// How it works:
//   1. The function checks if an environment variable with the given key exists.
//   2. If it exists and is not an empty string, it returns the value of that variable.
//   3. If the variable does not exist or is empty, it returns the default value provided
//      in the `def` parameter.
//
// Example usage:
//   port := custom.GetEnv("HOSTINFO_PORT", "8080")
//    If HOSTINFO_PORT is set to "9090" in the environment, port == "9090"
//    If HOSTINFO_PORT is not set, port == "8080"
//
// Notes:
//   - The function treats an empty string as "not set". If you explicitly want
//     to allow empty values, you need to handle that differently.
//   - This is useful for configuration settings in CLI apps, microservices,
//     Docker containers, or Kubernetes, where environment variables are commonly used
//     for passing runtime configuration.

func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
