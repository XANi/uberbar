package pomodoro

import (
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/uber"
	"github.com/XANi/uberstatus/util"
	"go.uber.org/zap"

	//	"gopkg.in/yaml.v1"
	"fmt"
	"time"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

// set up a pluginConfig struct
type pluginConfig struct {
	Prefix         string
	PomodoroTime   int
	ShortBreakTime int
	LongBreakTime  int
}

type plugin struct {
	l           *zap.SugaredLogger
	cfg         pluginConfig
	pomodoroEnd time.Time
	breakEnd    time.Time
	nextTs      time.Time
	state       int
	pomodoros   int
}

const (
	stopped = iota
	inPomodor
	inBreakStart
	inShortBreak
	inLongBreak
	inBreakEnd
)

func New(cfg uber.PluginConfig) (z uber.Plugin, err error) {
	p := &plugin{
		l: cfg.Logger,
	}
	p.cfg, err = loadConfig(cfg.Config)
	return p, nil
}
func (p *plugin) Init() error {
	return nil
}

func (p *plugin) GetUpdateInterval() int {
	return 999
}
func (p *plugin) UpdatePeriodic() uber.Update {
	var update uber.Update
	// TODO precompile and preallcate
	// example on how to allow UpdateFromEvent to display for some time
	// without being overwritten by periodic updates.
	// We set up ts in our plugin, update it in UpdateFromEvent() and just wait if it is in future via helper function

	util.WaitForTs(&p.nextTs)
	update.Markup = `pango`
	update.Color = `#66cc66`
	p.nextStateFromTime()
	switch p.state {
	case stopped:
		p.pomodoroEnd = time.Now().Add(time.Duration(p.cfg.PomodoroTime) * time.Minute)
		update.FullText = "stop, click to start"
		update.Color = `#cccc66`
	case inPomodor:
		diff := p.pomodoroEnd.Sub(time.Now())
		update.FullText = fmt.Sprintf(`<span foreground="#ff0000">🍅</span>: %s`, util.FormatDuration(diff))
		update.Color = `#ccccff`
	case inBreakStart:
		update.FullText = fmt.Sprintf(`<span foreground="#000000" background="#aa0000">%d🍅</span><span background="#aa0000">BREAK:</span>`, p.pomodoros)
	case inBreakEnd:
		update.FullText = `<span foreground="#000000" background="#aa0000">⌛</span><span background="#aa0000">END</span>`
	case inShortBreak:
		diff := p.breakEnd.Sub(time.Now())
		update.FullText = fmt.Sprintf("⌛: %s", util.FormatDuration(diff))
		update.Color = `#ccffcc`
	case inLongBreak:
		diff := p.breakEnd.Sub(time.Now())
		update.FullText = fmt.Sprintf("⏲: %s", util.FormatDuration(diff))
		update.Color = `#ccccff`
	default:
		update.FullText = "wtf"
	}
	return update
}

func (p *plugin) UpdateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	update.Markup = `pango`
	if e.Button == 1 {
		p.nextStateFromClick()
	}
	switch p.state {
	case stopped:
		update.FullText = "pomodoro stopped"
		update.Color = `#cccc66`
	case inPomodor:
		if time.Now().After(p.pomodoroEnd) {
			update.FullText = "pomodoro ended!"
			update.Color = `#cccc66`
		} else {
			update.FullText = fmt.Sprintf(`<span foreground="#ff0000">🍅</span>: %d`, p.pomodoros)
			update.Color = `#ccccff`
		}
	case inShortBreak:
		update.FullText = "short break start"
	case inLongBreak:
		update.FullText = "long break start"
	default:
		update.FullText = "wtf"
	}
	// set next TS updatePeriodic will wait to.
	p.nextTs = time.Now().Add(time.Second * 3)
	return update
}

func (p *plugin) nextStateFromClick() {
	switch p.state {
	case stopped:
		p.state = inPomodor
		p.pomodoroEnd = time.Now().Add(time.Duration(p.cfg.PomodoroTime) * time.Minute)
	case inPomodor:
	case inBreakStart:
		if (p.pomodoros % 4) == 3 {
			p.breakEnd = time.Now().Add(time.Duration(p.cfg.LongBreakTime) * time.Minute)
			p.state = inLongBreak
		} else {
			p.breakEnd = time.Now().Add(time.Duration(p.cfg.ShortBreakTime) * time.Minute)
			p.state = inShortBreak
		}
	case inBreakEnd:
		p.pomodoroEnd = time.Now().Add(time.Duration(p.cfg.PomodoroTime) * time.Minute)
		p.state = inPomodor
	case inShortBreak:
	case inLongBreak:
	default:
		p.l.Warnf("out of state machine: %d", p.state)
		p.state = stopped
	}
}
func (p *plugin) nextStateFromTime() {
	switch p.state {
	case stopped:
	case inPomodor:
		if time.Now().After(p.pomodoroEnd) {
			p.pomodoros++
			p.state = inBreakStart
		}
	case inShortBreak:
		if time.Now().After(p.breakEnd) {
			p.state = inBreakEnd
		}
	case inLongBreak:
		if time.Now().After(p.breakEnd) {
			p.state = inBreakEnd
		}
	}
}

// parse received structure into pluginConfig
func loadConfig(c config.PluginConfig) (pluginConfig, error) {
	var cfg pluginConfig
	cfg.Prefix = "ex: "
	cfg.PomodoroTime = 25
	cfg.ShortBreakTime = 5
	cfg.LongBreakTime = 15
	return cfg, c.GetConfig(&cfg)
}
