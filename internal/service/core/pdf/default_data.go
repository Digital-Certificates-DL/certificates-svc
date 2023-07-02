package pdf

var DefaultTemplateNormal = PDF{
	High:  595,
	Width: 842,
	Name: Field{
		X:        200,
		Y:        217,
		FontSize: 28,
		Font:     "semibold",
	},
	Course: Field{
		X:        61,
		Y:        259,
		FontSize: 14,
		Font:     "semibold",
	},
	Credits: Field{
		X:        70,
		Y:        56,
		FontSize: 12,
		Font:     "regular",
	},
	Points: Field{
		X:        70,
		Y:        79,
		FontSize: 12,
		Font:     "regular",
	},
	SerialNumber: Field{
		X:        641,
		Y:        56,
		FontSize: 12,
		Font:     "regular",
	},
	Date: Field{
		X:        641,
		Y:        79,
		FontSize: 12,
		Font:     "regular",
	},
	QR: Field{
		X:     658,
		Y:     106,
		High:  114,
		Width: 114,
	},
	Exam: Field{
		X:        300,
		Y:        300,
		FontSize: 15,
		Font:     "italic",
	},
	Level: Field{
		X:        300,
		Y:        277,
		FontSize: 14,
		Font:     "semibold",
	},
}

var DefaultTemplateTall = PDF{
	High:  1190,
	Width: 1684,
	Name: Field{
		Y:        434,
		FontSize: 56,
		Font:     "semibold",
	},
	Course: Field{
		Y:        518,
		FontSize: 28,
		Font:     "semibold",
	},
	Credits: Field{ //todo get from front and save to db
		X:        140,
		Y:        112,
		FontSize: 24,
		Font:     "regular",
	},
	Points: Field{
		X:        140,
		Y:        158,
		FontSize: 24,
		Font:     "regular",
	},
	SerialNumber: Field{
		X:        1250,
		Y:        112,
		FontSize: 24,
		Font:     "regular",
	},
	Date: Field{
		X:        1282,
		Y:        158,
		FontSize: 24,
		Font:     "regular",
	},
	QR: Field{
		X:     1316,
		Y:     212,
		High:  228,
		Width: 228,
	},
	Exam: Field{
		Y:        600,
		FontSize: 30,
		Font:     "italic",
	},
	Level: Field{
		Y:        554,
		FontSize: 28,
		Font:     "semibold",
	},
}

var DefaultData = PDFData{
	Name:         "Test Name",
	Course:       "Blockchain and Distributed Systems",
	Credits:      " 99",
	Points:       "100",
	SerialNumber: "694d0f5a7afe6fbc99cb",
	Date:         "30.05.2018",
	QR:           nil,
	Exam:         "passed",
	Level:        "graduated with honors",
	Note:         "************************************************",
}
