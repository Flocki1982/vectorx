package intents

import (
	"fmt"
	"strings"
)

// Wirepod supported locales

const LOCALE_ENGLISH = "en"
const LOCALE_ITALIAN = "it"
const LOCALE_SPANISH = "es"
const LOCALE_FRENCH = "fr"
const LOCALE_GERMAN = "de"

// Supported parameters

const PARAMETER_USERNAME = "PARAMETER_USERNAME"
const PARAMETER_LANGUAGE = "PARAMETER_LANGUAGE"

// Standard intents a production Vector understands

const STANDARD_INTENT_NAMES_USERNAME_EXTEND = "intent_names_username_extend"
const STANDARD_INTENT_WEATHER_EXTEND = "intent_weather_extend"
const STANDARD_INTENT_NAMES_ASK = "intent_names_ask"
const STANDARD_INTENT_IMPERATIVE_EYECOLOR = "intent_imperative_eyecolor"
const STANDARD_INTENT_CHARACTER_AGE = "intent_character_age"
const STANDARD_INTENT_EXPLORE_START = "intent_explore_start"
const STANDARD_INTENT_SYSTEM_CHARGER = "intent_system_charger"
const STANDARD_INTENTSYSTEM_SLEEP_ = "intent_system_sleep"
const STANDARD_INTENT_GREETING_GOODMORNING = "intent_greeting_goodmorning"
const STANDARD_INTENT_GREETING_GOODNIGHT = "intent_greeting_goodnight"
const STANDARD_INTENT_GREETING_GOODBYE = "intent_greeting_goodbye"
const STANDARD_INTENT_SEASONAL_HAPPYNEWYEAR = "intent_seasonal_happynewyear"
const STANDARD_INTENT_SEASONAL_HAPPY_HOLIDAYS = "intent_seasonal_happyholidays"
const STANDARD_INTENT_AMAZON_SIGNIN = "intent_amazon_signin"
const STANDARD_INTENT_AMAZON_SIGNOUT = "intent_amazon_signout"
const STANDARD_INTENT_IMPERATIVE_FORWARD = "intent_imperative_forward"
const STANDARD_INTENT_IMPERATIVE_TURNAROUND = "intent_imperative_turnaround"
const STANDARD_INTENT_IMPERATIVE_TURNLEFT = "intent_imperative_turnleft"
const STANDARD_INTENT_IMPERATIVE_TURNRIGHT = "intent_imperative_turnright"
const STANDARD_INTENT_PLAY_ROLLCUBE = "intent_play_rollcube"
const STANDARD_INTENT_PLAY_POPAWHEELIE = "intent_play_popawheelie"
const STANDARD_INTENT_PLAY_FISTBUMP = "intent_play_fistbump"
const STANDARD_INTENT_PLAY_BLACKJACK = "intent_play_blackjack"
const STANDARD_INTENT_IMPERATIVE_AFFIRMATIVE = "intent_imperative_affirmative"
const STANDARD_INTENT_IMPERATIVE_NEGATIVE = "intent_imperative_negative"
const STANDARD_INTENT_PHOTO_TAKE_EXTEND_ = "intent_photo_take_extend"
const STANDARD_INTENT_IMPERATIVE_PRAISE = "intent_imperative_praise"
const STANDARD_INTENT_IMPERATIVE_ABUSE = "intent_imperative_abuse"
const STANDARD_INTENT_IMPERATIVE_APOLOGIZE = "intent_imperative_apologize"
const STANDARD_INTENT_IMPERATIVE_BACKUP = "intent_imperative_backup"
const STANDARD_INTENT_IMPERATIVE_VOLUMEDOWN = "intent_imperative_volumedown"
const STANDARD_INTENT_IMPERATIVE_VOLUMEUP = "intent_imperative_volumeup"
const STANDARD_INTENT_IMPERATIVE_LOOKATME = "intent_imperative_lookatme"
const STANDARD_INTENT_IMPERATIVE_VOLUMELEVEL_EXTEND = "intent_imperative_volumelevel_extend"
const STANDARD_INTENT_IMPERATIVE_SHUTUP = "intent_imperative_shutup"
const STANDARD_INTENT_GREETING_HELLO = "intent_greeting_hello"
const STANDARD_INTENT_IMPERATIVE_COME = "intent_imperative_come"
const STANDARD_INTENT_IMPERATIVE_LOVE = "intent_imperative_love"
const STANDARD_INTENT_PROMPTQUESTION = "intent_knowledge_promptquestion"
const STANDARD_INTENT_CHECKTIMER = "intent_clock_checktimer"
const STANDARD_INTENT_GLOBAL_STOP_EXTEND = "intent_global_stop_extend"
const STANDARD_INTENT_SETTIMER_EXTEND = "intent_clock_settimer_extend"
const STANDARD_INTENT_CLOCK_TIME = "intent_clock_time"
const STANDARD_INTENT_IMPERATIVE_QUIET = "intent_imperative_quiet"
const STANDARD_INTENT_IMPERATIVE_DANCE = "intent_imperative_dance"
const STANDARD_INTENT_PLAY_PICKUPCUBE = "intent_play_pickupcube"
const STANDARD_INTENT_IMPERATIVE_FETCHCUBE = "intent_imperative_fetchcube"
const STANDARD_INTENT_IMPERATIVE_FINDCUBE = "intent_imperative_findcube"
const STANDARD_INTENT_PLAY_ANYTRICK = "intent_play_anytrick"
const STANDARD_INTENT_RECORDMESSAGE_EXTEND = "intent_message_recordmessage_extend"
const STANDARD_INTENT_PLAYMESSAGE_EXTEND = "intent_message_playmessage_extend"
const STANDARD_INTENTBLACKJACK_HIT = "intent_blackjack_hit"
const STANDARD_INTENT_BLACKJACK_STAND = "intent_blackjack_stand"
const STANDARD_INTENT_KEEPAWAY = "intent_play_keepaway"

type IntentParams struct {
	RobotName string
	Language  string
}

type IntentHandlerFunc func(IntentDef, IntentParams) string

type IntentDef struct {
	IntentName string
	Utterances map[string][]string
	Parameters []string
	Handler    IntentHandlerFunc
}

var intents []IntentDef

func RegisterIntents() {
	HelloWorld_Register(&intents)
	RollaDie_Register(&intents)
	RobotName_Register(&intents)
	ImageTest_Register(&intents)
	ChangeLanguage_Register(&intents)
}

func IntentMatch(speechText string, locale string) (IntentDef, error) {
	for _, intent := range intents {
		if hasPerfectMatch(intent.Utterances[locale], speechText) {
			return intent, nil
		}
	}
	for _, intent := range intents {
		if hasPartialMatch(intent.Utterances[locale], speechText) {
			return intent, nil
		}
	}
	return IntentDef{}, fmt.Errorf("Intent not found")
}

/**********************************************************************************************************************/
/*                                                PRIVATE FUNCTIONS                                                   */
/**********************************************************************************************************************/

func hasPerfectMatch(utterances []string, phrase string) bool {
	for _, s := range utterances {
		if strings.ToLower(s) == strings.ToLower(phrase) {
			return true
		}
	}
	return false
}

func hasPartialMatch(utterances []string, phrase string) bool {
	for _, s := range utterances {
		if strings.Contains(strings.ToLower(phrase), strings.ToLower(s)) {
			return true
		}
	}
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
