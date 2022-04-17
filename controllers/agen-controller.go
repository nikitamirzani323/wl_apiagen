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

const Fieldagen_home_redis = "LISTAGEN_MASTER_WL"

func Agenhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_agen)
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

	var obj entities.Model_agen
	var arraobj []entities.Model_agen
	var obj_listcurr entities.Model_agencurr
	var arraobj_listcurr []entities.Model_agencurr
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldagen_home_redis + "_" + strings.ToLower(client_idcompany))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	listcurr_RD, _, _, _ := jsonparser.Get(jsonredis, "listcurrency")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		agen_idagen, _ := jsonparser.GetString(value, "agen_idagen")
		agen_startjoin, _ := jsonparser.GetString(value, "agen_startjoin")
		agen_endjoin, _ := jsonparser.GetString(value, "agen_endjoin")
		agen_idcurr, _ := jsonparser.GetString(value, "agen_idcurr")
		agen_nmagen, _ := jsonparser.GetString(value, "agen_nmagen")
		agen_owneragen, _ := jsonparser.GetString(value, "agen_owneragen")
		agen_ownerphone, _ := jsonparser.GetString(value, "agen_ownerphone")
		agen_owneremail, _ := jsonparser.GetString(value, "agen_owneremail")
		agen_urlendpoint, _ := jsonparser.GetString(value, "agen_urlendpoint")
		agen_status, _ := jsonparser.GetString(value, "agen_status")
		agen_create, _ := jsonparser.GetString(value, "agen_create")
		agen_update, _ := jsonparser.GetString(value, "agen_update")

		obj.Agen_idagen = agen_idagen
		obj.Agen_startjoin = agen_startjoin
		obj.Agen_endjoin = agen_endjoin
		obj.Agen_idcurr = agen_idcurr
		obj.Agen_nmagen = agen_nmagen
		obj.Agen_ownername = agen_owneragen
		obj.Agen_ownerphone = agen_ownerphone
		obj.Agen_owneremail = agen_owneremail
		obj.Agen_urlendpoint = agen_urlendpoint
		obj.Agen_status = agen_status
		obj.Agen_create = agen_create
		obj.Agen_update = agen_update
		arraobj = append(arraobj, obj)
	})
	jsonparser.ArrayEach(listcurr_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		curr_idcurr, _ := jsonparser.GetString(value, "curr_idcurr")

		obj_listcurr.Curr_idcurr = curr_idcurr
		arraobj_listcurr = append(arraobj_listcurr, obj_listcurr)
	})
	if !flag {
		result, err := models.Fetch_agenHome(client_idcompany)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldagen_home_redis+"_"+strings.ToLower(client_idcompany), result, 30*time.Hour)
		log.Println("AGEN MYSQL")
		return c.JSON(result)
	} else {
		log.Println("AGEN CACHE")
		return c.JSON(fiber.Map{
			"status":       fiber.StatusOK,
			"message":      "Success",
			"record":       arraobj,
			"listcurrency": arraobj_listcurr,
			"time":         time.Since(render_page).String(),
		})
	}
}
func AgenSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_agensave)
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

	result, err := models.Save_agenHome(
		client_admin, client_idcompany,
		client.Agen_idagen, client.Agen_idcurr, client.Agen_nmagen,
		client.Agen_ownername, client.Agen_ownerphone, client.Agen_owneremail, client.Agen_urlendpoint,
		client.Agen_status, client.Sdata)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_agen(client_idcompany, "")
	return c.JSON(result)
}
func AgenListadmin(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_agenadmin)
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
	var obj entities.Model_agenadmin
	var arraobj []entities.Model_agenadmin
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldagen_home_redis + "_LISTADMIN_" + strings.ToLower(client_idcompany) + "_" + strings.ToLower(client.Idagen))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		agenadmin_username, _ := jsonparser.GetString(value, "agenadmin_username")
		agenadmin_type, _ := jsonparser.GetString(value, "agenadmin_type")
		agenadmin_name, _ := jsonparser.GetString(value, "agenadmin_name")
		agenadmin_phone, _ := jsonparser.GetString(value, "agenadmin_phone")
		agenadmin_email, _ := jsonparser.GetString(value, "agenadmin_email")
		agenadmin_status, _ := jsonparser.GetString(value, "agenadmin_status")
		agenadmin_lastlogin, _ := jsonparser.GetString(value, "agenadmin_lastlogin")
		agenadmin_lastipaddress, _ := jsonparser.GetString(value, "agenadmin_lastipaddress")
		agenadmin_create, _ := jsonparser.GetString(value, "agenadmin_create")
		agenadmin_update, _ := jsonparser.GetString(value, "agenadmin_update")

		obj.Agenadmin_username = agenadmin_username
		obj.Agenadmin_type = agenadmin_type
		obj.Agenadmin_name = agenadmin_name
		obj.Agenadmin_phone = agenadmin_phone
		obj.Agenadmin_email = agenadmin_email
		obj.Agenadmin_status = agenadmin_status
		obj.Agenadmin_lastlogin = agenadmin_lastlogin
		obj.Agenadmin_lastipaddress = agenadmin_lastipaddress
		obj.Agenadmin_create = agenadmin_create
		obj.Agenadmin_update = agenadmin_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_agenListAdmin(client_idcompany, client.Idagen)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldagen_home_redis+"_LISTADMIN_"+strings.ToLower(client_idcompany)+"_"+strings.ToLower(client.Idagen), result, 30*time.Hour)
		log.Println("COMPANY LISTADMIN AGEN MYSQL")
		return c.JSON(result)
	} else {
		log.Println("COMPANY LISTADMIN AGEN CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func AgenSavelistadmin(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_agenadminsave)
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

	result, err := models.Save_agenlistadmin(
		client_admin, client_idcompany,
		client.Agen_idagen, client.Agenadmin_username, client.Agenadmin_password,
		client.Agenadmin_name, client.Agenadmin_email, client.Agenadmin_phone, client.Agenadmin_status,
		client.Sdata)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_agen(client_idcompany, client.Agen_idagen)
	return c.JSON(result)
}
func _deleteredis_agen(idcompany, idagen string) {
	val_master := helpers.DeleteRedis(Fieldagen_home_redis + "_" + strings.ToLower(idcompany))
	log.Printf("REDIS DELETE MASTER AGEN : %d", val_master)

	val_superlog := helpers.DeleteRedis(Fieldlog_home_redis + "_" + strings.ToLower(idcompany))
	log.Printf("REDIS DELETE MASTER LOG : %d", val_superlog)

	if idagen != "" {
		val_agenadmin := helpers.DeleteRedis(Fieldagen_home_redis + "_LISTADMIN_" + strings.ToLower(idcompany) + "_" + strings.ToLower(idagen))
		log.Printf("REDIS DELETE MASTER AGEN LISTADMIN : %d", val_agenadmin)
	}

}
