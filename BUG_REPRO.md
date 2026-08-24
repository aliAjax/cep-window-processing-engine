# Bug Reproduction

## Bug

When no external variable provider is configured, the default provider carries uninitialized storage. A disabled provider can also survive as a typed-nil interface and pass validation, so the first variable write reaches a nil map and panics.

## Trigger

Start the engine without an external provider, build and validate the default provider, then write the first pattern variable. The same failure is observable by passing a typed-nil provider through the factory and validation path.

## Observed Errors

```text
default provider storage was not initialized
disabled provider must be a nil interface, got *provider.memoryProvider
typed-nil provider passed validation
panic: assignment to entry in nil map
```
