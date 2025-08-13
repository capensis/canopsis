L’export des fixtures permet de sauvegarder, l’intégralité ou seulement certaines parties de la configuration de Canopsis, afin de pouvoir la réutiliser ultérieurement, par exemple lors d’un déploiement sur une nouvelle instance.

Pour réaliser un export, il suffit d’appeler le endpoint d’API suivant :

```
POST /api/v4/export-configuration
```

Le corps de la requête doit contenir un objet JSON avec la clé `export` listant les collections à extraire.

Exemple minimal de payload :

```json
{
  "export": [
    "metaalarm",
    "pbehavior"
  ]
}
```

Exemple de résultat au format YAML :

```yaml
meta_alarm_rules:
  "0":
    _id: Ping_1
    alarm_pattern:
    - - cond:
          type: eq
          value: Ping_1
        field: v.resource
    author: root
    auto_resolve: false
    config:
      threshold_count: 2
      time_interval:
        unit: m
        value: 15
    corporate_alarm_pattern: ""
    corporate_alarm_pattern_title: ""
    corporate_entity_pattern: ""
    corporate_entity_pattern_title: ""
    corporate_total_entity_pattern: ""
    corporate_total_entity_pattern_title: ""
    created: 1684142212
    entity_pattern: []
    name: Ping_1
    output_template: "{{ .Children.Alarm.Value.State.Message }}"
    total_entity_pattern: []
    type: complex
    updated: 1684142212
pbehavior:
  "0":
    _id: df28d4a0-7e79-4a89-b957-0385d673f087
    alarm_count: 863
    author: root
    color: ""
    corporate_entity_pattern: ""
    corporate_entity_pattern_title: ""
    created: 1683795610
    enabled: true
    entity_pattern:
    - - cond:
          type: begin_with
          value: Joignabilité
        field: name
    exceptions: []
    exdates: []
    last_alarm_date: 1754982491
    name: Joignabilité
    reason: 5fee78e8-2519-4cc5-a9d5-19cbde6519e6
    rrule: FREQ=DAILY
    rrule_cstart: 1753394400
    tstart: 1683756000
    tstop: 1683842399
    type_: 5ea9d2d8-0f16-4e19-bcca-64b1e96e00fa
    updated: 1683795610

```

## Liste complète des paramètres disponibles

La liste des paramètres est disponible ici : [schemas_swagger.yaml](https://github.com/capensis/canopsis/blob/develop/community/go-engines-community/lib/api/docs/schemas_swagger.yaml#L2835)

Chaque paramètres correspond à une collection MongoDB associée.

Il est possible de sélectionner uniquement les collections nécessaires ou bien de toutes les exporter en les listant dans le tableau `export`.