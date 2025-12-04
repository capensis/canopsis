package validator

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/fr"
)

func RegisterTranslations(v *validator.Validate, uniTrans *ut.UniversalTranslator, invalidIDChars string) error {
	translators := make(map[string]ut.Translator, len(registerDefaultTranslations))
	for locale, f := range registerDefaultTranslations {
		trans, ok := uniTrans.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("%s translator found", locale)
		}

		err := f(v, trans)
		if err != nil {
			return fmt.Errorf("default translations for %q translator cannot be registred: %w", locale, err)
		}

		translators[locale] = trans
	}

	tagTranslations := getTagTranslations(invalidIDChars)
	for tag, c := range tagTranslations {
		for locale, s := range c.translations {
			trans, ok := translators[locale]
			if !ok {
				continue
			}

			err := v.RegisterTranslation(tag, trans, registrationFunc(tag, s, true), translateFunc)
			if err != nil {
				return fmt.Errorf("translation for %q tag and %q translator cannot be registred: %w", tag, locale, err)
			}
		}
	}

	return nil
}

var registerDefaultTranslations = map[string]func(v *validator.Validate, ut ut.Translator) error{
	types.LocaleEn: en.RegisterDefaultTranslations,
	types.LocaleFr: fr.RegisterDefaultTranslations,
}

type tagTransConfig struct {
	translations map[string]string
}

