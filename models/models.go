package Models

import (
	db "database/sql"

	tsgutils "github.com/timespacegroup/go-utils"
)

/*

 */
type User struct {
	Id       int64  `column:"id"`       //
	Name     string `column:"name"`     //
	Password string `column:"password"` //
	State    int64  `column:"state"`    //
	Token    string `column:"token"`    //
	Users    []User // This value is used for batch queries and inserts.
}

/*
func (user *User) RowToStruct(row *db.Row) error {
	builder := tsgutils.NewInterfaceBuilder()
	builder.Append(&user.Id)
	builder.Append(&user.Name)
	builder.Append(&user.Password)
	builder.Append(&user.State)
	builder.Append(&user.Token)
	err := row.Scan(builder.ToInterfaces()...)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) RowsToStruct(rows *db.Rows) error {
	var users []User
	builder := tsgutils.NewInterfaceBuilder()
	for rows.Next() {
		builder.Clear()
		builder.Append(&user.Id)
		builder.Append(&user.Name)
		builder.Append(&user.Password)
		builder.Append(&user.State)
		builder.Append(&user.Token)
		err := rows.Scan(builder.ToInterfaces()...)
		if err != nil {
			return err
		}
		users = append(users, *user)
	}
	if rows != nil {
		defer rows.Close()
	}
	user.Users = users
	return nil
}

func (user *User) Insert(client *DBClient, idSet bool) (int64, error) {
	structParam := *user
	sql := tsgutils.NewStringBuilder()
	qSql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("INSERT INTO ")
	sql.Append("user")
	sql.Append(" (")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(col).Append("`,")
		qSql.Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	qSql.RemoveLast()
	sql.Append(") VALUES (").Append(qSql.ToString()).Append(");")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (user *User) UpdateUserById(client *DBClient) (int64, error) {
	structParam := *user
	sql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("UPDATE ")
	sql.Append("user")
	sql.Append(" SET ")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	var id interface{}
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" {
			id = v
			continue
		}
		sql.Append(col).Append("=").Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	params.Append(id)
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (user *User) DeleteUserById(client *DBClient) (int64, error) {
	structParam := user
	sql := tsgutils.NewStringBuilder()
	sql.Append("DELETE FROM ")
	sql.Append("user")
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), structParam.Id)
}

func (user *User) BatchInsert(client *DBClient, idSet, returnIds bool) ([]int64, error) {
	structParam := *user
	list := structParam.Users
	var result []int64
	listLen := len(list)
	if listLen == 0 {
		return result, errors.New("no data needs to be inserted")
	}
	sql := tsgutils.NewStringBuilder()
	oneQSql := tsgutils.NewStringBuilder()
	batchQSql := tsgutils.NewStringBuilder()
	ks := reflect.TypeOf(structParam)
	fieldsNum := ks.NumField() - 1
	sql.Append("INSERT INTO ")
	sql.Append("user")
	sql.Append(" (")
	for i := 0; i < fieldsNum; i++ {
		iCol := ks.Field(i).Tag.Get("column")
		if iCol == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(iCol).Append("`,")
	}
	sql.RemoveLast().Append(") VALUES ")
	batchInsertColsLen := tsgutils.InterfaceToInt(tsgutils.IIIInterfaceOperator(idSet, fieldsNum, fieldsNum-1))
	oneQSql.Append("(")
	for j := 0; j < batchInsertColsLen; j++ {
		oneQSql.Append("?,")
	}
	oneQSql.RemoveLast().Append(")")
	if !returnIds {
		for j := 0; j < listLen; j++ {
			batchQSql.Append(oneQSql.ToString()).Append(",")
		}
		batchQSql.RemoveLast()
		batchSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(batchQSql.ToString()).Append(";").ToString()
		batchParams := tsgutils.NewInterfaceBuilder()
		for k := range list {
			item := list[k]
			kItem := reflect.ValueOf(item)
			for l := 0; l < fieldsNum; l++ {
				lCol := ks.Field(l).Tag.Get("column")
				if lCol == "id" && !idSet {
					continue
				}
				batchParams.Append(kItem.Field(l).Interface())
			}
		}
		id, err := client.Exec(batchSql, batchParams.ToInterfaces()...)
		if err != nil {
			return result, err
		}
		result = append(result, id)
	} else {
		oneSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(oneQSql.ToString()).Append(";").ToString()
		oneParams := tsgutils.NewInterfaceBuilder()
		tx, err := client.TxBegin()
		if err != nil {
			return result, err
		}
		for m := range list {
			oneParams.Clear()
			item := list[m]
			mItem := reflect.ValueOf(item)
			for n := 0; n < fieldsNum; n++ {
				nCol := ks.Field(n).Tag.Get("column")
				if nCol == "id" && !idSet {
					continue
				}
				oneParams.Append(mItem.Field(n).Interface())
			}
			id, err := client.TxExec(tx, oneSql, oneParams.ToInterfaces()...)
			if err != nil {
				client.TxRollback(tx)
				var resultTxRollback []int64
				return resultTxRollback, err
			}
			result = append(result, id)
		}
		if !client.TxCommit(tx) {
			return result, errors.New("batch insert (returnIds=true) tx commit failed")
		}
	}
	defer client.CloseConn()
	return result, nil
}
*/

