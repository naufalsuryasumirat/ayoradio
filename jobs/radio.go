package jobs

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/blang/mpv"
)

var ipcc *mpv.IPCClient
var c *mpv.Client

var skipDay bool = false
var muD sync.RWMutex

func PlayAudioLink(link string) {
    _, err := url.Parse(link)
    if err != nil {
        log.Println(err.Error())
        return
    }

    playAudio(link, true)
}

func PlayAudioFile(fpath string) {
    if !filepath.IsAbs(fpath) {
		var err error
		fpath, err = filepath.Abs(fpath)
		chk(err)
	}

    playAudio(fpath, true)
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

    // paused, _ := c.Pause()
    // fname, _ := c.Filename()
    idling, _ := c.Idle()
    if !idling {
        return
    }

    devices := ScanLocalDevices()
    if devices == nil {
        TurnOffRadio()
        return
    }

	// const linkDefault = "https://www.youtube.com/watch?v=jfKfPfyJRdk"
 //    playAudio(linkDefault, false)
    PlayAudioFile("./tmp/sample.wav")
}

// TEST: check if sets idle
func TurnOffRadio() {
    err := c.SetPause(true)
    chk(err)
    c.SetProperty("idle", true)
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
	PlayAudioFile("./tmp/sample.wap")
    time.Sleep(time.Duration(2*time.Second))
    PlayAudioFile("./tmp/sample.mp3")
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
	// startMpv()
    startClient()
	// dummy()
}
