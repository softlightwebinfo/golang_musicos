package models

import (
	"../settings"
	_ "database/sql"
	_ "errors"
	"time"
	_ "time"
)

type Package struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	UpdatedAt       time.Time `json:"updated_at"`
	PackagePrice    float32   `json:"package_price"`
	PackageAmount   int       `json:"package_amount"`
	PackageDiscount float32   `json:"package_discount"`
	Value           float32   `json:"value"`
}
type PackagePermissions struct {
	FkIdPackage    int      `json:"fk_id_package"`
	FkIdPermission int      `json:"fk_id_permission"`
	Value          *float32 `json:"value"`
	Price          *float32 `json:"price"`
	Active         bool     `json:"active"`
	Key            string   `json:"key"`
	Name           string   `json:"name"`
}
type PackageGroup struct {
	Package           Package              `json:"package"`
	PackagePermission []PackagePermissions `json:"package_permission"`
}

func Packages(business bool) (pack []Package, err error) {
	q := `
SELECT p.id,
       p.name,
       p.updated_at,
       pp.value as package_price,
       case WHEN po IS NULL then null else po.fk_amount_type END as package_amount,
       case WHEN po IS NULL then null else po.discount END as package_discount,
       CASE
           WHEN po IS NULL then pp.value
           else
               CASE
                   when po.fk_amount_type = 1 THEN ROUND(po.discount, 2)
                   when po.fk_amount_type = 2 THEN ROUND(pp.value - (pp.value * po.discount) / 100, 2)
                   when po.fk_amount_type = 3 THEN
                       CASE
                           WHEN pp.value > po.discount
                               THEN round(pp.value - po.discount, 2)
                           ELSE
                               ROUND(pp.value, 2)
                           END
                   ELSE ROUND(pp.value, 2)
                   END
           END                                                                as value
FROM packages p
         INNER JOIN packages_prices pp on p.id = pp.fk_id_package
         LEFT JOIN packages_offers po
                   on p.id = po.fk_id_package and po.start <= now() and (po."end" IS NULL OR po.end >= now())
WHERE p.active = true
AND p.business = true
ORDER BY p.id
`
	db := settings.InstanceDb
println(business)
	rows, err := db.Query(q)
	if err != nil {
		return
	}
	for rows.Next() {
		pk := Package{}
		_ = rows.Scan(
			&pk.Id,
			&pk.Name,
			&pk.UpdatedAt,
			&pk.PackagePrice,
			&pk.PackageAmount,
			&pk.PackageDiscount,
			&pk.Value,
		)
		pack = append(pack, pk)
	}
	defer rows.Close()
	return
}
func PackagesGroup(business bool) (packages []PackageGroup, err error) {
	packs, _ := Packages(business)
	for _, pack := range packs {
		pk := PackageGroup{
			Package:           pack,
			PackagePermission: PackagesPermissions(pack),
		}
		packages = append(packages, pk)
	}
	return
}
func PackagesPermissions(packa Package) (pack []PackagePermissions) {
	q := `
	SELECT t.fk_id_package,
       t.fk_id_permission,
       t.value,
       t.price,
       t.active,
	   t.name,
	   t.key
FROM (
         SELECT 
			CASE when pop is NULL then pp.fk_id_package ELSE po.fk_id_package END,
			CASE when pop is NULL then pp.fk_id_permission ELSE pop.fk_id_permission END,
			CASE when pop is NULL then pp.value ELSE pop.value END,
			CASE when pop is NULL then pp.price ELSE pop.price END,
			CASE when pop is NULL then pp.active ELSE pop.active END,
			ps.name,
			ps.key
         FROM permissions ps
                  LEFT JOIN packages_permissions pp
                            on ps.id = pp.fk_id_permission and pp.fk_id_package = $1 and pp.active = true
                  LEFT JOIN packages_offers_permissions pop on ps.id = pop.fk_id_permission and pop.active = true
                  LEFT JOIN packages_offers po on pop.fk_id_offer = po.id and po.fk_id_package=pp.fk_id_package
                  left join packages p on po.fk_id_package = p.id and p.active = true                  
         ORDER BY ps.id
     ) as t
WHERE t.active = true
  and fk_id_package is Not null
  and fk_id_permission is not null
`
	rows, err := settings.InstanceDb.Query(
		q,
		&packa.Id,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		pk := PackagePermissions{}
		_ = rows.Scan(
			&pk.FkIdPackage,
			&pk.FkIdPermission,
			&pk.Value,
			&pk.Price,
			&pk.Active,
			&pk.Name,
			&pk.Key,
		)
		pack = append(pack, pk)
	}
	defer rows.Close()
	return
}

