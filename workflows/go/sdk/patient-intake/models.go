package main

// PatientRecord represents a patient intake submitted by the front desk.
// In a real deployment the Name / DOB / MRN fields are protected health info
// and would be candidates for redaction when the record is propagated downstream
// this is a future add for history propagation.
type PatientRecord struct {
	PatientID  string  `json:"patientId"`
	Name       string  `json:"name"`
	DOB        string  `json:"dob"`
	MRN        string  `json:"mrn"`
	Condition  string  `json:"condition"`
	Medication string  `json:"medication"`
	Dosage     float64 `json:"dosage"`
}

// ComplianceResult is the output of the ComplianceAudit child workflow.
type ComplianceResult struct {
	Compliant  bool    `json:"compliant"`
	RiskScore  float64 `json:"riskScore"`
	Reason     string  `json:"reason"`
	EventCount int     `json:"eventCount"`
}

// DispenseResult is the output of the DispenseMedication activity.
type DispenseResult struct {
	DispenseID string `json:"dispenseId"`
	Status     string `json:"status"`
	EventCount int    `json:"eventCount"`
}
