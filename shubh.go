package main

import (
  "fmt"
  "github.com/kelvins/sunrisesunset"
  "io/ioutil"
  "log"
  "math"
  "os"
  "os/exec"
  "strconv"
  "time"
)

type Chowgadhiya int
type Phase int

const (
  Day Phase = iota
  Night
)

const (
  Chal Chowgadhiya = iota
  Amrit
  Kaal
  Labh
  Rog
  Shubh
  Udveg
)

// https://hinduism.stackexchange.com/questions/26242/how-is-the-first-choghadiya-decided
// Golang does not allow constant maps, but a literal map is close enough
var CHOWGADHIYA_LIST = map[Phase]map[time.Weekday][]Chowgadhiya{
  Day: map[time.Weekday][]Chowgadhiya{
    time.Sunday:    []Chowgadhiya{Udveg, Chal, Labh, Amrit, Kaal, Shubh, Rog, Udveg},
    time.Monday:    []Chowgadhiya{Amrit, Kaal, Shubh, Rog, Udveg, Chal, Labh, Amrit},
    time.Tuesday:   []Chowgadhiya{Rog, Udveg, Chal, Labh, Amrit, Kaal, Shubh, Rog},
    time.Wednesday: []Chowgadhiya{Labh, Amrit, Kaal, Shubh, Rog, Udveg, Chal, Labh},
    time.Thursday:  []Chowgadhiya{Shubh, Rog, Udveg, Chal, Labh, Amrit, Kaal, Shubh},
    time.Friday:    []Chowgadhiya{Chal, Labh, Amrit, Kaal, Shubh, Rog, Udveg, Chal},
    time.Saturday:  []Chowgadhiya{Kaal, Shubh, Rog, Udveg, Chal, Labh, Amrit, Kaal},
  },
  Night: map[time.Weekday][]Chowgadhiya{
    time.Sunday:    []Chowgadhiya{Shubh, Amrit, Chal, Rog, Kaal, Labh, Udveg, Shubh},
    time.Monday:    []Chowgadhiya{Chal, Rog, Kaal, Labh, Udveg, Shubh, Amrit, Chal},
    time.Tuesday:   []Chowgadhiya{Kaal, Labh, Udveg, Shubh, Amrit, Chal, Rog, Kaal},
    time.Wednesday: []Chowgadhiya{Udveg, Shubh, Amrit, Chal, Rog, Kaal, Labh, Udveg},
    time.Thursday:  []Chowgadhiya{Amrit, Chal, Rog, Kaal, Labh, Udveg, Shubh, Amrit},
    time.Friday:    []Chowgadhiya{Rog, Kaal, Labh, Udveg, Shubh, Amrit, Chal, Rog},
    time.Saturday:  []Chowgadhiya{Labh, Udveg, Shubh, Amrit, Chal, Rog, Kaal, Labh},
  },
}

/**
 * There is some confusion as to whether Chal is
 * considered Shubh or not, we are avoiding
 */
func isChowgadhiyaConsideredShubh(c Chowgadhiya) bool {
  debug("Picked Chowgadhiya", c)
  return (c == Amrit || c == Shubh || c == Labh)
}

/**
 * Returns the list of Chowgadhiyas in Order for Daytime
 */
func getChowgadhiyaListFromWeekday(day time.Weekday, phase Phase) []Chowgadhiya {
  return CHOWGADHIYA_LIST[phase][day]
}

/**
 * Takes time and returns the correct Chowgadhiya
 */
func getChowgadhiya(t time.Time) Chowgadhiya {
  sunrise, sunset, nextSunrise := getVedicDay(t)

  debug("Next sunrise:", nextSunrise)
  debug("Current time:", t)

  if t.Before(sunrise) || t.After(nextSunrise) {
    panic("current time does not fall between Sunrise and Sunset")
  }

  var baseTime time.Time
  var phase Phase
  var offsetInSeconds float64

  if t.Before(sunset) {
    // Daytime
    phase = Day
    baseTime = sunrise
    offsetInSeconds = (sunset.Sub(sunrise) / 8).Seconds()
  } else {
    // Nighttime
    phase = Night
    baseTime = sunset
    offsetInSeconds = (nextSunrise.Sub(sunset) / 8).Seconds()
    debug("time difference:", nextSunrise.Sub(t).Hours())
  }

  timePassedInCurrentPhase := t.Sub(baseTime).Seconds()
  debug("timePassedInCurrentPhase:", timePassedInCurrentPhase)
  debug("offsetInSeconds:", offsetInSeconds)
  numberOfChowgadhiyaPassed := timePassedInCurrentPhase / offsetInSeconds
  debug("numberOfChowgadhiyaPassed:", numberOfChowgadhiyaPassed)
  chowgadhiyaIndex := int(math.Floor(numberOfChowgadhiyaPassed))
  debug("chowgadhiyaIndex:", chowgadhiyaIndex)
  list := getChowgadhiyaListFromWeekday(sunrise.Weekday(), phase)
  debug("phase:", phase)
  debug("list:", list)
  return list[chowgadhiyaIndex]
}

