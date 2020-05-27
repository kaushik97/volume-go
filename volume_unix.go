// +build !windows

package volume

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func execCmd(cmdArgs []string) ([]byte, error) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = append(os.Environ(), cmdEnv()...)
	out, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf(`failed to execute "%v" (%+v)`, strings.Join(cmdArgs, " "), err)
	}
	return out, err
}

// GetVolume returns the current volume (0 to 100).
func GetVolume(outputdevice string) (int, error) {
	out, err := execCmd(getVolumeCmd(outputdevice))
	if err != nil {
		return 0, err
	}
	return parseVolume(string(out))
}

// SetVolume sets the sound volume to the specified value.
func SetVolume(volume int, outputdevice string) error {
	if volume < 0 || 100 < volume {
		return errors.New("out of valid volume range")
	}
	_, err := execCmd(setVolumeCmd(volume, outputdevice))
	return err
}

// IncreaseVolume increases (or decreases) the audio volume by the specified value.
func IncreaseVolume(diff int, outputdevice string) error {
	_, err := execCmd(increaseVolumeCmd(diff, outputdevice))
	return err
}

// GetMuted returns the current muted status.
func GetMuted(outputdevice string) (bool, error) {
	out, err := execCmd(getMutedCmd(outputdevice))
	if err != nil {
		return false, err
	}
	return parseMuted(string(out))
}

// Mute mutes the audio.
func Mute(outputdevice string) error {
	_, err := execCmd(muteCmd(outputdevice))
	return err
}

// Unmute unmutes the audio.
func Unmute(outputdevice string) error {
	_, err := execCmd(unmuteCmd(outputdevice))
	return err
}

// SidetoneMute mutes the sidetone
func SidetoneMute(inputdevice, cardnumber string) error {
	args := []string{"amixer", "-c", cardnumber, "set", inputdevice, "mute"}
	_, err := execCmd(args)
	return err
}

// SidetoneUnmute unmutes the sidetone
func SidetoneUnmute(inputdevice, cardnumber string) error {
	args := []string{"amixer", "-c", cardnumber, "set", inputdevice, "unmute"}
	_, err := execCmd(args)
	return err
}

// PulseMute mutes the speaker only for voice application
func PulseMute(sinkId string) error {
	args := []string{"pacmd", "set-sink-input-mute", sinkId, "true"}
	_, err := execCmd(args)
	return err
}

// PulseUnmute unmutes the speaker only for voice application
func PulseUnmute(sinkId string) error {
	args := []string{"pacmd", "set-sink-input-mute", sinkId, "false"}
	_, err := execCmd(args)
	return err
}
