package main

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/restartfu/roles/roles"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
)

func main() {
	err := roles.Load("assets/roles/")
	if err != nil {
		panic(err)
	}
	fmt.Println(roles.All())
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := server.DefaultConfig().Config(log)
	if err != nil {
		log.Fatalln(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()

	for srv.Accept(func(p *player.Player) {
		p.Handle(&handler{p: p})
	}) {

	}
}

type handler struct {
	player.NopHandler
	p *player.Player
}

func (h *handler) HandleChat(ctx *event.Context, message *string) {
	ctx.Cancel()

	owner, ok := roles.ByName("owner")
	if !ok {
		panic("role not found")
	}
	format := text.Colourf("<grey>[</grey>%s<grey>]</grey> %s<grey>:</grey> <white>%s</white>", owner.Coloured(owner.Name()), owner.Coloured(h.p.Name()), *message)
	_, _ = chat.Global.WriteString(format)
}
