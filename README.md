# Evan

Evan is like heaven for your deployment automation.

## Structure

- **config**: You set up Evan by configuring your applications and how they should be deployed.
  - `Application`:
    - `Repository`: Where the source code of the application is located.
    - `Environment`: Strings describing different places the application can be deployed: "production", "staging", etc.
    - `Strategy`: Strategies determine how the application should be deployed. A strategy consists *preconditions* and *phases*. Once all the preconditions are met it proceeds with executing the phases in order.
- **context**: Evan is stateless. Every time it is invoked it queries your ops infrastructure—GitHub, Heroku, etc.—to determine what's going on and where it may be in the process of executing a strategy.
- **http_handlers**: Provides common handlers conforming to the Go `net/http.Handler` interface for receiving various events and commands.
