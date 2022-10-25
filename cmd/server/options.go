package main

import (
	"errors"
	"fetch-me-if-you-read-me/imaginer"
	"fetch-me-if-you-read-me/server"
	"flag"
	"fmt"
	"image/color"
	"os"
)

var (
	host       = flag.String("host", "0.0.0.0", "Host where server will listen")
	port       = flag.String("port", "3000", "Port where server will listen")
	imageColor = flag.String("image-color", "#00FFFF", "Image color")
)

type Options struct {
	Imaginer *imaginer.ImaginerConfs
	Server   *server.ServerConfs
}

func parseOptions() (*Options, error) {
	flag.Parse()

	hostEnv, hostEnvSet := os.LookupEnv("HOST")
	portEnv, portEnvSet := os.LookupEnv("PORT")
	imageColorEnv, imageColorSet := os.LookupEnv("IMAGE_COLOR")

	if hostEnvSet {
		host = &hostEnv
	}

	if portEnvSet {
		port = &portEnv
	}

	if imageColorSet {
		imageColor = &imageColorEnv
	}

	rgbaColor, err := parseHexColor(*imageColor)

	if err != nil {
		return nil, errors.New("image color is not a valid hex value")
	}

	imaginerConf := imaginer.ImaginerConfs{
		Color: &rgbaColor,
	}

	serverConf := server.ServerConfs{
		Host: *host,
		Port: *port,
	}

	return &Options{
		Imaginer: &imaginerConf,
		Server:   &serverConf,
	}, nil
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}
