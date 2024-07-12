/*
The package tests provides functions to compare files and directories between them.

The compare algorithm used can be found at https://github.com/golang/go/tree/master/src/internal/diff/diff.go.
As it's not exposed by golang source code, it was just copied and will not be modified unless the source code is modified.

Additional functions are exposed like EqualFiles and EqualDirs which respectively compare two files and two directories.
When comparing directories, all subdirectories are read and compared.
*/
package tests