func getSunriseSunset(t time.Time) (time.Time, time.Time) {
  reference_time := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

  _, offset := t.Zone()

  fractional_offset := (float64(offset) / 60 / 60)

  if fractional_offset > 12 {
    fractional_offset = 12 - fractional_offset
  }

  latitude, _ := strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
  longitude, _ := strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)

  p := sunrisesunset.Parameters{
    Latitude:  latitude,
    Longitude: longitude,
    UtcOffset: fractional_offset,
    Date:      reference_time,
  }

  sunrise, sunset, err := p.GetSunriseSunset()

  if err == nil {
    return sunrise, sunset
  }
  panic("sunrise/sunset calculations failed")
}

/**
 * returns whether now is an auspicious time or not
 */
func isShubh(now time.Time) bool {
  chowgadhiya := getChowgadhiya(now)
  return isChowgadhiyaConsideredShubh(chowgadhiya)
}

func debug(strings ...interface{}) {
  Debug := log.New(os.Stdout,
    "DEBUG:",
    log.Ldate|log.Ltime|log.Lshortfile)

  _, debug := os.LookupEnv("DEBUG")

  if debug != true {
    Debug.SetOutput(ioutil.Discard)
  }

  Debug.Println(strings...)
}

func getVedicDay(now time.Time) (time.Time, time.Time, time.Time) {

  var sunrise, sunset, nextSunrise time.Time

  sunrise, sunset = getSunriseSunset(now)

  yesterday := now.AddDate(0, 0, -1)
  tomorrow := now.AddDate(0, 0, 1)

  loc := now.Location()
  tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.UTC)
  yesterday = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)

  sunrise = time.Date(now.Year(), now.Month(), now.Day(), sunrise.Hour(), sunrise.Minute(), sunrise.Second(), sunrise.Nanosecond(), loc)
  sunset = time.Date(now.Year(), now.Month(), now.Day(), sunset.Hour(), sunset.Minute(), sunset.Second(), sunset.Nanosecond(), loc)

  // Sun has not risen yet
  // So check the sunrise for yesterday
  if now.Before(sunrise) {
    debug("Sun is not yet up, go back to bed")
    nextSunrise, sunset = getSunriseSunset(yesterday)

    sunset = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), sunset.Hour(), sunset.Minute(), sunset.Second(), sunset.Nanosecond(), loc)
    nextSunrise = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), nextSunrise.Hour(), nextSunrise.Minute(), nextSunrise.Second(), nextSunrise.Nanosecond(), loc)
  } else {
    debug("Sun is up, rise and shine")
    // Calculate the sunrise time for tomorrow
    nextSunrise, _ = getSunriseSunset(tomorrow)
    nextSunrise = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), nextSunrise.Hour(), nextSunrise.Minute(), nextSunrise.Second(), nextSunrise.Nanosecond(), loc)
  }

  // Now we have a definite sunrise time for the "vedic day"

  debug("Sunrise:", sunrise)
  debug("Sunset:", sunset)
  debug("Next sunrise:", nextSunrise)

  return sunrise, sunset, nextSunrise
}

func printHelp() {
  // Replacing this with a proper parser is left
  // as an exercise for the reader
  fmt.Println("Usage: shubh command [args...]")
  fmt.Println("  Runs the command only if the time is auspicious")
  fmt.Println("  Exits with status 1 otherwise")
  fmt.Println("  Set SHUBH_WAIT environment variable to wait and run the command instead")
  fmt.Println("  Set DEBUG environment variable for debugging")
}

/**
 * Runs the command if the time is Shubh
 * and exits if it was ran
 */
func runCommand() {
  command := os.Args[1]
  argsWithoutProg := os.Args[2:]

  now := time.Now()

  if isShubh(now) {
    cmd := exec.Command(command, argsWithoutProg...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err == nil {
      os.Exit(0)
    } else {
      fmt.Println("error in executing command. Command:", os.Args[1:])
      os.Exit(255)
    }
  }
}

func main() {
  if len(os.Args) < 2 {
    printHelp()
    os.Exit(0)
  }

  _, wait := os.LookupEnv("SHUBH_WAIT")

  // Since our shubh times are ~90 minutes long
  // we are okay checking every minute
  runCommand()
  if wait {
    debug("Running in wait mode")
    for range time.Tick(10 * time.Second) {
      runCommand()
    }
  }
}
