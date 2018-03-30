# shubhcron

shubhcron is a auspicious cron runner. It is a drop-in replacement for your existing cron runners (cronie/anacron etc) and allows you to make sure that your jobs only run on auspicious timings.

The package comes with a `shubh` binary that you can use to wrap other programs and run them only if the current time is shubh.

For eg:

```bash
# Restart server only on auspicious timings
shubh shutdown --reboot
# Run package uprades only when the stars align
shubh apt-get upgrade --assume-yes
# Run deployments
shubh kubectl set deployments image app application=$VERSION
# Schedule emails to your customers to increase conversions
shubh php artisan email:send
# Schedule money transfers to your preferred charity
shubh sendtoaddress <bitcoinaddress> <amount>
# Run the JVM garbage collector
shubh /usr/bin/bin/jcmd GC.run
```

# Installation

TODO