# Bug Reproduction

## Bug

A quota reservation that exceeds the tenant's available budget loses its structured error identity across the repository, service, policy, and handler layers. The request is classified as retryable and returned as HTTP 503 instead of being rejected as quota exceeded.

## Trigger

Submit a reservation for a tenant with a requested quota greater than its available budget, then inspect the decision and returned status at the handler boundary.

## Error

```text
repository lost quota-exceeded sentinel: reserve quota for tenant "tenant-a": quota reservation exceeds available budget
service lost repository error chain: apply tenant quota "tenant-a": quota reservation exceeds available budget
quota excess classified as "retry", want reject
panic: quota excess mapped to status=503 decision=retry
```
