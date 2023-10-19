**Multithreading Challenge**

This app gets the postal code from the following APIs:
- https://brasilapi.com.br/api/cep/v1/
- https://viacep.com.br/ws/

It timesout if none of them returns a response in less than 1 second.

**Executing**

After executing the file below, the app will prompt the user to enter a Brazilian postal code (CEP). e.g. "93520-575".
```
go run main.go
```

**Testing**

The command below calls the app 10 times with a fixed input (e.g. "93520-575").
```
go test -v
```