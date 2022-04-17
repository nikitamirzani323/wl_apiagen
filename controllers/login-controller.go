package controllers

import (
	"log"
	"strconv"
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

const Field_login_redis = "LISTLOGINMASTER_MASTER_WL"

func CheckLogin(c *fiber.Ctx) error {
	msg := "Username and Password not register"
	render_page := time.Now()
	var errors []*helpers.ErrorResponse
	client := new(entities.Login)
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
	flag_login := false
	idcompany := ""
	typeadmin := ""
	ruleadmin := 0
	resultredis, flag := helpers.GetRedis(Field_login_redis + "_" + strings.ToLower(idcompany))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		login_username, _ := jsonparser.GetString(value, "login_username")
		login_password, _ := jsonparser.GetString(value, "login_password")
		login_typeadmin, _ := jsonparser.GetString(value, "login_typeadmin")
		Login_idcompany, _ := jsonparser.GetString(value, "login_idcompany")
		login_idrule, _ := jsonparser.GetInt(value, "login_idrule")

		if login_username == client.Username {
			hashpass := helpers.HashPasswordMD5(client.Password)
			if hashpass == login_password {
				flag_login = true
				idcompany = Login_idcompany
				typeadmin = login_typeadmin
				ruleadmin = int(login_idrule)
			}
		}
	})
	if !flag {
		result, flag_model, idcompany_model, typeadmin_model, rule_model, err := models.Login_Model(client.Username, client.Password)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		if flag_model {
			flag_login = true
			idcompany = idcompany_model
			typeadmin = typeadmin_model
			ruleadmin = rule_model
			helpers.SetRedis(Field_login_redis+"_"+strings.ToLower(idcompany), result, 30*time.Hour)
			log.Println("LIST LOGIN ADMIN SUPER MYSQL")

		}

	} else {
		log.Println("LIST LOGIN ADMIN SUPER CACHE")

	}
	temp_token := ""
	if flag_login {
		_deletelogin_admin(idcompany)
		models.Update_login(client.Username, client.Ipaddress, client.Timezone, idcompany)

		dataclient := client.Username + "==" + idcompany + "==" + typeadmin + "==" + strconv.Itoa(ruleadmin)
		dataclient_encr, keymap := helpers.Encryption(dataclient)
		dataclient_encr_final := dataclient_encr + "|" + strconv.Itoa(keymap)
		t, err := helpers.GenerateNewAccessToken(dataclient_encr_final)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		temp_token = t
		msg = ""
	}

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"token":   temp_token,
		"message": msg,
		"time":    time.Since(render_page).String(),
	})
}
func Home(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Home)
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
	client_username, client_idcompany, typeadmin, idruleadmin := helpers.Parsing_Decry(temp_decp, "==")
	log.Printf("USERNAME : %s", client_username)
	log.Printf("COMPANY : %s", client_idcompany)
	log.Printf("TYPE : %s", typeadmin)
	log.Printf("RULE : %s", idruleadmin)
	log.Printf("PAGE : %s", client.Page)
	if typeadmin != "MASTER" {
		idrule, _ := strconv.Atoi(idruleadmin)
		ruleadmin := models.Get_AdminRule("rulecomp", client_idcompany, idrule)
		flag := models.Get_listitemsearch(ruleadmin, ",", client.Page)
		if !flag {
			c.Status(fiber.StatusForbidden)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusForbidden,
				"message": "Anda tidak bisa akses halaman ini",
				"record":  nil,
			})
		} else {
			c.Status(fiber.StatusOK)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusOK,
				"message": "ADMIN",
				"record":  nil,
			})
		}
	} else {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "ADMIN",
			"record":  nil,
		})
	}

}

func _deletelogin_admin(idcompany string) {
	val_super := helpers.DeleteRedis(Fieldadmin_home_redis + "_" + strings.ToLower(idcompany))
	log.Printf("REDIS DELETE MASTER ADMIN : %d", val_super)

	val_superlog := helpers.DeleteRedis(Fieldlog_home_redis + "_" + strings.ToLower(idcompany))
	log.Printf("REDIS DELETE MASTER LOG : %d", val_superlog)
}
