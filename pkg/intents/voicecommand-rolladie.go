package intents

import (
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"image/color"
	"math/rand"
	"time"
)

/**********************************************************************************************************************/
/*                                                ROLL A DIE                                                          */
/**********************************************************************************************************************/

func RollaDie_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"roll a die", "roll the die", "roll the dice", "roll the die", "dice", "roll a number", "die"}
	utterances[LOCALE_ITALIAN] = []string{"tira un dado", "lancia un dado"}

	var intent = IntentDef{
		IntentName: "extended_intent_rolladie",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    rollADie,
	}
	*intentList = append(*intentList, intent)
	return nil
}

func rollADie(intent IntentDef, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	sdk_wrapper.MoveHead(3.0)
	sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})

	sdk_wrapper.UseVectorEyeColorInImages(true)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	die := r1.Intn(6) + 1
	dieImage := fmt.Sprintf("images/dice/%d.png", die)
	dieImage = sdk_wrapper.GetDataPath(dieImage)

	sdk_wrapper.DisplayAnimatedGif(sdk_wrapper.GetDataPath("images/dice/roll-the-dice.gif"), sdk_wrapper.ANIMATED_GIF_SPEED_FASTEST, 1, false)
	sdk_wrapper.DisplayImage(dieImage, 100, false)
	sdk_wrapper.PlaySystemSound(sdk_wrapper.SYSTEMSOUND_WIN)
	sdk_wrapper.SayText(fmt.Sprintf("You rolled a %d", die))
	sdk_wrapper.DisplayImageWithTransition(dieImage, 1000, sdk_wrapper.IMAGE_TRANSITION_FADE_OUT, 10)
	return returnIntent
}
