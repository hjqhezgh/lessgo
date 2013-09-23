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

func QueryMenus(w http.ResponseWriter, r *http.Request) {

	var menus []Menu
	data := make(map[string]interface {})

	username := r.FormValue("username")
	if username == "" {
		Log.Error("username is NULL!")
		return
	}

	db := GetMySQL()
	defer db.Close()

	sql := `select a.action_id,a.action_name,a.icon,a.url from action a where
				a.action_id in (select distinct(ra.action_id) from role_action ra where
					ra.role_id in (select ur.role_id from user_role ur where
						ur.user_id=(select u.user_id from users u where u.username=?)))`

	sql2 := "select a.action_id,a.action_name,a.icon,a.url from action a where parent_id=?"

	rows, err := db.Query(sql, username)
	if err != nil {
		Log.Error(err.Error())
		return
	}

	for rows.Next() {
		menu := new(Menu)
		err := rows.Scan(&menu.Id, &menu.Name, &menu.Icon, &menu.Url)
		if err != nil {
			Log.Error(err.Error())
			return
		}
		rows2, err := db.Query(sql2, menu.Id)
		if err != nil {
			Log.Error(err.Error())
			return
		}
		for rows2.Next() {
			child := new(Menu)
			err := rows2.Scan(&child.Id, &child.Name, &child.Icon, &child.Url)
			if err != nil {
				Log.Error(err.Error())
				return
			}
			menu.Children = append(menu.Children, *child)
		}
//		if len(menu.Children) == 0 {
//			menu.Children=[]
//		}
		Log.Warn(*menu)
		menus = append(menus, *menu)
	}
	data["menus"] = menus
	Log.Warn(data)
	outputJson(w, menus)
}

func outputJson(w http.ResponseWriter, object interface {}) {
	b, err := json.Marshal(object)
	if err != nil {
		Log.Error("error!", err.Error())
		return
	}
	w.Write(b)
}

