### Guidelines for use of 'logger' package
Makes it easier to stay coherent with logging along the project

* `logger.Log.Fatal` -> server startup goes wrong. Makes a call to `os.Exit(1)`!
* `logger.Log.Critical` -> error seriously compromises one or more parts of the server
* `logger.Log.Error` -> error on developer side, should not have happened
* `logger.Log.Warning` -> error on client side, should not have happened
* `logger.Log.Info` -> important info, useful for long-term logging
* `logger.Log.Debug` -> debugging purposes
