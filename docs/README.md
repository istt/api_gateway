# `api-gateway`

The API Gateway features is a lightweight replacement for JHipster API Gateway.

### API Endpoints:

- `POST api/authenticate`: Login with username, password. System will return a JWT bearer token to login.
- `GET api/account`: Retrieve the current logged in account information
- `POST api/account`: Save current account preferences


#### CLI for backend deployment


```bash
$ make api-gateway
$ 
$ rsync -avP build/package/api-gateway root@staging.soibe.xyz:/var/www/lab.soibe.xyz/apps/
$ 
$ systemctl restart lab.soibe.xyz.service
$ 
$ journalctl -u lab.soibe.xyz.service -f
$ 
```

#### Build Tool account-cli


```bash
$ # go build -o [outputFile] [srcFile]
$ go build  -o build/package/account-cli tools/account-cli.go
```

Upload to server to create new account for api-gateway


#### CLI for frontend deployment

```bash
$ cd Jhipster/jhipster
$ 
$ npm run build
$ 
$ rsync -avP target/classes/static/ root@staging.soibe.xyz:/var/www/lab.soibe.xyz/public/
```

#### Config for NGINX


```
 	# backend proxy
    location ~ (api|management|v2|services|schemas) {
        try_files $uri $uri/ @api;
    }


    # reverse proxy
    location @api {
        proxy_pass http://127.0.0.1:8080;
        include    nginxconfig.io/proxy.conf;
    }
```
