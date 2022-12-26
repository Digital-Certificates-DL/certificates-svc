// Copyright 2016 - 2022 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.
//
// Package excelize providing a set of functions that allow you to write to and
// read from XLAM / XLSM / XLSX / XLTM / XLTX files. Supports reading and
// writing spreadsheet documents generated by Microsoft Excel™ 2007 and later.
// Supports complex components by high compatibility, and provided streaming
// API for generating or reading data from a worksheet with huge amounts of
// data. This library needs Go version 1.15 or later.

package excelize

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/nfp"
)

// languageInfo defined the required fields of localization support for number format.
type languageInfo struct {
	apFmt      string
	tags       []string
	localMonth func(t time.Time, abbr int) string
}

// numberFormat directly maps the number format parser runtime required
// fields.
type numberFormat struct {
	section                                        []nfp.Section
	t                                              time.Time
	sectionIdx                                     int
	date1904, isNumeric, hours, seconds            bool
	number                                         float64
	ap, localCode, result, value, valueSectionType string
}

var (
	// supportedTokenTypes list the supported number format token types currently.
	supportedTokenTypes = []string{
		nfp.TokenSubTypeLanguageInfo,
		nfp.TokenTypeColor,
		nfp.TokenTypeCurrencyLanguage,
		nfp.TokenTypeDateTimes,
		nfp.TokenTypeElapsedDateTimes,
		nfp.TokenTypeGeneral,
		nfp.TokenTypeLiteral,
		nfp.TokenTypeTextPlaceHolder,
		nfp.TokenTypeZeroPlaceHolder,
	}
	// supportedLanguageInfo directly maps the supported language ID and tags.
	supportedLanguageInfo = map[string]languageInfo{
		"36":   {tags: []string{"af"}, localMonth: localMonthsNameAfrikaans, apFmt: apFmtAfrikaans},
		"445":  {tags: []string{"bn-IN"}, localMonth: localMonthsNameBangla, apFmt: nfp.AmPm[0]},
		"4":    {tags: []string{"zh-Hans"}, localMonth: localMonthsNameChinese1, apFmt: nfp.AmPm[2]},
		"7804": {tags: []string{"zh"}, localMonth: localMonthsNameChinese1, apFmt: nfp.AmPm[2]},
		"804":  {tags: []string{"zh-CN"}, localMonth: localMonthsNameChinese1, apFmt: nfp.AmPm[2]},
		"1004": {tags: []string{"zh-SG"}, localMonth: localMonthsNameChinese2, apFmt: nfp.AmPm[2]},
		"7C04": {tags: []string{"zh-Hant"}, localMonth: localMonthsNameChinese3, apFmt: nfp.AmPm[2]},
		"C04":  {tags: []string{"zh-HK"}, localMonth: localMonthsNameChinese2, apFmt: nfp.AmPm[2]},
		"1404": {tags: []string{"zh-MO"}, localMonth: localMonthsNameChinese3, apFmt: nfp.AmPm[2]},
		"404":  {tags: []string{"zh-TW"}, localMonth: localMonthsNameChinese3, apFmt: nfp.AmPm[2]},
		"9":    {tags: []string{"en"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"1000": {tags: []string{
			"aa", "aa-DJ", "aa-ER", "aa-ER", "aa-NA", "agq", "agq-CM", "ak", "ak-GH", "sq-ML",
			"gsw-LI", "gsw-CH", "ar-TD", "ar-KM", "ar-DJ", "ar-ER", "ar-IL", "ar-MR", "ar-PS",
			"ar-SO", "ar-SS", "ar-SD", "ar-001", "ast", "ast-ES", "asa", "asa-TZ", "ksf", "ksf-CM",
			"bm", "bm-Latn-ML", "bas", "bas-CM", "bem", "bem-ZM", "bez", "bez-TZ", "byn", "byn-ER",
			"brx", "brx-IN", "ca-AD", "ca-FR", "ca-IT", "ceb", "ceb-Latn", "ceb-Latn-PH", "tzm-Latn-MA",
			"ccp", "ccp-Cakm", "ccp-Cakm-BD", "ccp-Cakm-IN", "ce-RU", "cgg", "cgg-UG", "cu-RU", "swc",
			"swc-CD", "kw", "ke-GB", "da-GL", "dua", "dua-CM", "nl-AW", "nl-BQ", "nl-CW", "nl-SX",
			"nl-SR", "dz", "ebu", "ebu-KE", "en-AS", "en-AI", "en-AG", "en-AT", "en-BS", "en-BB",
			"en-BE", "en-BM", "en-BW", "en-IO", "en-VG", "en-BI", "en-CM", "en-KY", "en-CX", "en-CC",
			"en-CK", "en-CY", "en-DK", "en-DM", "en-ER", "en-150", "en-FK", "en-FI", "en-FJ", "en-GM",
			"en-DE", "en-GH", "en-GI", "en-GD", "en-GU", "en-GG", "en-GY", "en-IM", "en-IL", "en-JE",
			"en-KE", "en-KI", "en-LS", "en-LR", "en-MO", "en-MG", "en-MW", "en-MT", "en-MH", "en-MU",
			"en-FM", "en-MS", "en-NA", "en-NR", "en-NL", "en-NG", "en-NU", "en-NF", "en-MP", "en-PK",
			"en-PW", "en-PG", "en-PN", "en-PR", "en-RW", "en-KN", "en-LC", "en-VC", "en-WS", "en-SC",
			"en-SL", "en-SX", "en-SI", "en-SB", "en-SS", "en-SH", "en-SD", "en-SZ", "en-SE", "en-CH",
			"en-TZ", "en-TK", "en-TO", "en-TC", "en-TV", "en-UG", "en-UM", "en-VI", "en-VU", "en-001",
			"en-ZM", "eo", "eo-001", "ee", "ee-GH", "ee-TG", "ewo", "ewo-CM", "fo-DK", "fr-DZ",
			"fr-BJ", "fr-BF", "fr-BI", "fr-CF", "fr-TD", "fr-KM", "fr-CG", "fr-DJ", "fr-GQ", "fr-GF",
			"fr-PF", "fr-GA", "fr-GP", "fr-GN", "fr-MG", "fr-MQ", "fr-MR", "fr-MU", "fr-YT", "fr-NC",
			"fr-NE", "fr-RW", "fr-BL", "fr-MF", "fr-PM", "fr-SC", "fr-SY", "fr-TG", "fr-TN", "fr-VU",
			"fr-WF", "fur", "fur-IT", "ff-Latn-BF", "ff-CM", "ff-Latn-CM", "ff-Latn-GM", "ff-Latn-GH",
			"ff-GN", "ff-Latn-GN", "ff-Latn-GW", "ff-Latn-LR", "ff-MR", "ff-Latn-MR", "ff-Latn-NE",
			"ff-Latn-SL", "lg", "lg-UG", "de-BE", "de-IT", "el-CY", "guz", "guz-KE", "ha-Latn-GH",
			"ha-Latn-NG", "ia-FR", "ia-001", "it-SM", "it-VA", "jv", "jv-Latn", "jv-Latn-ID", "dyo",
			"dyo-SN", "kea", "kea-CV", "kab", "kab-DZ", "kkj", "kkj-CM", "kln", "kln-KE", "kam",
			"kam-KE", "ks-Arab-IN", "ki", "ki-KE", "sw-TZ", "sw-UG", "ko-KP", "khq", "khq-ML", "ses",
			"ses-ML", "nmg", "nmq-CM", "ku-Arab-IR", "lkt", "lkt-US", "lag", "lag-TZ", "ln", "ln-AO",
			"ln-CF", "ln-CD", "nds", "nds-DE", "nds-NL", "lu", "lu-CD", "luo", "luo", "luo-KE", "luy",
			"luy-KE", "jmc", "jmc-TZ", "mgh", "mgh-MZ", "kde", "kde-TZ", "mg", "mg-MG", "gv", "gv-IM",
			"mas", "mas-KE", "mas-TZ", "mas-IR", "mer", "mer-KE", "mgo", "mgo-CM", "mfe", "mfe-MU",
			"mua", "mua-CM", "nqo", "nqo-GN", "nqa", "naq-NA", "nnh", "nnh-CM", "jgo", "jgo-CM",
			"lrc-IQ", "lrc-IR", "nd", "nd-ZW", "nb-SJ", "nus", "nus-SD", "nus-SS", "nyn", "nyn-UG",
			"om-KE", "os", "os-GE", "os-RU", "ps-PK", "fa-AF", "pt-AO", "pt-CV", "pt-GQ", "pt-GW",
			"pt-LU", "pt-MO", "pt-MZ", "pt-ST", "pt-CH", "pt-TL", "prg-001", "ksh", "ksh-DE", "rof",
			"rof-TZ", "rn", "rn-BI", "ru-BY", "ru-KZ", "ru-KG", "ru-UA", "rwk", "rwk-TZ", "ssy",
			"ssy-ER", "saq", "saq-KE", "sg", "sq-CF", "sbp", "sbp-TZ", "seh", "seh-MZ", "ksb", "ksb-TZ",
			"sn", "sn-Latn", "sn-Latn-ZW", "xog", "xog-UG", "so-DJ", "so-ET", "so-KE", "nr", "nr-ZA",
			"st-LS", "es-BZ", "es-BR", "es-PH", "zgh", "zgh-Tfng-MA", "zgh-Tfng", "ss", "ss-ZA",
			"ss-SZ", "sv-AX", "shi", "shi-Tfng", "shi-Tfng-MA", "shi-Latn", "shi-Latn-MA", "dav",
			"dav-KE", "ta-MY", "ta-SG", "twq", "twq-NE", "teo", "teo-KE", "teo-UG", "bo-IN", "tig",
			"tig-ER", "to", "to-TO", "tr-CY", "uz-Arab", "us-Arab-AF", "vai", "vai-Vaii",
			"vai-Vaii-LR", "vai-Latn-LR", "vai-Latn", "vo", "vo-001", "vun", "vun-TZ", "wae",
			"wae-CH", "wal", "wae-ET", "yav", "yav-CM", "yo-BJ", "dje", "dje-NE",
		}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"C09":  {tags: []string{"en-AU"}, localMonth: localMonthsNameEnglish, apFmt: strings.ToLower(nfp.AmPm[0])},
		"2829": {tags: []string{"en-BZ"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"1009": {tags: []string{"en-CA"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"2409": {tags: []string{"en-029"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"3C09": {tags: []string{"en-HK"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"4009": {tags: []string{"en-IN"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"1809": {tags: []string{"en-IE"}, localMonth: localMonthsNameEnglish, apFmt: strings.ToLower(nfp.AmPm[0])},
		"2009": {tags: []string{"en-JM"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"4409": {tags: []string{"en-MY"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"1409": {tags: []string{"en-NZ"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"3409": {tags: []string{"en-PH"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"4809": {tags: []string{"en-SG"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"1C09": {tags: []string{"en-ZA"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"2C09": {tags: []string{"en-TT"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"4C09": {tags: []string{"en-AE"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"809":  {tags: []string{"en-GB"}, localMonth: localMonthsNameEnglish, apFmt: strings.ToLower(nfp.AmPm[0])},
		"409":  {tags: []string{"en-US"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"3009": {tags: []string{"en-ZW"}, localMonth: localMonthsNameEnglish, apFmt: nfp.AmPm[0]},
		"C":    {tags: []string{"fr"}, localMonth: localMonthsNameFrench, apFmt: nfp.AmPm[0]},
		"7":    {tags: []string{"de"}, localMonth: localMonthsNameGerman, apFmt: nfp.AmPm[0]},
		"C07":  {tags: []string{"de-AT"}, localMonth: localMonthsNameAustria, apFmt: nfp.AmPm[0]},
		"407":  {tags: []string{"de-DE"}, localMonth: localMonthsNameGerman, apFmt: nfp.AmPm[0]},
		"3C":   {tags: []string{"ga"}, localMonth: localMonthsNameIrish, apFmt: apFmtIrish},
		"83C":  {tags: []string{"ga-IE"}, localMonth: localMonthsNameIrish, apFmt: apFmtIrish},
		"10":   {tags: []string{"it"}, localMonth: localMonthsNameItalian, apFmt: nfp.AmPm[0]},
		"11":   {tags: []string{"ja"}, localMonth: localMonthsNameChinese3, apFmt: apFmtJapanese},
		"411":  {tags: []string{"ja-JP"}, localMonth: localMonthsNameChinese3, apFmt: apFmtJapanese},
		"12":   {tags: []string{"ko"}, localMonth: localMonthsNameKorean, apFmt: apFmtKorean},
		"412":  {tags: []string{"ko-KR"}, localMonth: localMonthsNameKorean, apFmt: apFmtKorean},
		"7C50": {tags: []string{"mn-Mong"}, localMonth: localMonthsNameTraditionalMongolian, apFmt: nfp.AmPm[0]},
		"850":  {tags: []string{"mn-Mong-CN"}, localMonth: localMonthsNameTraditionalMongolian, apFmt: nfp.AmPm[0]},
		"C50":  {tags: []string{"mn-Mong-MN"}, localMonth: localMonthsNameTraditionalMongolian, apFmt: nfp.AmPm[0]},
		"19":   {tags: []string{"ru"}, localMonth: localMonthsNameRussian, apFmt: nfp.AmPm[0]},
		"819":  {tags: []string{"ru-MD"}, localMonth: localMonthsNameRussian, apFmt: nfp.AmPm[0]},
		"419":  {tags: []string{"ru-RU"}, localMonth: localMonthsNameRussian, apFmt: nfp.AmPm[0]},
		"A":    {tags: []string{"es"}, localMonth: localMonthsNameSpanish, apFmt: apFmtSpanish},
		"2C0A": {tags: []string{"es-AR"}, localMonth: localMonthsNameSpanish, apFmt: apFmtSpanish},
		"200A": {tags: []string{"es-VE"}, localMonth: localMonthsNameSpanish, apFmt: apFmtSpanish},
		"400A": {tags: []string{"es-BO"}, localMonth: localMonthsNameSpanish, apFmt: apFmtSpanish},
		"340A": {tags: []string{"es-CL"}, localMonth: localMonthsNameSpanish, apFmt: apFmtSpanish},
		"240A": {tags: []string{"es-CO"}, localMonth: localMonthsNameSpanish, apFmt: apFmtSpanish},
		"140A": {tags: []string{"es-CR"}, localMonth: localMonthsNameSpanish, apFmt: apFmtSpanish},
		"5C0A": {tags: []string{"es-CU"}, localMonth: localMonthsNameSpanish, apFmt: apFmtCuba},
		"1C0A": {tags: []string{"es-DO"}, localMonth: localMonthsNameSpanish, apFmt: apFmtSpanish},
		"300A": {tags: []string{"es-EC"}, localMonth: localMonthsNameSpanish, apFmt: apFmtSpanish},
		"440A": {tags: []string{"es-SV"}, localMonth: localMonthsNameSpanish, apFmt: apFmtSpanish},
		"1E":   {tags: []string{"th"}, localMonth: localMonthsNameThai, apFmt: nfp.AmPm[0]},
		"41E":  {tags: []string{"th-TH"}, localMonth: localMonthsNameThai, apFmt: nfp.AmPm[0]},
		"51":   {tags: []string{"bo"}, localMonth: localMonthsNameTibetan, apFmt: apFmtTibetan},
		"451":  {tags: []string{"bo-CN"}, localMonth: localMonthsNameTibetan, apFmt: apFmtTibetan},
		"1F":   {tags: []string{"tr"}, localMonth: localMonthsNameTurkish, apFmt: apFmtTurkish},
		"41F":  {tags: []string{"tr-TR"}, localMonth: localMonthsNameTurkish, apFmt: apFmtTurkish},
		"52":   {tags: []string{"cy"}, localMonth: localMonthsNameWelsh, apFmt: apFmtWelsh},
		"452":  {tags: []string{"cy-GB"}, localMonth: localMonthsNameWelsh, apFmt: apFmtWelsh},
		"2A":   {tags: []string{"vi"}, localMonth: localMonthsNameVietnamese, apFmt: apFmtVietnamese},
		"42A":  {tags: []string{"vi-VN"}, localMonth: localMonthsNameVietnamese, apFmt: apFmtVietnamese},
		"88":   {tags: []string{"wo"}, localMonth: localMonthsNameWolof, apFmt: apFmtWolof},
		"488":  {tags: []string{"wo-SN"}, localMonth: localMonthsNameWolof, apFmt: apFmtWolof},
		"34":   {tags: []string{"xh"}, localMonth: localMonthsNameXhosa, apFmt: nfp.AmPm[0]},
		"434":  {tags: []string{"xh-ZA"}, localMonth: localMonthsNameXhosa, apFmt: nfp.AmPm[0]},
		"78":   {tags: []string{"ii"}, localMonth: localMonthsNameYi, apFmt: apFmtYi},
		"478":  {tags: []string{"ii-CN"}, localMonth: localMonthsNameYi, apFmt: apFmtYi},
		"35":   {tags: []string{"zu"}, localMonth: localMonthsNameZulu, apFmt: nfp.AmPm[0]},
		"435":  {tags: []string{"zu-ZA"}, localMonth: localMonthsNameZulu, apFmt: nfp.AmPm[0]},
	}
	// monthNamesBangla list the month names in the Bangla.
	monthNamesBangla = []string{
		"\u099C\u09BE\u09A8\u09C1\u09AF\u09BC\u09BE\u09B0\u09C0",
		"\u09AB\u09C7\u09AC\u09CD\u09B0\u09C1\u09AF\u09BC\u09BE\u09B0\u09C0",
		"\u09AE\u09BE\u09B0\u09CD\u099A",
		"\u098F\u09AA\u09CD\u09B0\u09BF\u09B2",
		"\u09AE\u09C7",
		"\u099C\u09C1\u09A8",
		"\u099C\u09C1\u09B2\u09BE\u0987",
		"\u0986\u0997\u09B8\u09CD\u099F",
		"\u09B8\u09C7\u09AA\u09CD\u099F\u09C7\u09AE\u09CD\u09AC\u09B0",
		"\u0985\u0995\u09CD\u099F\u09CB\u09AC\u09B0",
		"\u09A8\u09AD\u09C7\u09AE\u09CD\u09AC\u09B0",
		"\u09A1\u09BF\u09B8\u09C7\u09AE\u09CD\u09AC\u09B0",
	}
	// monthNamesAfrikaans list the month names in the Afrikaans.
	monthNamesAfrikaans = []string{"Januarie", "Februarie", "Maart", "April", "Mei", "Junie", "Julie", "Augustus", "September", "Oktober", "November", "Desember"}
	// monthNamesChinese list the month names in the Chinese.
	monthNamesChinese = []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二"}
	// monthNamesFrench list the month names in the French.
	monthNamesFrench = []string{"janvier", "février", "mars", "avril", "mai", "juin", "juillet", "août", "septembre", "octobre", "novembre", "décembre"}
	// monthNamesGerman list the month names in the German.
	monthNamesGerman = []string{"Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"}
	// monthNamesAustria list the month names in the Austria.
	monthNamesAustria = []string{"Jänner", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"}
	// monthNamesIrish list the month names in the Irish.
	monthNamesIrish = []string{"Eanáir", "Feabhra", "Márta", "Aibreán", "Bealtaine", "Meitheamh", "Iúil", "Lúnasa", "Meán Fómhair", "Deireadh Fómhair", "Samhain", "Nollaig"}
	// monthNamesItalian list the month names in the Italian.
	monthNamesItalian = []string{"gennaio", "febbraio", "marzo", "aprile", "maggio", "giugno", "luglio", "agosto", "settembre", "ottobre", "novembre", "dicembre"}
	// monthNamesRussian list the month names in the Russian.
	monthNamesRussian = []string{"январь", "февраль", "март", "апрель", "май", "июнь", "июль", "август", "сентябрь", "октябрь", "ноябрь", "декабрь"}
	// monthNamesSpanish list the month names in the Spanish.
	monthNamesSpanish = []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
	// monthNamesThai list the month names in the Thai.
	monthNamesThai = []string{
		"\u0e21\u0e01\u0e23\u0e32\u0e04\u0e21",
		"\u0e01\u0e38\u0e21\u0e20\u0e32\u0e1e\u0e31\u0e19\u0e18\u0e4c",
		"\u0e21\u0e35\u0e19\u0e32\u0e04\u0e21",
		"\u0e40\u0e21\u0e29\u0e32\u0e22\u0e19",
		"\u0e1e\u0e24\u0e29\u0e20\u0e32\u0e04\u0e21",
		"\u0e21\u0e34\u0e16\u0e38\u0e19\u0e32\u0e22\u0e19",
		"\u0e01\u0e23\u0e01\u0e0e\u0e32\u0e04\u0e21",
		"\u0e2a\u0e34\u0e07\u0e2b\u0e32\u0e04\u0e21",
		"\u0e01\u0e31\u0e19\u0e22\u0e32\u0e22\u0e19",
		"\u0e15\u0e38\u0e25\u0e32\u0e04\u0e21",
		"\u0e1e\u0e24\u0e28\u0e08\u0e34\u0e01\u0e32\u0e22\u0e19",
		"\u0e18\u0e31\u0e19\u0e27\u0e32\u0e04\u0e21",
	}
	// monthNamesTibetan list the month names in the Tibetan.
	monthNamesTibetan = []string{
		"\u0f5f\u0fb3\u0f0b\u0f56\u0f0b\u0f51\u0f44\u0f0b\u0f54\u0f7c\u0f0b",
		"\u0f5f\u0fb3\u0f0b\u0f56\u0f0b\u0f42\u0f49\u0f72\u0f66\u0f0b\u0f54\u0f0b",
		"\u0f5f\u0fb3\u0f0b\u0f56\u0f0b\u0f42\u0f66\u0f74\u0f58\u0f0b\u0f54\u0f0b",
		"\u0f5f\u0fb3\u0f0b\u0f56\u0f0b\u0f56\u0f5e\u0f72\u0f0b\u0f54\u0f0b",
		"\u0f5f\u0fb3\u0f0b\u0f56\u0f0b\u0f63\u0f94\u0f0b\u0f54\u0f0b",
		"\u0f5f\u0fb3\u0f0b\u0f56\u0f0b\u0f51\u0fb2\u0f74\u0f42\u0f0b\u0f54\u0f0b",
		"\u0f5f\u0fb3\u0f0b\u0f56\u0f0b\u0f56\u0f51\u0f74\u0f53\u0f0b\u0f54\u0f0b",
		"\u0f5f\u0fb3\u0f0b\u0f56\u0f0b\u0f56\u0f62\u0f92\u0fb1\u0f51\u0f0b\u0f54\u0f0b",
		"\u0f5f\u0fb3\u0f0b\u0f56\u0f0b\u0f51\u0f42\u0f74\u0f0b\u0f54\u0f0b",
		"\u0f66\u0fa4\u0fb1\u0f72\u0f0b\u0f5f\u0fb3\u0f0b\u0f56\u0f45\u0f74\u0f0b\u0f54\u0f0d",
		"\u0f5f\u0fb3\u0f0b\u0f56\u0f0b\u0f56\u0f45\u0f74\u0f0b\u0f42\u0f45\u0f72\u0f42\u0f0b\u0f54\u0f0b",
		"\u0f5f\u0fb3\u0f0b\u0f56\u0f0b\u0f56\u0f45\u0f74\u0f0b\u0f42\u0f49\u0f72\u0f66\u0f0b\u0f54\u0f0b",
	}
	// monthNamesTurkish list the month names in the Turkish.
	monthNamesTurkish = []string{"Ocak", "Şubat", "Mart", "Nisan", "Mayıs", "Haziran", "Temmuz", "Ağustos", "Eylül", "Ekim", "Kasım", "Aralık"}
	// monthNamesWelsh list the month names in the Welsh.
	monthNamesWelsh = []string{"Ionawr", "Chwefror", "Mawrth", "Ebrill", "Mai", "Mehefin", "Gorffennaf", "Awst", "Medi", "Hydref", "Tachwedd", "Rhagfyr"}
	// monthNamesWolof list the month names in the Wolof.
	monthNamesWolof = []string{"Samwiye", "Fewriye", "Maars", "Awril", "Me", "Suwe", "Sullet", "Ut", "Septàmbar", "Oktoobar", "Noowàmbar", "Desàmbar"}
	// monthNamesXhosa list the month names in the Xhosa.
	monthNamesXhosa = []string{"Januwari", "Febuwari", "Matshi", "Aprili", "Meyi", "Juni", "Julayi", "Agasti", "Septemba", "Oktobha", "Novemba", "Disemba"}
	// monthNamesYi list the month names in the Yi.
	monthNamesYi = []string{"\ua2cd", "\ua44d", "\ua315", "\ua1d6", "\ua26c", "\ua0d8", "\ua3c3", "\ua246", "\ua22c", "\ua2b0", "\ua2b0\ua2aa", "\ua2b0\ua44b"}
	// monthNamesZulu list the month names in the Zulu.
	monthNamesZulu = []string{"Januwari", "Febhuwari", "Mashi", "Ephreli", "Meyi", "Juni", "Julayi", "Agasti", "Septemba", "Okthoba", "Novemba", "Disemba"}
	// apFmtAfrikaans defined the AM/PM name in the Afrikaans.
	apFmtAfrikaans = "vm./nm."
	// apFmtCuba defined the AM/PM name in the Cuba.
	apFmtCuba = "a.m./p.m."
	// apFmtIrish defined the AM/PM name in the Irish.
	apFmtIrish = "r.n./i.n."
	// apFmtJapanese defined the AM/PM name in the Japanese.
	apFmtJapanese = "午前/午後"
	// apFmtKorean defined the AM/PM name in the Korean.
	apFmtKorean = "오전/오후"
	// apFmtSpanish defined the AM/PM name in the Spanish.
	apFmtSpanish = "a. m./p. m."
	// apFmtTibetan defined the AM/PM name in the Tibetan.
	apFmtTibetan = "\u0f66\u0f94\u0f0b\u0f51\u0fb2\u0f7c\u0f0b/\u0f55\u0fb1\u0f72\u0f0b\u0f51\u0fb2\u0f7c\u0f0b"
	// apFmtTurkish defined the AM/PM name in the Turkish.
	apFmtTurkish = "\u00F6\u00F6/\u00F6\u0053"
	// apFmtVietnamese defined the AM/PM name in the Vietnamese.
	apFmtVietnamese = "SA/CH"
	// apFmtWolof defined the AM/PM name in the Wolof.
	apFmtWolof = "Sub/Ngo"
	// apFmtYi defined the AM/PM name in the Yi.
	apFmtYi = "\ua3b8\ua111/\ua06f\ua2d2"
	// apFmtWelsh defined the AM/PM name in the Welsh.
	apFmtWelsh = "yb/yh"
)

// prepareNumberic split the number into two before and after parts by a
// decimal point.
func (nf *numberFormat) prepareNumberic(value string) {
	if nf.isNumeric, _ = isNumeric(value); !nf.isNumeric {
		return
	}
}

// format provides a function to return a string parse by number format
// expression. If the given number format is not supported, this will return
// the original cell value.
func format(value, numFmt string, date1904 bool) string {
	p := nfp.NumberFormatParser()
	nf := numberFormat{section: p.Parse(numFmt), value: value, date1904: date1904}
	nf.number, nf.valueSectionType = nf.getValueSectionType(value)
	nf.prepareNumberic(value)
	for i, section := range nf.section {
		nf.sectionIdx = i
		if section.Type != nf.valueSectionType {
			continue
		}
		if nf.isNumeric {
			switch section.Type {
			case nfp.TokenSectionPositive:
				return nf.positiveHandler()
			case nfp.TokenSectionNegative:
				return nf.negativeHandler()
			default:
				return nf.zeroHandler()
			}
		}
		return nf.textHandler()
	}
	return value
}

// positiveHandler will be handling positive selection for a number format
// expression.
func (nf *numberFormat) positiveHandler() (result string) {
	nf.t, nf.hours, nf.seconds = timeFromExcelTime(nf.number, nf.date1904), false, false
	for i, token := range nf.section[nf.sectionIdx].Items {
		if inStrSlice(supportedTokenTypes, token.TType, true) == -1 || token.TType == nfp.TokenTypeGeneral {
			result = nf.value
			return
		}
		if token.TType == nfp.TokenTypeCurrencyLanguage {
			if err := nf.currencyLanguageHandler(i, token); err != nil {
				result = nf.value
				return
			}
		}
		if token.TType == nfp.TokenTypeDateTimes {
			nf.dateTimesHandler(i, token)
		}
		if token.TType == nfp.TokenTypeElapsedDateTimes {
			nf.elapsedDateTimesHandler(token)
		}
		if token.TType == nfp.TokenTypeLiteral {
			nf.result += token.TValue
			continue
		}
		if token.TType == nfp.TokenTypeZeroPlaceHolder && token.TValue == strings.Repeat("0", len(token.TValue)) {
			if isNum, precision := isNumeric(nf.value); isNum {
				if nf.number < 1 {
					nf.result += "0"
					continue
				}
				if precision > 15 {
					nf.result += roundPrecision(nf.value, 15)
				} else {
					nf.result += fmt.Sprintf("%.f", nf.number)
				}
				continue
			}
		}
	}
	result = nf.result
	return
}

// currencyLanguageHandler will be handling currency and language types tokens for a number
// format expression.
func (nf *numberFormat) currencyLanguageHandler(i int, token nfp.Token) (err error) {
	for _, part := range token.Parts {
		if inStrSlice(supportedTokenTypes, part.Token.TType, true) == -1 {
			err = ErrUnsupportedNumberFormat
			return
		}
		if _, ok := supportedLanguageInfo[strings.ToUpper(part.Token.TValue)]; !ok {
			err = ErrUnsupportedNumberFormat
			return
		}
		nf.localCode = strings.ToUpper(part.Token.TValue)
	}
	return
}

// localAmPm return AM/PM name by supported language ID.
func (nf *numberFormat) localAmPm(ap string) string {
	if languageInfo, ok := supportedLanguageInfo[nf.localCode]; ok {
		return languageInfo.apFmt
	}
	return ap
}

// localMonthsNameEnglish returns the English name of the month.
func localMonthsNameEnglish(t time.Time, abbr int) string {
	if abbr == 3 {
		return t.Month().String()[:3]
	}
	if abbr == 4 {
		return t.Month().String()
	}
	return t.Month().String()[:1]
}

// localMonthsNameAfrikaans returns the Afrikaans name of the month.
func localMonthsNameAfrikaans(t time.Time, abbr int) string {
	if abbr == 3 {
		month := monthNamesAfrikaans[int(t.Month())-1]
		if len([]rune(month)) <= 3 {
			return month
		}
		return string([]rune(month)[:3]) + "."
	}
	if abbr == 4 {
		return monthNamesAfrikaans[int(t.Month())-1]
	}
	return monthNamesAfrikaans[int(t.Month())-1][:1]
}

// localMonthsNameAustria returns the Austria name of the month.
func localMonthsNameAustria(t time.Time, abbr int) string {
	if abbr == 3 {
		return string([]rune(monthNamesAustria[int(t.Month())-1])[:3])
	}
	if abbr == 4 {
		return monthNamesAustria[int(t.Month())-1]
	}
	return monthNamesAustria[int(t.Month())-1][:1]
}

// localMonthsNameBangla returns the German name of the month.
func localMonthsNameBangla(t time.Time, abbr int) string {
	if abbr == 3 || abbr == 4 {
		return monthNamesBangla[int(t.Month())-1]
	}
	return string([]rune(monthNamesBangla[int(t.Month())-1])[:1])
}

// localMonthsNameFrench returns the French name of the month.
func localMonthsNameFrench(t time.Time, abbr int) string {
	if abbr == 3 {
		month := monthNamesFrench[int(t.Month())-1]
		if len([]rune(month)) <= 4 {
			return month
		}
		return string([]rune(month)[:4]) + "."
	}
	if abbr == 4 {
		return monthNamesFrench[int(t.Month())-1]
	}
	return monthNamesFrench[int(t.Month())-1][:1]
}

// localMonthsNameIrish returns the Irish name of the month.
func localMonthsNameIrish(t time.Time, abbr int) string {
	if abbr == 3 {
		switch int(t.Month()) {
		case 1, 4, 8:
			return string([]rune(monthNamesIrish[int(t.Month())-1])[:3])
		case 2, 3, 6:
			return string([]rune(monthNamesIrish[int(t.Month())-1])[:5])
		case 9, 10:
			return string([]rune(monthNamesIrish[int(t.Month())-1])[:1]) + "Fómh"
		default:
			return string([]rune(monthNamesIrish[int(t.Month())-1])[:4])
		}
	}
	if abbr == 4 {
		return monthNamesIrish[int(t.Month())-1]
	}
	return string([]rune(monthNamesIrish[int(t.Month())-1])[:1])
}

// localMonthsNameItalian returns the Italian name of the month.
func localMonthsNameItalian(t time.Time, abbr int) string {
	if abbr == 3 {
		return monthNamesItalian[int(t.Month())-1][:3]
	}
	if abbr == 4 {
		return monthNamesItalian[int(t.Month())-1]
	}
	return monthNamesItalian[int(t.Month())-1][:1]
}

// localMonthsNameGerman returns the German name of the month.
func localMonthsNameGerman(t time.Time, abbr int) string {
	if abbr == 3 {
		return string([]rune(monthNamesGerman[int(t.Month())-1])[:3])
	}
	if abbr == 4 {
		return monthNamesGerman[int(t.Month())-1]
	}
	return string([]rune(monthNamesGerman[int(t.Month())-1])[:1])
}

// localMonthsNameChinese1 returns the Chinese name of the month.
func localMonthsNameChinese1(t time.Time, abbr int) string {
	if abbr == 3 {
		return strconv.Itoa(int(t.Month())) + "月"
	}
	if abbr == 4 {
		return monthNamesChinese[int(t.Month())-1] + "月"
	}
	return monthNamesChinese[int(t.Month())-1]
}

// localMonthsNameChinese2 returns the Chinese name of the month.
func localMonthsNameChinese2(t time.Time, abbr int) string {
	if abbr == 3 || abbr == 4 {
		return monthNamesChinese[int(t.Month())-1] + "月"
	}
	return monthNamesChinese[int(t.Month())-1]
}

// localMonthsNameChinese3 returns the Chinese name of the month.
func localMonthsNameChinese3(t time.Time, abbr int) string {
	if abbr == 3 || abbr == 4 {
		return strconv.Itoa(int(t.Month())) + "月"
	}
	return strconv.Itoa(int(t.Month()))
}

// localMonthsNameKorean returns the Korean name of the month.
func localMonthsNameKorean(t time.Time, abbr int) string {
	if abbr == 3 || abbr == 4 {
		return strconv.Itoa(int(t.Month())) + "월"
	}
	return strconv.Itoa(int(t.Month()))
}

// localMonthsNameTraditionalMongolian returns the Traditional Mongolian name of the month.
func localMonthsNameTraditionalMongolian(t time.Time, abbr int) string {
	if abbr == 5 {
		return "M"
	}
	return fmt.Sprintf("M%02d", int(t.Month()))
}

// localMonthsNameRussian returns the Russian name of the month.
func localMonthsNameRussian(t time.Time, abbr int) string {
	if abbr == 3 {
		month := monthNamesRussian[int(t.Month())-1]
		if len([]rune(month)) <= 4 {
			return month
		}
		return string([]rune(month)[:3]) + "."
	}
	if abbr == 4 {
		return monthNamesRussian[int(t.Month())-1]
	}
	return string([]rune(monthNamesRussian[int(t.Month())-1])[:1])
}

// localMonthsNameSpanish returns the Spanish name of the month.
func localMonthsNameSpanish(t time.Time, abbr int) string {
	if abbr == 3 {
		return monthNamesSpanish[int(t.Month())-1][:3]
	}
	if abbr == 4 {
		return monthNamesSpanish[int(t.Month())-1]
	}
	return monthNamesSpanish[int(t.Month())-1][:1]
}

// localMonthsNameThai returns the Thai name of the month.
func localMonthsNameThai(t time.Time, abbr int) string {
	if abbr == 3 {
		r := []rune(monthNamesThai[int(t.Month())-1])
		return string(r[:1]) + "." + string(r[len(r)-2:len(r)-1]) + "."
	}
	if abbr == 4 {
		return monthNamesThai[int(t.Month())-1]
	}
	return string([]rune(monthNamesThai[int(t.Month())-1])[:1])
}

// localMonthsNameTibetan returns the Tibetan name of the month.
func localMonthsNameTibetan(t time.Time, abbr int) string {
	if abbr == 3 {
		return "\u0f5f\u0fb3\u0f0b" + []string{"\u0f21", "\u0f22", "\u0f23", "\u0f24", "\u0f25", "\u0f26", "\u0f27", "\u0f28", "\u0f29", "\u0f21\u0f20", "\u0f21\u0f21", "\u0f21\u0f22"}[int(t.Month())-1]
	}
	if abbr == 5 {
		if t.Month() == 10 {
			return "\u0f66"
		}
		return "\u0f5f"
	}
	return monthNamesTibetan[int(t.Month())-1]
}

// localMonthsNameTurkish returns the Turkish name of the month.
func localMonthsNameTurkish(t time.Time, abbr int) string {
	if abbr == 3 {
		return string([]rune(monthNamesTurkish[int(t.Month())-1])[:3])
	}
	if abbr == 4 {
		return monthNamesTurkish[int(t.Month())-1]
	}
	return string([]rune(monthNamesTurkish[int(t.Month())-1])[:1])
}

// localMonthsNameWelsh returns the Welsh name of the month.
func localMonthsNameWelsh(t time.Time, abbr int) string {
	if abbr == 3 {
		switch int(t.Month()) {
		case 2, 7:
			return string([]rune(monthNamesWelsh[int(t.Month())-1])[:5])
		case 8, 9, 11, 12:
			return string([]rune(monthNamesWelsh[int(t.Month())-1])[:4])
		default:
			return string([]rune(monthNamesWelsh[int(t.Month())-1])[:3])
		}
	}
	if abbr == 4 {
		return monthNamesWelsh[int(t.Month())-1]
	}
	return string([]rune(monthNamesWelsh[int(t.Month())-1])[:1])
}

// localMonthsNameVietnamese returns the Vietnamese name of the month.
func localMonthsNameVietnamese(t time.Time, abbr int) string {
	if abbr == 3 {
		return "Thg " + strconv.Itoa(int(t.Month()))
	}
	if abbr == 5 {
		return "T " + strconv.Itoa(int(t.Month()))
	}
	return "Tháng " + strconv.Itoa(int(t.Month()))
}

// localMonthsNameWolof returns the Wolof name of the month.
func localMonthsNameWolof(t time.Time, abbr int) string {
	if abbr == 3 {
		switch int(t.Month()) {
		case 3, 6:
			return string([]rune(monthNamesWolof[int(t.Month())-1])[:3])
		case 5, 8:
			return string([]rune(monthNamesWolof[int(t.Month())-1])[:2])
		case 9:
			return string([]rune(monthNamesWolof[int(t.Month())-1])[:4]) + "."
		case 11:
			return "Now."
		default:
			return string([]rune(monthNamesWolof[int(t.Month())-1])[:3]) + "."
		}
	}
	if abbr == 4 {
		return monthNamesWolof[int(t.Month())-1]
	}
	return string([]rune(monthNamesWolof[int(t.Month())-1])[:1])
}

// localMonthsNameXhosa returns the Xhosa name of the month.
func localMonthsNameXhosa(t time.Time, abbr int) string {
	if abbr == 3 {
		switch int(t.Month()) {
		case 4:
			return "uEpr."
		case 8:
			return "u" + string([]rune(monthNamesXhosa[int(t.Month())-1])[:2]) + "."
		default:
			return "u" + string([]rune(monthNamesXhosa[int(t.Month())-1])[:3]) + "."
		}
	}
	if abbr == 4 {
		return "u" + monthNamesXhosa[int(t.Month())-1]
	}
	return "u"
}

// localMonthsNameYi returns the Yi name of the month.
func localMonthsNameYi(t time.Time, abbr int) string {
	if abbr == 3 || abbr == 4 {
		return string(monthNamesYi[int(t.Month())-1]) + "\ua1aa"
	}
	return string([]rune(monthNamesYi[int(t.Month())-1])[:1])
}

// localMonthsNameZulu returns the Zulu name of the month.
func localMonthsNameZulu(t time.Time, abbr int) string {
	if abbr == 3 {
		if int(t.Month()) == 8 {
			return string([]rune(monthNamesZulu[int(t.Month())-1])[:4])
		}
		return string([]rune(monthNamesZulu[int(t.Month())-1])[:3])
	}
	if abbr == 4 {
		return monthNamesZulu[int(t.Month())-1]
	}
	return string([]rune(monthNamesZulu[int(t.Month())-1])[:1])
}

// localMonthName return months name by supported language ID.
func (nf *numberFormat) localMonthsName(abbr int) string {
	if languageInfo, ok := supportedLanguageInfo[nf.localCode]; ok {
		return languageInfo.localMonth(nf.t, abbr)
	}
	return localMonthsNameEnglish(nf.t, abbr)
}

// dateTimesHandler will be handling date and times types tokens for a number
// format expression.
func (nf *numberFormat) dateTimesHandler(i int, token nfp.Token) {
	if idx := inStrSlice(nfp.AmPm, strings.ToUpper(token.TValue), false); idx != -1 {
		if nf.ap == "" {
			nextHours := nf.hoursNext(i)
			aps := strings.Split(nf.localAmPm(token.TValue), "/")
			nf.ap = aps[0]
			if nextHours > 12 {
				nf.ap = aps[1]
			}
		}
		nf.result += nf.ap
		return
	}
	if strings.Contains(strings.ToUpper(token.TValue), "M") {
		l := len(token.TValue)
		if l == 1 && !nf.hours && !nf.secondsNext(i) {
			nf.result += strconv.Itoa(int(nf.t.Month()))
			return
		}
		if l == 2 && !nf.hours && !nf.secondsNext(i) {
			nf.result += fmt.Sprintf("%02d", int(nf.t.Month()))
			return
		}
		if l == 3 {
			nf.result += nf.localMonthsName(3)
			return
		}
		if l == 4 || l > 5 {
			nf.result += nf.localMonthsName(4)
			return
		}
		if l == 5 {
			nf.result += nf.localMonthsName(5)
			return
		}
	}
	nf.yearsHandler(i, token)
	nf.daysHandler(i, token)
	nf.hoursHandler(i, token)
	nf.minutesHandler(token)
	nf.secondsHandler(token)
}

// yearsHandler will be handling years in the date and times types tokens for a
// number format expression.
func (nf *numberFormat) yearsHandler(i int, token nfp.Token) {
	years := strings.Contains(strings.ToUpper(token.TValue), "Y")
	if years && len(token.TValue) <= 2 {
		nf.result += strconv.Itoa(nf.t.Year())[2:]
		return
	}
	if years && len(token.TValue) > 2 {
		nf.result += strconv.Itoa(nf.t.Year())
		return
	}
}

// daysHandler will be handling days in the date and times types tokens for a
// number format expression.
func (nf *numberFormat) daysHandler(i int, token nfp.Token) {
	if strings.Contains(strings.ToUpper(token.TValue), "D") {
		switch len(token.TValue) {
		case 1:
			nf.result += strconv.Itoa(nf.t.Day())
			return
		case 2:
			nf.result += fmt.Sprintf("%02d", nf.t.Day())
			return
		case 3:
			nf.result += nf.t.Weekday().String()[:3]
			return
		default:
			nf.result += nf.t.Weekday().String()
			return
		}
	}
}

// hoursHandler will be handling hours in the date and times types tokens for a
// number format expression.
func (nf *numberFormat) hoursHandler(i int, token nfp.Token) {
	nf.hours = strings.Contains(strings.ToUpper(token.TValue), "H")
	if nf.hours {
		h := nf.t.Hour()
		ap, ok := nf.apNext(i)
		if ok {
			nf.ap = ap[0]
			if h > 12 {
				h -= 12
				nf.ap = ap[1]
			}
		}
		if nf.ap != "" && nf.hoursNext(i) == -1 && h > 12 {
			h -= 12
		}
		switch len(token.TValue) {
		case 1:
			nf.result += strconv.Itoa(h)
			return
		default:
			nf.result += fmt.Sprintf("%02d", h)
			return
		}
	}
}

// minutesHandler will be handling minutes in the date and times types tokens
// for a number format expression.
func (nf *numberFormat) minutesHandler(token nfp.Token) {
	if strings.Contains(strings.ToUpper(token.TValue), "M") {
		nf.hours = false
		switch len(token.TValue) {
		case 1:
			nf.result += strconv.Itoa(nf.t.Minute())
			return
		default:
			nf.result += fmt.Sprintf("%02d", nf.t.Minute())
			return
		}
	}
}

// secondsHandler will be handling seconds in the date and times types tokens
// for a number format expression.
func (nf *numberFormat) secondsHandler(token nfp.Token) {
	nf.seconds = strings.Contains(strings.ToUpper(token.TValue), "S")
	if nf.seconds {
		switch len(token.TValue) {
		case 1:
			nf.result += strconv.Itoa(nf.t.Second())
			return
		default:
			nf.result += fmt.Sprintf("%02d", nf.t.Second())
			return
		}
	}
}

// elapsedDateTimesHandler will be handling elapsed date and times types tokens
// for a number format expression.
func (nf *numberFormat) elapsedDateTimesHandler(token nfp.Token) {
	if strings.Contains(strings.ToUpper(token.TValue), "H") {
		nf.result += fmt.Sprintf("%.f", nf.t.Sub(excel1900Epoc).Hours())
		return
	}
	if strings.Contains(strings.ToUpper(token.TValue), "M") {
		nf.result += fmt.Sprintf("%.f", nf.t.Sub(excel1900Epoc).Minutes())
		return
	}
	if strings.Contains(strings.ToUpper(token.TValue), "S") {
		nf.result += fmt.Sprintf("%.f", nf.t.Sub(excel1900Epoc).Seconds())
		return
	}
}

// hoursNext detects if a token of type hours exists after a given tokens list.
func (nf *numberFormat) hoursNext(i int) int {
	tokens := nf.section[nf.sectionIdx].Items
	for idx := i + 1; idx < len(tokens); idx++ {
		if tokens[idx].TType == nfp.TokenTypeDateTimes {
			if strings.Contains(strings.ToUpper(tokens[idx].TValue), "H") {
				t := timeFromExcelTime(nf.number, false)
				return t.Hour()
			}
		}
	}
	return -1
}

// apNext detects if a token of type AM/PM exists after a given tokens list.
func (nf *numberFormat) apNext(i int) ([]string, bool) {
	tokens := nf.section[nf.sectionIdx].Items
	for idx := i + 1; idx < len(tokens); idx++ {
		if tokens[idx].TType == nfp.TokenTypeDateTimes {
			if strings.Contains(strings.ToUpper(tokens[idx].TValue), "H") {
				return nil, false
			}
			if i := inStrSlice(nfp.AmPm, tokens[idx].TValue, false); i != -1 {
				return strings.Split(nf.localAmPm(tokens[idx].TValue), "/"), true
			}
		}
	}
	return nil, false
}

// secondsNext detects if a token of type seconds exists after a given tokens
// list.
func (nf *numberFormat) secondsNext(i int) bool {
	tokens := nf.section[nf.sectionIdx].Items
	for idx := i + 1; idx < len(tokens); idx++ {
		if tokens[idx].TType == nfp.TokenTypeDateTimes {
			return strings.Contains(strings.ToUpper(tokens[idx].TValue), "S")
		}
	}
	return false
}

// negativeHandler will be handling negative selection for a number format
// expression.
func (nf *numberFormat) negativeHandler() (result string) {
	for _, token := range nf.section[nf.sectionIdx].Items {
		if inStrSlice(supportedTokenTypes, token.TType, true) == -1 || token.TType == nfp.TokenTypeGeneral {
			result = nf.value
			return
		}
		if token.TType == nfp.TokenTypeLiteral {
			nf.result += token.TValue
			continue
		}
		if token.TType == nfp.TokenTypeZeroPlaceHolder && token.TValue == strings.Repeat("0", len(token.TValue)) {
			if isNum, precision := isNumeric(nf.value); isNum {
				if math.Abs(nf.number) < 1 {
					nf.result += "0"
					continue
				}
				if precision > 15 {
					nf.result += strings.TrimLeft(roundPrecision(nf.value, 15), "-")
				} else {
					nf.result += fmt.Sprintf("%.f", math.Abs(nf.number))
				}
				continue
			}
		}
	}
	result = nf.result
	return
}

// zeroHandler will be handling zero selection for a number format expression.
func (nf *numberFormat) zeroHandler() string {
	return nf.value
}

// textHandler will be handling text selection for a number format expression.
func (nf *numberFormat) textHandler() (result string) {
	for _, token := range nf.section[nf.sectionIdx].Items {
		if token.TType == nfp.TokenTypeLiteral {
			result += token.TValue
		}
		if token.TType == nfp.TokenTypeTextPlaceHolder {
			result += nf.value
		}
	}
	return result
}

// getValueSectionType returns its applicable number format expression section
// based on the given value.
func (nf *numberFormat) getValueSectionType(value string) (float64, string) {
	isNum, _ := isNumeric(value)
	if !isNum {
		return 0, nfp.TokenSectionText
	}
	number, _ := strconv.ParseFloat(value, 64)
	if number > 0 {
		return number, nfp.TokenSectionPositive
	}
	if number < 0 {
		return number, nfp.TokenSectionNegative
	}
	return number, nfp.TokenSectionZero
}
