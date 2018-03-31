# shubhcron

Shubhcron is an auspicious cron runner. It is a drop-in replacement for your existing cron runners (cronie/anacron etc) and allows you to make sure that your jobs only run on auspicious timings.

The package comes with a `shubh` binary that you can use to wrap other programs and run them only if the current time is shubh.

For eg:

```bash
# Restart server, but if the time is auspicious
shubh shutdown --reboot
# Run package uprades, but only if the stars align
shubh apt-get upgrade --assume-yes
# Run deployments, but only if Saturn is in the fifth house
shubh kubectl set deployments image app application=$VERSION
# Schedule emails, but only if the Vedas indicate high conversions
shubh php artisan email:send
# Schedule money transfers to your preferred charity, but only if luck favours
shubh sendtoaddress <bitcoinaddress> <amount>
# Run the JVM garbage collector, but only if all is well
shubh /usr/bin/bin/jcmd GC.run
# Sign documents, but wait till the time is right
SHUBH_WAIT=1 shubh gpg --sign
```

## Cron Usage

If you have the `shubh` binary available, you can prepend all your important jobs with `shubh` in your crontab to run them only if the time is right.

```
# Try this job every 15 minutes on the 31st of March after 6pm
*/15 18-23 31 MAR * /usr/bin/shubh rake finance:closing
```

If you are using the `shubhcron` package, you can omit the `/usr/bin/shubh`:

```
# Set the correct location of your server
LATITUDE=26.7880
LONGITUDE=82.1986
# Attempts to send a mail every 5 minutes, only runs if the time is shubh
*/5 * * * * sendmail --subject "shubh labh"
```

You can also pass an extra environment variable `SHUBH_WAIT=1` to sleep till the time is shubh instead of exiting.

## Installation

### Stand-Alone

If you'd prefer to not replace your system's `cron` package, you can download the `shubh` binary from the [releases page][releases] for your prefered OS and use it directly as a prefix in your cronjobs.

### Ubuntu

We have a ready [Debian and Ubuntu package][releases] available that replaces the system `cron` package.

### Alpine Linux

An Alpine Linux package is in the Docker pipeline.

### Arch Linux

Package for Arch Linux (based on `cronie`) is available at [AUR][aur].

### Kubernetes

We plan to release the kubernetes controller that lets you create ShubhCronJob resources soon.

## Reading the Stars

In order to figure out the best shubh muhurat, `shubh` calculates the vedic day (which starts from sunrise and lasts till the next sunrise), calculates the correct चौगाड़िया and then uses a vedic lookup table (memoized) to determine if the time is auspiciousness or not.

The current system time is used for all calculations. The latitude and longitude, required for calculation of sunrise and sunset times, are picked up from the environment. If not available, the coordinates of Ayodhya are used by default.
The [library]((https://github.com/kelvins/sunrisesunset)) used for calculating sunrise/sunset times implements the same algorithm published by the [National Oceanic & Atmospheric Administration (NOAA)][noaa] and includes corrections for Atmospheric refraction.

## Web

We offer a web-service as well to calculate the auspiciousness of a moment. Just make a request to the [ShubhCron Pandit API endpoint](https://shubhcron-pandit.herokuapp.com/chowgadhiya) to get the current chowgadhiya (along with the next few).
Note that the webservice assumes Ayodhya to be the location.

# License

`shubh` and `shubhcron` are licensed under the MIT License. See `LICENSE` for more details. All code in this repo is released to the spiritual domain.

[releases]: https://github.com/razorpay/shubhcron/releases/latest
[aur]: https://aur.archlinux.org/packages/shubhcron
[noaa]: http://www.noaa.gov/