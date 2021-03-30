# mongo-go-driver-wrapper

Wrapper for the MongoDB Go driver to enable unit testing and dependency injection.

The official MongoDB driver for Go doesn't use interfaces making it tricky to unit test and use dependency injection or mocking when working with MongoDB databases.

This wrapper is intended to provide an abstraction layer around the driver to enable DI and more effective unit testing to avoid needing to use integration testing and a running DB instance.
