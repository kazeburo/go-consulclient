# go-consulclient

tiny consul (service) client

## Usage

```
cc := consulclient.New("http://consul.service.example.consul:8500", timeout)
len, err := cc.PassingNodeLen("example")
```




