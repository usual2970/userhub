# UserHub

UserHub是一个通用的用户中心或用户服务，使用它可以快速的使你的应用拥有注册、登录等能力，从而让你集中精力处理业务逻辑。

## Why UserHub

每次想做个业余项目或者创业者开发一个新的项目，我们可能在考虑业务的同时都必不可少的需要考虑用户中心、消息中心等等这些通用且基础的服务，从而导致重复做一些不必要的劳动，然而这些基础服务事实上是可以通用、复用的，可以做到一次开发多次使用。

UserHub就是这样一个项目，使你不必考虑用户中心，集中精力处理自己的业务逻辑。UserHub提供了用户登录、注册、用户管理等等作为用户中心必须的功能。登录、注册的方式除了内置的几种方式，还支持使用者自行扩展。

## 安装

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

- [ ] UserHub的整体设计，架构设计包括要实现的功能等
- [ ] 框架的搭建（配置、日志、orm等基础组件的选型）
- [ ] 基础功能的实现（注册、登录、绑定、登出、查询等）
- [ ] 对外接口实现（http、grpc、甚至go mod)
- [ ] 优化，比如分库分表、docker等

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)