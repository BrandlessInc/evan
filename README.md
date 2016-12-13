# Evan

Evan is like [heaven][] for your deployment automation.

[heaven]: https://github.com/atmos/heaven

## Structure

- **config**: You set up Evan by configuring your applications and how they should be deployed.
  - `Application`:
    - `Repository`: Where the source code of the application is located.
    - `Environments`: List of strings describing different places the application can be deployed: "production", "staging", etc.
    - `DeployEnvironment`: Function that returns a strategy for deploying the application to a given environment.
  - `Strategy`: Strategies determine how the application should be deployed. A strategy consists *preconditions* and *phases*. It executes all the preconditions first; if they all pass then it executes the phases. Preconditions are repeatable whereas phases execute the actual deployment.
- **preconditions**: Built-in preconditions.
- **phases**: Built-in phases.
- **stores**: Persist the state of deployments during and after their execution. This allows the system to report the progress of deployments and keep track of deployments after they go out.
- **http_handlers**: Provides common handlers conforming to the Go `net/http.Handler` interface for receiving various events and commands.
  - **rest_json**: These handlers make it easy to build a REST'ish JSON API for creating, managing, and querying deployments.
- **common**: Shared protocols for communicating information and functionality between the subsystems that make up Evan.

## License

Released under the MIT license, see [LICENSE](LICENSE) for details.
