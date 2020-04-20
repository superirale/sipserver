package sip

import (
	"strconv"
	"strings"
)

type CSeqStruct struct {
	CSeqText   string
	CSeqNumber int
	SIPMethod  string
}

func BuildCSeq(c string) *CSeqStruct {
	cseq := new(CSeqStruct)
	cArr := strings.Split(c, " ")
	if len(cArr) == 3 {
		cseq.CSeqText = c
		n, _ := strconv.Atoi(cArr[1])
		cseq.CSeqNumber = n
		cseq.SIPMethod = cArr[2]
	}
	return cseq
}