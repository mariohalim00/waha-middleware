package models

import (
	"time"
)

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
	Username             string `json:"username"`
	IsInternal           bool   `json:"isInternal"`
	Voucher              string `json:"voucher"`
}

type Job struct {
	Pic      Pic      `json:"pic"`
	Customer Customer `json:"customer"`
	Success  bool     `json:"success"`
}

type JobResponse struct {
	CustomerNumber string `json:"customerNumber"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	Voucher        string `json:"voucher"`
}

type PromoToken struct {
	UserName  string    `json:"userName"`
	PromoCode string    `json:"promoCOde"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type CustdatumDto struct {
	Username                string    `json:"username"`
	Userid                  string    `json:"userid"`
	Suspend                 bool      `json:"suspend"`
	StatusActive            bool      `json:"statusActive"`
	Name                    string    `json:"name"`
	Phone                   string    `json:"phone"`
	Email                   string    `json:"email"`
	MailSubs                bool      `json:"mailSubs"`
	FundMethod              string    `json:"fundMethod"`
	Balance                 int64     `json:"balance"`
	BankGroup               string    `json:"bankGroup"`
	RegisDate               time.Time `json:"regisDate"`
	FirstPayDate            time.Time `json:"firstPayDate"`
	LastPayDate             time.Time `json:"lastPayDate"`
	LastLoginDate           time.Time `json:"lastLoginDate"`
	LastDetailDate          time.Time `json:"lastDetailDate"`
	DateCreated             time.Time `json:"dateCreated"`
	SocialContact           string    `json:"socialContact"`
	HomeAddress             string    `json:"homeAddress"`
	ReferralLink            string    `json:"referralLink"`
	LastWithdrawalDate      time.Time `json:"lastWithdrawalDate"`
	MemberLevel             string    `json:"memberLevel"`
	LevelValidity           string    `json:"levelValidity"`
	LevelStartDate          time.Time `json:"levelStartDate"`
	LevelExpiredDate        time.Time `json:"levelExpiredDate"`
	CurrentRequiredDeposit  string    `json:"currentRequiredDeposit"`
	CurrentRequiredTurnover string    `json:"currentRequiredTurnover"`
	TotalDeposit            int64     `json:"totalDeposit"`
	TotalWithdrawal         int64     `json:"totalWithdrawal"`
	TotalPromotion          string    `json:"totalPromotion"`
	RefComm                 int64     `json:"refComm"`
	TotalSum                int64     `json:"totalSum"`
	RegisterDomain          string    `json:"registerDomain"`
	IsInactive              bool      `json:"isInactive"`
	LastBlastDate           time.Time `json:"lastBlastDate"`
	Priority                int32     `json:"priority"`
	InProcess               bool      `json:"inProcess"`
	Bucket                  int32     `json:"bucket"`
	ProcessStartDate        time.Time `json:"processStartDate"`
	PhoneNumberExists       bool      `json:"phoneNumberExists"`
}
