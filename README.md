debgo
======

dpkg-like functionality for building debs & source packages in Go.


**Warning: v0.2 is currently in a state of flux. Please use `go get github.com/laher/debgo-v0/deb` for the time-being. Ta**

Introduction
------------

debgo can produce 3 types of artifact:

 * The 'Binary debs' - per-architecture `.deb` packages, usually containing compiled artifacts..
 * The 'source packages' - a .dsc file plus 2 archives. Contains sources and build information.
 * The '-dev' package - a `.deb` file containing only sources. Commonly used as build dependencies

debgo has extra features for packaging 'go' applications, but in theory it could be used for various other tools.

*Note that Binary debs should normally be built from source debs. But, for Go programs in particular it's convenient to skip this step - especially when cross-compiling. The reason for this is because the Go cross-compiler is very straightforward, whereas the standard dpkg toolchain is not portable AFAIK.*

Libary use
----------

Once the API firms up, the preferred way to use debgo will be as a library.

 * The default behaviour uses text/template to generate files such as 'debian/control'.
 * Packages can be built with custom logic, by overriding the 'default' BuildFunc.
 * At the very least you can just use debgo to generate the final .deb files, having generated the contents elsewhere.

Basic commands
--------------

debgo comes with a few basic commands for building Debian packages. For the most part, each takes the same arguments.

`go get github.com/laher/debgo-v0.2/cmd/...`

 * debgo-deb produces .deb files for each architecture
 * debgo-source produces 3 'source package' files.
 * debgo-dev produces one '-dev.deb' file

