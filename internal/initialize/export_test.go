package initialize

import "io"

// SetReader overrides the input reader used by Run function (to reader user inputs from os.Stdin) to facilitate tests.
func SetReader(r io.Reader) {
	reader = r
}
