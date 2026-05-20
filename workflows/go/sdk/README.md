# Dapr workflows — Go SDK quickstarts

This directory contains Go quickstart examples for the Dapr Workflows. 
Each example is self-contained: it has its own `README`, `dapr.yaml`, and `makefile`.

## Examples

### [`order-processor/`](./order-processor) — Intro

The canonical "getting started" example. A single `OrderProcessingWorkflow`
that purchases items from a store, using 5 activities (notify, verify
inventory, request approval, process payment, update inventory) against a
Redis state store.

Run it:

```sh
cd order-processor
dapr run -f .
```

### [`patient-intake/`](./patient-intake) — Workflow history propagation

Requirements: 
- Dapr v1.18+
- go-sdk v1.15+

A healthcare scenario demonstrating workflow history propagation: a root
`PatientIntake` workflow orders a prescription via a child workflow, which
in turn runs a `ComplianceAudit` child workflow and a `DispenseMedication`
activity. Downstream steps inspect the propagated execution history of
their callers to verify the required upstream checks actually ran before
acting.

## Layout

```
sdk/
├── README.md         # this landing page
├── makefile          # fans `make validate` out to each example
├── order-processor/  # intro example
└── patient-intake/   # workflow history propagation
```
