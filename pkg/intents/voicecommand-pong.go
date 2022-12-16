package intents

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/fforchino/vector-go-sdk/pkg/vectorpb"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"math/rand"
	"os"
	"time"
	opencv_ifc "vectorx/pkg/opencv-ifc"
)

/**********************************************************************************************************************/
/*                                          FOLLOW (INDEX) FINGER                                                     */
/**********************************************************************************************************************/

func Pong_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"let's play a classic"}
	utterances[LOCALE_ITALIAN] = []string{"giochiamo a pong"}
	utterances[LOCALE_SPANISH] = []string{"juguemos a pong"}
	utterances[LOCALE_FRENCH] = []string{"jouons à pong"}
	utterances[LOCALE_GERMAN] = []string{"lass uns pong spielen"}

	var intent = IntentDef{
		IntentName: "extended_intent_play_pong",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    playPong,
	}
	*intentList = append(*intentList, intent)
	addLocalizedString("STR_PONG_I_WON", []string{"I won! ", "ho vinto!", "yo gané!", "j'ai gagné!", "ich habe gewonnen!"})
	addLocalizedString("STR_PONG_YOU_WON", []string{"you won!", "hai vinto!", "ganaste tu!", "tu as gagné!", "du hast gewonnen!"})

	return nil
}

func playPong(intent IntentDef, speechText string, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	doPong(true)
	return returnIntent
}

