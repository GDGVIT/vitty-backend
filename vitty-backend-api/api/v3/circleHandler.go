package v3

import (
	"errors"
	"log"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/middleware"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/serializers"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func circleHandler(api fiber.Router) {
	group := api.Group("/circles")
	group.Use(middleware.JWTAuthMiddleware)
	group.Get("/", getCircles)
	group.Get("/:circleId", getUsersOfCircle)
	group.Get("/leisure/:circleId", getLeisureTime)
	group.Get("/requests/received", getReceivedCircleRequests)
	group.Get("/requests/sent", getSentCircleRequests)
	group.Post("/create/:circleName", createCircle)
	group.Post("/sendRequest/:circleId/:username", sendCircleRequestToUser)
	group.Post("/:circleId/generateJoinCode", generateCircleJoinCode)
	group.Post("/acceptRequest/:circleId", acceptCircleRequest)
	group.Post("/join", joinCircleByCode)
	group.Post("/declineRequest/:circleId", declineCircleRequest)
	group.Patch("/", updateCircleName)
	group.Get("/:circleId/:username", getCircleMemberTimetable)
	group.Delete("/:circleId", deleteCircle)
	group.Delete("/remove/:circleId/:username", removeUserFromCircle)
	group.Delete("/leave/:circleId", leaveCircle)
	group.Delete("/unsendRequest/:circleId/:username", unsendCircleRequestToUser)
}

func getCircles(c *fiber.Ctx) error {
	var user_circle models.UsersCirclesJoin

	username := c.Locals("user").(models.User).Username

	user_circle.Uname = username

	err, circles := user_circle.GetCirclesofUser()

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "circles not fetched",
		})
	}

	if len(circles) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": "You are not part of any circle",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": serializers.CirclesListSerializer(circles),
	})
}

func getUsersOfCircle(c *fiber.Ctx) error {
	var userCircle models.UsersCirclesJoin

	circleId := c.Params("circleId")
	username := c.Locals("user").(models.User).Username

	userCircle.CID = circleId
	userCircle.Uname = username
	err, isPartOfCircle := userCircle.IsUserOfCircle()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid circle id",
			})
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user authorization could not be verified",
		})
	}

	if isPartOfCircle {
		userCircle.Uname = ""
		err, users := userCircle.GetUsersofCircle()

		if err == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"data": serializers.UsersListCircleSerializer(users),
			})
		}
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "users list fetch failed",
	})
}

func getCircleMemberTimetable(c *fiber.Ctx) error {
	var userCircle models.UsersCirclesJoin

	circleId := c.Params("circleId")
	username := c.Params("username")
	requestUser := c.Locals("user").(models.User).Username

	userCircle.CID = circleId
	userCircle.Uname = requestUser
	err, isPartOfCircle := userCircle.IsUserOfCircle()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid circle id",
			})
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user authorization could not be verified",
		})
	}

	if !isPartOfCircle {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you are not part of this circle",
		})
	}

	if !utils.CheckUserExists(username) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"detail": "User not found",
		})
	}

	userCircle.Uname = username
	err, isTargetUserInCircle := userCircle.IsUserOfCircle()
	if err != nil || !isTargetUserInCircle {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user is not part of this circle",
		})
	}

	user := utils.GetUserByUsername(username)
	campus := "vellore"
	if user.Campus != nil {
		campus = string(*user.Campus)
	}

	return c.Status(fiber.StatusOK).JSON(serializers.TimetableSerializer(user.GetTimeTable(), campus))
}

func getLeisureTime(c *fiber.Ctx) error {
	var circle models.Circles

	circleId := c.Params("circleId")

	circle.CircleId = circleId

	err, circleSlotMap := circle.GetCircleSlots()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid circle id",
			})
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "leisure times couldn't be fetched",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": circleSlotMap,
	})

}

func getReceivedCircleRequests(c *fiber.Ctx) error {
	var circleRequest models.CircleRequest

	username := c.Locals("user").(models.User).Username

	err, circleRequests := circleRequest.GetReceivedRequests(username)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "requests not able to be fetched",
		})
	}

	if len(circleRequests) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": "your inbox is empty",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": serializers.CircleRequestsSerializer(circleRequests),
	})

}

func getSentCircleRequests(c *fiber.Ctx) error {
	var circleRequest models.CircleRequest

	username := c.Locals("user").(models.User).Username

	err, circleRequests := circleRequest.GetSentRequests(username)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "requests not able to be fetched",
		})
	}

	if len(circleRequests) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": "your outbox is empty",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": serializers.CircleRequestsSerializer(circleRequests),
	})

}

func createCircle(c *fiber.Ctx) error {
	var circle models.Circles
	var user_circle models.UsersCirclesJoin

	circleName := c.Params("circleName")
	username := c.Locals("user").(models.User).Username

	if circleName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	} else {
		circle.CircleId = utils.UUIDWithPrefix("circle")
		circle.CircleName = circleName
		circle.Uname = username
	}

	err := circle.CreateCircle()
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"detail": "circle name already exists",
			})
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "circle creation failed",
		})
	}

	request_user := c.Locals("user").(models.User)

	user_circle.CID = circle.CircleId
	user_circle.CircleRole = "admin"
	user_circle.Uname = request_user.Username

	err = user_circle.AddUserToCircle()

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "user error circle creation failed",
		})
	}

	joinCode := utils.GenerateJoinCode(10)
	err = circle.CreateCircleJoinCode(joinCode)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "circle created but join code generation failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail":    "circle created successfully",
		"circle_id": circle.CircleId,
		"join_code": joinCode,
	})
}

