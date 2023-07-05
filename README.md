# short-url

### Start app
##### Specify which database will be used
##### Redis will be used by default
```
sudo docker compose build --build-args DB="postgres" && sudo docker compose up
```

### Endpoints
##### Create short url
```
POST
127.0.0.1:9000/create

data : {
    "url": "github.com/myhorizonn/short-url"
}

response: {
    "url": "jO0j5jyCyj"
}
```
##### Get original url
```
GET
127.0.0.1:9000/get_original

data: {
    "url": "jO0j5jyCyj"
}

response: {
     "url": "github.com/myhorizonn/short-url"
}
```