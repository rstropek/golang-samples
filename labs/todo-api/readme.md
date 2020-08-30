# ToDo List API

## Anforderung

Ihre Aufgabe ist die Implementierung einer RESTful Web API **mit Go**.

Die technische Spezifikation der Web API finden Sie im *Open API* Format in [*api-spec.yaml*](api-spec.yaml). Hinweis: Sie können die Spezifikation in einem angenehm lesbaren Format ansehen, indem Sie [https://editor.swagger.io/](https://editor.swagger.io/) im Browser öffnen und den Inhalt der Datei [*api-spec.yaml*](api-spec.yaml) in den Eingabebereich auf der linken Seite kopieren.

**Lesen Sie die API Spezifikation genau.** Auch scheinbare Kleinigkeiten wie zum Beispiel geforderte *Response Codes* (z.B. *404* für *Not Found*, *201* für *Created*) sind wichtig und sollten beachtet werden.

## Levels

Je nach Vorwissen wird diese Aufgabe manchen schwerer und manchen leichter fallen. Daher hier Anregungen, wie man Schritt für Schritt das Beispiel lösen könnte. Jeder kann basierend auf der bestehenden Programmierpraxis die Anzahl an Levels abarbeiten, die für sie oder ihn passen.

### Level 1 - *Get Items*

* Befüllen Sie die Todo-Liste im Hauptspeicher im Programm mit ein paar Beispieldatensätzen
* Implementieren Sie den Basis-Webserver und die Operation *getItems* (`GET /api/toDoItems`)

### Level 2 - *Get Single Item*

* Erweitern Sie die API um eine Funktion zum Abfragen eines einzelnen Items (*getItem* (`GET /api/toDoItems/{id}`))

### Level 3 - *Add Items*

* Erweitern Sie die API um eine Funktion zum Hinzufügen von Items (*addItem* (`POST /api/toDoItems`))

### Level 4 - *Nested Package*

* Gliedern Sie die Geschäftslogik zum Verwalten der Todo-Liste in ein eigenes *nested package* aus und trennen Sie diesen Code damit vom Code der Web API.
* Schreiben Sie zumindest drei sinnvolle Unit Tests für das Package.

### Level 5 - Docker

* Schreiben Sie ein Dockerfile zum Erstellen eines Docker Image für die API (`docker build`)
* Testen Sie Ihre API in einem Docker Container (`docker run`)
* Veröffentlichen Sie das Docker Image am *Docker Hub* (`docker push`)

## Testfälle

### ToDo Items abfragen

Request:

```http
GET http://localhost:8080/api/toDoItems
```

### ToDo Item einfügen

Request:

```http
POST http://localhost:8080/api/toDoItems
Content-Type: text/plain

Einkaufen
```

### Einzelnes ToDo Item abfragen

Request (ID muss durch eigenen Wert ersetzt werden):

```http
GET http://localhost:8080/api/toDoItems/MXdzY2yO5rnZoZeg
```

### ToDo Item auf erledigt setzen

Request (ID muss durch eigenen Wert ersetzt werden):

```http
POST http://localhost:8080/api/toDoItems/MXdzY2yO5rnZoZeg/setDone
```
