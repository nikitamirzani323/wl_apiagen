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

const Fieldadmin_home_redis = "LISTADMIN_MASTER_WL"

func Adminhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_admin)
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

	var obj entities.Model_admin
	var arraobj []entities.Model_admin
	var obj_listruleadmin entities.Model_adminrule
	var arraobj_listruleadmin []entities.Model_adminrule
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldadmin_home_redis + "_" + strings.ToLower(client_idcompany))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	listruleadmin_RD, _, _, _ := jsonparser.Get(jsonredis, "listruleadmin")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		admin_username, _ := jsonparser.GetString(value, "admin_username")
		admin_type, _ := jsonparser.GetString(value, "admin_type")
		admin_idrule, _ := jsonparser.GetInt(value, "admin_idrule")
		admin_rule, _ := jsonparser.GetString(value, "admin_rule")
		admin_nama, _ := jsonparser.GetString(value, "admin_nama")
		admin_phone, _ := jsonparser.GetString(value, "admin_phone")
		admin_email, _ := jsonparser.GetString(value, "admin_email")
		admin_joindate, _ := jsonparser.GetString(value, "admin_joindate")
		admin_lastlogin, _ := jsonparser.GetString(value, "admin_lastlogin")
		admin_lastipaddres, _ := jsonparser.GetString(value, "admin_lastipaddres")
		admin_status, _ := jsonparser.GetString(value, "admin_status")
		admin_create, _ := jsonparser.GetString(value, "admin_create")
		admin_update, _ := jsonparser.GetString(value, "admin_update")

		obj.Admin_username = admin_username
		obj.Admin_type = admin_type
		obj.Admin_idrule = int(admin_idrule)
		obj.Admin_rule = admin_rule
		obj.Admin_nama = admin_nama
		obj.Admin_phone = admin_phone
		obj.Admin_email = admin_email
		obj.Admin_joindate = admin_joindate
		obj.Admin_lastlogin = admin_lastlogin
		obj.Admin_lastIpaddress = admin_lastipaddres
		obj.Admin_status = admin_status
		obj.Admin_create = admin_create
		obj.Admin_update = admin_update
		arraobj = append(arraobj, obj)
	})
	jsonparser.ArrayEach(listruleadmin_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		adminrule_idrule, _ := jsonparser.GetInt(value, "adminrule_idrule")
		adminrule_nmrule, _ := jsonparser.GetString(value, "adminrule_nmrule")

		obj_listruleadmin.Admin_idrule = int(adminrule_idrule)
		obj_listruleadmin.Admin_nmrule = adminrule_nmrule
		arraobj_listruleadmin = append(arraobj_listruleadmin, obj_listruleadmin)
	})
	if !flag {
		result, err := models.Fetch_adminHome(client_idcompany)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldadmin_home_redis+"_"+strings.ToLower(client_idcompany), result, 30*time.Minute)
		log.Println("ADMIN MYSQL")
		return c.JSON(result)
	} else {
		log.Println("ADMIN CACHE")
		return c.JSON(fiber.Map{
			"status":        fiber.StatusOK,
			"message":       "Success",
			"record":        arraobj,
			"listruleadmin": arraobj_listruleadmin,
			"time":          time.Since(render_page).String(),
		})
	}
}
func AdminSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_adminsave)
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

	//admin, idcompany, username, password, nama, email, phone, status, sData string, idrule int
	result, err := models.Save_adminHome(
		client_admin, client_idcompany,
		client.Username, client.Password, client.Nama, client.Email, client.Phone,
		client.Status, client.Sdata, client.Idrule)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_admin(client_idcompany)
	return c.JSON(result)
}
func _deleteredis_admin(idcompany string) {
	val_super := helpers.DeleteRedis(Fieldadmin_home_redis + "_" + strings.ToLower(idcompany))
	log.Printf("REDIS DELETE MASTER ADMIN : %d", val_super)

	val_superlog := helpers.DeleteRedis(Fieldlog_home_redis + "_" + strings.ToLower(idcompany))
	log.Printf("REDIS DELETE MASTER LOG : %d", val_superlog)
}
