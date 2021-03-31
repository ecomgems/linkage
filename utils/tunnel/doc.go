// Package tunnel opens a tunnel using provided server and tunnel configuration and
// keeps it until something will break it or close it. As soon as the tunnel closed,
// it returns a signal to executing the routine using the channel flag.
package tunnel
