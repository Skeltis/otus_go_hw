package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const tagName = "validate"

type ruleType int

const (
	Regexp ruleType = iota
	Min
	Max
	Len
	InString
	InInt
)

const (
	NotFoundInArrayErrorText              = "'%v' doesn't exist in the provided array"
	ShouldBeLessOrEqualErrorText          = "%v should be less or equal to %v"
	ShouldBeMoreOrEqualErrorText          = "%v should be more or equal to %v"
	DoesntMatchRegularExpressionErrorText = "'%s' doesn't match to given regular expression"
	LengthDoesntMatchErrorText            = "length of '%s' doesn't match length = %v"
)

var ErrBaseValidationError = errors.New("validation error")

type ValidationError struct {
	Field string
	Err   error
}

type validationRule struct {
	validationType ruleType
	regexpArg      *regexp.Regexp
	intArgs        []int
	stringArgs     []string
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var stringBuilder strings.Builder
	for _, validationError := range v {
		stringBuilder.WriteString(fmt.Sprintf("%v", validationError.Error()))
	}
	return stringBuilder.String()
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field %v: %v \n", v.Field, v.Err.Error())
}

func wrapBaseValidationError(text string) error {
	return fmt.Errorf("%w: %s", ErrBaseValidationError, text)
}

func Validate(v interface{}) error {
	if v == nil {
		return nil
	}

	validationErrors, err := validate(v)
	if err != nil {
		return err
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return validationErrors
}

func validate(v interface{}) (ValidationErrors, error) {
	validationErrors := make([]ValidationError, 0)
	reflection := reflect.ValueOf(v)
	if reflection.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%T is not a pointer to a struct", v)
	}

	for fieldNumber := 0; fieldNumber < reflection.NumField(); fieldNumber++ {
		fieldInfo := reflection.Type().Field(fieldNumber)
		field := reflection.Field(fieldNumber)
		tag := fieldInfo.Tag.Get(tagName)

		if tag == "nested" && reflection.Kind() == reflect.Struct {
			iterationValError, err := validate(field.Interface())
			if err != nil {
				return nil, err
			}
			validationErrors = append(validationErrors, iterationValError...)
			continue
		}

		if !isEmptyOrWhiteSpace(tag) {
			extractedRules, err := extractRules(tag, fieldInfo.Type)
			if err != nil {
				return nil, fmt.Errorf("error on validation rules: %w", err)
			}

			validationErrors = append(validationErrors, validateField(field, fieldInfo.Name, extractedRules)...)
		}
	}

	return validationErrors, nil
}

func validateField(field reflect.Value, fieldName string, validationRules []validationRule) ValidationErrors {
	var err error
	validationErrors := make([]ValidationError, 0)

	switch field.Kind() { //nolint:exhaustive
	case reflect.Int:
		err = validateIntegerValue(int(field.Int()), validationRules)
		if err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}

	case reflect.String:
		err = validateStringValue(field.String(), validationRules)
		if err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}

	case reflect.Slice:
		switch field.Type().Elem().Kind() { //nolint:exhaustive
		case reflect.String:
			for _, value := range field.Interface().([]string) {
				err = validateStringValue(value, validationRules)
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldName,
						Err:   err,
					})
				}
			}
		case reflect.Int:
			for _, value := range field.Interface().([]int) {
				err = validateIntegerValue(value, validationRules)
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldName,
						Err:   err,
					})
				}
			}
		}
	}
	return validationErrors
}

func extractRules(validationExpression string, refType reflect.Type) ([]validationRule, error) {
	wrappedRules := strings.Split(validationExpression, "|")
	unwrappedRules := make([]validationRule, 0, len(wrappedRules))

	for _, rule := range wrappedRules {
		ruleAndArgs := strings.Split(rule, ":")
		if len(ruleAndArgs) < 2 {
			return nil, fmt.Errorf("validation format should be [ruleType]:[ruleArguments,...]: %s", rule)
		}

		reflectSingleType := refType
		if reflectSingleType == reflect.SliceOf(reflect.TypeOf("string")) {
			reflectSingleType = reflect.TypeOf("string")
		} else if reflectSingleType == reflect.SliceOf(reflect.TypeOf(12)) {
			reflectSingleType = reflect.TypeOf(12)
		}

		extractedRule, err := extractValidationRule(ruleAndArgs[0], ruleAndArgs[1], reflectSingleType)
		if err != nil {
			return nil, err
		}

		unwrappedRules = append(unwrappedRules, *extractedRule)
	}
	return unwrappedRules, nil
}

