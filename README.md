# Webserver Infrastructure

The Webserver Infrastructure is responsible for managing common webserver
behavior so other Wrecker Labs projects can focus on business value. The code
in this repository may be included into other projects as a go package.

#### Links
* [Confluence Project Management](https://wreckingball.atlassian.net/wiki/display/IN/) - Project plan, meetings, etc.
* [JIRA Issue Tracker](https://wreckingball.atlassian.net/browse/IN/) - All infrastructure work tickets are managed through a single Infrastructure JIRA Project
* [Jenkins Build Server](http://home.aarongreenlee.com:8080/) - CI Server
* [Stash Git Version Control](http://git.wreckerlabs.com/projects/IN) - Source Code Repository

# Known Bugs
A version has not yet shipped.

# Developing the Webserver Infrastructure

#### Before You Start

* Verify you have the required software installed on your system.
* The Webserver Infrastructure is intended to run on Linux, Mac, and Windows and may be run on
either x86 (32/64bit) or ARM architecture. Download the correct Go distribution
for your environment.
* Ensure you have access to the source repository and database environments.

#### Required Software
1. [Go Distribution v1.3](http://golang.org/doc/install) distribution for your
system.
2. A working ``git`` installation and a registered ``ssh-key`` for the
repository.

## Download the Source

**Step 1:**
Use GO's get command to download the source code. The go get
command will use your git installation to download the source code and download
any dependent GO packages from their repositories. The source code will be
downloaded into your $GO_PATH directory. See the GO documentation if you
run into trouble. To download the source code run the following command:
``go get git.wreckerlabs.com/in/webserver``

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
