1) Criar o Docker com o comando: docker build -t stress-test .

2) Rodar ele com o comando: docker run stress-test --url=http://google.com --requests=1000 --concurrency=10
