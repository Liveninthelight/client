package client

import (
	"context"
	"errors"
	"strconv"

	"github.com/keybase/cli"
	"github.com/keybase/client/go/libcmdline"
	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/chat1"
)

type CmdChatDownload struct {
	libkb.Contextified
	tlf        string
	public     bool
	messageID  uint64
	outputFile string
	preview    bool
}

func newCmdChatDownload(cl *libcmdline.CommandLine, g *libkb.GlobalContext) cli.Command {
	return cli.Command{
		Name:         "download",
		Usage:        "Download an attachment from a conversation",
		ArgumentHelp: "<conversation> <attachment id> [-o filename]",
		Action: func(c *cli.Context) {
			cmd := &CmdChatDownload{Contextified: libkb.NewContextified(g)}
			cl.ChooseCommand(cmd, "download", c)
		},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "p, preview",
				Usage: "Download preview",
			},
			cli.BoolFlag{
				Name:  "public",
				Usage: "Download attachment from public conversation",
			},
			cli.StringFlag{
				Name:  "o, outfile",
				Usage: "Specify output file (stdout default)",
			},
		},
	}
}

func (c *CmdChatDownload) ParseArgv(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return errors.New("usage: keybase chat download <conversation> <attachment id> [-o filename]")
	}
	c.tlf = ctx.Args()[0]
	id, err := strconv.ParseUint(ctx.Args()[1], 10, 64)
	if err != nil {
		return err
	}
	c.messageID = id
	c.outputFile = ctx.String("outfile")
	c.preview = ctx.Bool("preview")
	c.public = ctx.Bool("public")

	return nil
}

func (c *CmdChatDownload) Run() error {
	opts := downloadOptionsV1{
		Channel: ChatChannel{
			Name:   c.tlf,
			Public: c.public,
		},
		MessageID: chat1.MessageID(c.messageID),
		Output:    c.outputFile,
		Preview:   c.preview,
	}
	api := &CmdChatAPI{
		Contextified: libkb.NewContextified(c.G()),
	}
	reply := api.DownloadV1(context.Background(), opts)
	if reply.Error != nil {
		return errors.New(reply.Error.Message)
	}
	return nil
}

func (c *CmdChatDownload) GetUsage() libkb.Usage {
	return libkb.Usage{
		API:       true,
		KbKeyring: true,
		Config:    true,
	}
}
