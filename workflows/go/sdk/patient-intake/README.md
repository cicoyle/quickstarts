# Dapr Workflow History Propagation — Patient Intake

This example demonstrates how Dapr workflows can propagate their execution
history to child workflows and activities, so downstream consumers can
inspect the full (or partial) execution context of their caller.

The scenario is a patient intake / e-prescribing pipeline: a compliance
audit and a pharmacy dispense step refuse to act unless they can see
proof — in the propagated history — that the required upstream checks
(insurance, allergies, drug interactions) actually ran.

## Workflow architecture

```
PatientIntake (workflow)
├── VerifyInsurance (activity, no propagation)
└── PrescribeMedication (child workflow, PropagateLineage)
    ├── CheckAllergies (activity, no propagation)
    ├── ScreenDrugInteractions (activity, no propagation)
    ├── ComplianceAudit (child workflow, PropagateLineage)
    │     → sees PatientIntake + PrescribeMedication events
    └── DispenseMedication (activity, PropagateOwnHistory)
          → sees PrescribeMedication events only
```

### Propagation scope

| Mode | What it sends | Use case |
|------|---------------|----------|
| `PropagateLineage()` | Caller's own events + any ancestor events it received | Full chain-of-custody verification (compliance audits) |
| `PropagateOwnHistory()` | Caller's own events only (no ancestor chain) | Trust boundary — downstream only sees the immediate caller (pharmacy dispense) |

### Key demonstration

- **ComplianceAudit** receives the full lineage via `PropagateLineage()` —
  it verifies that `VerifyInsurance` ran in the grandparent workflow
  (PatientIntake), plus `CheckAllergies` and `ScreenDrugInteractions`
  ran in PrescribeMedication.

- **DispenseMedication** receives only PrescribeMedication's history via
  `PropagateOwnHistory()`. The PatientIntake ancestral history is excluded
  — the pharmacy system doesn't need (or get to see) the upstream chain.

## Running this example

1. Build and run with `dapr run`:

```bash
cd workflows/go/sdk/patient-intake
go build -o patient-app .
dapr run --app-id patient-app --resources-path ../../../components -- ./patient-app
```

Build the binary and run it directly (rather than `go run .`) for clean
signal handling. The app runs the demo once and exits on its own — no
Ctrl+C needed.

You'll see lines like:

```
[ComplianceAudit] Received propagated history: 15 events (scope: LINEAGE)
[ComplianceAudit] APPROVED (risk=0.10)
[DispenseMedication] Dispensing amoxicillin 500mg ... (propagated history: 12 events, scope=OWN_HISTORY)
[DispenseMedication] DISPENSED: rx-P-1042-...
```

In standalone mode the sidecar will log
`propagating unsigned workflow history to ...` warnings — these are
expected. Without `WorkflowHistorySigning` enabled, propagated history
chunks aren't cryptographically signed, which is fine for a local
`dapr run` demo. Signing the chunks within an mTLS trust boundary is a
production concern handled at the cluster/control-plane level and is out
of scope for this quickstart. 


## Files

```
patient-intake/
├── README.md      # this file
├── main.go        # registry + worker setup, schedules one workflow run
├── models.go      # PatientRecord, ComplianceResult, DispenseResult
├── workflow.go    # workflow + activity definitions, history helpers
└── go.mod         # module + deps
```