func getTagTranslations(invalidIDChars string) map[string]tagTransConfig {
	return map[string]tagTransConfig{
		"gtfield": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be greater than the specified starting value",
				types.LocaleFr: "{0} doit être supérieur à la valeur de départ spécifiée",
			},
		},
		"gtefield": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be greater than or equal to the specified starting value",
				types.LocaleFr: "{0} doit être supérieur ou égal à la valeur de départ spécifiée",
			},
		},
		"ltfield": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be less than the specified end value",
				types.LocaleFr: "{0} doit être inférieur à la valeur de fin spécifiée",
			},
		},
		"ltefield": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be less than or equal to the specified end value",
				types.LocaleFr: "{0} doit être inférieur ou égal à la valeur de fin spécifiée",
			},
		},
		"latitude": {
			translations: map[string]string{
				types.LocaleEn: "{0} must contain valid latitude coordinates",
				types.LocaleFr: "{0} doit contenir des coordonnées latitude valides",
			},
		},
		"longitude": {
			translations: map[string]string{
				types.LocaleEn: "{0} must contain valid longitude coordinates",
				types.LocaleFr: "{0} doit contenir des coordonnées longitudes valides",
			},
		},
		"notblank": {
			translations: map[string]string{
				types.LocaleEn: "{0} is a required field",
				types.LocaleFr: "{0} est un champ obligatoire",
			},
		},
		"required_or": {
			translations: map[string]string{
				types.LocaleEn: "{0} is a required field",
				types.LocaleFr: "{0} est un champ obligatoire",
			},
		},
		"required_not_both": {
			translations: map[string]string{
				types.LocaleEn: "{0} is an excluded field",
				types.LocaleFr: "{0} est un champ exclu",
			},
		},
		"must_be_empty": {
			translations: map[string]string{
				types.LocaleEn: "{0} is an excluded field",
				types.LocaleFr: "{0} est un champ exclu",
			},
		},
		"exist": {
			translations: map[string]string{
				types.LocaleEn: "{0} already exists",
				types.LocaleFr: "{0} existe déjà",
			},
		},
		"not_exist": {
			translations: map[string]string{
				types.LocaleEn: "{0} doesn't exist",
				types.LocaleFr: "{0} n'existe pas",
			},
		},
		"oneoforempty": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be one of [{1}] or empty",
				types.LocaleFr: "{0} doit être l'un des choix suivants [{1}] ou vide",
			},
		},
		"iscolororempty": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid color or empty",
				types.LocaleFr: "{0} doit être une couleur valide ou vide",
			},
		},
		"alarm_pattern": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid alarm pattern",
				types.LocaleFr: "{0} doit être un modèle d'alarme valide",
			},
		},
		"entity_pattern": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid entity pattern",
				types.LocaleFr: "{0} doit être un modèle d'entité valide",
			},
		},
		"pbehavior_pattern": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid pbehavior pattern",
				types.LocaleFr: "{0} doit être un modèle de comportement périodique valide",
			},
		},
		"event_pattern": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid event pattern",
				types.LocaleFr: "{0} doit être un modèle d'événement valide",
			},
		},
		"weather_service_pattern": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid weather service pattern",
				types.LocaleFr: "{0} doit être un modèle de météo de service valide",
			},
		},
		"id": {
			translations: map[string]string{
				types.LocaleEn: "{0} cannot contain '" + invalidIDChars + "'",
				types.LocaleFr: "{0} ne doit pas contenir '" + invalidIDChars + "'",
			},
		},
		"template": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid template",
				types.LocaleFr: "{0} doit être un modèle valide",
			},
		},
		"entity_not_exist": {
			translations: map[string]string{
				types.LocaleEn: "{0} has no associated entity",
				types.LocaleFr: "{0} n'a aucune entité associée",
			},
		},
		"alarm_not_exist": {
			translations: map[string]string{
				types.LocaleEn: "{0} has no associated alarm",
				types.LocaleFr: "{0} n'a aucune alarme associée",
			},
		},
		"not_json": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid JSON",
				types.LocaleFr: "{0} doit être un JSON valide",
			},
		},
		"multi_sort_invalid": {
			translations: map[string]string{
				types.LocaleEn: "{0} must contain a field and a sort direction",
				types.LocaleFr: "{0} doit contenir un champ et une direction de tri",
			},
		},
		"timezone": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid timezone",
				types.LocaleFr: "{0} doit être un fuseau horaire valide",
			},
		},
		"time_format": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid time format",
				types.LocaleFr: "{0} doit être dans un format d'heure valide",
			},
		},
		"rrule": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid recurrent rule",
				types.LocaleFr: "{0} doit être une règle récurrente valide",
			},
		},
		"table_name": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid table name",
				types.LocaleFr: "{0} doit être un nom de table valide",
			},
		},
		"value_string": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid string value",
				types.LocaleFr: "{0} doit être une valeur de type chaîne valide",
			},
		},
		"value_number": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid number value",
				types.LocaleFr: "{0} doit être une valeur de type nombre valide",
			},
		},
		"value_boolean": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid boolean value",
				types.LocaleFr: "{0} doit être une valeur de type booléen valide",
			},
		},
		"value_string_array": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid string array value",
				types.LocaleFr: "{0} doit être une valeur de type tableau de chaînes valide",
			},
		},
		"value_timestamp": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid timestamp value",
				types.LocaleFr: "{0} doit être une valeur de timestamp valide",
			},
		},
		"info_value": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid integer, string, boolean or array of strings value",
				types.LocaleFr: "{0} doit être une valeur entière, une chaîne, un booléen ou un tableau de chaînes valide",
			},
		},
		"regexp": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid regular expression",
				types.LocaleFr: "{0} doit être une expression régulière valide",
			},
		},
		"filemask": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid filemask",
				types.LocaleFr: "{0} doit être un masque de fichier valide",
			},
		},
		"not_api_perm": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid API permission",
				types.LocaleFr: "{0} doit être une permission API valide",
			},
		},
		"not_ui_perm": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid UI permission",
				types.LocaleFr: "{0} doit être une permission UI valide",
			},
		},
		"not_approver": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a valid profile with approval permissions",
				types.LocaleFr: "{0} doit être un profil valide disposant des permissions d'approbation",
			},
		},
		"invalid": {
			translations: map[string]string{
				types.LocaleEn: "{0} contains an invalid value",
				types.LocaleFr: "{0} contient une valeur invalide",
			},
		},
		"not_supported": {
			translations: map[string]string{
				types.LocaleEn: "{0} contains an unsupported value",
				types.LocaleFr: "{0} contient une valeur non prise en charge",
			},
		},
		"not_applicable": {
			translations: map[string]string{
				types.LocaleEn: "{0} contains a value that is not applicable",
				types.LocaleFr: "{0} contient une valeur qui n'est pas applicable",
			},
		},
		"not_accessible": {
			translations: map[string]string{
				types.LocaleEn: "{0} contains a value that is not accessible",
				types.LocaleFr: "{0} contient une valeur qui n'est pas accessible",
			},
		},
		"filesize": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a file of {1} bytes or less",
				types.LocaleFr: "{0} doit être un fichier de {1} bytes ou moins",
			},
		},
		"filetype": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a file of one of the following types [{1}]",
				types.LocaleFr: "{0} doit être un fichier d'un des types suivants [{1}]",
			},
		},
		"type": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be of a valid type",
				types.LocaleFr: "{0} doit être d'un type valide",
			},
		},
		"rulematch": {
			translations: map[string]string{
				types.LocaleEn: "{0} must match the rule",
				types.LocaleFr: "{0} doit correspondre à la règle",
			},
		},
		"allchildren": {
			translations: map[string]string{
				types.LocaleEn: "{0} must not contain all linked alarms",
				types.LocaleFr: "{0} ne doit pas contenir toutes les alarmes liées",
			},
		},
		"disabled": {
			translations: map[string]string{
				types.LocaleEn: "{0} must not contain a disabled object",
				types.LocaleFr: "{0} ne doit pas contenir d'objet désactivé",
			},
		},
		"unchangeable": {
			translations: map[string]string{
				types.LocaleEn: "{0} must not be updated",
				types.LocaleFr: "{0} ne doit pas être mis à jour",
			},
		},
		"slicelen": {
			translations: map[string]string{
				types.LocaleEn: "{0} must contain {1}",
				types.LocaleFr: "{0} doit contenir {1}",
			},
		},
		"strmax": {
			translations: map[string]string{
				types.LocaleEn: "{0} must be a maximum of {1} in length",
				types.LocaleFr: "{0} doit faire une taille maximum de {1}",
			},
		},
		"invalidcols": {
			translations: map[string]string{
				types.LocaleEn: "{0} contains invalid columns [{1}]",
				types.LocaleFr: "{0} contient des colonnes invalides [{1}]",
			},
		},
		"missingcols": {
			translations: map[string]string{
				types.LocaleEn: "{0} must contain columns [{1}]",
				types.LocaleFr: "{0} doit contenir toutes les colonnes [{1}]",
			},
		},
		"invalidpayload": {
			translations: map[string]string{
				types.LocaleEn: "{0} contains a payload that cannot be executed",
				types.LocaleFr: "{0} contient une charge utile qui ne peut pas être exécutée",
			},
		},
	}
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, override)
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
	if err != nil {
		return fe.Error()
	}

	return t
}
