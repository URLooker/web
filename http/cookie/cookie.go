package cookie

import (
	"github.com/gorilla/securecookie"
	"github.com/peng19940915/urlooker/web/g"
	"github.com/gin-gonic/gin"
	"fmt"
)

var SecureCookie *securecookie.SecureCookie

func Init() {
	var hashKey = []byte(g.Config.Http.Secret)
	var blockKey = []byte(nil)
	SecureCookie = securecookie.New(hashKey, blockKey)
}

func ReadUser(c *gin.Context) (id int64, name string, cnName string, found bool) {
	if cookieValue, err := c.Cookie("u"); err == nil {
		value := make(map[string]interface{})
		if err = SecureCookie.Decode("u", cookieValue, &value); err == nil {
			id = value["id"].(int64)
			name = value["name"].(string)
			if value["cnName"] != nil {
				cnName = value["cnName"].(string)
			}else {
				cnName = ""
			}


			if id == 0 || name == "" {
				return
			} else {
				found = true
				return
			}
		}
	}else {

			fmt.Println("Not Find User Info From Cookie")
	}
	return
}

func WriteUser(c *gin.Context, id int64, name string, cnName string) error {
	value := make(map[string]interface{})
	value["id"] = id
	value["name"] = name
	value["cnName"] = cnName
	encoded, err := SecureCookie.Encode("u", value)
	if err != nil {
		return err
	}

	c.SetCookie("u", encoded, 3600 * 24 * 7, "/","", false, true)

	return nil
}

func RemoveUser(c *gin.Context) error {
	value := make(map[string]interface{})
	value["id"] = ""
	value["name"] = ""
	encoded, err := SecureCookie.Encode("u", value)
	if err != nil {
		return err
	}
	c.SetCookie("u", encoded, -1, "/","", false, true)

	return nil
}
