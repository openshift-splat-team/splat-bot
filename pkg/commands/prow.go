package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	prowv1 "k8s.io/test-infra/prow/apis/prowjobs/v1"
)

const (
	prowJobsUrl      = "https://prow.ci.openshift.org/prowjobs.js?omit=annotations,decoration_config,pod_spec"
	retrievalTimeout = time.Hour + time.Minute*30
	retrieveEvery    = time.Minute * 29
)

var (
	prowTimer    *time.Timer
	prowJobList  *prowv1.ProwJobList
	periodicJobs []prowv1.ProwJob
	mu           sync.Mutex
)

var ProwAttibutes = Attributes{
	Regex:          `\bprow\b`,
	RequireMention: true,
	Callback: func(evt *slackevents.MessageEvent, args []string) ([]slack.MsgOption, error) {
		startProwRetrievalTimers()

		results, err := queryProwResults(args[1], args[2], prowv1.ProwJobState(args[3]))
		if err != nil {
			return nil, err
		}

		return StringToBlock(results, false), nil
	},
	RequiredArgs: 4,
	HelpMarkdown: "retrieve prow results: `prow [platform] [version] [state]`",
}

func fetchProwJobs() (*prowv1.ProwJobList, error) {
	var prowJobList prowv1.ProwJobList

	// Send HTTP GET request
	response, err := http.Get(prowJobsUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Decode JSON response
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&prowJobList); err != nil {
		return nil, err
	}

	return &prowJobList, nil
}

func startProwRetrievalTimers() {
	var err error

	if prowTimer == nil {
		mu.Lock()
		log.Print("mu.Lock()")
		prowJobList, err = fetchProwJobs()
		if err != nil {
			log.Fatal(err)
		}
		mu.Unlock()
		log.Print("mu.Unlock()")

		prowTimer = time.AfterFunc(retrieveEvery, func() {
			log.Print("time.AfterFunc")

			mu.Lock()
			log.Print("AfterFunc mu.Lock()")
			prowJobList, err = fetchProwJobs()
			if err != nil {
				log.Fatal(err)
			}
			mu.Unlock()
			log.Print("AfterFunc mu.Unlock()")
		})

		time.AfterFunc(retrievalTimeout, func() {
			mu.Lock()
			log.Print("prowTimer Stop()")
			prowTimer.Stop()
			prowJobList = nil
			mu.Unlock()
		})
	}
}

func queryProwResults(platform, version string, prowJobState prowv1.ProwJobState) (string, error) {
	var resultsBuilder strings.Builder
	numToRetrieve := 10

	mu.Lock()
	if prowJobList != nil {
		periodicJobs = make([]prowv1.ProwJob, 0)

		periodicJobs = slices.DeleteFunc(prowJobList.Items, func(p prowv1.ProwJob) bool {
			return p.Spec.Type != prowv1.PeriodicJob || !strings.Contains(p.Spec.Job, "nightly") || !strings.Contains(p.Spec.Job, platform) || !p.Complete() || !strings.Contains(p.Spec.Job, version)
		})

		periodicJobs = slices.DeleteFunc(periodicJobs, func(p prowv1.ProwJob) bool {
			return p.Status.State != prowJobState
		})

		periodicJobLength := len(periodicJobs)
		if periodicJobLength > numToRetrieve {
			periodicJobLength = numToRetrieve
		}

		periodicJobs = periodicJobs[0:periodicJobLength:periodicJobLength]

		for _, j := range periodicJobs {
			fmt.Fprintf(&resultsBuilder, "<%s|%s>\n", j.Status.URL, j.Spec.Job)
		}
	}
	mu.Unlock()
	return resultsBuilder.String(), nil
}
