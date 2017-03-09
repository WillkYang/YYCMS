package models


func FormatStartEndTimeParams(starttime, endtime *string) {
	FormatStartTime(starttime)
	FormatEndTime(endtime)
}

func FormatStartTime(starttime *string) {
	if len(*starttime) <= 0 {
		*starttime = "-1"
	}
}

func FormatEndTime(endtime *string) {
	if len(*endtime) <= 0 {
		*endtime = "999999999999999"
	}
}