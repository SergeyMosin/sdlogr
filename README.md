## sdlogr: `logr` implementation

Another implementation of [logr interface](https://github.com/go-logr/logr) with addition of `systemd` specific prefixes (severity levels). This logger is meant for apps/services that are started via `systemd` and want to send their logs to the system `journal`.

`logr.Info` and `logr.Error` messages are prefixed by `SD_INFO` and `SD_ERR` respectively and therefor should be easy to distinguish and filter with `journalctl` command.

```
 // from systemd/sd-daemon.h
 #define SD_EMERG   "<0>"  /* system is unusable */
 #define SD_ALERT   "<1>"  /* action must be taken immediately */
 #define SD_CRIT    "<2>"  /* critical conditions */
 #define SD_ERR     "<3>"  /* error conditions */
 #define SD_WARNING "<4>"  /* warning conditions */
 #define SD_NOTICE  "<5>"  /* normal but significant condition */
 #define SD_INFO    "<6>"  /* informational */
 #define SD_DEBUG   "<7>"  /* debug-level messages */
```