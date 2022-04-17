package controllers

import (
	"log"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/wl_apiagen/entities"
	"github.com/nikitamirzani323/wl_apiagen/helpers"
	"github.com/nikitamirzani323/wl_apiagen/models"
)

const Fieldadminrule_home_redis = "LISTADMINRULE_MASTER_WL"

func Adminrulehome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_adminrule)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, client_idcompany, _, _ := helpers.Parsing_Decry(temp_decp, "==")
	var obj entities.Model_adminruleall
	var arraobj []entities.Model_adminruleall
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldadminrule_home_redis + "_" + strings.ToLower(client_idcompany))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		adminrule_idrule, _ := jsonparser.GetInt(value, "adminrule_idrule")
		adminrule_nmrule, _ := jsonparser.GetString(value, "adminrule_nmrule")
		adminrule_rule, _ := jsonparser.GetString(value, "adminrule_rule")
		adminrule_create, _ := jsonparser.GetString(value, "adminrule_create")
		adminrule_update, _ := jsonparser.GetString(value, "adminrule_update")

		obj.Adminrule_idrule = int(adminrule_idrule)
		obj.Adminrule_nmrule = adminrule_nmrule
		obj.Adminrule_rule = adminrule_rule
		obj.Adminrule_create = adminrule_create
		obj.Adminrule_update = adminrule_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_adminruleHome(client_idcompany)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldadminrule_home_redis+"_"+strings.ToLower(client_idcompany), result, 60*time.Minute)
		log.Println("ADMIN RULE MYSQL")
		return c.JSON(result)
	} else {
		log.Println("ADMIN RULE CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func AdminruleSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_adminrulesave)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	client_admin, client_idcompany, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	//admin, idcompany, name, rule, sData string, idrule int
	result, err := models.Save_adminrule(client_admin, client_idcompany, client.Nmrule, client.Rule, client.Sdata, client.Idrule)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_adminrule(client_idcompany)
	return c.JSON(result)
}
func _deleteredis_adminrule(idcompany string) {
	val_super := helpers.DeleteRedis(Fieldadminrule_home_redis + "_" + strings.ToLower(idcompany))
	log.Printf("REDIS DELETE MASTER ADMIN RULE : %d", val_super)

	val_superlog := helpers.DeleteRedis(Fieldlog_home_redis + "_" + strings.ToLower(idcompany))
	log.Printf("REDIS DELETE MASTER LOG : %d", val_superlog)
}