type Goods struct {
	Describe string  `column:"describe"` //
	Id       int64   `column:"id"`       //
	Photo    string  `column:"photo"`    //
	Price    float64 `column:"price"`    //
	Quantity int64   `column:"quantity"` //
	Userid   int64   `column:"userid"`   //
	Goodss   []Goods // This value is used for batch queries and inserts.
}

/*
func (goods *Goods) RowToStruct(row *db.Row) error {
	builder := tsgutils.NewInterfaceBuilder()
	builder.Append(&goods.Describe)
	builder.Append(&goods.Id)
	builder.Append(&goods.Photo)
	builder.Append(&goods.Price)
	builder.Append(&goods.Quantity)
	builder.Append(&goods.Userid)
	err := row.Scan(builder.ToInterfaces()...)
	if err != nil {
		return err
	}
	return nil
}

func (goods *Goods) RowsToStruct(rows *db.Rows) error {
	var goodss []Goods
	builder := tsgutils.NewInterfaceBuilder()
	for rows.Next() {
		builder.Clear()
		builder.Append(&goods.Describe)
		builder.Append(&goods.Id)
		builder.Append(&goods.Photo)
		builder.Append(&goods.Price)
		builder.Append(&goods.Quantity)
		builder.Append(&goods.Userid)
		err := rows.Scan(builder.ToInterfaces()...)
		if err != nil {
			return err
		}
		goodss = append(goodss, *goods)
	}
	if rows != nil {
		defer rows.Close()
	}
	goods.Goodss = goodss
	return nil
}

func (goods *Goods) Insert(client *DBClient, idSet bool) (int64, error) {
	structParam := *goods
	sql := tsgutils.NewStringBuilder()
	qSql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("INSERT INTO ")
	sql.Append("goods")
	sql.Append(" (")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(col).Append("`,")
		qSql.Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	qSql.RemoveLast()
	sql.Append(") VALUES (").Append(qSql.ToString()).Append(");")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (goods *Goods) UpdateGoodsById(client *DBClient) (int64, error) {
	structParam := *goods
	sql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("UPDATE ")
	sql.Append("goods")
	sql.Append(" SET ")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	var id interface{}
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" {
			id = v
			continue
		}
		sql.Append(col).Append("=").Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	params.Append(id)
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (goods *Goods) DeleteGoodsById(client *DBClient) (int64, error) {
	structParam := goods
	sql := tsgutils.NewStringBuilder()
	sql.Append("DELETE FROM ")
	sql.Append("goods")
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), structParam.Id)
}

func (goods *Goods) BatchInsert(client *DBClient, idSet, returnIds bool) ([]int64, error) {
	structParam := *goods
	list := structParam.Goodss
	var result []int64
	listLen := len(list)
	if listLen == 0 {
		return result, errors.New("no data needs to be inserted")
	}
	sql := tsgutils.NewStringBuilder()
	oneQSql := tsgutils.NewStringBuilder()
	batchQSql := tsgutils.NewStringBuilder()
	ks := reflect.TypeOf(structParam)
	fieldsNum := ks.NumField() - 1
	sql.Append("INSERT INTO ")
	sql.Append("goods")
	sql.Append(" (")
	for i := 0; i < fieldsNum; i++ {
		iCol := ks.Field(i).Tag.Get("column")
		if iCol == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(iCol).Append("`,")
	}
	sql.RemoveLast().Append(") VALUES ")
	batchInsertColsLen := tsgutils.InterfaceToInt(tsgutils.IIIInterfaceOperator(idSet, fieldsNum, fieldsNum-1))
	oneQSql.Append("(")
	for j := 0; j < batchInsertColsLen; j++ {
		oneQSql.Append("?,")
	}
	oneQSql.RemoveLast().Append(")")
	if !returnIds {
		for j := 0; j < listLen; j++ {
			batchQSql.Append(oneQSql.ToString()).Append(",")
		}
		batchQSql.RemoveLast()
		batchSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(batchQSql.ToString()).Append(";").ToString()
		batchParams := tsgutils.NewInterfaceBuilder()
		for k := range list {
			item := list[k]
			kItem := reflect.ValueOf(item)
			for l := 0; l < fieldsNum; l++ {
				lCol := ks.Field(l).Tag.Get("column")
				if lCol == "id" && !idSet {
					continue
				}
				batchParams.Append(kItem.Field(l).Interface())
			}
		}
		id, err := client.Exec(batchSql, batchParams.ToInterfaces()...)
		if err != nil {
			return result, err
		}
		result = append(result, id)
	} else {
		oneSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(oneQSql.ToString()).Append(";").ToString()
		oneParams := tsgutils.NewInterfaceBuilder()
		tx, err := client.TxBegin()
		if err != nil {
			return result, err
		}
		for m := range list {
			oneParams.Clear()
			item := list[m]
			mItem := reflect.ValueOf(item)
			for n := 0; n < fieldsNum; n++ {
				nCol := ks.Field(n).Tag.Get("column")
				if nCol == "id" && !idSet {
					continue
				}
				oneParams.Append(mItem.Field(n).Interface())
			}
			id, err := client.TxExec(tx, oneSql, oneParams.ToInterfaces()...)
			if err != nil {
				client.TxRollback(tx)
				var resultTxRollback []int64
				return resultTxRollback, err
			}
			result = append(result, id)
		}
		if !client.TxCommit(tx) {
			return result, errors.New("batch insert (returnIds=true) tx commit failed")
		}
	}
	defer client.CloseConn()
	return result, nil
}
*/

