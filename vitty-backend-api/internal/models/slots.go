package models

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
}

type Timings struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

var TheoryTimings = []Timings{
	{StartTime: "08:00", EndTime: "08:50"},
	{StartTime: "09:00", EndTime: "09:50"},
	{StartTime: "10:00", EndTime: "10:50"},
	{StartTime: "11:00", EndTime: "11:50"},
	{StartTime: "12:00", EndTime: "12:50"},
	{StartTime: "14:00", EndTime: "14:50"},
	{StartTime: "15:00", EndTime: "15:50"},
	{StartTime: "16:00", EndTime: "16:50"},
	{StartTime: "17:00", EndTime: "17:50"},
	{StartTime: "18:00", EndTime: "18:50"},
	{StartTime: "19:00", EndTime: "19:50"},
}

var LabTimings = []Timings{
	{StartTime: "08:00", EndTime: "08:50"},
	{StartTime: "08:51", EndTime: "09:40"},
	{StartTime: "09:51", EndTime: "10:40"},
	{StartTime: "10:41", EndTime: "11:30"},
	{StartTime: "11:40", EndTime: "12:30"},
	{StartTime: "12:31", EndTime: "13:20"},
	{StartTime: "14:00", EndTime: "14:50"},
	{StartTime: "14:51", EndTime: "15:40"},
	{StartTime: "15:51", EndTime: "16:40"},
	{StartTime: "16:41", EndTime: "17:30"},
	{StartTime: "17:40", EndTime: "18:30"},
	{StartTime: "18:31", EndTime: "19:20"},
}
