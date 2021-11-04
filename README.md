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
)

func main() {
    log := logrus.New()
    newLogger := logger.New(log, gormLogger.Config{
        SlowThreshold:  time.Second,    // Slow SQL threshold
        SkipErrRecordNotFound: true,    // Skip ErrRecordNotFound error for logger
        SourceField:    "source",    //  Source field in config which is recorded file name and line number of the current error
        ModuleName:     "my-gorm-logger",    // Name of the module in the log. Default value: "gorm"
        LogLevel:   logger.Error, // Log level. Default value: gormLogger.Info
    })

    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
        Logger: newLogger,
    })
}
```

## License

© CryptoMath, since 2021

Released under the [MIT License](https://github.com/mathandcrypto/cryptomath-gorm-logger/blob/master/LICENSE)