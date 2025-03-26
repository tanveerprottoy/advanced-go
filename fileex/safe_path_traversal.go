package fileex

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

// Path traversal attacks
// “Path traversal” covers a number of related attacks following a common pattern:
// A program attempts to open a file in some known location, but an attacker causes
// it to open a file in a different location.
// If the attacker controls part of the filename, they may be able to use relative
// directory components ("..") to escape the intended location:

// f, err := os.Open(filepath.Join(trustedLocation, "../../../../etc/passwd"))
// On Windows systems, some names have special meaning:

// // f will print to the console.
// f, err := os.Create(filepath.Join(trustedLocation, "CONOUT$"))
// If the attacker controls part of the local filesystem, they may be able to use
// symbolic links to cause a program to access the wrong file:

// // Attacker links /home/user/.config to /home/otheruser/.config:
// err := os.WriteFile("/home/user/.config/foo", config, 0o666)

// If the program defends against symlink traversal by first verifying that the intended
// file does not contain any symlinks, it may still be vulnerable to time-of-check/time-of-use
// (TOCTOU) races, where the attacker creates a symlink after the program’s check:
// Another variety of TOCTOU race involves moving a directory that forms part of a
// path mid-traversal. For example, the attacker provides a path such as
// “a/b/c/../../etc/passwd”, and renames “a/b/c” to “a/b” while the open operation
// is in progress.
// symlinkCheck vulnerable to time-of-check/time-of-use
func symlinkCheck() error {
	unsafePath := "/home/user/.config/foo"
	// Validate the path before use.
	cleaned, err := filepath.EvalSymlinks(unsafePath)
	if err != nil {
		return err
	}

	if !filepath.IsLocal(cleaned) {
		return errors.New("unsafe path")
	}

	// Attacker replaces part of the path with a symlink.
	// The Open call follows the symlink:
	f, err := os.Open(cleaned)
	if err != nil {
		return err
	}

	log.Println(f)

	return nil
}

// Path sanitization
// When a program’s threat model does not include attackers with access to the local file
// system, it can be sufficient to validate untrusted input paths before use.
// Unfortunately, sanitizing paths can be surprisingly tricky, especially for portable
// programs that must handle both Unix and Windows paths.
// For example, on Windows filepath.IsAbs(`\foo`) reports false, because the path “\foo” is
// relative to the current drive.
// In Go 1.20, the path/filepath.IsLocal function was added, which reports whether a
// path is “local”. A “local” path is one which:
// does not escape the directory in which it is evaluated ("../etc/passwd" is not allowed);
// is not an absolute path ("/etc/passwd" is not allowed);
// is not empty ("" is not allowed);
// on Windows, is not a reserved name (“COM1” is not allowed).
// In Go 1.23, the path/filepath.Localize function was added, which converts a /-separated
// path into a local operating system path. Programs that accept and operate on potentially
// attacker-controlled paths should almost always use filepath.IsLocal
// or filepath.Localize to validate or sanitize those paths.
func localLocalize() {
	path := "directoy/file.txt"

	if !filepath.IsLocal(path) {
		log.Println("path is not local")
	}

	// TODO: filepath.Localize example, needs Go 1.23
}

// Beyond sanitization
// Path sanitization is not sufficient when attackers may have access to part of the local filesystem.
// Multi-user systems are uncommon these days, but attacker access to the filesystem can still occur
// in a variety of ways. An unarchiving utility that extracts a tar or zip file may be induced
// to extract a symbolic link and then extract a file name that traverses that link.
// A container runtime may give untrusted code access to a portion of the local filesystem.
// Programs may defend against unintended symlink traversal by using the
// path/filepath.EvalSymlinks function to resolve links in untrusted names before validation,
// but as described above this two-step process is vulnerable to TOCTOU races.

// Before Go 1.24, the safer option was to use a package such as github.com/google/safeopen,
// that provides path traversal-resistant functions for opening a potentially-untrusted
// filename within a specific directory.