func doPong(useFx bool) {
	opencv_ifc.CreateClient()
	// Run opencv server on my pc to be faster
	// opencv_ifc.SetServerAddress("http://192.168.43.65:8090")

	s1 := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(s1)

	sdk_wrapper.MoveHead(3.0)
	sdk_wrapper.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
	sdk_wrapper.UseVectorEyeColorInImages(true)
	paddleFile, _ := os.Open(sdk_wrapper.GetDataPath("images/pong/paddle.png"))
	ballFile, _ := os.Open(sdk_wrapper.GetDataPath("images/pong/ball.png"))
	playFieldFile, _ := os.Open(sdk_wrapper.GetDataPath("images/pong/playfield.png"))
	paddle, _, _ := image.Decode(paddleFile)
	ball, _, _ := image.Decode(ballFile)
	playField, _, _ := image.Decode(playFieldFile)

	var scores []image.Image = []image.Image{}

	for i := 0; i <= 9; i++ {
		dFile, _ := os.Open(sdk_wrapper.GetDataPath(fmt.Sprintf("images/pong/digits_%d.png", i)))
		score, _, _ := image.Decode(dFile)
		scores = append(scores, score)
	}

	humanScore := 0
	vectorScore := 0
	const WIDTH = sdk_wrapper.VECTOR_SCREEN_WIDTH
	const HEIGHT = sdk_wrapper.VECTOR_SCREEN_HEIGHT
	PADDLE_WIDTH := paddle.Bounds().Dx()
	PADDLE_HEIGHT := paddle.Bounds().Dy()
	SPEED := 5
	VECTOR_PADDLE_SPEED := 3
	BALL_WIDTH := ball.Bounds().Dx()
	BALL_HEIGHT := ball.Bounds().Dy()

	// Human paddle coordinates
	humanPaddle := image.Point{0, HEIGHT / 2}
	// Vector paddle coordinates
	vectorPaddle := image.Point{WIDTH - PADDLE_WIDTH, HEIGHT / 2}
	// Ball coordinates
	ballObj := image.Point{WIDTH / 2, HEIGHT / 2}
	// Ball speed
	bSpeed := image.Point{X: -1 * SPEED, Y: 0}

	fx := ""

	// Play audio asynchronously
	go func() {
		for true {
			if useFx && fx != "" {
				println(fx)
				sdk_wrapper.PlaySound(fx)
				fx = ""
			}
		}
	}()

	// Read input asynchronously
	go func() {
		for true {
			//tx := time.Now().UnixMilli()
			img, err := sdk_wrapper.GetStaticCameraPicture(false)
			if nil == err {
				//println(fmt.Sprintf("T1: %d", time.Now().UnixMilli()-tx))
				var handInfo map[string]interface{}
				jsonData := opencv_ifc.SendImageToImageServer(&img)
				//println("OpenCV server response: " + jsonData)
				json.Unmarshal([]byte(jsonData), &handInfo)
				index_x := int(handInfo["index_x"].(float64))
				if index_x != -1 {
					// Increment human paddle position
					humanPaddle.Y = HEIGHT * index_x / img.Bounds().Dx()
					if humanPaddle.Y < PADDLE_HEIGHT/2 {
						humanPaddle.Y = PADDLE_HEIGHT / 2
					} else if humanPaddle.Y > HEIGHT-PADDLE_HEIGHT/2 {
						humanPaddle.Y = HEIGHT - PADDLE_HEIGHT/2
					}
				}
				//println(fmt.Sprintf("T2: %d", time.Now().UnixMilli()-tx))
			}
		}
	}()

	for humanScore < 9 && vectorScore < 9 {
		// Increment ball position
		ballObj.X += bSpeed.X
		ballObj.Y += bSpeed.Y

		// Increment Vector's paddle position, and check bounds
		if (vectorPaddle.Y + PADDLE_WIDTH) < ballObj.Y {
			vectorPaddle.Y += VECTOR_PADDLE_SPEED
		}
		if (vectorPaddle.Y + PADDLE_WIDTH) > ballObj.Y {
			vectorPaddle.Y -= VECTOR_PADDLE_SPEED
		}

		if vectorPaddle.Y < PADDLE_HEIGHT/2 {
			vectorPaddle.Y = PADDLE_HEIGHT / 2
		} else if vectorPaddle.Y > HEIGHT-PADDLE_HEIGHT/2 {
			vectorPaddle.Y = HEIGHT - PADDLE_HEIGHT/2
		}

		// Check bouncing
		if ballObj.X <= PADDLE_WIDTH && bSpeed.X < 0 {
			// Ball hits Human's wall
			if ballObj.Y >= humanPaddle.Y-PADDLE_HEIGHT/2 && ballObj.Y <= humanPaddle.Y+PADDLE_HEIGHT/2 {
				// Paddle hits the ball
				fx = sdk_wrapper.GetDataPath("audio/pong/ball_ping.pcm")
				bSpeed.X = bSpeed.X * -1
				bSpeed.Y = (humanPaddle.Y - ballObj.Y) / 4
				println(fmt.Sprintf(">> HUMAN HITS THE BALL, new speed %d,%d", bSpeed.X, bSpeed.Y))
			} else {
				// Ball lost
				if vectorScore < 9 {
					vectorScore++
				}
				fx = sdk_wrapper.GetDataPath("audio/pong/ball_out.pcm")
				ballObj.X = WIDTH / 2
				ballObj.Y = HEIGHT/2 + -1*HEIGHT/5 + rnd.Intn(HEIGHT/5*2)
				bSpeed.X = -1 * SPEED
				bSpeed.Y = -3 + rnd.Intn(6)
				println(fmt.Sprintf(">> BALL LOST, new speed %d,%d", bSpeed.X, bSpeed.Y))
			}
		} else if ballObj.X+BALL_WIDTH >= WIDTH-PADDLE_WIDTH && bSpeed.X > 0 {
			if ballObj.Y >= vectorPaddle.Y-PADDLE_HEIGHT/2 && ballObj.Y <= vectorPaddle.Y+PADDLE_HEIGHT/2 {
				// Paddle hits the ball
				fx = sdk_wrapper.GetDataPath("audio/pong/ball_pong.pcm")
				bSpeed.X = bSpeed.X * -1
				bSpeed.Y = (vectorPaddle.Y - ballObj.Y) / 4
				println(fmt.Sprintf(">> VECTOR HITS THE BALL, new speed %d,%d", bSpeed.X, bSpeed.Y))
			} else {
				// Ball lost
				if humanScore < 9 {
					humanScore++
				}
				fx = sdk_wrapper.GetDataPath("audio/pong/ball_out.pcm")
				ballObj.X = WIDTH / 2
				ballObj.Y = HEIGHT/2 + -1*HEIGHT/5 + rnd.Intn(HEIGHT/5*2)
				bSpeed.X = -1 * SPEED
				bSpeed.Y = -3 + rnd.Intn(6)
				println(fmt.Sprintf(">> BALL LOST, new speed %d,%d", bSpeed.X, bSpeed.Y))
			}
		} else if ballObj.Y <= BALL_HEIGHT || ballObj.Y+BALL_HEIGHT >= HEIGHT {
			// Ball hits top or bottom part of the screen, bounce back
			println(fmt.Sprintf(">> BALL BOUNCE, new speed %d,%d", bSpeed.X, bSpeed.Y))
			fx = sdk_wrapper.GetDataPath("audio/pong/ball_bounce.pcm")
			bSpeed.Y = bSpeed.Y * -1
		}

		// Draw
		dc := gg.NewContext(WIDTH, HEIGHT)
		dc.DrawImage(playField, 0, 0)
		dc.DrawImage(paddle, humanPaddle.X, humanPaddle.Y-PADDLE_HEIGHT/2)
		dc.DrawImage(paddle, vectorPaddle.X, vectorPaddle.Y-PADDLE_HEIGHT/2)
		dc.DrawImage(ball, ballObj.X, ballObj.Y)
		dc.DrawImage(scores[humanScore], WIDTH/4-5, 0)
		dc.DrawImage(scores[vectorScore], WIDTH/4*3-5, 0)

		buf := new(bytes.Buffer)
		bitmap := sdk_wrapper.ConvertPixelsToRawBitmap(dc.Image(), 100)
		for _, ui := range bitmap {
			binary.Write(buf, binary.LittleEndian, ui)
		}
		_, _ = sdk_wrapper.Robot.Conn.DisplayFaceImageRGB(
			context.Background(),
			&vectorpb.DisplayFaceImageRGBRequest{
				FaceData:         buf.Bytes(),
				DurationMs:       100,
				InterruptRunning: true,
			},
		)
		//println(fmt.Sprintf("Step %d/%d. User pos @ %d,%d, ball pos %d,%d ballspeed @ %d,%d", i, numSteps, humanPaddle.X, humanPaddle.Y, ballObj.X, ballObj.Y, bSpeed.X, bSpeed.Y))
	}

	// Game over. Let's see who won
	if vectorScore > humanScore {
		sdk_wrapper.SayText(getText("STR_PONG_I_WON"))
	} else {
		sdk_wrapper.SayText(getText("STR_PONG_YOU_WON"))
	}
}
