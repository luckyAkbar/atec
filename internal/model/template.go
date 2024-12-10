package model

// Template is the ATEC template in which contains several known subtest's group
type Template struct {
	SubTest SubTest
}

// SubtestDetail each questionnaire groups along with its detail
type SubtestDetail struct {
	Name          string
	OptionCount   int
	QuestionCount int
}

// SubTest is each questionnaire group with its respective details
type SubTest map[int]SubtestDetail

// DefaultATECTemplate is the constant of default ATEC template.
// ATEC template itself is unlikely to change, thus using variable as constant
// will be faster than storing it in the database.
// Also if on the future will use different template,
// will be easier to refactor / implement new feature to handle it.
// each map key, represented by int is the subtest group id.
//
//nolint:mnd
var DefaultATECTemplate = Template{
	SubTest: SubTest{
		0: {
			Name:          "Speech/Language/Communication",
			OptionCount:   3,
			QuestionCount: 14,
		},
		1: {
			Name:          "Sociability",
			OptionCount:   3,
			QuestionCount: 20,
		},
		2: {
			Name:          "Sensory/Cognitive Awareness",
			OptionCount:   3,
			QuestionCount: 18,
		},
		3: {
			Name:          "Health/Physical/Behavior",
			OptionCount:   4,
			QuestionCount: 25,
		},
	},
}
