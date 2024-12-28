package models

var TimetableSlots = []string{"A1", "A2", "B1", "B2", "C1", "C2", "D1", "D2", "E1", "E2", "F1", "F2", "G1", "G2",
	"TA1", "TA2", "TAA1", "TAA2", "TB1", "TB2", "TBB2", "TC1", "TC2", "TCC1", "TCC2", "TD1", "TD2", "TDD2", "TE1", "TE2", "TF1", "TF2",
	"TG1", "TG2", "V1", "V2", "V3", "V4", "V5", "V6", "V7", "V8", "V9", "V10", "V11", "W21", "W22", "X11", "X12", "X21", "Y11", "Y12", "Y21", "Z21"}

var DailySlots = map[string]map[string][]string{
	"Monday": {
		"Theory": {"A1", "F1", "D1", "TB1", "TG1", "A2", "F2", "D2", "TB2", "TG2", "V3"},
		"Lab":    {"L1", "L2", "L3", "L4", "L5", "L6", "L31", "L32", "L33", "L34", "L35", "L36"},
	},
	"Tuesday": {
		"Theory": {"B1", "G1", "E1", "TC1", "TAA1", "B2", "G2", "E2", "TC2", "TAA2", "V4"},
		"Lab":    {"L7", "L8", "L9", "L10", "L11", "L12", "L37", "L38", "L39", "L40", "L41", "L42"},
	},
	"Wednesday": {
		"Theory": {"C1", "A1", "F1", "V1", "V2", "C2", "A2", "F2", "TD2", "TBB2", "V5"},
		"Lab":    {"L13", "L14", "L15", "L16", "L17", "L18", "L43", "L44", "L45", "L46", "L47", "L48"},
	},
	"Thursday": {
		"Theory": {"D1", "B1", "G1", "TE1", "TCC1", "D2", "B2", "G2", "TE2", "TCC2", "V6"},
		"Lab":    {"L19", "L20", "L21", "L22", "L23", "L24", "L49", "L50", "L51", "L52", "L53", "L54"},
	},
	"Friday": {
		"Theory": {"E1", "C1", "TA1", "TF1", "TD1", "E2", "C2", "TA2", "TF2", "TDD2", "V7"},
		"Lab":    {"L25", "L26", "L27", "L28", "L29", "L30", "L55", "L56", "L57", "L58", "L59", "L60"},
	},
	"Saturday": {
		"Theory": {"V8", "X11", "X12", "Y11", "Y12", "X21", "Z21", "Y21", "W21", "W22", "V9"},
		"Lab":    {"L71", "L72", "L73", "L74", "L75", "L76", "L77", "L78", "L79", "L80", "L81", "L82"},
	},
	"Sunday": {
		"Theory": {"V10", "Y11", "Y12", "X11", "X12", "Y21", "Z21", "X21", "W21", "W22", "V11"},
		"Lab":    {"L83", "L84", "L85", "L86", "L87", "L88", "L89", "L90", "L91", "L92", "L93", "L94"},
	},
}

type Timings struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

var TheoryTimings = []Timings{
	{StartTime: "0000-01-01T09:00", EndTime: "0000-01-01T09:50"},
	{StartTime: "0000-01-01T10:00", EndTime: "0000-01-01T10:50"},
	{StartTime: "0000-01-01T11:00", EndTime: "0000-01-01T11:50"},
	{StartTime: "0000-01-01T12:00", EndTime: "0000-01-01T12:50"},
	{StartTime: "0000-01-01T08:00", EndTime: "0000-01-01T08:50"},
	{StartTime: "0000-01-01T14:00", EndTime: "0000-01-01T14:50"},
	{StartTime: "0000-01-01T15:00", EndTime: "0000-01-01T15:50"},
	{StartTime: "0000-01-01T16:00", EndTime: "0000-01-01T16:50"},
	{StartTime: "0000-01-01T17:00", EndTime: "0000-01-01T17:50"},
	{StartTime: "0000-01-01T18:00", EndTime: "0000-01-01T18:50"},
	{StartTime: "0000-01-01T19:00", EndTime: "0000-01-01T19:50"},
}

var LabTimings = []Timings{
	{StartTime: "0000-01-01T08:00", EndTime: "0000-01-01T08:50"},
	{StartTime: "0000-01-01T08:51", EndTime: "0000-01-01T09:40"},
	{StartTime: "0000-01-01T09:51", EndTime: "0000-01-01T10:40"},
	{StartTime: "0000-01-01T10:41", EndTime: "0000-01-01T11:30"},
	{StartTime: "0000-01-01T11:40", EndTime: "0000-01-01T12:30"},
	{StartTime: "0000-01-01T12:31", EndTime: "0000-01-01T13:20"},
	{StartTime: "0000-01-01T14:00", EndTime: "0000-01-01T14:50"},
	{StartTime: "0000-01-01T14:51", EndTime: "0000-01-01T15:40"},
	{StartTime: "0000-01-01T15:51", EndTime: "0000-01-01T16:40"},
	{StartTime: "0000-01-01T16:41", EndTime: "0000-01-01T17:30"},
	{StartTime: "0000-01-01T17:40", EndTime: "0000-01-01T18:30"},
	{StartTime: "0000-01-01T18:31", EndTime: "0000-01-01T19:20"},
}
