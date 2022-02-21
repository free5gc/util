package logger

import (
	"reflect"

	"github.com/asaskevich/govalidator"
)

type Logger struct {
	AMF   *LogSetting `yaml:"AMF" valid:"optional"`
	AUSF  *LogSetting `yaml:"AUSF" valid:"optional"`
	N3IWF *LogSetting `yaml:"N3IWF" valid:"optional"`
	NRF   *LogSetting `yaml:"NRF" valid:"optional"`
	NSSF  *LogSetting `yaml:"NSSF" valid:"optional"`
	PCF   *LogSetting `yaml:"PCF" valid:"optional"`
	SMF   *LogSetting `yaml:"SMF" valid:"optional"`
	UDM   *LogSetting `yaml:"UDM" valid:"optional"`
	UDR   *LogSetting `yaml:"UDR" valid:"optional"`
	UPF   *LogSetting `yaml:"UPF" valid:"optional"`
	NEF   *LogSetting `yaml:"NEF" valid:"optional"`
	BSF   *LogSetting `yaml:"BSF" valid:"optional"`
	CHF   *LogSetting `yaml:"CHF" valid:"optional"`
	UDSF  *LogSetting `yaml:"UDSF" valid:"optional"`
	NWDAF *LogSetting `yaml:"NWDAF" valid:"optional"`
	WEBUI *LogSetting `yaml:"WEBUI" valid:"optional"`

	Aper *LogSetting `yaml:"Aper" valid:"optional"`
	FSM  *LogSetting `yaml:"FSM" valid:"optional"`
	NAS  *LogSetting `yaml:"NAS" valid:"optional"`
	NGAP *LogSetting `yaml:"NGAP" valid:"optional"`
	PFCP *LogSetting `yaml:"PFCP" valid:"optional"`
}

func (l *Logger) Validate() (bool, error) {
	logger := reflect.ValueOf(l).Elem()
	for i := 0; i < logger.NumField(); i++ {
		if logSetting := logger.Field(i).Interface().(*LogSetting); logSetting != nil {
			result, err := logSetting.validate()
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(l)
	return result, err
}

type LogSetting struct {
	DebugLevel   string `yaml:"debugLevel" valid:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller" valid:"type(bool)"`
}

func (l *LogSetting) validate() (bool, error) {
	govalidator.TagMap["debugLevel"] = govalidator.Validator(func(str string) bool {
		if str == "panic" || str == "fatal" || str == "error" || str == "warn" ||
			str == "info" || str == "debug" || str == "trace" {
			return true
		} else {
			return false
		}
	})

	result, err := govalidator.ValidateStruct(l)
	return result, err
}
