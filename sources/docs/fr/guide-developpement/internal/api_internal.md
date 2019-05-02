# API Internal

### Modification de l'edition et de la stack de Canopsis

Des valeurs par défaut sont mises en place pour l'édition et la stack de Canopsis. L'édition vaut soit `"cat"`, soit `"core"`. La stack est soit `"go"`, soit `"python"`.

Pour modifier ces valeurs, il faut passer l'API.

**URL** : `/api/internal/properties`

**Méthode** : `POST`

**Authentification requise** : Oui

**Permissions requise** : Aucune

**Exemple de corps de requête** :
```json
{
    "edition":"cat",
    "stack":"go
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le Json ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
    "edition":"cat",
    "stack":"go
}' 'http://<Canopsis_URL>/api/internal/properties'
```

#### Réponse en cas de réussite

**Condition** : l'édition et/ou la stack ont été modifiées.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```json
{}
```

#### Réponse en cas d'erreur

**Condition** : Si le corps de la requête n'est pas valide.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "Invalid JSON"
}
```

---

**Condition** : Si `edition` ou `stack` sont définis dans le JSON mais ne correspondant pas à des valeurs valies.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
    "name": "",
    "description": "Invalid value(s)."
}
```

---

**Condition** : En cas d'erreur avec la base de données.

**Code** : `400 BAD REQUEST`

**Exemple du corps de la réponse** :

```json
{
  "description": "Error while updating edition/stack values"
}
```


