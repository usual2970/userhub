# UserHub

UserHub is a general-purpose user center or user service that allows your app to have registration, login, and other capabilities, allowing you to focus on your business logic.

## Why UserHub

Every time we want to do a hobby project or a startup develops a new project, we may need to consider user centers, message centers, and other common and basic services while considering the business, which leads to repeated unnecessary work. However, these basic services are actually reusable and can be developed once and used multiple times.

UserHub is such a project that allows you to focus on your business logic without having to consider the user center. UserHub provides user login, registration, user management, and other necessary functions for the user center. In addition to the built-in methods for login and registration, UserHub also supports users to extend their own functions.

## Installation

Use the package manager [pip](https://pip.pypa.io/en/stable/) to install foobar.

```bash
pip install foobar
```

## Usage

```python
import foobar

# returns 'words'
foobar.pluralize('word')

# returns 'geese'
foobar.pluralize('goose')

# returns 'phenomenon'
foobar.singularize('phenomena')
```

## Roadmap

- [ ] UserHub's overall design and architecture includes the functions to be implemented.
- [ ] The construction of the framework (the selection of configuration, log, orm and other basic components)
- [ ] The implementation of basic functions (registration, login, binding, logout, query, etc.)
- [ ] Implementation of external interfaces (http, grpc, or even go mod)
- [ ] Optimization, such as sharding, docker, etc.

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)