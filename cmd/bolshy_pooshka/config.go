package main

import (
	"math/rand"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/satori/go.uuid"
)

type Config struct {
	DbDriver   string
	DataSource string
	Stages     []Stage
}

type Stage struct {
	StageName   string        `yaml:"stage"` // used as a part of metric name
	RPS         float32       // 0 - infinity
	Concurrency int           // How many repeatable requests must be run in parallel, 0 is 1
	Duration    time.Duration /*
		0 - end as soon as all the RunOnce queries done
		duration - obvious
		set a huge duration to run until interrupted
	*/
	RunOnce []*Query    // executed one by one
	Repeat  []*Scenario // executed in parallel according to their probability
	Pause   bool        // Do not step to the next stage automatically

}

type Scenario struct {
	ScenarioName string   `yaml:"scenario"` // used as a part of metric name
	Queries      []*Query // Queries, run sequentially
	Probability  float32  // 0 - never, 1 - each time, ignored for RunOnce
}

type Query struct {
	QueryName  string   `yaml:"query"` // used as a part of metric name
	SQL        string   // SQL itself
	Params     []*Param // Parameters for query placeholders
	Update     bool     // This query is DB update
	RandRepeat int      // Repeat randomly when used in a scenario, ignore otherwise
}

type Param struct {
	ParamName string `yaml:"param"`
	Type      string
	Generator string
}

var globalConfig Config

type Template struct {
	*template.Template
}

func (t *Template) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var strVal string

	err := unmarshal(&strVal)
	if err != nil {
		return err
	}

	newT, err := template.New(strVal).Parse(strVal)
	if err != nil {
		return err
	}

	*t = Template{newT}
	return nil
}

type QueryData struct {
	Rand1    int64
	Rand2    int64
	UserName string
	Inc1     int64
}

var inc1 = int64(0)

func (d *QueryData) Init() *QueryData {
	d.Rand1 = rand.Int63()
	d.Rand2 = rand.Int63()
	d.UserName = uuid.NewV4().String()
	d.Inc1 = atomic.AddInt64(&inc1, 1)
	return d
}
