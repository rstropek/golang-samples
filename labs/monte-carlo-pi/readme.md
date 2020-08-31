# PI Monte Carlo

## Anforderung

Ihre Aufgabe ist die Näherungsweise Berechnung von PI mit Hilfe der [Monte Carlo](https://en.wikipedia.org/wiki/Monte_Carlo_method) Methode. Es geht bei diesem Beispiel nicht darum, PI möglichst genau zu berechnen. Ziel ist das Üben von Parallelverarbeitung mit Go. Wenn das Ergebnis in etwa PI entspricht (auf 2-3 Stellen genau) ist das ausreichend.

Verwenden Sie zur Berechnung von PI möglichst alle in Ihrem PC vorhandenen Prozessoren parallel.

## Berechnungsmethode

Eine gute Beschreibung der näherungsweisen Berechnung von PI findet man [hier](https://www.zum.de/Faecher/Inf/RP/Java/java_1.htm).

Erstellen Sie viele, zufällige Punkte in einem Quadrat mit Seitenlänge 1.

![Punkte](https://www.zum.de/Faecher/Inf/RP/Java/Bild/Monte_Carlo.jpg)

Zählen Sie die Punkte, die im Einheitskreis (Radius = 1) liegen.

![Einheitskreis](https://www.zum.de/Faecher/Inf/RP/Java/Bild/java_15.gif)

PI ergibt sich aus dem Verhältnis von Punkten innerhalb des Einheitskreises und der Anzahl an zufälligen Punkten.

![Formel](https://www.zum.de/Faecher/Inf/RP/Java/Bild/java_11.gif)

Hier Beispielcode, der zeigt, wie die näherungsweise Berechnung von PI in Go aussehen könnte:

```go
const ITERATIONS = 10000000

// Create a new random number generate. Attention! r is not thread safe.
// You need to call `rand.New` in each goroutine.
r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

inside := uint(0)
for i := 0; i < ITERATIONS; i++ {
    a := r.Float32()
    b := r.Float32()
    c := a*a + b*b
    if c <= float32(1) {
        inside++
    }
}

pi := float32(inside) / float32(ITERATIONS) * float32(4)
fmt.Printf("%.6f", pi)
```

## Tipps

* Die Anzahl an logischen CPUs erhalten Sie mit `runtime.NumCPU()`

## Levels

Je nach Vorwissen wird diese Aufgabe manchen schwerer und manchen leichter fallen. Daher hier Anregungen, wie man Schritt für Schritt das Beispiel lösen könnte. Jeder kann basierend auf der bestehenden Programmierpraxis die Anzahl an Levels abarbeiten, die für sie oder ihn passen.

### Level 1 - Rechenlogik

* Nehmen Sie den Beispielcode von oben und wandeln Sie ihn in ein lauffähiges Go-Programm um.
* Kontrollieren Sie, ob ein Ergebnis in der Nähe von PI geliefert wird.

### Level 2 - Parallelverarbeitung

* Ändern Sie den Code so, dass er parallel auf jeder CPU ausgeführt wird.

### Level 3 - Hilfsmethode

* Gliedern Sie den Algorithmus, der eine Funktion auf allen CPUs parallel ausführt, in eine eigene Funktion aus. Der Code zur Parallelisierung und der zur Berechnung von PI sollte logisch getrennt sein. Hier ein Beispiel, wie das aussehen könnte:

```go
func runInParallel(body func(...), goroutines int, ...) {
    ...
    for i := 0; i < goroutines; i++ {
        go body(...)
    }
    ...
}

func main2() {
    cpus := runtime.NumCPU()

    go runInParallel(func(...) {
        ...
    }, ...)

    ...
    fmt.Printf("%.6f", pi)
}
```

### Level 4 - Unit Testing

* Schreiben Sie mindestens einen sinnvollen Unit Test ([Kurzanleitung](https://gobyexample.com/testing)) für die Methode zur Parallelisierung von Funktionen (im obigen Beispiel `runInParallel`).
