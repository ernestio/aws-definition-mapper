# AWS definition mapper

This service will validate and map a user service definition into a valid ernest service. It service will respond to nats endpoints *definition.map.creation.aws* & *definition.map.deletion.aws*

## Build status

* master: [![CircleCI](https://circleci.com/gh/r3labs/aws-definition-mapper/tree/master.svg?style=svg)](https://circleci.com/gh/r3labs/aws-definition-mapper/tree/master)
* develop: [![CircleCI](https://circleci.com/gh/ErnestIO/aws-definition-mapper/tree/develop.svg?style=svg)](https://circleci.com/gh/r3labs/aws-definition-mapper/tree/develop)

## Installation

```
make deps
make install
```

## Running Tests

```
make test
```

## Contributing

Please read through our
[contributing guidelines](CONTRIBUTING.md).
Included are directions for opening issues, coding standards, and notes on
development.

Moreover, if your pull request contains patches or features, you must include
relevant unit tests.

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2015 r3labs.io authors.

Code released under
[the Mozilla Public License Version 2.0](LICENSE).

