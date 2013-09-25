/**
 * FileDes: 
 * User: Samurai
 * Date: 13-9-23
 * Time: 上午10:06
 */
package lessgo

import(
	"net/http"
	"encoding/json"
	"fmt"
)

type Menu struct {
	Id			string			`json:"id"`
	Name		string		`json:"name"`
	Icon		string		`json:"icon"`
	Url			string	    `json:"url"`
	Children 	[]Menu		`json:"children"`
}


func GetMenus(username string) []Menu{

	var menus []Menu
	var parent_id string
	db := GetMySQL()
	defer db.Close()
	fmt.Println("uname:", username)
	sql := `select a.action_id,a.action_name,a.icon,a.parent_id,a.url from action a where
				a.action_id in (select distinct(ra.action_id) from role_action ra where
					ra.role_id in (select er.role_id from employee_role er where
						er.user_id=(select e.user_id from employee e where e.username=?))) order by parent_id`

	rows, err := db.Query(sql, username)
	if err != nil {
		Log.Error(err.Error())
		return nil
	}
	for rows.Next() {
		menu := new(Menu)
		err := rows.Scan(&menu.Id, &menu.Name, &menu.Icon, &parent_id, &menu.Url)
		if err != nil {
			Log.Error(err.Error())
			return nil
		}

		if parent_id > "0" {
			for i := 0; i < len(menus); i++ {
				if menus[i].Id == parent_id {
					menus[i].Children = append(menus[i].Children, *menu)
				}
			}
		}else {
			menus = append(menus, *menu)
		}
	}
	fmt.Println("menus: ", menus)
	return menus
}

func QueryMenusAction(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface {})

	username := r.FormValue("username")
	if username == "" {
		Log.Error("username is NULL!")
		return
	}

	data["menus"] = GetMenus(username)

	Log.Warn(data)
	fmt.Println(data)
	outputJson(w, data)
}

func outputJson(w http.ResponseWriter, object interface {}) {
	b, err := json.Marshal(object)
	if err != nil {
		Log.Error("error!", err.Error())
		return
	}
	w.Write(b)
}

