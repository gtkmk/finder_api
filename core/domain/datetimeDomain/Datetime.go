package datetimeDomain

import (
	"github.com/gtkmk/finder_api/core/domain/helper"
	"time"
)

const (
	DATETIME_FORMAT              = "2006-01-02 15:04:05"
	DATETIME_FORMAT_MILLISECONDS = "2006-01-02 15:04:05.000"
	DATE_REVERTED                = "2006-01-02"
	LAYOUT_PARSER                = "02/01/2006"
	LAYOUT_REVERTED              = "2006-01-02 15:04:05"
	DateTimeFormatWithSlashConst = "2006/01/02 15:04:05"
)

func CreateFormattedNow() string {
	formattedNow := (time.Now()).Format(DATETIME_FORMAT)
	return formattedNow
}

func FormatDateAsString(date time.Time) string {
	return date.Format(DATETIME_FORMAT)
}

func CreateNow() (time.Time, error) {
	timezoneLocation, err := time.LoadLocation("America/Sao_Paulo")

	if err != nil {
		return time.Now(),
			helper.ErrorBuilder(
				helper.ErrorWithCodeConst,
				helper.ErrorLoadingSaoPauloLocationCodeConst,
			)
	}

	formattedTime, err := time.Parse(
		DATETIME_FORMAT_MILLISECONDS,
		(time.Now().In(timezoneLocation)).Format(DATETIME_FORMAT_MILLISECONDS),
	)

	if err != nil {
		return time.Now(),
			helper.ErrorBuilder(
				helper.ErrorWithCodeConst,
				helper.ErrorGeneratingFormattedTimestampCodeConst,
			)
	}

	return formattedTime, nil
}

func FormatDateAsTime(date string) (time.Time, error) {
	formatted, err := time.Parse(LAYOUT_PARSER, date)

	if err != nil {
		return time.Time{}, err
	}

	return formatted, err
}

func FormatDateAsTimeReverted(date string) (time.Time, error) {
	formatted, err := time.Parse(LAYOUT_REVERTED, date)

	if err != nil {
		return time.Time{}, err
	}

	return formatted, err
}

func FormatDate(unformattedDate string) (string, error) {
	formatted, err := time.Parse(LAYOUT_PARSER, unformattedDate)

	if err != nil {
		return "", err
	}

	formattedDate := formatted.Format(DATETIME_FORMAT)
	return formattedDate, err
}

func FormatDateFromDefaultToDateTimeWithSlash(dateStr interface{}) string {
	dataOriginalString, ok := dateStr.(string)
	if !ok {
		return ""
	}

	formatedDate, err := time.Parse(DATETIME_FORMAT, dataOriginalString)
	if err != nil {
		return ""
	}

	return formatedDate.Format(DateTimeFormatWithSlashConst)
}
