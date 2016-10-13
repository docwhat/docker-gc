Change Log
==========

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

[Unreleased]
------------

### Changed

-   Updated docker support.
-   Logger now logs docker actions correctly.
-   Fixed regression: Events are now listened for again.

[1.0.5] - 2016-05-26
--------------------

### Changed

-   Binaries are smaller now. About 8MiB to 5MiB.

[1.0.4] - 2016-05-26
--------------------

### Added

-   Better documentation for [contributors](CONTRIBUTING.md).
-   Released [docker container](https://hub.docker.com/r/docwhat/docker-gc/).

### Changed

-   All logging goes to `stderr` except for fatals.
-   `--quiet` doesn't hide errors anymore.
-   Removed timestamps from logging output.

[1.0.3] - 2016-05-26
--------------------

### Changed

-   Force all binaries to be statically compiled.

[1.0.2] - 2016-05-26
--------------------

### Added

-   Clean up dangling images too.
-   Added `--dangle-safe-duration=DUR` flag

### Changed

-   Replaced `--verbosity=NUM` flag with simpler `--quiet`

[1.0.1] - 2016-05-26
----------------------

### Changed

-   Nothing; version changed to fix build error.

[1.0.0] - 2016-05-26
----------------------

### Added

-   Everything!


[Unreleased]: https://github.com/docwhat/docker-gc/compare/1.0.5...HEAD
[1.0.5]: https://github.com/docwhat/docker-gc/compare/1.0.4...1.0.5
[1.0.4]: https://github.com/docwhat/docker-gc/compare/1.0.3...1.0.4
[1.0.3]: https://github.com/docwhat/docker-gc/compare/1.0.2...1.0.3
[1.0.2]: https://github.com/docwhat/docker-gc/compare/1.0.1...1.0.2
[1.0.1]: https://github.com/docwhat/docker-gc/compare/1.0.0...1.0.1
[1.0.0]: https://github.com/docwhat/docker-gc/commits/1.0.0
