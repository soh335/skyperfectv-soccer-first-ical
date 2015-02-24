[![wercker status](https://app.wercker.com/status/68cd212ef7d65141e1abf1d7dce2c433/s/master "wercker status")](https://app.wercker.com/project/bykey/68cd212ef7d65141e1abf1d7dce2c433)

# skyperfetct-soccer-first-ical

generate ical data from http://soccer.skyperfectv.co.jp/static/first/

## HOW TO INSTALL

```
$ go get github.com/soh335/skyperfectv-soccer-first-ical
```

## USAGE

### print categories

```
$ skyperfectv-soccer-first-ical category
```

### print ical

```
$ skyperfectv-soccer-first-ical ical --category=ucl,serie,liga,swiss,uel,premier,bundes,eredivisie --channel=BS242,BS243,BS244,BS245,CS250,BS238,BS241,CS296,CS800,CS801,CS802,CS805 --calname=soccer.skyperfect --liveonly
```

## LICENSE

MIT
