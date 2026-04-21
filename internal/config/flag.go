package config

import "flag"

// SetConfigFlag содержит значения конфигурации, считанные из флагов CLI.
type SetConfigFlag struct {
	// RunAddress — адрес запуска сервера (флаг: -a, по умолчанию "localhost:8000").
	RunAddress string
	// DatabaseURI — строка подключения к базе данных (флаг: -d).
	DatabaseURI string
	// AccrualSystemAddress — адрес сервиса начислений (флаг: -r).
	AccrualSystemAddress string
}

// ParseFlag разбирает флаги командной строки и возвращает заполненный SetConfigFlag.
// Если флаги уже были разобраны ранее, повторный вызов flag.Parse() пропускается.
func ParseFlag() SetConfigFlag {
	var configFlag = SetConfigFlag{}
	flag.StringVar(&configFlag.RunAddress, "a", "localhost:8000", "Run address")
	flag.StringVar(&configFlag.DatabaseURI, "d", "", "Database URI")
	flag.StringVar(&configFlag.AccrualSystemAddress, "r", "", "System address")
	if !flag.Parsed() {
		flag.Parse()
	}
	return configFlag
}
