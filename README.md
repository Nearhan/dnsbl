# DNSBL

A little go application that checks ips against `zen.spamhaus.org` to see if they are listed for spam. The api is implemented through graphql.

Read more about how [dnsbls](https://en.wikipedia.org/wiki/Domain_Name_System-based_Blackhole_List) work.

Major Technolgies used:

- graphql via gqlgen
- sqlite3 for protable database

I tried to use a few libraries a possible, leveraging the std lib as much as possible and hand rolling most aspects of the application. The exceptions are listed below.

Libraries used:

	- github.com/stretchr/testify v1.4.0 // adds some great testing assertions for go
	- github.com/satori/go.uuid v1.2.0   // to generate UUID's
	- golang.org/x/crypto v0.0.0....     // generate hash and salt for passwords


## Deving

This is a graphql project that heavily relies on [gqlen](https://gqlgen.com/) to generate a ton of the graphql portion. 


If you want to modify the the schema it is located at  `schema.graphqls`, then run `go generate`, this will generate resolvers in the home directory for you to implement. 


The config for generation is located at `gqlgen.yml`

```
├── Dockerfile                    <-- basic docker file
├── auth.go                       <-- auth middle wear for graphql 
├── buildachart                   <-- helm charts
|
├── cmd
│   └── server
│       └── main.go                 <--- Entrypoint for app
├── generate.go
├── go.mod
├── go.sum
├── gqlgen.yml           
├── integration_test.go            <--- integration tests
├── models.go                      <--- where sqlite models live
├── models_test.go 
├── resolver.go                    <--- where majority of the implmentation gql live
├── resolver_test.go
├── schema.graphqls                <--- graphql schema
├── schema.resolvers.go            <--- generated resolvers
├── server_gen.go                  <--- server generated
└── testdata                       <--- testdata for integration test
```


Current Schema.

```
type IPDetail {
  id: ID!
  createdAt: Time!
  updatedAt: Time!
  responseCode: String!
  ipAddress: String!
}

type Query {
  getIPDetails(ipAddress: String!): IPDetail!
}

type Mutation {
  enqueue(ipAddresses: [String!]): String!
}

scalar Time
```


## Running

You can run this application in a couple of ways.
Just locally.

```
PORT=9999 go run cmd/server/main.go

In your browser go to localhost:9999
```

In Docker

```
docker build . -t nearhan/dnsbl
docker run -p 8222:9000 -e PORT=9000 nearhan/dnsbl

In your browser go to localhost: 8222
```

With helm and minikube.

```
helm install dnsbl buildachart/ --values buildachart/values.yaml
minikube service --url dnsbl

follow the url obtained from the output
```

### Basic Auth
Once you have the application running you go can to the url and you'll enter the grapql playground webpage. This will give you a 401 unless you fill in basic authorization information. In the playground home page near the bottom where it says http headers add the following.


```
{
  "Authorization": "Basic c2VjdXJld29ya3M6cGFzc3dvcmQ="
}
```

### Example


```
query {
  getIPDetails(ipAddress:"1.1.1.1") {
    id
	ipAddress
    createdAt
    updatedAt
    responseCode
  }
}

mutation {
  enqueue(ipAddresses:["1.1.1.1"])
}

```

curl
 
```
getIpDetails

curl 'http://localhost:9999/query' -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' -H 'Connection: keep-alive' -H 'DNT: 1' -H 'Origin: http://localhost:9999' -H 'Authorization: Basic c2VjdXJld29ya3M6cGFzc3dvcmQ=' --data-binary '{"query":"query {\n  getIPDetails(ipAddress:\"1.1.1.1\") {\n    id\n  }\n}"}' --compressed



enqueue

curl 'http://localhost:9999/query' -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' -H 'Connection: keep-alive' -H 'DNT: 1' -H 'Origin: http://localhost:9999' -H 'Authorization: Basic c2VjdXJld29ya3M6cGFzc3dvcmQ=' --data-binary '{"query":"mutation {\n  enqueue(ipAddresses:[\"1.1.1.1\"])\n}"}' --compressed
```



###

## Tests

There are unit tests and integration tests.
To run the integration tests you'll need an internet connection.

At repo root.

For unit tests.

```
go test -v ./...
```

For integration tests

```
go test -v ./... --tags=integration
```



## Things to improve


- Enqueue 
	- For each ip that you pass in, it spins off a go routine to do the look up against zen.spamhaus.org. Would be smarter to implement a resource bounded worker pool instead of N go routines.
	
	

