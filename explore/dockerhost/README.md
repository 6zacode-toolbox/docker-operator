# How the data was extract

```bash 
docker info --format '{{json .}}'

docker ps -a --format "{{json .}}" | jq -s
```


# Translation 

https://www.json2yaml.com/
https://mholt.github.io/json-to-go/