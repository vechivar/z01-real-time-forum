package rtfServer

func GetUsernameFromId(user_id int) (string, error) {
	var username string

	err := db.QueryRow("SELECT username FROM user WHERE user_id = ?", user_id).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func GetCategoryFromId(category_id int) (string, error) {
	var category string

	err := db.QueryRow("SELECT name FROM category WHERE category_id = ?", category_id).Scan(&category)
	if err != nil {
		return "", err
	}

	return category, nil
}
