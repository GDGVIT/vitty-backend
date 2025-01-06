package serializers

import "github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"

func NotesSerializer(notes []models.Notes) []map[string]interface{} {
	var result []map[string]interface{}

	for _, note := range notes {
		out := map[string]interface{}{
			"note_id":      note.NoteID,
			"note_name":    note.NoteName,
			"note_content": note.NoteContent,
			"course_id":    note.CourseID,
			"course_name":  note.Courses.CourseName,
		}
		result = append(result, out)
	}
	return result
}
