# ToDo List API mit Go Kit

## Anforderung

[Auf GitHub](https://github.com/rstropek/golang-samples/tree/master/go-microservices/05-simple-web-api) finden Sie ein einfaches Beispiel für eine Web API zum Verwalten von Personen:

* [Go Sourcecode](https://github.com/rstropek/golang-samples/blob/master/go-microservices/05-simple-web-api/web.go)
* [Beispielanfragen](https://github.com/rstropek/golang-samples/blob/master/go-microservices/05-simple-web-api/demo.http)

Ihre Aufgabe ist es, diese Web API mit Go Kit umzusetzen. Als Grundlage können Sie dafür das besprochene [Go Kit](https://github.com/rstropek/golang-samples/tree/master/go-kit) Beispiel verwenden.

## Levels

Je nach Vorwissen wird diese Aufgabe manchen schwerer und manchen leichter fallen. Daher hier Anregungen, wie man Schritt für Schritt das Beispiel lösen könnte. Jeder kann basierend auf der bestehenden Programmierpraxis die Anzahl an Levels abarbeiten, die für sie oder ihn passen.

### Level 1 - *Service für Geschäftslogik*

* Extrahieren Sie die Geschäftslogik zum Verwalten der Personen in ein [Go Kit Service](https://gokit.io/faq/#services-mdash-what-is-a-go-kit-service).
* Wenn Sie sich den Code des Ausgangsservice ansehen, werden Sie sehen, dass er keine Mutexe zur Steuerung des gleichzeitigen Zugriffs enthält. Korrigieren Sie das im Rahmen der Serviceentwicklung.

### Level 2 - *Endpoints* und *Transports* für `GetPerson`

* Entwickeln Sie einen [Go Kit Endpoint](https://gokit.io/faq/#endpoints-mdash-what-are-go-kit-endpoints) und [Transport](https://gokit.io/faq/#transports-mdash-what-are-go-kit-transports) für die Operation *Get Person*.
* Schreiben Sie eine *Main* Konsolenanwendung zum Betrieb des Go Kit Microservice.

### Level 3 - *Endpoints* und *Transports* vervollständigen

* Fügen Sie Endpoints und Transports für die restlichen Operationen (*Get People*, *Create Person*, *Delete Person*) hinzu.

### Level 4 - *Logging Middleware*

* Fügen Sie eine [*Logging Middleware*](https://gokit.io/faq/#middlewares-mdash-what-are-middlewares-in-go-kit) hinzu.

### Level 5 - Docker

* Schreiben Sie ein Dockerfile zum Erstellen eines Docker Image für die API (`docker build`)
* Testen Sie Ihre API in einem Docker Container (`docker run`)
* Veröffentlichen Sie das Docker Image am *Docker Hub* (`docker push`)