/*

 */
type Communications struct {
	Email           string           `column:"email"`   //
	Id              int64            `column:"id"`      //
	Phone           string           `column:"phone"`   //
	Userid          int64            `column:"userid"`  //
	Weichat         string           `column:"weichat"` //
	Communicationss []Communications // This value is used for batch queries and inserts.
}

func (communications *Communications) RowToStruct(row *db.Row) error {
	builder := tsgutils.NewInterfaceBuilder()
	builder.Append(&communications.Email)
	builder.Append(&communications.Id)
	builder.Append(&communications.Phone)
	builder.Append(&communications.Userid)
	builder.Append(&communications.Weichat)
	err := row.Scan(builder.ToInterfaces()...)
	if err != nil {
		return err
	}
	return nil
}

/*
func (communications *Communications) RowsToStruct(rows *db.Rows) error {
	var communicationss []Communications
	builder := tsgutils.NewInterfaceBuilder()
	for rows.Next() {
		builder.Clear()
		builder.Append(&communications.Email)
		builder.Append(&communications.Id)
		builder.Append(&communications.Phone)
		builder.Append(&communications.Userid)
		builder.Append(&communications.Weichat)
		err := rows.Scan(builder.ToInterfaces()...)
		if err != nil {
			return err
		}
		communicationss = append(communicationss, *communications)
	}
	if rows != nil {
		defer rows.Close()
	}
	communications.Communicationss = communicationss
	return nil
}

func (communications *Communications) Insert(client *DBClient, idSet bool) (int64, error) {
	structParam := *communications
	sql := tsgutils.NewStringBuilder()
	qSql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("INSERT INTO ")
	sql.Append("communications")
	sql.Append(" (")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(col).Append("`,")
		qSql.Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	qSql.RemoveLast()
	sql.Append(") VALUES (").Append(qSql.ToString()).Append(");")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (communications *Communications) UpdateCommunicationsById(client *DBClient) (int64, error) {
	structParam := *communications
	sql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("UPDATE ")
	sql.Append("communications")
	sql.Append(" SET ")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	var id interface{}
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" {
			id = v
			continue
		}
		sql.Append(col).Append("=").Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	params.Append(id)
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (communications *Communications) DeleteCommunicationsById(client *DBClient) (int64, error) {
	structParam := communications
	sql := tsgutils.NewStringBuilder()
	sql.Append("DELETE FROM ")
	sql.Append("communications")
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), structParam.Id)
}

func (communications *Communications) BatchInsert(client *DBClient, idSet, returnIds bool) ([]int64, error) {
	structParam := *communications
	list := structParam.Communicationss
	var result []int64
	listLen := len(list)
	if listLen == 0 {
		return result, errors.New("no data needs to be inserted")
	}
	sql := tsgutils.NewStringBuilder()
	oneQSql := tsgutils.NewStringBuilder()
	batchQSql := tsgutils.NewStringBuilder()
	ks := reflect.TypeOf(structParam)
	fieldsNum := ks.NumField() - 1
	sql.Append("INSERT INTO ")
	sql.Append("communications")
	sql.Append(" (")
	for i := 0; i < fieldsNum; i++ {
		iCol := ks.Field(i).Tag.Get("column")
		if iCol == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(iCol).Append("`,")
	}
	sql.RemoveLast().Append(") VALUES ")
	batchInsertColsLen := tsgutils.InterfaceToInt(tsgutils.IIIInterfaceOperator(idSet, fieldsNum, fieldsNum-1))
	oneQSql.Append("(")
	for j := 0; j < batchInsertColsLen; j++ {
		oneQSql.Append("?,")
	}
	oneQSql.RemoveLast().Append(")")
	if !returnIds {
		for j := 0; j < listLen; j++ {
			batchQSql.Append(oneQSql.ToString()).Append(",")
		}
		batchQSql.RemoveLast()
		batchSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(batchQSql.ToString()).Append(";").ToString()
		batchParams := tsgutils.NewInterfaceBuilder()
		for k := range list {
			item := list[k]
			kItem := reflect.ValueOf(item)
			for l := 0; l < fieldsNum; l++ {
				lCol := ks.Field(l).Tag.Get("column")
				if lCol == "id" && !idSet {
					continue
				}
				batchParams.Append(kItem.Field(l).Interface())
			}
		}
		id, err := client.Exec(batchSql, batchParams.ToInterfaces()...)
		if err != nil {
			return result, err
		}
		result = append(result, id)
	} else {
		oneSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(oneQSql.ToString()).Append(";").ToString()
		oneParams := tsgutils.NewInterfaceBuilder()
		tx, err := client.TxBegin()
		if err != nil {
			return result, err
		}
		for m := range list {
			oneParams.Clear()
			item := list[m]
			mItem := reflect.ValueOf(item)
			for n := 0; n < fieldsNum; n++ {
				nCol := ks.Field(n).Tag.Get("column")
				if nCol == "id" && !idSet {
					continue
				}
				oneParams.Append(mItem.Field(n).Interface())
			}
			id, err := client.TxExec(tx, oneSql, oneParams.ToInterfaces()...)
			if err != nil {
				client.TxRollback(tx)
				var resultTxRollback []int64
				return resultTxRollback, err
			}
			result = append(result, id)
		}
		if !client.TxCommit(tx) {
			return result, errors.New("batch insert (returnIds=true) tx commit failed")
		}
	}
	defer client.CloseConn()
	return result, nil
}
*/
/*

 */
