Changelog
=========

## 1.4.0
### Unreleased
- Implement ``url`` in action, An URL to "GET" in case service is up/down, for example 'hipchat'.
- debug (print all the headers)
- kill -USR1 shows cleaner info.
- Implement ``Retry`` default to 3, with 1 second interval
- Fix ``Emoji`` and ``msg`` implementation to behave like a list
- Implement timestamp ``when`` RFC3339

## 1.3.0
- ``Insecure`` feature to ignore SSL / self signed certificates.
- ``Stop`` establish a limit on how many times to retry a cmd, ``-1`` loops for ever.
- ``Emoji`` support per action, add/remove icons on email subject.

## 1.2.0
- Improve expect/header match.
- Fix service notification to avoid spamming recipients.

## 1.1.0
- Added -d debug flag.
- Added ``Follow`` setting to avoid/allow following redirects.
