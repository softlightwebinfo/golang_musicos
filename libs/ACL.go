package libs

import (
	"database/sql"
)

type UserPackage struct {
	FkIdPackageOffer int     `json:"fk_id_package_offer"`
	FkIdPackage      int     `json:"fk_id_package"`
	FkAmountType     int     `json:"fk_amount_type"`
	Discount         float32 `json:"discount"`
	Name             string  `json:"name"`
	PackagePrice     float32 `json:"package_price"`
	PriceTotal       float32 `json:"price_total"`
}
type PackagePermission struct {
	FkIdPackage    int      `json:"fk_id_package"`
	FkIdPermission int      `json:"fk_id_permission"`
	Value          *float32 `json:"value"`
	Price          *float32 `json:"price"`
	Active         bool     `json:"active"`
}
type ACLResponse struct {
	Package           UserPackage         `json:"package"`
	PackagePermission []PackagePermission `json:"package_permission"`
}
type ACL struct {
	Id                 int64
	userPackage        UserPackage
	packagePermissions []PackagePermission
	db                 *sql.DB
}

func (then *ACL) Init(id int64) {
	then.db = GetConnection()
	then.Id = id
	then.getPackageSelected()
	then.getPackagePermissions()
	defer then.db.Close()
}
func (then *ACL) getPackageSelected() {
	q := `SELECT 
				up.fk_id_package_offer,
       			po.fk_id_package,
				po.fk_amount_type,
				po.discount,
				p.name,
    			pp.value,       		
       			CASE 
       			    when po.fk_amount_type=1 THEN ROUND(po.discount, 2)
       			    when po.fk_amount_type=2 THEN ROUND(pp.value - (pp.value*po.discount)/100, 2)
       			    when po.fk_amount_type=3 THEN 
       			    	CASE 
       			    	    WHEN pp.value>po.discount
       			    	    	THEN round(pp.value - po.discount, 2)
							ELSE
       			    	    	ROUND(pp.value, 2)
						END 
					ELSE ROUND(pp.value, 2)
				END
			from users_packages up 
			INNER JOIN packages_offers po ON po.id=up.fk_id_package_offer
			    INNER JOIN packages_prices pp on po.fk_id_package=pp.fk_id_package
			INNER JOIN packages p ON po.fk_id_package=p.id
			WHERE (up.expired_at IS NULL OR up.expired_at >= now()) and up.active=true AND up.fk_id_user=$1`
	err := then.db.QueryRow(q, then.Id).Scan(
		&then.userPackage.FkIdPackageOffer,
		&then.userPackage.FkIdPackage,
		&then.userPackage.FkAmountType,
		&then.userPackage.Discount,
		&then.userPackage.Name,
		&then.userPackage.PackagePrice,
		&then.userPackage.PriceTotal,
	)
	if err != nil {
		println(err.Error())
		return
	}
}
func (then *ACL) getPackagePermissions() {
	q := `
		SELECT t.fk_id_package,
       t.fk_id_permission,
       t.value,
       t.price,
       t.active
FROM (
         SELECT CASE
                    when upp is null then case when pop is NULL then pp.fk_id_package ELSE p.id END
                    ELSE pou.id END as fk_id_package,
                CASE
                    when upp is null then case when pop is NULL then pp.fk_id_permission ELSE pop.fk_id_permission END
                    ELSE upp.fk_id_permission END,
                CASE
                    when upp is null then case when pop is NULL then pp.value ELSE pop.value END
                    ELSE upp.value END,
                CASE
                    when upp is null then case when pop is NULL then pp.price ELSE pop.price END
                    ELSE upp.price END,
                CASE
                    when upp is null then case when pop is NULL then pp.active ELSE pop.active END
                    ELSE upp.active END
         FROM permissions ps
                  LEFT JOIN packages_permissions pp
                            on ps.id = pp.fk_id_permission and pp.fk_id_package = $1 and pp.active = true
                  LEFT JOIN packages_offers_permissions pop on ps.id = pop.fk_id_permission and pop.active = true
                  LEFT JOIN packages_offers po on pop.fk_id_offer = po.id
                  left join packages p on po.fk_id_package = p.id and p.active = true
                  LEFT JOIN users_packages_permissions upp
                            on upp.fk_id_user = $2 AND fk_id_package_offer = $3 and upp.fk_id_permission = ps.id and
                               upp.active = true
                  LEFT JOIN packages_offers pou on upp.fk_id_package_offer = pou.id
         ORDER BY ps.id
     ) as t
WHERE t.active = true
  and fk_id_package is Not null
  and fk_id_permission is not null
`
	rows, err := then.db.Query(
		q,
		&then.userPackage.FkIdPackage,
		&then.Id,
		&then.userPackage.FkIdPackageOffer,
	)
	if err != nil {
		println(err.Error())
		return
	}

	for rows.Next() {
		pack := PackagePermission{}
		_ = rows.Scan(
			&pack.FkIdPackage,
			&pack.FkIdPermission,
			&pack.Value,
			&pack.Price,
			&pack.Active,
		)
		then.packagePermissions = append(then.packagePermissions, pack)
	}
}
func (then *ACL) GetData() ACLResponse {
	return ACLResponse{
		Package:           then.userPackage,
		PackagePermission: then.packagePermissions,
	}
}
