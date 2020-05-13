package ffmpeg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
)

// Constants for video stream
const (
	FrameX      = 1920
	FrameY      = 1080
	FrameSize   = FrameX * FrameY * 3
	MinimumArea = 3000
)

// CreateH264Pipe creates an ffmpeg pipe.
func CreateH264Pipe(src string) io.ReadCloser {
	ffmpeg := exec.Command(
		"ffmpeg",
		"-f", src,
		// "-capture_cursor", "1",
		// "-i", "1:none",
		"-i", ":0.0",
		"-s", strconv.Itoa(FrameX)+"x"+strconv.Itoa(FrameY),
		"-c:v", "libx264",
		// "-vsync", "2",
		"-framerate", "30",
		"-preset", "veryfast",
		"-tune", "zerolatency",
		"-f", "h264",
		"-pix_fmt", "yuv420p",
		"pipe:1",
	)

	ffmpegOut, err := ffmpeg.StdoutPipe()
	if err != nil {
		panic(err)
	}

	ffmpegErr, err := ffmpeg.StderrPipe()
	if err != nil {
		panic(err)
	}

	// Log any errors
	go func() {
		scanner := bufio.NewScanner(ffmpegErr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	// Kill the spawned process if this program receives an interrupt signal
	go func() {
		<-sigs
		if err := ffmpeg.Process.Kill(); err != nil {
			panic(err)
		}
	}()

	if err := ffmpeg.Start(); err != nil {
		panic(err)
	}

	return ffmpegOut
}