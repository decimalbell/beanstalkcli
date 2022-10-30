package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		host string
		port int
	)

	var (
		priority, delay, ttr int
		body                 string
	)

	var (
		id int64
	)

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Usage:       "host",
				Value:       "127.0.0.1",
				Destination: &host,
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "port",
				Value:       11300,
				Destination: &port,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "put",
				Usage: "put",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "priority",
						Usage:       "priority",
						Destination: &priority,
						Required:    true,
					},
					&cli.IntFlag{
						Name:        "delay",
						Usage:       "delay",
						Destination: &delay,
						Required:    true,
					},
					&cli.IntFlag{
						Name:        "ttr",
						Usage:       "ttr",
						Destination: &ttr,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "body",
						Usage:       "body",
						Destination: &body,
						Required:    true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					conn, err := beanstalk.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
					if err != nil {
						return err
					}

					id, err := conn.Put([]byte(body), uint32(priority),
						time.Duration(delay)*time.Second, time.Duration(ttr)*time.Second)
					if err != nil {
						return nil
					}

					fmt.Println(id)
					return nil
				},
			},
			{
				Name:  "peek",
				Usage: "peek",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:        "id",
						Usage:       "id",
						Destination: &id,
						Required:    true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					conn, err := beanstalk.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
					if err != nil {
						return err
					}

					body, err := conn.Peek(uint64(id))
					if err != nil {
						return nil
					}

					fmt.Println(string(body))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
