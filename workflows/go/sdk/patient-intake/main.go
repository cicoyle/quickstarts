// This quickstart demonstrates workflow history propagation in a
// patient intake / e-prescribing scenario. A root PatientIntake workflow
// orders a prescription via a child PrescribeMedication workflow, which
// in turn runs a ComplianceAudit child workflow and a DispenseMedication
// activity. The compliance audit and dispensing steps inspect the
// propagated execution history of their callers to verify that the
// required upstream checks (insurance, allergies, drug interactions)
// actually ran before they make a decision.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dapr/durabletask-go/workflow"
	"github.com/dapr/go-sdk/client"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func main() {
	r := workflow.NewRegistry()

	// Step 1: PatientIntake (root workflow)
	//   Step 1.1: VerifyInsurance (activity, no propagation)
	//   Step 1.2: PrescribeMedication (child wf, propagate lineage)
	//     Step 1.2.1: CheckAllergies (activity, no propagation)
	//     Step 1.2.2: ScreenDrugInteractions (activity, no propagation)
	//     Step 1.2.3: ComplianceAudit (grandchild wf, propagate lineage)
	//     Step 1.2.4: DispenseMedication (activity, propagate own history)
	for _, add := range []func() error{
		func() error { return r.AddWorkflow(PatientIntake) },
		func() error { return r.AddActivity(VerifyInsurance) },
		func() error { return r.AddWorkflow(PrescribeMedication) },
		func() error { return r.AddActivity(CheckAllergies) },
		func() error { return r.AddActivity(ScreenDrugInteractions) },
		func() error { return r.AddWorkflow(ComplianceAudit) },
		func() error { return r.AddActivity(DispenseMedication) },
	} {
		if err := add(); err != nil {
			logger.Fatal(err)
		}
	}

	wfClient, err := client.NewWorkflowClient()
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err = wfClient.StartWorker(ctx, r); err != nil {
		logger.Fatal(err)
	}

	fmt.Println(banner("WORKFLOW HISTORY PROPAGATION DEMO — PATIENT INTAKE"))
	fmt.Println()
	fmt.Println("  Flow: PatientIntake -> VerifyInsurance")
	fmt.Println("           -> PrescribeMedication (child wf, lineage)")
	fmt.Println("               -> CheckAllergies -> ScreenDrugInteractions")
	fmt.Println("               -> ComplianceAudit (child wf, lineage)     <-- sees PatientIntake + PrescribeMedication events")
	fmt.Println("               -> DispenseMedication (activity, own only) <-- sees only PrescribeMedication events")
	fmt.Println()

	id, err := wfClient.ScheduleWorkflow(ctx, "PatientIntake",
		workflow.WithInstanceID("intake-001"),
		workflow.WithInput(PatientRecord{
			PatientID:  "P-1042",
			Name:       "Jane Doe",
			DOB:        "1985-06-12",
			MRN:        "MRN-77231",
			Condition:  "bacterial sinusitis",
			Medication: "amoxicillin",
			Dosage:     500,
		}),
	)
	if err != nil {
		logger.Fatalf("failed to start workflow: %v", err)
	}
	fmt.Printf("  [main] Started workflow: %s\n", id)

	waitCtx, waitCancel := context.WithTimeout(ctx, 30*time.Second)
	_, err = wfClient.WaitForWorkflowCompletion(waitCtx, id)
	waitCancel()
	if err != nil {
		logger.Fatalf("workflow failed: %v", err)
	}

	if err = wfClient.PurgeWorkflowState(ctx, id); err != nil {
		logger.Printf("failed to purge: %v", err)
	}

	fmt.Println()
	fmt.Println(banner("COMPLETE"))

	// The workflow has completed and its state was purged above, so return.
	// The deferred cancel() stops the worker and lets `dapr run` exit on its
	// own — no Ctrl+C needed.
}

func banner(msg string) string {
	line := strings.Repeat("=", len(msg)+4)
	return fmt.Sprintf("%s\n= %s =\n%s", line, msg, line)
}
