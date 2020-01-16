package text

import (
	"strings"
	"testing"
)

func TestReadingTime(t *testing.T) {

	text := `В блистательном сериале «Семнадцать мгновений весны» есть два примечательных момента, на
  которые редко обращают внимание — это просто фоновые детали.
  Но, тем не менее, они весьма интересны. Первая деталь. В одной из сцен 1-й серии
  Штирлиц крутит настройки приёмника и находит лёгенький джаз. Голос за кадром вещает: «Лондон
  передавал веселую музыку. Оркестр американца Гленна Миллера играл композицию из «Серенады солнечной долины».`

	if ReadingTime(strings.NewReader(text)) != 18 {
		t.Fatal("Invalid reading time")
	}

	text = `Date and time formats cause a lot of confusion and interoperability
   problems on the Internet.  This document addresses many of the
   problems encountered and makes recommendations to improve consistency
   and interoperability when representing and using date and time in
   Internet protocols.

   This document includes an Internet profile of the ISO 8601 [ISO8601]
   standard for representation of dates and times using the Gregorian
   calendar.

   There are many ways in which date and time values might appear in
   Internet protocols:  this document focuses on just one common usage,
   viz. timestamps for Internet protocol events.  This limited
   consideration has the following consequences:

   o  All dates and times are assumed to be in the "current era",
      somewhere between 0000AD and 9999AD.

   o  All times expressed have a stated relationship (offset) to
      Coordinated Universal Time (UTC).  (This is distinct from some
      usage in scheduling applications where a local time and location
      may be known, but the actual relationship to UTC may be dependent
      on the unknown or unknowable actions of politicians or
      administrators.  The UTC time corresponding to 17:00 on 23rd March
      2005 in New York may depend on administrative decisions about
      daylight savings time.  This specification steers well clear of
      such considerations.)

   o  Timestamps can express times that occurred before the introduction
      of UTC.  Such timestamps are expressed relative to universal time,
      using the best available practice at the stated time.

   o  Date and time expressions indicate an instant in time.
      Description of time periods, or intervals, is not covered here.

2. Definitions

   The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
   "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
   document are to be interpreted as described in RFC 2119 [RFC2119].`

	if ReadingTime(strings.NewReader(text)) != 80 {
		t.Fatal("Invalid reading time")
	}

	text = `В блистательном сериале «Семнадцать мгновений весны»
  есть два примечательных момента, на
  которые редко обращают внимание — это просто фоновые детали.
  Но, тем не менее, они весьма интересны. Первая деталь. В одной
  из сцен 1-й серии Штирлиц крутит настройки приёмника и находит
  лёгенький джаз`

	if ReadingTime(strings.NewReader(text)) != 12 {
		t.Fatal("Invalid reading time")
	}

	text = ``

	if ReadingTime(strings.NewReader(text)) != 1 {
		t.Fatal("Invalid reading time")
	}
}