func sendCircleRequestToUser(c *fiber.Ctx) error {
	var circleRequest models.CircleRequest
	var usersCirclesJoin models.UsersCirclesJoin

	circleId := c.Params("circleId")
	from_user := c.Locals("user").(models.User).Username
	to_user := c.Params("username")

	usersCirclesJoin.CID = circleId
	usersCirclesJoin.Uname = from_user

	err, isCircleAdmin := usersCirclesJoin.IsUserCircleAdmin()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "circle does not exist",
			})
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "circle search failed",
		})
	}

	if !isCircleAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"detail": "you cannot send requests",
		})
	}

	circleRequest.CID = circleId
	circleRequest.FromUsername = from_user
	circleRequest.ToUsername = to_user

	err = circleRequest.CreateRequest()
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "request pending",
			})
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "request not sent",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "request sent successfully",
	})
}

func generateCircleJoinCode(c *fiber.Ctx) error {
	var circle models.Circles

	circleId := c.Params("circleId")

	circle.CircleId = circleId

	joinCode := utils.GenerateJoinCode(10)
	err := circle.CreateCircleJoinCode(joinCode)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "circle share code generation failed",
		})
	}

	circle.CircleJoinCode = joinCode

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"circle_id": circle.CircleId,
		"join_code": circle.CircleJoinCode,
		"detail":    "circle code generated",
	})
}

func acceptCircleRequest(c *fiber.Ctx) error {
	var circleRequest models.CircleRequest

	circleId := c.Params("circleId")
	to_user := c.Locals("user").(models.User).Username

	circleRequest.CID = circleId
	circleRequest.ToUsername = to_user

	err := circleRequest.AcceptRequest()

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "you do not have any requests to accept",
			})
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "request not accepted",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "request accepted successfully",
	})
}

func joinCircleByCode(c *fiber.Ctx) error {
	var circle models.Circles
	var userCircle models.UsersCirclesJoin

	joinCode := c.Query("code")
	reqUser := c.Locals("user").(models.User).Username

	err := circle.GetCircleByJoinCode(joinCode)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "invalid join code",
			})
		}
		log.Println(err)
	}

	userCircle.CID = circle.CircleId
	userCircle.CircleRole = "member"
	userCircle.Uname = reqUser

	err = userCircle.AddUserToCircle()

	if err != nil {
		log.Println(err)

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "you are already part of the circle",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to join circle",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "joined circle successfully",
	})
}

func updateCircleName(c *fiber.Ctx) error {
	var circle models.Circles

	err := c.BodyParser(&circle)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}

	err = circle.UpdateCircleName(circle.CircleName)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "circle name not updated",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "circle name updated",
	})
}

func deleteCircle(c *fiber.Ctx) error {
	var circle models.Circles

	circleId := c.Params("circleId")
	circle.CircleId = circleId

	err := circle.DeleteCircle()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "circle was not deleted",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "circle deleted successfully",
	})
}

func removeUserFromCircle(c *fiber.Ctx) error {
	var ucj models.UsersCirclesJoin

	circleId := c.Params("circleId")
	username := c.Params("username")
	requestUser := c.Locals("user").(models.User)

	ucj.CID = circleId
	ucj.Uname = requestUser.Username

	err, isCircleAdmin := ucj.IsUserCircleAdmin()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "circle does not exist",
			})
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "circle search failed",
		})
	}

	if !isCircleAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"detail": "you cannot remove user",
		})
	}

	ucj.Uname = username
	err = ucj.DeleteUserFromCircle()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user not removed from the circle",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "user removed from circle",
	})
}

func leaveCircle(c *fiber.Ctx) error {
	var ucj models.UsersCirclesJoin

	circleId := c.Params("circleId")
	username := c.Locals("user").(models.User).Username

	ucj.CID = circleId
	ucj.Uname = username

	err, isCircleAdmin := ucj.IsUserCircleAdmin()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": "circle does not exist",
			})
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "circle search failed",
		})
	}

	err = ucj.DeleteUserFromCircle()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "you didn't leave the circle",
		})
	}

	if isCircleAdmin {
		ucj.Uname = ""
		err, users := ucj.GetUsersofCircle()
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "you didn't leave the circle;admin error",
			})
		}

		if len(users) != 0 {
			var circle models.Circles
			circle.CircleId = circleId
			err = circle.UpdateCircleUsername(users[0].Username)
			if err != nil {
				log.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "you didn't leave the circle;admin update error",
				})
			}
			ucj.Uname = users[0].Username
			err = ucj.UpdateUserCircleRole("admin")
			if err != nil {
				log.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "you didn't leave the circle;admin update error",
				})
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "you left the circle ",
	})
}

func unsendCircleRequestToUser(c *fiber.Ctx) error {
	var circleRequest models.CircleRequest

	circleId := c.Params("circleId")
	from_user := c.Locals("user").(models.User).Username
	to_user := c.Params("username")

	circleRequest.CID = circleId
	circleRequest.FromUsername = from_user
	circleRequest.ToUsername = to_user

	err := circleRequest.DeleteRequest()

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "request could not be unsent",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "request unsent successfully",
	})
}

func declineCircleRequest(c *fiber.Ctx) error {
	var circleRequest models.CircleRequest

	circleId := c.Params("circleId")
	to_user := c.Locals("user").(models.User).Username

	circleRequest.CID = circleId
	circleRequest.ToUsername = to_user

	err := circleRequest.DeclineRequest()
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "request not declined",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "request declined successfully",
	})
}
