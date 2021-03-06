.. _platforms:

Supported Platforms
===================

Goiardi has been built and run with the native 6g compiler on Mac OS X (10.7 and above), Debian squeeze, wheezy, and jessie, a fairly recent Arch Linux, FreeBSD 9.2, Ubuntu 14.04, Solaris, and Raspbian (on both the original Raspberry Pi and the Raspberry Pi 2). Using Go's cross compiling capabilities, goiardi builds for all of Go's supported platforms except plan9 (because of issues with the postgres client library). Windows support has not been tested extensively, but a cross compiled binary has been tested successfully on Windows.

At one point goiardi was able to be built and run with gccgo (using the ``-compiler gccgo`` option with the ``go`` command) on Arch Linux. Unfortunately while recent gccgo versions include the ``go`` command, so building go programs with gccgo is theoretically much easier than before, it currently doesn't actually work because some dependencies blow up under gccgo.
