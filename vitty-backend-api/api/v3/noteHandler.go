package v3

import (
	"strings"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/middleware"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/serializers"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func noteHandler(api fiber.Router) {
	group := api.Group("/notes")
	group.Use(middleware.JWTAuthMiddleware)
	group.Get("/:courseId", getNotes)
	group.Post("/save", saveNote)
	group.Delete("/:noteId?", deleteNote)
}

func getNotes(c *fiber.Ctx) error {
	var notes models.Notes
	courseId := c.Params("courseId")
	notes.CourseID = courseId
	notes.UserName = c.Locals("user").(models.User).Username
	err, userNotes := notes.GetNotesByCourseId()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Notes fetch failed",
		})
	}

	if len(userNotes) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": "No notes to display",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": serializers.NotesSerializer(userNotes),
	})
}

func saveNote(c *fiber.Ctx) error {
	var note models.Notes

	if err := c.BodyParser(&note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if !strings.Contains(note.NoteID, "note_") || len(note.NoteID) < 32 {
		note.NoteID = utils.UUIDWithPrefix("note")
	}

	note.SaveNote()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Note Saved Successfully",
	})
}

func deleteNote(c *fiber.Ctx) error {
	var note models.Notes

	noteId := c.Params("noteId")
	if noteId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "Note id is missing",
		})
	}

	note.NoteID = noteId
	note.UserName = c.Locals("user").(models.User).Username

	err := note.DeleteNote()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Note deleted successfully",
	})
}
