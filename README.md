KeySig
---

Extract useful statistics from keyboard events.

![Single Tile](docs/TimeOfPressReport.png)

Some initial results
---

Generated signatures from two runs of typing the same text (just a random wikipedia article). The metric used to generate these signatures was the time taken from releasing a key to pressing the next one (time-to-next).

##### Run 1
![Run 1](docs/run1.jpeg)

##### Run 2
![Run 2](docs/run2.jpeg)

Goals
---

* Extract metrics from typing data.
* Generate a verifiable signature from the metrics.
* Get the typing data using a KeyLogger that runs in the background.
* Create a mechanism to authenticate a user based on a previous signature.
* Create alert when user fails authentication.

Ideas for metrics
---

1. Histogram of keys typed.
2. Stream of events (this could be the base data actually) such as keydown, keyup.
3. Time between different keys. h -> e -> l -> l -> o, etc
4. Length of line.
5. Wrong spelling.
6. Keypress duration.
7. Different transition periods (word-to-word, char-to-char, etc).

Program structure
---

* Keylogger is running on main thread and is blocking.
* Keylogger provides a channel of events.
* MetricsAnalyzers are consuming this channel (however we need this channel to be broadcast to all of them)
Each metric analyzer is specialized for a certain type of metric.

References
---

* In Linux: https://github.com/MarinX/keylogger
* In Mac: can't keylog (must be active process with Window atm)
* OSX Keylogger in C: https://github.com/caseyscarborough/keylogger
* http://osxbook.com/book/bonus/chapter2/alterkeys/
