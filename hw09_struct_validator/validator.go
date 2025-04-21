package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field %s: %v", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for i, err := range v {
		sb.WriteString(err.Error())
		if i < len(v)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

var (
	ErrNotStruct         = errors.New("input is not a struct")
	ErrValidateRule      = errors.New("invalid validation rule")
	ErrRegexpCompile     = errors.New("failed to compile regexp")
	ErrInvalidType       = errors.New("unsupported type for validation")
	ErrStringLength      = errors.New("string length mismatch")
	ErrStringMinLength   = errors.New("string length is less than min length")
	ErrStringMaxLength   = errors.New("string length is greater than max length")
	ErrStringRegexp      = errors.New("string does not match regexp")
	ErrStringNotInSet    = errors.New("string is not in allowed set")
	ErrNumberMin         = errors.New("number is less than min")
	ErrNumberMax         = errors.New("number is greater than max")
	ErrNumberNotInSet    = errors.New("number is not in allowed set")
	ErrInvalidNestedType = errors.New("invalid nested type")
)

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ValidationErrors{
			ValidationError{
				Field: "", // no field name for top-level error
				Err:   ErrNotStruct,
			},
		}
	}

	var valErrors ValidationErrors
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		if !fieldType.IsExported() {
			continue
		}

		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		if validateTag == "nested" {
			valErrors = append(valErrors, handleNestedField(field, fieldType)...)
			continue
		}

		rules := strings.Split(validateTag, "|")
		fieldErrs := validateField(field, fieldType.Name, rules)
		valErrors = append(valErrors, fieldErrs...)
	}

	if len(valErrors) > 0 {
		return valErrors
	}

	return nil
}

func handleNestedField(
	field reflect.Value,
	fieldType reflect.StructField,
) ValidationErrors {
	var valErrors ValidationErrors

	if field.Kind() == reflect.Struct {
		if err := Validate(field.Interface()); err != nil {
			var nestedValErrors ValidationErrors
			if errors.As(err, &nestedValErrors) {
				for _, e := range nestedValErrors {
					valErrors = append(valErrors, ValidationError{
						Field: fieldType.Name + "." + e.Field,
						Err:   e.Err,
					})
				}
			}
		}
	} else {
		valErrors = append(valErrors, ValidationError{
			Field: fieldType.Name,
			Err:   ErrInvalidNestedType,
		})
	}

	return valErrors
}

func validateField(field reflect.Value, fieldName string, rules []string) ValidationErrors {
	var valErrors ValidationErrors

	if field.Kind() == reflect.Slice {
		for i := 0; i < field.Len(); i++ {
			elemErrs := validateValue(field.Index(i), fieldName, rules)
			for _, err := range elemErrs {
				err.Field = fmt.Sprintf("%s[%d]", err.Field, i)
				valErrors = append(valErrors, err)
			}
		}
		return valErrors
	}

	return validateValue(field, fieldName, rules)
}

func validateValue(value reflect.Value, fieldName string, rules []string) ValidationErrors {
	var valErrors ValidationErrors

	for _, rule := range rules {
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			valErrors = append(valErrors, ValidationError{
				Field: fieldName,
				Err:   ErrValidateRule,
			})
			continue
		}

		ruleName := parts[0]
		ruleParam := parts[1]

		//nolint:exhaustive // We are not validating all possible types here, just the ones we need.
		switch value.Kind() {
		case reflect.String:
			err := validateString(value.String(), ruleName, ruleParam)
			if err != nil {
				valErrors = append(valErrors, ValidationError{Field: fieldName, Err: err})
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err := validateInt(value.Int(), ruleName, ruleParam)
			if err != nil {
				valErrors = append(valErrors, ValidationError{Field: fieldName, Err: err})
			}
		default:
			valErrors = append(valErrors, ValidationError{
				Field: fieldName,
				Err:   ErrInvalidType,
			})
		}
	}

	return valErrors
}

func validateString(value string, ruleName string, ruleParam string) error {
	switch ruleName {
	case "len":
		length, err := strconv.Atoi(ruleParam)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrValidateRule, err)
		}
		if len(value) != length {
			return ErrStringLength
		}
	case "regexp":
		re, err := regexp.Compile(ruleParam)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrRegexpCompile, err)
		}
		if !re.MatchString(value) {
			return ErrStringRegexp
		}
	case "in":
		allowedValues := strings.Split(ruleParam, ",")
		for _, allowed := range allowedValues {
			if value == allowed {
				return nil
			}
		}
		return ErrStringNotInSet
	case "minlength":
		minLength, err := strconv.Atoi(ruleParam)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrValidateRule, err)
		}
		if len(value) < minLength {
			return ErrStringMinLength
		}
	case "maxlength":
		maxLength, err := strconv.Atoi(ruleParam)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrValidateRule, err)
		}
		if len(value) > maxLength {
			return ErrStringMaxLength
		}
	default:
		return fmt.Errorf("%w: unknown rule %s", ErrValidateRule, ruleName)
	}
	return nil
}

func validateInt(value int64, ruleName string, ruleParam string) error {
	switch ruleName {
	case "min":
		minValue, err := strconv.ParseInt(ruleParam, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrValidateRule, err)
		}
		if value < minValue {
			return ErrNumberMin
		}
	case "max":
		maxValue, err := strconv.ParseInt(ruleParam, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrValidateRule, err)
		}
		if value > maxValue {
			return ErrNumberMax
		}
	case "in":
		allowedValues := strings.Split(ruleParam, ",")
		for _, allowed := range allowedValues {
			allowedInt, err := strconv.ParseInt(allowed, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrValidateRule, err)
			}
			if value == allowedInt {
				return nil
			}
		}
		return ErrNumberNotInSet
	default:
		return fmt.Errorf("%w: unknown rule %s", ErrValidateRule, ruleName)
	}
	return nil
}