type Label struct {
	Goodsid int64   `column:"goodsid"` //
	Id      int64   `column:"id"`      //
	Kind    string  `column:"kind"`    //
	Labels  []Label // This value is used for batch queries and inserts.
}

/*
func (label *Label) RowToStruct(row *db.Row) error {
	builder := tsgutils.NewInterfaceBuilder()
	builder.Append(&label.Goodsid)
	builder.Append(&label.Id)
	builder.Append(&label.Kind)
	err := row.Scan(builder.ToInterfaces()...)
	if err != nil {
		return err
	}
	return nil
}

func (label *Label) RowsToStruct(rows *db.Rows) error {
	var labels []Label
	builder := tsgutils.NewInterfaceBuilder()
	for rows.Next() {
		builder.Clear()
		builder.Append(&label.Goodsid)
		builder.Append(&label.Id)
		builder.Append(&label.Kind)
		err := rows.Scan(builder.ToInterfaces()...)
		if err != nil {
			return err
		}
		labels = append(labels, *label)
	}
	if rows != nil {
		defer rows.Close()
	}
	label.Labels = labels
	return nil
}

func (label *Label) Insert(client *DBClient, idSet bool) (int64, error) {
	structParam := *label
	sql := tsgutils.NewStringBuilder()
	qSql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("INSERT INTO ")
	sql.Append("label")
	sql.Append(" (")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(col).Append("`,")
		qSql.Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	qSql.RemoveLast()
	sql.Append(") VALUES (").Append(qSql.ToString()).Append(");")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (label *Label) UpdateLabelById(client *DBClient) (int64, error) {
	structParam := *label
	sql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("UPDATE ")
	sql.Append("label")
	sql.Append(" SET ")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	var id interface{}
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" {
			id = v
			continue
		}
		sql.Append(col).Append("=").Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	params.Append(id)
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (label *Label) DeleteLabelById(client *DBClient) (int64, error) {
	structParam := label
	sql := tsgutils.NewStringBuilder()
	sql.Append("DELETE FROM ")
	sql.Append("label")
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), structParam.Id)
}

func (label *Label) BatchInsert(client *DBClient, idSet, returnIds bool) ([]int64, error) {
	structParam := *label
	list := structParam.Labels
	var result []int64
	listLen := len(list)
	if listLen == 0 {
		return result, errors.New("no data needs to be inserted")
	}
	sql := tsgutils.NewStringBuilder()
	oneQSql := tsgutils.NewStringBuilder()
	batchQSql := tsgutils.NewStringBuilder()
	ks := reflect.TypeOf(structParam)
	fieldsNum := ks.NumField() - 1
	sql.Append("INSERT INTO ")
	sql.Append("label")
	sql.Append(" (")
	for i := 0; i < fieldsNum; i++ {
		iCol := ks.Field(i).Tag.Get("column")
		if iCol == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(iCol).Append("`,")
	}
	sql.RemoveLast().Append(") VALUES ")
	batchInsertColsLen := tsgutils.InterfaceToInt(tsgutils.IIIInterfaceOperator(idSet, fieldsNum, fieldsNum-1))
	oneQSql.Append("(")
	for j := 0; j < batchInsertColsLen; j++ {
		oneQSql.Append("?,")
	}
	oneQSql.RemoveLast().Append(")")
	if !returnIds {
		for j := 0; j < listLen; j++ {
			batchQSql.Append(oneQSql.ToString()).Append(",")
		}
		batchQSql.RemoveLast()
		batchSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(batchQSql.ToString()).Append(";").ToString()
		batchParams := tsgutils.NewInterfaceBuilder()
		for k := range list {
			item := list[k]
			kItem := reflect.ValueOf(item)
			for l := 0; l < fieldsNum; l++ {
				lCol := ks.Field(l).Tag.Get("column")
				if lCol == "id" && !idSet {
					continue
				}
				batchParams.Append(kItem.Field(l).Interface())
			}
		}
		id, err := client.Exec(batchSql, batchParams.ToInterfaces()...)
		if err != nil {
			return result, err
		}
		result = append(result, id)
	} else {
		oneSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(oneQSql.ToString()).Append(";").ToString()
		oneParams := tsgutils.NewInterfaceBuilder()
		tx, err := client.TxBegin()
		if err != nil {
			return result, err
		}
		for m := range list {
			oneParams.Clear()
			item := list[m]
			mItem := reflect.ValueOf(item)
			for n := 0; n < fieldsNum; n++ {
				nCol := ks.Field(n).Tag.Get("column")
				if nCol == "id" && !idSet {
					continue
				}
				oneParams.Append(mItem.Field(n).Interface())
			}
			id, err := client.TxExec(tx, oneSql, oneParams.ToInterfaces()...)
			if err != nil {
				client.TxRollback(tx)
				var resultTxRollback []int64
				return resultTxRollback, err
			}
			result = append(result, id)
		}
		if !client.TxCommit(tx) {
			return result, errors.New("batch insert (returnIds=true) tx commit failed")
		}
	}
	defer client.CloseConn()
	return result, nil
}
*/
/*

 */
