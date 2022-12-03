package analyzer

const QUERY_PERPAGE = 30

/* 取得已授權的youtuber列表
 */
func GetAllYtbers(offset int, limit int) ([]*Youtuber, error) {

	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var ytbers []*Youtuber
	query := "SELECT * FROM `youtuber` LIMIT :limit OFFSET :offset"
	rows, err := db.NamedQuery(query, map[string]interface{}{
		"offset": offset,
		"limit":  limit,
	})
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ytb Youtuber
		err := rows.StructScan(&ytb)
		if err != nil {
			return nil, err
		}

		ytbers = append(ytbers, &ytb)
	}

	return ytbers, nil
}