func extractValidationRule(ruleName string, ruleArgs string, refType reflect.Type) (*validationRule, error) {
	lowerRuleName := strings.ToLower(ruleName)
	switch lowerRuleName {
	case "regexp":
		if refType.Kind() != reflect.String {
			return nil, fmt.Errorf("incompatible field type for rule '%s': %s", lowerRuleName, ruleArgs)
		}

		regexArg, err := regexp.Compile(ruleArgs)
		if err != nil {
			return nil, fmt.Errorf("error in regular expression: %s, %w", ruleArgs, err)
		}

		return &validationRule{
			validationType: Regexp,
			regexpArg:      regexArg,
		}, nil

	case "len":
		if refType.Kind() != reflect.String {
			return nil, fmt.Errorf("incompatible field type for rule '%s': %s", lowerRuleName, ruleArgs)
		}

		intValue, err := strconv.Atoi(ruleArgs)
		if err != nil {
			return nil, fmt.Errorf("incorrect digit format: %s; %w", ruleArgs, err)
		}

		return &validationRule{
			validationType: Len,
			intArgs:        []int{intValue},
		}, nil

	case "min", "max":
		if refType.Kind() != reflect.Int {
			return nil, fmt.Errorf("incompatible field type for rule '%s': %s", lowerRuleName, ruleArgs)
		}

		intValue, err := strconv.Atoi(ruleArgs)
		if err != nil {
			return nil, fmt.Errorf("incorrect digit format: %s; %w", ruleArgs, err)
		}

		outRule := validationRule{intArgs: []int{intValue}}

		if lowerRuleName == "min" {
			outRule.validationType = Min
		} else {
			outRule.validationType = Max
		}

		return &outRule, nil

	case "in":
		stringValues := strings.Split(ruleArgs, ",")
		if len(stringValues) < 1 {
			return nil, fmt.Errorf("nowhere to search for value: %s", ruleArgs)
		}

		if refType.Kind() == reflect.String {
			return &validationRule{
				validationType: InString,
				stringArgs:     stringValues,
			}, nil
		}

		intValues := make([]int, len(stringValues))
		for i := 0; i < len(stringValues); i++ {
			intVal, err := strconv.Atoi(stringValues[i])
			if err == nil {
				intValues[i] = intVal
			} else {
				return nil, fmt.Errorf("error parsing integers: %s", ruleArgs)
			}
		}
		return &validationRule{
			validationType: InInt,
			intArgs:        intValues,
		}, nil

	default:
		return nil, fmt.Errorf("unknown validation rule: %s", lowerRuleName)
	}
}

func isEmptyOrWhiteSpace(s string) bool {
	if len(s) == 0 {
		return true
	}

	if len(strings.TrimSpace(s)) == 0 {
		return true
	}

	return false
}

func validateStringValue(value string, validationRules []validationRule) error {
	for _, rule := range validationRules {
		switch rule.validationType { //nolint:exhaustive
		case Len:
			if len([]rune(value)) != rule.intArgs[0] {
				return wrapBaseValidationError(fmt.Sprintf(LengthDoesntMatchErrorText, value, rule.intArgs[0]))
			}

		case Regexp:
			doesMatch := rule.regexpArg.MatchString(value)
			if !doesMatch {
				return wrapBaseValidationError(fmt.Sprintf(DoesntMatchRegularExpressionErrorText, value))
			}

		case InString:
			found := false
			for _, elementToMatch := range rule.stringArgs {
				if elementToMatch == value {
					found = true
					break
				}
			}
			if !found {
				return wrapBaseValidationError(fmt.Sprintf(NotFoundInArrayErrorText, value))
			}

		default:
		}
	}
	return nil
}

func validateIntegerValue(value int, validationRules []validationRule) error {
	for _, rule := range validationRules {
		switch rule.validationType { //nolint:exhaustive
		case Min:
			if value < rule.intArgs[0] {
				return wrapBaseValidationError(fmt.Sprintf(ShouldBeMoreOrEqualErrorText, value, rule.intArgs[0]))
			}

		case Max:
			if value > rule.intArgs[0] {
				return wrapBaseValidationError(fmt.Sprintf(ShouldBeLessOrEqualErrorText, value, rule.intArgs[0]))
			}

		case InInt:
			found := false
			for _, elementToMatch := range rule.intArgs {
				if elementToMatch == value {
					found = true
					break
				}
			}
			if !found {
				return wrapBaseValidationError(fmt.Sprintf(NotFoundInArrayErrorText, value))
			}

		default:
		}
	}
	return nil
}
