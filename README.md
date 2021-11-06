# CryptoMath GORM Logger

A customized [GORM](https://gorm.io/) logger that implements the [appropriate interface](https://gorm.io/docs/logger.html#Customize-Logger) and uses [Logrus](https://github.com/sirupsen/logrus) to output logs.

## Install
```shell
go get github.com/mathandcrypto/cryptomath-gorm-logger
```

## Basic usage
```go
package main

import (
    "github.com/mathandcrypto/cryptomath-gorm-logger"
    "github.com/sirupsen/logrus"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    gormLogger "gorm.io/gorm/logger"
    "time"
)

func main() {
    log := logrus.New()
    newLogger := logger.New(log, logger.Config{
        SlowThreshold:  time.Second,
        SkipErrRecordNotFound: true,
        LogLevel:   gormLogger.Error,
    })

    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
        Logger: newLogger,
    })
}
```

## Logger сonfiguration
When creating in the `logger.New` method, the second parameter specifies the configuration of the logger.
The `logger.Config` structure has the following fields:

| Parameter               | Type                  | Default value     | Description |
| ----------------------- | --------------------- | ----------------- | ----------- |
| SlowThreshold           | `time.Duration`       |                   | If the sql query time exceeds this value, a warning log about the slow sql query time will be output |
| SkipErrRecordNotFound   | `bool`                | `false`            | Skip `ErrRecordNotFound` error for logger |
| SourceField             | `string`              |                   | Source field in config which is recorded file name and line number of the current error |
| ModuleName              | `string`              | `"gorm"`          | Name of the module in the log |
| LogLevel                | `gormLogger.LogLevel` | `gormLogger.Info` | [Log level](https://gorm.io/docs/logger.html#Log-Levels) |

## License

© CryptoMath, since 2021

Released under the [MIT License](https://github.com/mathandcrypto/cryptomath-gorm-logger/blob/master/LICENSE)