type Messages struct {
	Contain   string     `column:"contain"` //
	Fromid    int64      `column:"fromid"`  //
	Id        int64      `column:"id"`      //
	Toid      int64      `column:"toid"`    //
	Messagess []Messages // This value is used for batch queries and inserts.
}

/*
func (messages *Messages) RowToStruct(row *db.Row) error {
	builder := tsgutils.NewInterfaceBuilder()
	builder.Append(&messages.Contain)
	builder.Append(&messages.Fromid)
	builder.Append(&messages.Id)
	builder.Append(&messages.Toid)
	err := row.Scan(builder.ToInterfaces()...)
	if err != nil {
		return err
	}
	return nil
}

func (messages *Messages) RowsToStruct(rows *db.Rows) error {
	var messagess []Messages
	builder := tsgutils.NewInterfaceBuilder()
	for rows.Next() {
		builder.Clear()
		builder.Append(&messages.Contain)
		builder.Append(&messages.Fromid)
		builder.Append(&messages.Id)
		builder.Append(&messages.Toid)
		err := rows.Scan(builder.ToInterfaces()...)
		if err != nil {
			return err
		}
		messagess = append(messagess, *messages)
	}
	if rows != nil {
		defer rows.Close()
	}
	messages.Messagess = messagess
	return nil
}

func (messages *Messages) Insert(client *DBClient, idSet bool) (int64, error) {
	structParam := *messages
	sql := tsgutils.NewStringBuilder()
	qSql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("INSERT INTO ")
	sql.Append("messages")
	sql.Append(" (")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(col).Append("`,")
		qSql.Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	qSql.RemoveLast()
	sql.Append(") VALUES (").Append(qSql.ToString()).Append(");")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (messages *Messages) UpdateMessagesById(client *DBClient) (int64, error) {
	structParam := *messages
	sql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("UPDATE ")
	sql.Append("messages")
	sql.Append(" SET ")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	var id interface{}
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" {
			id = v
			continue
		}
		sql.Append(col).Append("=").Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	params.Append(id)
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (messages *Messages) DeleteMessagesById(client *DBClient) (int64, error) {
	structParam := messages
	sql := tsgutils.NewStringBuilder()
	sql.Append("DELETE FROM ")
	sql.Append("messages")
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), structParam.Id)
}

func (messages *Messages) BatchInsert(client *DBClient, idSet, returnIds bool) ([]int64, error) {
	structParam := *messages
	list := structParam.Messagess
	var result []int64
	listLen := len(list)
	if listLen == 0 {
		return result, errors.New("no data needs to be inserted")
	}
	sql := tsgutils.NewStringBuilder()
	oneQSql := tsgutils.NewStringBuilder()
	batchQSql := tsgutils.NewStringBuilder()
	ks := reflect.TypeOf(structParam)
	fieldsNum := ks.NumField() - 1
	sql.Append("INSERT INTO ")
	sql.Append("messages")
	sql.Append(" (")
	for i := 0; i < fieldsNum; i++ {
		iCol := ks.Field(i).Tag.Get("column")
		if iCol == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(iCol).Append("`,")
	}
	sql.RemoveLast().Append(") VALUES ")
	batchInsertColsLen := tsgutils.InterfaceToInt(tsgutils.IIIInterfaceOperator(idSet, fieldsNum, fieldsNum-1))
	oneQSql.Append("(")
	for j := 0; j < batchInsertColsLen; j++ {
		oneQSql.Append("?,")
	}
	oneQSql.RemoveLast().Append(")")
	if !returnIds {
		for j := 0; j < listLen; j++ {
			batchQSql.Append(oneQSql.ToString()).Append(",")
		}
		batchQSql.RemoveLast()
		batchSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(batchQSql.ToString()).Append(";").ToString()
		batchParams := tsgutils.NewInterfaceBuilder()
		for k := range list {
			item := list[k]
			kItem := reflect.ValueOf(item)
			for l := 0; l < fieldsNum; l++ {
				lCol := ks.Field(l).Tag.Get("column")
				if lCol == "id" && !idSet {
					continue
				}
				batchParams.Append(kItem.Field(l).Interface())
			}
		}
		id, err := client.Exec(batchSql, batchParams.ToInterfaces()...)
		if err != nil {
			return result, err
		}
		result = append(result, id)
	} else {
		oneSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(oneQSql.ToString()).Append(";").ToString()
		oneParams := tsgutils.NewInterfaceBuilder()
		tx, err := client.TxBegin()
		if err != nil {
			return result, err
		}
		for m := range list {
			oneParams.Clear()
			item := list[m]
			mItem := reflect.ValueOf(item)
			for n := 0; n < fieldsNum; n++ {
				nCol := ks.Field(n).Tag.Get("column")
				if nCol == "id" && !idSet {
					continue
				}
				oneParams.Append(mItem.Field(n).Interface())
			}
			id, err := client.TxExec(tx, oneSql, oneParams.ToInterfaces()...)
			if err != nil {
				client.TxRollback(tx)
				var resultTxRollback []int64
				return resultTxRollback, err
			}
			result = append(result, id)
		}
		if !client.TxCommit(tx) {
			return result, errors.New("batch insert (returnIds=true) tx commit failed")
		}
	}
	defer client.CloseConn()
	return result, nil
}
*/
// The generated tabs:  user goods communications label messages
