package templatemocks

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	m "github.com/expandorg/requester-service/pkg/model"
)

// Populate mock data
func Populate(d *mgo.Database) {
	err := taskTemplates(d)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Inserted test data")
	err = onboardingTemplates(d)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Inserted onboarding test data")
}

func taskTemplates(db *mgo.Database) error {
	t1 := &m.TaskTemplate{
		ID:    bson.NewObjectId(),
		Name:  "Region Bounding Box",
		Order: 1,
		Eligibility: &m.DraftEligibility{
			Module: m.EligibilityModuleAll,
		},
		Assignment: &m.DraftAssignment{
			Module: m.AssignmentModuleAll,
			Limit:  1,
			Repeat: true,
		},
		Funding: &m.DraftFunding{
			Module:      m.FundingModuleRequirement,
			Requirement: 0,
			Balance:     0,
		},
		Verification: &m.DraftVerification{
			Module:               m.VerificationModuleRequester,
			AgreementCount:       0,
			ScoreThreshold:       1,
			MinimumExecutionTime: 0,
		},
		Onboarding: &m.DraftOnboarding{
			Enabled: true,
			Steps:   make([]m.DraftOnboardingStep, 0),
		},
		Variables: []string{"image", "test2", "test3"},
		DataSample: map[string]string{
			"image": "https://expand.org/images/image-annotation.png",
			"test2": "Find logo",
			"test3": "sample variable",
		},
		TaskForm: unmarshalForm([]byte(`
			{
				"modules": [
					{
						"name": "regionSelect",
						"type": "regionSelect",
						"image": "$(image)"
					},
					{
						"name": "test",
						"type": "text",
						"style": "body",
						"content": "$(test2)"
					},					
					{
						"name": "submit",
						"caption": "Submit",
						"type": "submit"
					}
				]
			}	
	`)),
		VerificationForm: unmarshalForm([]byte(`
			{
				"modules": [
					{
						"name": "regionSelect",
						"type": "regionSelect",
						"image": "https://expand.org/images/image-annotation.png"
					},
					{
						"name": "submit",
						"caption": "Submit",
						"type": "submit"
					}
				]
			}	
	`)),
	}

	return db.C("taskTemplates").Insert(t1)
}

func onboardingTemplates(db *mgo.Database) error {
	o1 := &m.OnboardingTemplate{
		ID:      bson.NewObjectId(),
		Name:    "Region Bounding Box Onboarding",
		IsGroup: false,
		TaskForm: unmarshalForm([]byte(`
			{
				"modules": [
					{
						"name": "regionSelect",
						"type": "regionSelect",
						"image": "https://expand.org/images/image-annotation.png"
					},
					{
						"name": "submit",
						"caption": "Submit",
						"type": "submit"
					}
				]
			}	
		`)),
	}

	q1 := &m.OnboardingTemplate{
		ID:             bson.NewObjectId(),
		Name:           "quiz",
		IsGroup:        true,
		ScoreThreshold: 1,
		Retries:        3,
		FailureMessage: "Quiz failed",
		Data: &m.OnboardingGroupData{
			Answer: &m.OnboardingGroupDataAnswer{
				Field: "input",
			},
			Columns: []*m.OnboardingGroupDataVariable{
				&m.OnboardingGroupDataVariable{Name: "varname1", Type: "text"},
				&m.OnboardingGroupDataVariable{Name: "varname2", Type: "number"},
				&m.OnboardingGroupDataVariable{Name: "varname3", Type: "bool"},
			},
			Steps: []*m.OnboardingGroupDataStep{
				&m.OnboardingGroupDataStep{"1", []string{"Type your name", "1", ""}},
				&m.OnboardingGroupDataStep{"2", []string{"Your favorite book", "2", "true"}},
				&m.OnboardingGroupDataStep{"3", []string{"Your favorite color", "3", ""}},
				&m.OnboardingGroupDataStep{"4", []string{"Last question", "4", ""}},
			},
		},
		TaskForm: unmarshalForm([]byte(`
			{
				"modules": [
					{
						"name": "input",
						"type": "input",
						"inputType": "text",
						"placeholder": "$(varname1)..."
					},
					{
						"name": "test",
						"type": "text",
						"style": "body",
						"align": "left",
						"content": "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut"
					},
					{
						"name": "submit",
						"caption": "Submit",
						"type": "submit"
					}
				]
			}		
		`)),
	}

	q2 := &m.OnboardingTemplate{
		ID:             bson.NewObjectId(),
		Name:           "Cities quiz",
		IsGroup:        true,
		ScoreThreshold: 1,
		Retries:        3,
		FailureMessage: "Quiz failed",
		Data: &m.OnboardingGroupData{
			Answer: &m.OnboardingGroupDataAnswer{
				Field: "select-10b3",
			},
			Columns: []*m.OnboardingGroupDataVariable{
				&m.OnboardingGroupDataVariable{Name: "question", Type: "text"},
				&m.OnboardingGroupDataVariable{Name: "answer1", Type: "text"},
				&m.OnboardingGroupDataVariable{Name: "answer2", Type: "text"},
				&m.OnboardingGroupDataVariable{Name: "answer3", Type: "text"},
				&m.OnboardingGroupDataVariable{Name: "answer4", Type: "text"},
			},
			Steps: []*m.OnboardingGroupDataStep{
				&m.OnboardingGroupDataStep{"London", []string{"Select capital of Great Britain", "Washington", "London", "Moscow", "Edinburgh"}},
				&m.OnboardingGroupDataStep{"Washington", []string{"Select capital of USA", "Washington", "London", "Moscow", "Edinburgh"}},
				&m.OnboardingGroupDataStep{"Moscow", []string{"Select capital of Russia", "Washington", "London", "Moscow", "Edinburgh"}},
			},
		},
		TaskForm: unmarshalForm([]byte(`
			{
				"modules": [{
					"content": "<h1>Quiz</h1>",
					"name": "richText-b916",
					"type": "richText"
				}, {
					"content": "<p>$(question)</p>",
					"name": "richText-b0d3",
					"type": "richText"
				}, {
					"options": ["$(answer1)", "$(answer2)", "$(answer3)", "$(answer4)"],
					"type": "select",
					"columns": 2,
					"idType": "numerals",
					"name": "select-10b3"
				}, {
					"caption": "Answer",
					"name": "submit",
					"type": "submit"
				}]
			}		
		`)),
	}

	return db.C("onboardingTemplates").Insert(o1, q1, q2)
}

func unmarshalForm(data []byte) *m.Form {
	form := new(m.Form)
	err := json.Unmarshal(data, form)
	if err != nil {
		log.Fatalln(err)
	}
	return form
}
