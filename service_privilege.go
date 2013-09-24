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
)

type Menu struct {
	Id			int			`json:"id"`
	Name		string		`json:"name"`
	Icon		string		`json:"icon"`
	Url			string	    `json:"url"`
	Children 	[]Menu		`json:"children"`
}

func queryMenus(username string, menus *[]Menu) bool{

	db := GetMySQL()
	defer db.Close()

	sql := `select a.action_id,a.action_name,a.icon,a.url from action a where
				a.action_id in (select distinct(ra.action_id) from role_action ra where
					ra.role_id in (select er.role_id from employee_role er where
						er.user_id=(select e.user_id from employee e where e.username=?)))`

	sql2 := "select a.action_id,a.action_name,a.icon,a.url from action a where parent_id=?"

	rows, err := db.Query(sql, username)
	if err != nil {
		Log.Error(err.Error())
		return false
	}

	for rows.Next() {
		menu := new(Menu)
		err := rows.Scan(&menu.Id, &menu.Name, &menu.Icon, &menu.Url)
		if err != nil {
			Log.Error(err.Error())
			return false
		}
		rows2, err := db.Query(sql2, menu.Id)
		if err != nil {
			Log.Error(err.Error())
			return false
		}
		for rows2.Next() {
			child := new(Menu)
			err := rows2.Scan(&child.Id, &child.Name, &child.Icon, &child.Url)
			if err != nil {
				Log.Error(err.Error())
				return false
			}
			menu.Children = append(menu.Children, *child)
		}
		*menus = append(*menus, *menu)
	}
	return true
}

func GetMenus(username string) []Menu{

	var menus []Menu
	db := GetMySQL()
	defer db.Close()

	sql := `select a.action_id,a.action_name,a.icon,a.url from action a where
				a.action_id in (select distinct(ra.action_id) from role_action ra where
					ra.role_id in (select er.role_id from employee_role er where
						er.user_id=(select e.user_id from employee e where e.username=?)))`

	sql2 := "select a.action_id,a.action_name,a.icon,a.url from action a where parent_id=?"

	rows, err := db.Query(sql, username)
	if err != nil {
		Log.Error(err.Error())
		return nil
	}

	for rows.Next() {
		menu := new(Menu)
		err := rows.Scan(&menu.Id, &menu.Name, &menu.Icon, &menu.Url)
		if err != nil {
			Log.Error(err.Error())
			return nil
		}
		rows2, err := db.Query(sql2, menu.Id)
		if err != nil {
			Log.Error(err.Error())
			return nil
		}
		for rows2.Next() {
			child := new(Menu)
			err := rows2.Scan(&child.Id, &child.Name, &child.Icon, &child.Url)
			if err != nil {
				Log.Error(err.Error())
				return nil
			}
			menu.Children = append(menu.Children, *child)
		}
		menus = append(menus, *menu)
	}
	return menus
}

func QueryMenusAction(w http.ResponseWriter, r *http.Request) {

	var menus []Menu
	data := make(map[string]interface {})

	username := r.FormValue("username")
	if username == "" {
		Log.Error("username is NULL!")
		return
	}
	ret := queryMenus(username, &menus)
	if ret {
		data["menus"] = menus
	}
	Log.Warn(data)
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