func GetOfferPackageActive(id int64) (item AclOffer, err error) {
	db := settings.InstanceDb
	err = db.QueryRow(`SELECT id FROM packages_offers po WHERE po.fk_id_package=$1 and po.start <= now() and (po."end" >= now() OR po."end" is NULL)`, id).
		Scan(
			&item.Id,
		)
	return
}

func CreateOrder(fkIdUser int64, fkIdOffer int64, value float32, invoice string) (err error) {
	db := settings.InstanceDb
	q := `INSERT INTO users_packages_orders(fk_id_user, fk_id_offer, value, invoice) VALUES($1, $2, $3, $4)`
	_, err = db.Exec(
		q,
		fkIdUser,
		fkIdOffer,
		value,
		invoice,
	)
	return
}

func GetOrder(invoice string) (item AclOfferActivate, err error) {
	db := settings.InstanceDb
	err = db.QueryRow(`SELECT fk_id_user,fk_id_offer,invoice FROM users_packages_orders WHERE invoice=$1`, invoice).
		Scan(
			&item.FkIdUser,
			&item.FkIdOffer,
			&item.Invoice,
		)
	return
}

func ActivateOrder(aclOffer AclOfferActivate, pack AclPackageActivate) (err error) {
	db := settings.InstanceDb
	q := `INSERT INTO 
    users_packages_orders_records(fk_id_user, fk_id_offer, invoice, txn_type, subscr_id, item_name, recurring, payer_status, payer_email, subscr_date, recur_times, custom, period1, mc_amount1, period3, mc_amount3, ipn_track_id) 
    VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`
	_, err = db.Exec(
		q,
		&aclOffer.FkIdUser,
		&aclOffer.FkIdOffer,
		&aclOffer.Invoice,
		&pack.TxnType,
		&pack.SubscrId,
		&pack.ItemName,
		&pack.Recurring,
		&pack.PayerStatus,
		&pack.PayerEmail,
		&pack.SubscrDate,
		&pack.RecurTimes,
		&pack.Custom,
		&pack.Period1,
		&pack.McAmount1,
		&pack.Period3,
		&pack.McAmount3,
		&pack.IpnTrackId,
	)
	return
}

func DeleteOrder(item AclOfferActivate) (err error) {
	db := settings.InstanceDb
	q := `DELETE FROM users_packages_orders WHERE fk_id_user=$1 AND fk_id_offer=$2`
	_, err = db.Exec(
		q,
		&item.FkIdUser,
		&item.FkIdOffer,
	)
	return
}

func ChangeUserPackage(item AclOfferActivate) (err error) {
	db := settings.InstanceDb
	now := time.Now()
	after := now.AddDate(0, 12, 0)
	q := `INSERT INTO users_packages(fk_id_user, fk_id_package_offer, expired_at) VALUES($1, $2, $3)`
	_, err = db.Exec(
		q,
		&item.FkIdUser,
		&item.FkIdOffer,
		&after,
	)
	return
}

func DeleteUserPackages(item AclOfferActivate) (err error) {
	db := settings.InstanceDb
	q := `UPDATE users_packages SET active=false WHERE fk_id_user=$1`
	_, err = db.Exec(
		q,
		&item.FkIdUser,
	)
	return
}
