// This file is safe to edit. Once it exists it will not be overwritten.

package api

// ServerShutdown is called when the HTTP(S) server is shut down and done.
// handling all active connections and does not accept connections any more.
func ServerShutdown() {}

// PreServerShutdown is called before the HTTP(S) server is shutdown.
// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic.
func PreServerShutdown() {}

// ServerStartup is called once at server start, just before any request can be handled.
func ServerStartup() {}
