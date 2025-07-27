package models

import "time"

type ChatDetails struct {
	Session string `json:"session"`
	ChatId  string `json:"chatId"`
}

type SendTextDetails struct {
	Session                string `json:"session"`
	ChatId                 string `json:"chatId"`
	ReplyTo                string `json:"reply_to"`
	LinkPreview            bool   `json:"linkPreview"`
	LinkPreviewHighQuality bool   `json:"linkPreviewHighQuality"`
	Text                   string `json:"text"`
}
type Pic struct {
	FormattedPhoneNumber string `json:"formattedPhoneNumber"`
	Session              string `json:"session"`
	IsInternal           bool   `json:"isInternal"`
}

type Customer struct {
	FormattedPhoneNumber string `json:"formattedPhoneNumber"`
	Name                 string `json:"name"`
	IsInternal           bool   `json:"isInternal"`
	// ChatId               string `json:"chatId"`
}

type Job struct {
	Pic      Pic      `json:"pic"`
	Customer Customer `json:"customer"`
	Success  bool     `json:"success"`
	Voucher  string   `json:"voucher"`
}

type JobResponse struct {
	CustomerNumber string `json:"customerNumber"`
	Name           string `json:"name"`
	Voucher        string `json:"voucher"`
}

type PromoToken struct {
	UserName  string    `json:"userName"`
	PromoCode string    `json:"promoCOde"`
	ExpiresAt time.Time `json:"expiresAt"`
}
