package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/log/level"
)

const (
	pollAttempts  = 3
	pollInterval  = 10 * time.Second
	pollRetryWait = 3 * time.Second
	rawstatusURL  = "http://%s/RAWSTATUS"
)

type herpstat struct {
	NextAllowedPoll time.Time
	info            *info
}

// newHerpstat returns a new instance of the herpstat struct.
// [exporter.herpstat.nextAllowedPoll] is set to 10 seconds in the past to ensure that the first [exporter.herpstat.pollingTooQuickly()]
// call will return true
func newHerpstat() *herpstat {
	return &herpstat{
		NextAllowedPoll: time.Now().Add(-pollInterval),
		info: &info{
			system:  &system{},
			outputs: &[]output{},
		},
	}
}

// Polls the Herpstat SpyderWeb device to retrieve its status info, storing it in [herpstat.exporter.info]. They can
// sometimes be a little finicky and come back with invalid JSON data. If that happens, we'll try polling a total of
// three times, waiting three seconds in between each poll.
func (h *herpstat) poll() bool {
	if h.pollingTooQuickly() {
		level.Warn(logger).Log("msg", fmt.Sprintf("Polling too quickly! Please set polling interval to %.0f seconds.", pollInterval.Seconds()))
		level.Warn(logger).Log("msg", fmt.Sprintf("See http://%s/handleAdminControls for more information.", *herpstatAddress))

		return true
	}

	retried := false

	// herpstats can sometimes come back with weird data. we'll retry a total of three times, waiting
	// 3 seconds in between attempts if that happens.
	for i := 1; i <= pollAttempts; i++ {
		level.Debug(logger).Log("msg", fmt.Sprintf("poll attempt %d/%d", i, pollAttempts))

		rawstatus := h.getRawstatus()
		if rawstatus == nil {
			retried = true

			maybeWait(i, pollAttempts)

			continue
		}

		if err := json.Unmarshal(*rawstatus, &h.info); err != nil {
			retried = true

			level.Warn(logger).Log("msg", fmt.Sprintf("unable to unmarshal JSON: %s", err.Error()))
			level.Warn(logger).Log("msg", *rawstatus)
			maybeWait(i, pollAttempts)

			continue
		}

		if retried {
			level.Info(logger).Log("msg", "Successfully unmarshalled JSON after %d attempts", i)
		}

		// success! we can poll again in 10 seconds
		h.NextAllowedPoll = time.Now().Add(pollInterval)

		return true
	}

	level.Error(logger).Log("msg", fmt.Sprintf("unable to get data from device after %d attempts", pollAttempts))

	return false
}

// Waits for [herpstat.exporter.pollRetryWait] (3 seconds), trying again if we're not done with the loop.
func maybeWait(i, limit int) {
	if i > limit {
		return
	}

	level.Warn(logger).Log("msg", fmt.Sprintf("Waiting %.0f seconds before trying again...", pollRetryWait.Seconds()))
	time.Sleep(pollRetryWait)
}

// checks whether the next allowed poll time [herpstat.exporter.nextAllowedPoll] is after the current time
// [time.Now]. If the next allowed poll time is in the future, it means that the polling is happening too
// quickly and should be limited to a certain interval. Per the Herpstat SpyderWeb Admin Page:
// "Polling this page should be limited to intervals of 10 seconds or greater to prevent blocking other tasks"
func (h *herpstat) pollingTooQuickly() bool {
	return h.NextAllowedPoll.After(time.Now())
}

// Performs an HTTP request to the `/RAWSTATUS` endpoint of the Herpstat SpyderWeb and returns its raw, byte-encoded
// body.
func (h *herpstat) getRawstatus() *[]byte {
	level.Debug(logger).Log("msg", "getting data from herpstat")

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf(rawstatusURL, *herpstatAddress), http.NoBody)
	if err != nil {
		level.Error(logger).Log("msg", "unable to make new request object:", "err", err)
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		level.Error(logger).Log("msg", "problem making request", "err", err)
		return nil
	}
	defer resp.Body.Close()

	rawStatus, err := io.ReadAll(resp.Body)
	if err != nil {
		level.Error(logger).Log("msg", "problem reading response body", err)
		return nil
	}

	level.Debug(logger).Log("rawstatus", rawStatus)

	return &rawStatus
}
