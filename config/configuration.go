package config

import (
	"os"
)

const (
	LOCAL       = "local"
	DEVELOPMENT = "development"
	PRODUCTION  = "production"
)

const ENVIRONMENT string = LOCAL

var config_env = map[string]map[string]string{
	"local": {
		"PORT":    "8000",
		"PQS_URL": "http://127.0.0.1:8000",

		"MYSQL_HOST":   "localhost",
		"MYSQL_PORT":   "3306",
		"MYSQL_USER":   "root",
		"MYSQL_PASS":   "",
		"MYSQL_SCHEMA": "erp_salju",

		"TOKEN_KEY": "ftr$;C3Uck=2AH/xe(q;}Ak=#%2#@M?BNTrKPP[+zyP.B@G25@%L#AUQ}cvM[ZJ(7}hCNF;qrc$zPz?TB$YT+;BMK6!,SV?PzYXKUvG{:B-XKtL)(awL3ic$AjSzmq9bZ(3WTYrU_V8q*prA._pm;iv_=.FiD+LH+!&U-tpa}/ZzQ:RQ?U?uy75j6v*m[.!t$9UccH+j",
	},
}

var CONFIG = config_env[ENVIRONMENT]

func GetConfig(key string, config string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return config
}

func InitConfig() {
	for key := range CONFIG {
		CONFIG[key] = GetConfig(key, CONFIG[key])
		os.Setenv(key, CONFIG[key])
	}
}
