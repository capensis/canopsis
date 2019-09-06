# Post-traitement du moteur Axe

Les plugins de post-processing permettent d'appliquer des traitements aux alarmes après leur modification par le moteur [axe](../../guide-administration/moteurs/moteur-axe.md).

Les plugins de post-processing suivants sont disponibles dans Canopsis :
  - [Webhooks](../../guide-administration/moteurs/moteur-axe-webhooks.md)

Un plugin de post-processing est un module go exportant une variable
`AxePostProcessor` qui implémente l'interface `AxePostProcessor` (définie
dans
`git.canopsis.net/canopsis/go-revolution/cmd/engine-axe/plugins.go`).

Les plugins sont compilés avec `go build -buildmode=plugin`, qui crée un
fichier `.so`. Ils doivent être recompilés à chaque version de Canopsis.

Le code ci-dessous est un exemple minimal de plugin :

```go
package main

import (
    "log"

    "git.canopsis.net/canopsis/go-revolution/lib/canopsis/alarm"
    "git.canopsis.net/canopsis/go-revolution/lib/canopsis/types"
)

type dummyPostProcessor struct{}

// Init est appelé une fois au démarrage du moteur axe. Le moteur plante si
// cette méthode renvoie une erreur.
func (d dummyPostProcessor) Init() error {
    log.Println("dummy.Init")
    return nil
}

// Beat est appelé à chaque beat du moteur axe. Si cette méthode renvoie une
// erreur, elle est loggée, et l'exécution du moteur continue.
func (d dummyPostProcessor) Beat() error {
    log.Println("dummy.Beat")
    return nil
}

// ProcessEvent est appelé à chaque fois que le moteur axe a terminé le
// traitement d'un événement. Si cette méthode renvoie une erreur, elle est
// loggée, et l'exécution du moteur continue.
func (d dummyPostProcessor) ProcessEvent(alarmService alarm.Service, event types.Event, alarmChange types.AlarmChange) error {
    log.Printf("dummy.ProcessEvent: %+v\n", event)
    return nil
}

// ProcessAlarms est appelé quand plusieurs alarmes sont modifiées par lot. Si
// cette méthode renvoie une erreur, elle est loggée, et l'exécution du moteur
// continue.
func (d dummyPostProcessor) ProcessAlarms(alarmService alarm.Service, alarms []types.Alarm, change types.AlarmChangeType) error {
    log.Println("dummy.ProcessAlarms")
    return nil
}

// Crée un objet AxePostProcessor de type dummyPostProcessor.
var AxePostProcessor dummyPostProcessor
```
