package localization_test

import (
	"fmt"
	"testing"

	"github.com/aarongreenlee/go-infrastructure/localization"
)

func Example() {

	_ = localization.Init([]string{
		"./testdata/en-US.all.json",
		"./testdata/es-ES.all.json"}, "en-us")

	fmt.Printf("Default translator-alpha: %s\n", localization.T("alpha"))
	fmt.Printf("Default translator-beta: %s\n", localization.T("beta"))
	fmt.Printf("Default translator-omega: %s\n", localization.T("omega"))
	fmt.Printf("Default translator-Missing: %s\n", localization.T("missing"))

	t := localization.NewTranslationFunc("es-es")
	fmt.Printf("Spanish translator-alpha: %s\n", t("alpha"))
	fmt.Printf("Spanish translator-beta: %s\n", t("beta"))
	fmt.Printf("Spanish translator-omega: %s\n", t("omega"))
	fmt.Printf("Spanish translator-Missing: %s\n", t("missing"))

	fmt.Printf("Localized Error Key: %s\n", localization.NewError("alpha").Error())
	fmt.Printf("Localized Error Message: %s\n", localization.T(localization.NewError("alpha").Error()))

	// Output:
	//Default translator-alpha: Rocket
	//Default translator-beta: Ship
	//Default translator-omega: Trip
	//Default translator-Missing: missing
	//Spanish translator-alpha: Cohete
	//Spanish translator-beta: Barco
	//Spanish translator-omega: Viaje
	//Spanish translator-Missing: missing
	//Localized Error Key: alpha
	//Localized Error Message: Rocket
}

func Test_NewError(t *testing.T) {
	key := "errorkey"
	err := localization.NewError(key)

	if err == nil {
		t.Error("Nil error returned by factory")
	}

	if err.Error() != key {
		t.Errorf("Expected %s but got %s", key, err.Error())
	}
}

func Test_InitFailure(t *testing.T) {
	err := localization.Init([]string{}, "en-us")

	switch err {
	case nil:
		t.Error("Without files Init should error")
	case localization.ErrNoTranslationFiles:
		return
	}

	t.Error("Error returned is of an unexpected type")

}
