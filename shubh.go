package main

import (
  "fmt"
  "github.com/kelvins/sunrisesunset"
  "io/ioutil"
  "log"
  "math"
  "os"
  "os/exec"
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

/**
 * There is some confusion as to whether Chal is
 * considered Shubh or not, we are avoiding
 */
func IsChowgadhiyaConsideredShubh(c Chowgadhiya) bool {
  if c == Amrit || c == Shubh || c == Labh {
    return true
  }
  return false
}

/**
 * Returns the list of Chowgadhiyas in Order for Daytime
 */
func getChowgadhiyaListFromWeekday(day time.Weekday, phase Phase) []Chowgadhiya {
  var daytime = make(map[time.Weekday][]Chowgadhiya, 7)
  var nighttime = make(map[time.Weekday][]Chowgadhiya, 7)
  // https://hinduism.stackexchange.com/questions/26242/how-is-the-first-choghadiya-decided
  daytime[time.Sunday] = []Chowgadhiya{Udveg, Chal, Labh, Amrit, Kaal, Shubh, Rog, Udveg}
  daytime[time.Monday] = []Chowgadhiya{Amrit, Kaal, Shubh, Rog, Udveg, Chal, Labh, Amrit}
  daytime[time.Tuesday] = []Chowgadhiya{Rog, Udveg, Chal, Labh, Amrit, Kaal, Shubh, Rog}
  daytime[time.Wednesday] = []Chowgadhiya{Labh, Amrit, Kaal, Shubh, Rog, Udveg, Chal, Labh}
  daytime[time.Thursday] = []Chowgadhiya{Shubh, Rog, Udveg, Chal, Labh, Amrit, Kaal, Shubh}
  daytime[time.Friday] = []Chowgadhiya{Chal, Labh, Amrit, Kaal, Shubh, Rog, Udveg, Chal}
  daytime[time.Saturday] = []Chowgadhiya{Kaal, Shubh, Rog, Udveg, Chal, Labh, Amrit, Kaal}

  nighttime[time.Sunday] = []Chowgadhiya{Shubh, Amrit, Chal, Rog, Kaal, Labh, Udveg, Shubh}
  nighttime[time.Monday] = []Chowgadhiya{Chal, Rog, Kaal, Labh, Udveg, Shubh, Amrit, Chal}
  nighttime[time.Tuesday] = []Chowgadhiya{Kaal, Labh, Udveg, Shubh, Amrit, Chal, Rog, Kaal}
  nighttime[time.Wednesday] = []Chowgadhiya{Udveg, Shubh, Amrit, Chal, Rog, Kaal, Labh, Udveg}
  nighttime[time.Thursday] = []Chowgadhiya{Amrit, Chal, Rog, Kaal, Labh, Udveg, Shubh, Amrit}
  nighttime[time.Friday] = []Chowgadhiya{Rog, Kaal, Labh, Udveg, Shubh, Amrit, Chal, Rog}
  nighttime[time.Saturday] = []Chowgadhiya{Labh, Udveg, Shubh, Amrit, Chal, Rog, Kaal, Labh}

  if phase == Day {
    return daytime[day]
  }

  return nighttime[day]
}

/**
 * Takes time and returns the correct Chowgadhiya
 */
func getChowgadhiya(t time.Time) Chowgadhiya {
  sunrise, sunset, nextSunrise := GetVedicDay(t)

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
    offsetInSeconds = (nextSunrise.Sub(t) / 8).Seconds()
  }

  timePassedInCurrentPhase := t.Sub(baseTime).Seconds()
  numberOfChowgadhiyaPassed := timePassedInCurrentPhase / offsetInSeconds
  chowgadhiyaIndex := int(math.Floor(numberOfChowgadhiyaPassed))
  list := getChowgadhiyaListFromWeekday(sunrise.Weekday(), phase)
  return list[chowgadhiyaIndex]
}

func GetSunriseSunset(t time.Time) (time.Time, time.Time) {
  Error := log.New(os.Stderr,
    "ERROR: ",
    log.Ldate|log.Ltime|log.Lshortfile)

  t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

  p := sunrisesunset.Parameters{
    // Hardcoded to Bangalore right now
    // Should be changed to Ayodhya
    Latitude:  12.9716,
    Longitude: 77.5946,
    // TODO: Use t.Zone() instead
    // And make sure that the sign is correct
    UtcOffset: 5.5,
    Date:      t,
  }
  sunrise, sunset, err := p.GetSunriseSunset()

  if err == nil {
    return sunrise, sunset
  }

  _, debug := os.LookupEnv("DEBUG")
  if debug {
    Error.Panicln("Sunrise/Sunset calculations failed")
  }
  panic("")
}

/**
 * returns whether now is an auspicious time or not
 */
func isShubh(now time.Time) bool {
  chowgadhiya := getChowgadhiya(now)
  return IsChowgadhiyaConsideredShubh(chowgadhiya)
}

func GetVedicDay(now time.Time) (time.Time, time.Time, time.Time) {
  Debug := log.New(os.Stdout,
    "DEBUG: ",
    log.Ldate|log.Ltime|log.Lshortfile)

  _, debug := os.LookupEnv("DEBUG")

  if debug != true {
    Debug.SetOutput(ioutil.Discard)
  }

  var sunrise, sunset, nextSunrise time.Time

  sunrise, sunset = GetSunriseSunset(now)

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
    Debug.Println("sun is not yet up, go back to bed")
    nextSunrise, sunset = GetSunriseSunset(yesterday)

    sunset = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), sunset.Hour(), sunset.Minute(), sunset.Second(), sunset.Nanosecond(), loc)
    nextSunrise = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), nextSunrise.Hour(), nextSunrise.Minute(), nextSunrise.Second(), nextSunrise.Nanosecond(), loc)
  } else {
    Debug.Println("sun is up, rise and shine")
    // Calculate the sunrise time for tomorrow
    nextSunrise, _ = GetSunriseSunset(tomorrow)
    nextSunrise = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), nextSunrise.Hour(), nextSunrise.Minute(), nextSunrise.Second(), nextSunrise.Nanosecond(), loc)
  }

  // Now we have a definite sunrise time for the "vedic day"

  Debug.Println("Sunrise     :", sunrise)
  Debug.Println("Sunset      :", sunset)
  Debug.Println("Next Sunrise:", nextSunrise)

  return sunrise, sunset, nextSunrise
}

func printHelp() {
  // Replacing this with a proper parser is left
  // as an exercise for the reader
  fmt.Println("Usage: shubh command [args...]")
  fmt.Println("  runs the command only if the time is auspicious")
  fmt.Println("  exits with status 1 otherwise")
  fmt.Println("  set SHUBH_WAIT environment variable to wait and run the command instead")
  fmt.Println("  set DEBUG environment variable for debugging")
}

func main() {
  if len(os.Args) < 2 {
    printHelp()
    os.Exit(0)
  }

  command := os.Args[1]
  argsWithoutProg := os.Args[2:]
  now := time.Now()

  if isShubh(now) {
    cmd := exec.Command(command, argsWithoutProg...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run()
  } else {
    os.Exit(1)
  }

}
