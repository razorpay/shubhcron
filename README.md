# shubhcron

shubhcron is a auspicious cron runner. It is a drop-in replacement for your existing cron runners (cronie/anacron etc) and allows you to make sure that your jobs only run on auspicious timings.

The syntax is fairly simple:

`* * * * * AUSPICIOUS_FLAG`

The first 5 values are respected as per the standard crontab format. The AUSPICIOUS_FLAG accepts the following values:

## Choghadiya

There are four good Choghadiya, Amrit, Shubh, Labh, and Char to start auspicious work. The cronjob will only run if one of these 4 matches.

## Lagna

You can pick which Lagna you want by using Lagna_Zodiac as the flag.

The auspicious flag must always be in capitals.

- You can use multiple flags as a comma-separated string.
- You can negate a specific flag by using a - before it