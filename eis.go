package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"time"
)

type Eis struct {
	Count   int
	Updated int
}

func (t *Eis) run() {
	defer SaveStack()
	db, err := sql.Open("mysql", Dsn)
	if err != nil {
		Logging("Ошибка подключения к БД", err)
		return
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Second * 3600)
	stmt, _ := db.Prepare(fmt.Sprintf("SELECT COUNT(id_tender) FROM %stender WHERE type_fz = ?", Prefix))
	res, err := stmt.Query(Typefz)
	stmt.Close()
	if err != nil {
		Logging("Ошибка выполения запроса", err)
		return
	}
	t.Count = 0
	if res.Next() {
		err = res.Scan(&t.Count)
		if err != nil {
			Logging("Ошибка чтения результата запроса", err)
			return
		}
		res.Close()
	} else {
		res.Close()
		Logging("Count tenders return none", err)
		return
	}
	page := t.Count / 1000
	if page == 0 {
		u := fmt.Sprintf("SELECT id_tender FROM %stender WHERE type_fz = ?", Prefix)
		t.GetTenders(db, u)
	} else {
		for i := 0; i <= page; i++ {
			u := fmt.Sprintf("SELECT id_tender FROM %stender WHERE type_fz = ? LIMIT 1000 OFFSET %d", Prefix, i*1000)
			t.GetTenders(db, u)
		}
	}

}

func (t *Eis) GetTenders(db *sql.DB, u string) {
	defer SaveStack()
	stmt, _ := db.Prepare(u)
	res, err := stmt.Query(Typefz)
	stmt.Close()
	if err != nil {
		Logging("Ошибка выполения запроса", err)
		return
	}
	for res.Next() {
		var idTen sql.NullInt64
		err = res.Scan(&idTen)
		if err != nil {
			Logging("Ошибка чтения результата запроса", err)
			return
		}
		if idTen.Valid {
			idT := int(idTen.Int64)
			err = t.TenderKwords(db, idT)
			if err != nil {
				Logging("Ошибка TenderKwords", err)
				return
			}
		}

	}
	res.Close()

}

func (t *Eis) TenderKwords(db *sql.DB, idTender int) error {
	resString := ""
	stmt, _ := db.Prepare(fmt.Sprintf("SELECT DISTINCT po.name, po.okpd_name FROM %spurchase_object AS po LEFT JOIN %slot AS l ON l.id_lot = po.id_lot WHERE l.id_tender = ?", Prefix, Prefix))
	rows, err := stmt.Query(idTender)
	stmt.Close()
	if err != nil {
		Logging("Ошибка выполения запроса", err)
		return err
	}
	for rows.Next() {
		var name sql.NullString
		var okpdName sql.NullString
		err = rows.Scan(&name, &okpdName)
		if err != nil {
			Logging("Ошибка чтения результата запроса", err)
			return err
		}
		if name.Valid {
			resString = fmt.Sprintf("%s %s ", resString, name.String)
		}
		if okpdName.Valid {
			resString = fmt.Sprintf("%s %s ", resString, okpdName.String)
		}
	}
	rows.Close()
	stmt1, _ := db.Prepare(fmt.Sprintf("SELECT DISTINCT file_name FROM %sattachment WHERE id_tender = ?", Prefix))
	rows1, err := stmt1.Query(idTender)
	stmt1.Close()
	if err != nil {
		Logging("Ошибка выполения запроса", err)
		return err
	}
	for rows1.Next() {
		var attName sql.NullString
		err = rows1.Scan(&attName)
		if err != nil {
			Logging("Ошибка чтения результата запроса", err)
			return err
		}
		if attName.Valid {
			resString = fmt.Sprintf("%s %s ", resString, attName.String)
		}
	}
	rows1.Close()
	idOrg := 0
	stmt2, _ := db.Prepare(fmt.Sprintf("SELECT purchase_object_info, id_organizer FROM %stender WHERE id_tender = ?", Prefix))
	rows2, err := stmt2.Query(idTender)
	stmt2.Close()
	if err != nil {
		Logging("Ошибка выполения запроса", err)
		return err
	}
	for rows2.Next() {
		var idOrgNull sql.NullInt64
		var purOb sql.NullString
		err = rows2.Scan(&purOb, &idOrgNull)
		if err != nil {
			Logging("Ошибка чтения результата запроса", err)
			return err
		}
		if idOrgNull.Valid {
			idOrg = int(idOrgNull.Int64)
		}
		if purOb.Valid {
			resString = fmt.Sprintf("%s %s ", resString, purOb.String)
		}

	}
	rows2.Close()
	if idOrg != 0 {
		stmt3, _ := db.Prepare(fmt.Sprintf("SELECT full_name, inn FROM %sorganizer WHERE id_organizer = ?", Prefix))
		rows3, err := stmt3.Query(idOrg)
		stmt3.Close()
		if err != nil {
			Logging("Ошибка выполения запроса", err)
			return err
		}
		for rows3.Next() {
			var innOrg sql.NullString
			var nameOrg sql.NullString
			err = rows3.Scan(&nameOrg, &innOrg)
			if err != nil {
				Logging("Ошибка чтения результата запроса", err)
				return err
			}
			if innOrg.Valid {

				resString = fmt.Sprintf("%s %s ", resString, innOrg.String)
			}
			if nameOrg.Valid {
				resString = fmt.Sprintf("%s %s ", resString, nameOrg.String)
			}

		}
		rows3.Close()
	}
	stmt4, _ := db.Prepare(fmt.Sprintf("SELECT DISTINCT cus.inn, cus.full_name FROM %scustomer AS cus LEFT JOIN %spurchase_object AS po ON cus.id_customer = po.id_customer LEFT JOIN %slot AS l ON l.id_lot = po.id_lot WHERE l.id_tender = ?", Prefix, Prefix, Prefix))
	rows4, err := stmt4.Query(idTender)
	stmt4.Close()
	if err != nil {
		Logging("Ошибка выполения запроса", err)
		return err
	}
	for rows4.Next() {
		var innC sql.NullString
		var fullNameC sql.NullString
		err = rows4.Scan(&innC, &fullNameC)
		if err != nil {
			Logging("Ошибка чтения результата запроса", err)
			return err
		}
		if innC.Valid {

			resString = fmt.Sprintf("%s %s ", resString, innC.String)
		}
		if fullNameC.Valid {
			resString = fmt.Sprintf("%s %s ", resString, fullNameC.String)
		}
	}
	rows4.Close()
	re := regexp.MustCompile(`\s+`)
	resString = re.ReplaceAllString(resString, " ")
	stmtr, _ := db.Prepare(fmt.Sprintf("UPDATE %stender SET tender_kwords = ? WHERE id_tender = ?", Prefix))
	_, errr := stmtr.Exec(resString, idTender)
	stmtr.Close()
	if errr != nil {
		Logging("Ошибка вставки TenderKwords", errr)
		return err
	}
	t.Updated++
	return nil
}
