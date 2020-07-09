# minio-lifecycle
minio-lifecycle

Lifecycle your minio bucket for objects older than X days.

# Build

Install go.

`make build`

# Usage

`./minio-lifecycle -server http://<minio server ip>:<port> -key "<access key>" -secret '<access secret>' -bucket '<target bucket>' -days <max age of object>`

Example (deleting objects older than a week):

`./minio-lifecycle -server http://10.0.0.10:9768 -key "xxxxxxxxxxxxxxxx" -secret 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' -bucket 'dump' -days 7`