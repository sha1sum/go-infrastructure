# Webserver Infrastructure

This package has been open sourced in hopes some other gopher finds it helpful.
This package was not open sourced in hopes of becoming a really popular framework.

**What is this package all about?**
The Webserver provides infrastructure by supporting common webserver behavior
so gophers can focus on value.

This Webserver does not support the standard GO http interface but you can
easily work with handlers that do if you'd like. What this Webserver does do is
make it easier for gophers to build the kind of applications and platforms expected in my experieinces working with companies like Adobe, Dolby, TimeWarner Cable, and more. The basic building blocks and objectives of this package are to provide the following:

* Several convenience methods for working with input and output
* Means of passing information between handlers without requiring locking
* The ability to execute pre-handlers, post-handlers, and wrap handlers for custom behavior.
* Building blocks to help you auto-document handler/endpoint behavior.

# Developing the Webserver Infrastructure

#### Before You Start

* Verify you have the required software installed on your system.
* The Webserver Infrastructure is intended to run on Linux, Mac, and Windows and may be run on
either x86 (32/64bit) or ARM architecture. Download the correct Go distribution
for your environment.
* Ensure you have access to the source repository and database environments.

#### Required Software
1. [Go Distribution v1.4](http://golang.org/doc/install) distribution for your
system.
2. A working ``git`` installation and a registered ``ssh-key`` for the
repository.

## Download the Source

**Step 1:**
Use GO's get command to download the source code. The go get
command will use your git installation to download the source code and download
any dependent GO packages from their repositories. The source code will be
downloaded into your $GO_PATH directory. See the GO documentation if you
run into trouble.

## Compiling

This is a support library and will not compile on it's own.

## Writing Automated Tests
This paragraph still needs to be written.

### Go BDD Testing using github.com/onsi/ginkgo
Please familiarize yourself with the [Golang Ginkgo BDD Testing Framework](https://github.com/onsi/ginkgo).

## Questions?
Please contact Aaron Greenlee!

## Change Log
A version has not yet shipped.
