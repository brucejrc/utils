package id

import "time"

type CodeOptions struct {
	chars []rune
	n1    int
	n2    int
	l     int
	salt  uint64
}

func WithCodeChars(arr []rune) func(opions *CodeOptions) {
	return func(options *CodeOptions) {
		if len(arr) > 0 {
			getCodeOptionsOrSetDefault(options).chars = arr
		}
	}
}

func WithCodeN1(n1 int) func(options *CodeOptions) {
	return func(options *CodeOptions) {
		getCodeOptionsOrSetDefault(options).n1 = n1
	}
}

func WithCodeN2(n2 int) func(options *CodeOptions) {
	return func(options *CodeOptions) {
		getCodeOptionsOrSetDefault(options).n2 = n2
	}
}

func WithCodeL(length int) func(options *CodeOptions) {
	return func(options *CodeOptions) {
		if length > 0 {
			getCodeOptionsOrSetDefault(options).l = length
		}
	}
}

func WithCodeSalt(salt uint64) func(options *CodeOptions) {
	return func(options *CodeOptions) {
		if salt > 0 {
			getCodeOptionsOrSetDefault(options).salt = salt
		}
	}
}

func getCodeOptionsOrSetDefault(options *CodeOptions) *CodeOptions {
	if options == nil {
		return &CodeOptions{
			chars: []rune{
				'2', '3', '4', '5', '6',
				'7', '8', '9', 'A', 'B',
				'C', 'D', 'E', 'F', 'G',
				'H', 'J', 'K', 'L', 'M',
				'N', 'P', 'Q', 'R', 'S',
				'T', 'V', 'W', 'X', 'Y',
			},
			n1:   17,
			n2:   5,
			l:    8,
			salt: 123567369,
		}
	}

	return options
}

type SonyflakeOptions struct {
	machineId uint16
	startTime time.Time
}

func WithSnoyflakeMachineId(id uint16) func(options *SonyflakeOptions) {
	return func(options *SonyflakeOptions) {
		if id > 0 {
			getSonyflakeOptionsOrSetDefault(options).machineId = id
		}
	}
}

func WithSnoyflakeStartTime(time time.Time) func(options *SonyflakeOptions) {
	return func(options *SonyflakeOptions) {
		getSonyflakeOptionsOrSetDefault(options).startTime = time
	}
}

func getSonyflakeOptionsOrSetDefault(options *SonyflakeOptions) *SonyflakeOptions {
	if options == nil {
		return &SonyflakeOptions{
			machineId: 1,
			startTime: time.Date(2022, 03, 22, 0, 0, 0, 0, time.UTC),
		}
	}
	return options
}
