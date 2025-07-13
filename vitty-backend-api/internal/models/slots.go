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

var BhopalTimetableSlots = []string{"A11", "B11", "C11", "A21", "A14", "B21", "C21", "D11", "E11", "F11", "D21", "E14", "E21", "F21",
	"A12", "B12", "C12", "A22", "B14", "B22", "A24", "D12", "E12", "F12", "D22", "F14", "E22", "F22",
	"A13", "B13", "C13", "A23", "C14", "B23", "B24", "D13", "E13", "F13", "D23", "D14", "D24", "E23"}

var BhopalDailySlots = map[string]map[string][]string{
	"Monday": {
		"Theory": {"A11", "B11", "C11", "A21", "A14", "B21", "C21"},
		"Lab":    {},
	},
	"Tuesday": {
		"Theory": {"D11", "E11", "F11", "D21", "E14", "E21", "F21"},
		"Lab":    {},
	},
	"Wednesday": {
		"Theory": {"A12", "B12", "C12", "A22", "B14", "B22", "A24"},
		"Lab":    {},
	},
	"Thursday": {
		"Theory": {"D12", "E12", "F12", "D22", "F14", "E22", "F22"},
		"Lab":    {},
	},
	"Friday": {
		"Theory": {"A13", "B13", "C13", "A23", "C14", "B23", "B24"},
		"Lab":    {},
	},
	"Saturday": {
		"Theory": {"D13", "E13", "F13", "D23", "D14", "D24", "E23"},
		"Lab":    {},
	},
	"Sunday": {
		"Theory": {},
		"Lab":    {},
	},
}

// Bhopal-specific timings (7 slots per day)
var BhopalTheoryTimings = []Timings{
	{StartTime: "0000-01-01T08:30", EndTime: "0000-01-01T10:00"},
	{StartTime: "0000-01-01T10:05", EndTime: "0000-01-01T11:35"},
	{StartTime: "0000-01-01T11:40", EndTime: "0000-01-01T13:10"},
	{StartTime: "0000-01-01T13:15", EndTime: "0000-01-01T14:45"},
	{StartTime: "0000-01-01T14:50", EndTime: "0000-01-01T16:20"},
	{StartTime: "0000-01-01T16:25", EndTime: "0000-01-01T17:55"},
	{StartTime: "0000-01-01T18:00", EndTime: "0000-01-01T19:30"},
}

var BhopalLabTimings = []Timings{}

// Bhopal slot to timing mapping
var BhopalSlotToTimingMap = map[string]int{
	// Slot 0: 08:30-10:00
	"A11": 0, "D11": 0, "A12": 0, "D12": 0, "A13": 0, "D13": 0,
	// Slot 1: 10:05-11:35
	"B11": 1, "E11": 1, "B12": 1, "E12": 1, "B13": 1, "E13": 1,
	// Slot 2: 11:40-13:10
	"C11": 2, "F11": 2, "C12": 2, "F12": 2, "C13": 2, "F13": 2,
	// Slot 3: 13:15-14:45
	"A21": 3, "D21": 3, "A22": 3, "D22": 3, "A23": 3, "D23": 3,
	// Slot 4: 14:50-16:20
	"A14": 4, "E14": 4, "B14": 4, "F14": 4, "C14": 4, "D14": 4,
	// Slot 5: 16:25-17:55
	"B21": 5, "E21": 5, "B22": 5, "E22": 5, "B23": 5, "D24": 5,
	// Slot 6: 18:00-19:30
	"C21": 6, "F21": 6, "A24": 6, "F22": 6, "B24": 6, "E23": 6,
}

func GetDailySlotsForCampus(campus string) map[string]map[string][]string {
	if campus == "bhopal" {
		return BhopalDailySlots
	}
	return DailySlots
}

func GetTimetableSlotsForCampus(campus string) []string {
	if campus == "bhopal" {
		return BhopalTimetableSlots
	}
	return TimetableSlots
}

func GetTheoryTimingsForCampus(campus string) []Timings {
	if campus == "bhopal" {
		return BhopalTheoryTimings
	}
	return TheoryTimings
}

func GetLabTimingsForCampus(campus string) []Timings {
	if campus == "bhopal" {
		return BhopalLabTimings
	}
	return LabTimings
}

func GetBhopalSlotTimingIndex(slot string) (int, bool) {
	index, exists := BhopalSlotToTimingMap[slot]
	return index, exists
}
