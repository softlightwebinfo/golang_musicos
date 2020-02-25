package models

type AclPackageInvoice struct {
	Id    int64   `json:"id"`
	Value float32 `json:"value"`
}

type AclPackageInvoiceSuccess struct {
	Invoice string `json:"invoice"`
}

type AclOffer struct {
	Id int64 `json:"id"`
}

type AclPackageActivate struct {
	TxnType     string `json:"txn_type"`
	SubscrId    string `json:"subscr_id"`
	ItemName    string `json:"item_name"`
	Recurring   string `json:"recurring"`
	PayerStatus string `json:"payer_status"`
	PayerEmail  string `json:"payer_email"`
	Invoice     string `json:"invoice"`
	SubscrDate  string `json:"subscr_date"`
	RecurTimes  string `json:"recur_times"`
	Custom      string `json:"custom"`
	Period1     string `json:"period_1"`
	McAmount1   string `json:"mc_amount_1"`
	Period3     string `json:"period_3"`
	McAmount3   string `json:"mc_amount_3"`
	IpnTrackId  string `json:"ipn_track_id"`
}

type AclOfferActivate struct {
	FkIdUser  int64  `json:"fk_id_user"`
	FkIdOffer int64  `json:"fk_id_offer"`
	Invoice   string `json:"invoice"`
}