// Introducing os.Root
// In Go 1.24, new APIs were introduced in the os package to safely open a file in a location
// in a traversal-resistent fashion.
// The new os.Root type represents a directory somewhere in the local filesystem. Open a root
// with the os.OpenRoot function:

// root, err := os.OpenRoot("/some/root/directory")
// if err != nil {
//   return err
// }
// defer root.Close()
// Root provides methods to operate on files within the root. These methods all accept
// filenames relative to the root, and disallow any operations that would escape from the
// root either using relative path components ("..") or symlinks.

// f, err := root.Open("path/to/file")
// Root permits relative path components and symlinks that do not escape the root.
// For example, root.Open("a/../b") is permitted. Filenames are resolved using the
// semantics of the local platform: On Unix systems, this will follow any symlink in
// “a” (so long as that link does not escape the root); while on Windows systems this
// will open “b” (even if “a” does not exist).

// Root currently provides the following set of operations:

// func (*Root) Create(string) (*File, error)
// func (*Root) Lstat(string) (fs.FileInfo, error)
// func (*Root) Mkdir(string, fs.FileMode) error
// func (*Root) Open(string) (*File, error)
// func (*Root) OpenFile(string, int, fs.FileMode) (*File, error)
// func (*Root) OpenRoot(string) (*Root, error)
// func (*Root) Remove(string) error
// func (*Root) Stat(string) (fs.FileInfo, error)
// In addition to the Root type, the new os.OpenInRoot function provides a simple way to open a potentially-untrusted filename within a specific directory:

// f, err := os.OpenInRoot("/some/root/directory", untrustedFilename)
// The Root type provides a simple, safe, portable API for operating with untrusted filenames.

// Caveats and considerations
// Unix
// On Unix systems, Root is implemented using the openat family of system calls.
// A Root contains a file descriptor referencing its root directory and will track
// that directory across renames or deletion.
// Root defends against symlink traversal but does not limit traversal of mount points.
// For example, Root does not prevent traversal of Linux bind mounts. Our threat model
// is that Root defends against filesystem constructs that may be created by ordinary
// users (such as symlinks), but does not handle ones that require root privileges
// to create (such as bind mounts).

// Windows
// On Windows, Root opens a handle referencing its root directory. The open handle
// prevents that directory from being renamed or deleted until the Root is closed.
// Root prevents access to reserved Windows device names such as NUL and COM1.

// Performance
// Root operations on filenames containing many directory components can be much more
// expensive than the equivalent non-Root operation. Resolving “..” components can also
// be expensive. Programs that want to limit the cost of filesystem operations can use
// filepath.Clean to remove “..” components from input filenames, and may want to limit
// the number of directory components.

// Who should use os.Root?
// You should use os.Root or os.OpenInRoot if:

// you are opening a file in a directory; AND
// the operation should not access a file outside that directory.
// For example, an archive extractor writing files to an output directory should use os.Root,
// because the filenames are potentially untrusted and it would be incorrect to write a file
// outside the output directory.

// However, a command-line program that writes output to a user-specified location should not use os.Root, because the filename is not untrusted and may refer to anywhere on the filesystem.

// As a good rule of thumb, code which calls filepath.Join to combine a fixed directory and an externally-provided filename should probably use os.Root instead.

// // This might open a file not located in baseDirectory.
// f, err := os.Open(filepath.Join(baseDirectory, filename))

// // This will only open files under baseDirectory.
// f, err := os.OpenInRoot(baseDirectory, filename)
func osRoot() error {
	root, err := os.OpenRoot("/some/root/directory")
	if err != nil {
		return err
	}

	defer root.Close()

	f, err := root.Open("path/to/file")
	if err != nil {
		return err
	}

	defer f.Close()

	untrustedFilename := "path/to/file"

	f, err = os.OpenInRoot("/some/root/directory", untrustedFilename)
	if err != nil {
		return err
	}

	defer f.Close()

	return nil
}