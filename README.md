# Router GO Library

# Install
```yaml
import:
- package: github.com/best-expendables/router
  version: x.x.x
```

# Middlewares
- `middleware.Authentication`
- `middleware.Authorized`
- `middleware.OnlyRoles`
- `middleware.ACL`
- `middleware.AccessLog`
- `middleware.Prometheus`
- `middleware.Recoverer`
- `middleware.OpenTrace`
- `middleware.RequestID` [DEPRECATED]


# Initialization
Function `router.New()` returns router with configured list of middlewares.

```go
func main() {
	loggerFactory := logger.NewLoggerFactory(logger.InfoLevel),
	
	config := router.Configuration{
		LoggerFactory: loggerFactory,
		PanicHandler: panicHandler,
	}
	
	mux, err := router.New(config)
	...
}
```


### Configuration parameters:

**Namespace** [REQUIRED]   
Service name (namespace) from k8n.

**LoggerFactory** [REQUIRED]  
This factory using for "context-logger".

**PanicHandler**  
Panic will be handled by this handler.

**Tracer**   
Tracer for OpenTrace. Should be nil. By default middleware will use GlobalTracer and this is enough for most of our use-cases. 

**DisableAccessLog**  
Disable access log


### Custom router
For exiting code you can initialize middlewares by yourself.   
You have to follow this order of middlewares on  for router initialization:

* Recoverer
* RequestID
* Authentication
* ContextLogger
* AccessLog


### Panic recovery
For the panic case, recoverer does not responsible for error response, it only sends HTTP response code 500.
For response with error messages you have to define panic handler and pass it as argument.
[Example](https://github.com/best-expendables/repos/router/browse/_examples/01_recoverer/main.go)   


### AccessLog
Improved ReqResLogger which write log into JSON format and ignore binary(non-text) body. 


### Authentication & Authorization
Authentication middleware receive user information only and set it into the context.  
ACL middleware using for access control.