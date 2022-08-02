package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	TestStruct struct {
		testName    string
		in          interface{}
		expectedErr error
	}

	User struct {
		ID     string   `json:"id" validate:"len:36"`
		Name   string   `json:"name"`
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
	}

	Principal struct {
		User   User     `validate:"nested"`
		Claims []string `validate:"in:api,portal,admin_panel"`
	}

	IntegerTestStruct struct {
		Minimal  int   `validate:"min:-18"`
		Maximal  int   `validate:"max:50"`
		SingleIn int   `validate:"in:1,2,3"`
		MultiIn  []int `validate:"in:1,2,3,4,5"`
	}

	StringTestStruct struct {
		MailRegexp string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		ExactLen   string   `validate:"len:0"`
		SingleIn   string   `validate:"in:alfa,beta,gamma"`
		MultiIn    []string `validate:"in:alfa,beta,gamma,zeta,omega"`
	}

	EmptyStruct struct{}
)

func getValidTestStructures() []TestStruct {
	return []TestStruct{
		{
			testName:    "nil to validate",
			in:          nil,
			expectedErr: nil,
		},
		{
			testName:    "empty struct to validate",
			in:          EmptyStruct{},
			expectedErr: nil,
		},
		{
			testName: "valid structure",
			in: User{
				ID:    "d5ee6174-d0b6-43e2-be47-f4dff983ff48",
				Age:   19,
				Email: "someone@one.ru",
				Role:  "stuff",
				Phones: []string{
					"89161012035",
					"89053401943",
				},
			},
			expectedErr: nil,
		},
		{
			testName: "integer corner cases",
			in: IntegerTestStruct{
				Minimal:  -18,
				Maximal:  50,
				SingleIn: 2,
				MultiIn:  []int{1, 2, 3, 4, 5},
			},
			expectedErr: nil,
		},
		{
			testName: "string corner cases",
			in: StringTestStruct{
				MailRegexp: "a@a.a",
				ExactLen:   "",
				SingleIn:   "alfa",
				MultiIn: []string{
					"alfa",
					"beta",
					"gamma",
					"zeta",
					"omega",
				},
			},
			expectedErr: nil,
		},
		{
			testName: "valid structure with nested",
			in: Principal{
				User: User{
					ID:    "d5ee6174-d0b6-43e2-be47-f4dff983ff48",
					Age:   19,
					Email: "someone@one.ru",
					Role:  "stuff",
					Phones: []string{
						"89161012035",
						"89053401943",
					},
				},
				Claims: []string{
					"api",
					"portal",
					"admin_panel",
				},
			},
			expectedErr: nil,
		},
	}
}

func getInvalidTestStructures() []TestStruct {
	return []TestStruct{
		{
			testName: "integer invalid fields",
			in: IntegerTestStruct{
				Minimal:  -20,
				Maximal:  52,
				SingleIn: 5,
				MultiIn:  []int{1, 2, 3, 4, 5, 6},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Minimal",
					Err:   wrapBaseValidationError(fmt.Sprintf(ShouldBeMoreOrEqualErrorText, -20, -18)),
				},
				ValidationError{
					Field: "Maximal",
					Err:   wrapBaseValidationError(fmt.Sprintf(ShouldBeLessOrEqualErrorText, 52, 50)),
				},
				ValidationError{
					Field: "SingleIn",
					Err:   wrapBaseValidationError(fmt.Sprintf(NotFoundInArrayErrorText, 5)),
				},
				ValidationError{
					Field: "MultiIn",
					Err:   wrapBaseValidationError(fmt.Sprintf(NotFoundInArrayErrorText, 6)),
				},
			},
		},
		{
			testName: "string invalid fields",
			in: StringTestStruct{
				MailRegexp: "aaa.a",
				ExactLen:   "32",
				SingleIn:   "theta",
				MultiIn: []string{
					"alfa",
					"beta",
					"gamma",
					"zeta",
					"omega",
					"theta",
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "MailRegexp",
					Err:   wrapBaseValidationError(fmt.Sprintf(DoesntMatchRegularExpressionErrorText, "aaa.a")),
				},
				ValidationError{
					Field: "ExactLen",
					Err:   wrapBaseValidationError(fmt.Sprintf(LengthDoesntMatchErrorText, "32", 0)),
				},
				ValidationError{
					Field: "SingleIn",
					Err:   wrapBaseValidationError(fmt.Sprintf(NotFoundInArrayErrorText, "theta")),
				},
				ValidationError{
					Field: "MultiIn",
					Err:   wrapBaseValidationError(fmt.Sprintf(NotFoundInArrayErrorText, "theta")),
				},
			},
		},
	}
}

func TestValidate(t *testing.T) {
	tests := make([]TestStruct, 0)
	tests = append(tests, getValidTestStructures()...)
	tests = append(tests, getInvalidTestStructures()...)

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.testName), func(t *testing.T) {
			tt := tt
			t.Parallel()

			resultErr := Validate(tt.in)
			require.Equal(t, tt.expectedErr, resultErr)
		})
	}
}
