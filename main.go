package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	modifiedImgWidth = 340
	loop             = true
)

var (
	captcha  Captcha
	tipY     int
	fp       string
	detail   string
	unsolved bool
)

/*
 Retrieve the images that will be used to calculate the captcha.
*/
func getCaptchaImage() {
	fmt.Print("\r[⚙️] retrieving captcha images...")

	// Grab the image URLs needed and captcha metadata.
	resp, _ := http.Get("https://us.tiktok.com/captcha/get?lang=en&app_name=&h5_sdk_version=2.26.17&sdk_version=3.6.1&iid=0&did=0&device_id=0&ch=web_text&aid=1459&os_type=2&mode=&tmp=1664410177244&platform=pc&webdriver=false&fp=" + fp + "&type=verify&detail=" + detail + "&server_sdk_env={\"idc\":\"useast5\",\"region\":\"US-TTP\",\"server_type\":\"passport\"}&subtype=slide&challenge_code=3058&os_name=windows&h5_check_version=3.6.1")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	boardUrl := strings.Split(strings.Split(string(body), "\"url1\":\"")[1], "\"")[0]
	pieceUrl := strings.Split(strings.Split(string(body), "\"url2\":\"")[1], "\"")[0]
	tipY, _ = strconv.Atoi(strings.Split(strings.Split(string(body), "\"tip_y\":")[1], "}")[0])
	captcha.ID = strings.Split(strings.Split(string(body), "\"id\":\"")[1], "\"")[0]

	// Get the puzzle board and write it to a file.
	boardOut, _ := os.Create("board.jpg")
	defer boardOut.Close()
	boardResp, _ := http.Get(boardUrl)
	defer boardResp.Body.Close()
	io.Copy(boardOut, boardResp.Body)

	// Get the puzzle piece and write it to a file.
	pieceOut, _ := os.Create("piece.png")
	defer pieceOut.Close()
	pieceResp, _ := http.Get(pieceUrl)
	defer pieceResp.Body.Close()
	io.Copy(pieceOut, pieceResp.Body)

	fmt.Println("\r[✅️] retrieving captcha images... downloaded!")
}

func processCaptcha() {
	fmt.Print("\r[⚙️] detecting puzzle piece...")

	// Set metadata
	captcha.ModifiedImgWidth = modifiedImgWidth
	captcha.Mode = "slide"

	// Open image
	boardJPG, err := os.Open("./board.jpg")
	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}
	defer boardJPG.Close()

	board, err := jpeg.Decode(boardJPG)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Iterate over the picture to find the corner of the puzzle piece
	puzzleX := 0
	puzzleY := 0
	for x := 0; x < board.Bounds().Max.X; x++ {
		if puzzleX != 0 {
			break
		}
		for y := 0; y < board.Bounds().Max.Y; y++ {
			// Grab the next pixel
			r, g, b, _ := board.At(x, y).RGBA()
			r /= 257
			g /= 257
			b /= 257

			// Check if the red value of the pixel exceeds our threshold here
			if r > 200 {
				// Grab the corner pixels
				r2, g2, b2, _ := board.At(x+1, y).RGBA()
				r2 /= 257
				g2 /= 257
				b2 /= 257

				r3, g3, b3, _ := board.At(x, y+1).RGBA()
				r3 /= 257
				g3 /= 257
				b3 /= 257

				r4temp, g4temp, b4temp, _ := board.At(x+1, y+1).RGBA()
				r4temp /= 257
				g4temp /= 257
				b4temp /= 257

				// Check the difference in color between the supposed corner pixels and the pixel that is though to be in the puzzle
				if (int(r)-int(r4temp) > 150 && int(g)-int(g4temp) > 150 && int(b)-int(b4temp) > 150) && (int(r2)-int(r4temp) > 150 && int(g2)-int(g4temp) > 150 && int(b2)-int(b4temp) > 150) && (int(r3)-int(r4temp) > 150 && int(g3)-int(g4temp) > 150 && int(b3)-int(b4temp) > 150) {
					puzzleX = x
					puzzleY = y
					break
				}
			}
		}
	}
	fmt.Println("\r[✅️] detecting puzzle piece... found at (" + strconv.Itoa(puzzleX) + ", " + strconv.Itoa(puzzleY) + ")!")

	// do something idk lol
	fmt.Print("[⚙️] calculating humanized captcha data...")

	startTime := time.Now()
	time.Sleep(200 * time.Millisecond)

	// Linear calculation for 'Reply'
	for iterations := 1; iterations <= puzzleX+40; iterations++ {
		time.Since(startTime)
		c := struct {
			X            int `json:"x"`
			Y            int `json:"y"`
			RelativeTime int `json:"relative_time"`
		}{
			iterations,
			tipY,
			int(time.Since(startTime).Milliseconds()),
		}
		captcha.Reply = append(captcha.Reply, c)

		// Omitting this as it seems to be unnecessary
		/*if iterations%3 == 0 {
			//captcha.Reply2 = append(captcha.Reply2, c)
		}*/
		iterations++
		time.Sleep(5 * time.Millisecond)
	}
	fmt.Println("\r[✅️] calculating humanized captcha data... done!")
}

func submitCaptcha() {
	fmt.Print("\r[⚙️] submitting captcha to TikTok...")

	solvedCaptcha, _ := json.Marshal(captcha)
	resp, _ := http.Post("https://us.tiktok.com/captcha/verify?lang=en&app_name=&h5_sdk_version=2.26.17&sdk_version=3.6.1&iid=0&did=0&device_id=0&ch=web_text&aid=1459&os_type=2&mode=slide&tmp="+strconv.FormatInt(time.Now().UnixMilli()-10000, 10)+"&platform=pc&webdriver=false&fp="+fp+"&type=verify&detail="+detail+"&server_sdk_env=%7B%22idc%22:%22useast5%22,%22region%22:%22US-TTP%22,%22server_type%22:%22passport%22%7D&subtype=slide&challenge_code=99999&os_name=windows&h5_check_version=3.6.1", "application/json", bytes.NewBuffer(solvedCaptcha))
	body, _ := ioutil.ReadAll(resp.Body)

	if strings.Contains(string(body), "Verification complete") {
		unsolved = false
		fmt.Println("\r[✅️] submitting captcha to TikTok... solved!")
		/*for _, cookie := range resp.Cookies() {
			if cookie.Name == "_abck" {
				fmt.Println("[✅️] _abck: " + cookie.Value)
			} else if cookie.Name == "bm_sz" {
				fmt.Println("[✅️] bm_sz: " + cookie.Value)
			}
		}*/
	} else {
		fmt.Println("\r[❌️] submitting captcha to TikTok... failed!")
		fmt.Println("[❌️] Error: " + string(body))
	}
}

func main() {
	fmt.Println("[✨️] ttCaptchaSolver // v1.0-slide // made with ❤️ by luma")

	fmt.Print("[❓️] fp parameter: ")
	fmt.Scanln(&fp)
	fmt.Print("[❓️] detail parameter: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	detail = scanner.Text()

	unsolved = true

	if loop {
		for unsolved {
			getCaptchaImage()
			processCaptcha()
			submitCaptcha()
		}
	} else {
		getCaptchaImage()
		processCaptcha()
		submitCaptcha()
	}
}
