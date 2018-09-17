package schedule

import (
	"sync"
	"time"

	com "github.com/att/deadline/common"
	"github.com/att/deadline/config"
	"github.com/att/deadline/dao"
	"github.com/sirupsen/logrus"
)

var manager *ScheduleManager
var once sync.Once
var log = logrus.New()

// GetManagerInstance will return the singleton of the ScheduleManager object
func GetManagerInstance(cfg *config.Config) *ScheduleManager {
	once.Do(func() {
		manager = &ScheduleManager{
			db:     dao.NewScheduleDAO(cfg),
			rwLock: &sync.RWMutex{},
		}

		manager.schedules = make(map[string]*Schedule)
		manager.subscriptionTable = make(map[string][]*Schedule)
		manager.blueprints = make(chan com.ScheduleBlueprint)

		manager.loadAllSchedules()
	})
	return manager
}

// part of the initialization cycle, this function should only be called once per instance of the
// ScheduleManager.  It is not likely thread safe at this time.
func (manager *ScheduleManager) loadAllSchedules() {

	blueprints, err := manager.db.LoadScheduleBlueprints()
	if err != nil {
		//log that you couldn't load blueprints. return?
	}

	for _, blueprint := range blueprints {
		if err := manager.AddSchedule(blueprint); err != nil {
			//TODO log error
		}
	}

	events, err := manager.db.LoadEvents()
	if err != nil {
		//log that you couldn't load events
	} else {
		for _, event := range events {
			schedules := manager.subscriptionTable[event.Name]
			for _, schedule := range schedules {
				schedule.EventOccurred(&event)
			}
		}
	}
}

// Update updates any schedule currently alive with the event that you pass in
func (manager *ScheduleManager) Update(e *com.Event) {
	manager.rwLock.Lock()
	defer manager.rwLock.Unlock()

	scheds := manager.subscriptionTable[e.Name]

	if scheds == nil {
		// TODO log
	}

	for _, schedule := range scheds {
		schedule.EventOccurred(e)
	}
}

// GetBlueprint gets a blueprint for a schedule given the name of the blueprint
func (manager *ScheduleManager) GetBlueprint(name string) (*com.ScheduleBlueprint, error) {
	return manager.db.GetByName(name)
}

func (manager *ScheduleManager) AddScheduleAndSave(blueprint *com.ScheduleBlueprint) error {
	// TODO rollback the save if the other errors out
	if err := manager.db.Save(blueprint); err != nil {
		return err
	} else if err := manager.AddSchedule(*blueprint); err != nil {
		return err
	} else {
		return nil
	}
}

// AddSchedule adds the schedule to the current list of schedules. If the schedule's start time
// it will become live and the manager will start to evaluate it. Otherwise it will be scheduled
// to become live at that time
func (manager *ScheduleManager) AddSchedule(blueprint com.ScheduleBlueprint) error {

	if schedule, err := FromBlueprint(&blueprint); err != nil {
		return err
	} else if nextTime, err := timingToDuration(blueprint.Timing); err != nil {
		return err
	} else if startTime, err := normailizeTime(blueprint.StartsAt, nextTime); err != nil {
		return err
	} else {
		blueprint.StartsAt = startTime

		// TODO check and log duplicates entries
		manager.rwLock.Lock()
		defer manager.rwLock.Unlock()

		timer := time.NewTimer(nextTime)
		go func() {
			// TODO:bug - what happens when you remove the blueprint/stop the schedule?
			<-timer.C
			manager.AddSchedule(blueprint)
		}()

		manager.schedules[schedule.Name] = schedule
		for subscription := range schedule.SubscribesTo() {
			entry := manager.subscriptionTable[subscription]
			manager.subscriptionTable[subscription] = append(entry, schedule)
		}
	}

	return nil
}

// GetSchedule gets the current running schedule by the given name. If it exists, it'll
// return it, if not, it will return nil.
func (manager *ScheduleManager) GetSchedule(name string) *Schedule {
	manager.rwLock.RLock()
	defer manager.rwLock.Unlock()

	if s, ok := manager.schedules[name]; !ok {
		return nil
	} else {
		return s
	}
}

func normailizeTime(startTime string, timing time.Duration) (string, error) {
	var start time.Time
	var err error

	if start, err = time.Parse(time.RFC3339, startTime); err != nil {
		return "", err
	}

	now := time.Now()
	next := start.Add(timing)
	last := start

	// TODO not the best way to do this if startTime is very far in the past
	for next.Unix() < now.Unix() {
		last = next
		next = last.Add(timing)
	}

	return last.Format(ExpectedTimeLayout), nil
}

func timingToDuration(timing string) (time.Duration, error) {
	if alias, found := TimingAilias[timing]; found {
		return time.ParseDuration(alias)
	} else {
		return time.ParseDuration(timing)
	}
}
