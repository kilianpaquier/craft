/*
The fs package exposes some useful function around files, for instance a simple function as `func Exists(src) bool` to verify a file existence easily.

It also exposes `CopyFile(src, dest) error` which copies a given src file to dst with either specific permissions or not.

It also exposes `CopyDir(srcdir, destdir)` to copy a full directory at another place. The destination directory will be created if it doesn't already.

The package also exposes some constants around permissions.
*/
package fs
