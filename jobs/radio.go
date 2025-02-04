package jobs

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/blang/mpv"
)

var ipcc *mpv.IPCClient
var c *mpv.Client

var skipDay bool = false
var muD sync.RWMutex

// lofi girl link
const linkDefault string = "https://www.youtube.com/watch?v=jfKfPfyJRdk"

var fnameDefault string = func() string {
	url, _ := url.Parse(linkDefault)
	return url.RawQuery
}()

const volumeIncrement float64 = 2.5

func PlayLink(link string) {
	_, err := url.Parse(link)
	if err != nil {
		log.Println(err)
		return
	}

	playAudio(link, false)
}

func PlayFile(fpath string) {
	if !filepath.IsAbs(fpath) {
		var err error
		fpath, err = filepath.Abs(fpath)
		if err != nil {
			log.Println(err)
		}
	}

	playAudio(fpath, false)
}

func QueueLink(link string) {
	_, err := url.Parse(link)
	if err != nil {
		log.Println(err)
		return
	}

	playAudio(link, true)
}

func QueueFile(fpath string) {
	if !filepath.IsAbs(fpath) {
		var err error
		fpath, err = filepath.Abs(fpath)
		if err != nil {
			log.Println(err)
		}
	}

	playAudio(fpath, true)
}

func PlayNext() {
	chk(c.PlaylistNext())
}

func PlayPrev() {
	chk(c.PlaylistPrevious())
}

func VolumeIncrease() {
	v, err := c.Volume()
	chk(err)
	c.SetProperty("volume", clampvol(v+volumeIncrement))
}

func VolumeDecrease() {
	v, err := c.Volume()
	chk(err)
	c.SetProperty("volume", clampvol(v-volumeIncrement))
}

func Volume() float64 {
	v, err := c.Volume()
    if (err != nil) {
        log.Println(err)
        return 30.
    }
	return v
}

func clampvol(v float64) float64 {
	return min(max(v, 15.), 75.)
}

func playAudio(in string, queue bool) bool {
	var err error
	if queue {
		err = c.Loadfile(in, mpv.LoadFileModeAppendPlay)
	} else {
		err = c.Loadfile(in, mpv.LoadFileModeReplace)
	}

	if errors.Is(err, mpv.ErrInvalidType) {
		return false
	} else if errors.Is(err, mpv.ErrTimeoutRecv) || errors.Is(err, mpv.ErrTimeoutSend) {
		log.Panic(err)
	}

	c.SetPause(false)

	return true
}

// scheduled turn on radio default, won't run if skipped for the day
func TurnOnRadio() {
	muD.RLock()
	defer muD.RUnlock()
	if skipDay {
		return
	}

	fname, _ := c.Filename()
	nofile := fname == "<nil>"

	devices := ScanLocalDevices()
	noCon := len(devices) == 0
	if noCon {
		if !nofile {
			TurnOffRadio()
		}
		return
	}

	// defaultPlaying := false
	// if idx := strings.IndexRune(fname, '?'); idx >= 0 {
	//     qname := fname[idx+1:]
	//     defaultPlaying = qname == fnameDefault
	// }

	if !nofile { // file playing, and device connected
		return
	}

	log.Printf("Device connected, playing default link: %s\n", linkDefault)
	playAudio(linkDefault, false)
}

func TurnOffRadio() {
	err := c.SetPause(true)
	chk(err)
	c.Loadfile("", mpv.LoadFileModeReplace)
}

func SkipToday() {
	if skipDay {
		return
	}

	muD.Lock()
	skipDay = true
	muD.Unlock()
}

func ResetSkipDay() {
	muD.Lock()
	skipDay = false
	muD.Unlock()
}

func dummy() {
	time.Sleep(time.Duration(2 * time.Second))
	PlayLink(linkDefault)
}

func CurrentPlaying() string {
    fname, _ := c.Filename()

	if idx := strings.IndexRune(fname, '='); idx >= 0 {
		qname := fname[idx+1:]
        qname = strings.TrimSuffix(qname, "\"")
        return fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", qname)
	}

    return ""
}

func startMpv() {
	cmd := exec.Command(
		"mpv",
		"--idle",
		"--no-video",
		fmt.Sprintf("--input-ipc-server=%s", os.Getenv("MPVSOCKET_PATH")),
	)
	go func() {
		chk(cmd.Start())
		chk(cmd.Wait())
		chk(cmd.Process.Release())
	}()
}

func startClient() {
	time.Sleep(2 * time.Second)
	ipcc = mpv.NewIPCClient(os.Getenv("MPVSOCKET_PATH")) // Lowlevel client
	c = mpv.NewClient(ipcc)                              // Highlevel client, can also use RPCClient
}

func init() {
	startMpv()
	go startClient()
	// dummy()
}
