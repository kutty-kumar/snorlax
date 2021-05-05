package main

import "github.com/spf13/pflag"

const (
	defaultConfigDirectory = "../../deploy/"
	defaultConfigFile      = "config"
	defaultSecretFile      = ""
	defaultApplicationID   = "snorlax"

	// Heartbeat
	defaultKeepaliveTime    = 10
	defaultKeepaliveTimeout = 20
)

var (
	// define flag overrides
	flagConfigDirectory = pflag.String("config.source", defaultConfigDirectory, "directory of the configuration file")
	flagConfigFile      = pflag.String("config.file", defaultConfigFile, "directory of the configuration file")
	flagSecretFile      = pflag.String("config.secret.file", defaultSecretFile, "directory of the secrets configuration file")
	flagApplicationID   = pflag.String("app.id", defaultApplicationID, "identifier for the application")

	flagKeepaliveTime    = pflag.Int("config.keepalive.time", defaultKeepaliveTime, "default value, in seconds, of the keepalive time")
	flagKeepaliveTimeout = pflag.Int("config.keepalive.timeout", defaultKeepaliveTimeout, "default value, in seconds, of the keepalive timeout")
)
