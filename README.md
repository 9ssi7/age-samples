# Apache Age Golang Driver Samples

This repository contains samples for using the Apache Age Golang driver.

## Before All, Install Apache Age with Docker

```bash
docker run \
--name age  \
-p 5455:5432 \
-e POSTGRES_USER=postgres \
-e POSTGRES_PASSWORD=postgres \
-e POSTGRES_DB=postgres \
-d \
apache/age
```

## Running the samples

To run the samples, you need to have Go installed on your machine. You can download and install Go from [here](https://golang.org/dl/).

After installing Go, you can run the samples by executing the following commands:

```bash
cd basic
go run main.go
```

## License

This project is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.
