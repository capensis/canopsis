# Paramètres de calcul d'état/sévérité

L'état ou la sévérité d'un Service est dépendant de règles de calcul d'état qui s'appliquent sur les dépendances de celui-ci.  

Lorsqu’aucune règle ne s'applique, l'état du service vaut le pire état de ses dépendances.

### Les bases

| Paramètre                          | Description                                                  |
| ---------------------------------- | ------------------------------------------------------------ |
| **Nom**                            | Nom de la règle à créer                                      |
| **Priorité**                       | Priorité d'application de la règle                           |
| **Activée**                        | La règle est-elle activée ou non ?                           |
| **Appliqué pour le type d'entité** | Cette règle de calcul s'applique t-elle aux `Composants` ou aux `Services` ? |
| **Méthode de calcul d'état**       | Choix de la méthode de calcul d'état                         |

### Définir les entités cibles

| Paramètre               | Description                                                  |
| ----------------------- | ------------------------------------------------------------ |
| **Modèles des entités** | La règle de calcul d'état s'applique aux entités ciblées par ce modèle |



### Ajouter des conditions

#### L'État est hérité des dépendances

|                             | Description                                                  |
| --------------------------- | ------------------------------------------------------------ |
| **Modèles des dépendances** | Quelle(s) dépendance(s) de l'entité ciblée sera(ont) responsable(s) de l'état final ?<br />Si plusieurs dépendances sont sélectionnées par le modèle alors le pire état de celles-ci sera utilisé. |

#### L'État est défini par un calcul (pourcentage ou nombre) appliqué sur les états des dépendances

Dans ce mode, il est possible de définir l'état d'un service à partir de conditions basées sur un pourcentage ou un nombre d'états des dépendances du service.

Nous pourrions par exemple exprimer le fait que le service sera en état :

* Critique si plus de 50% de ses dépendances sont en état critique ou

* Majeur si 3 de ses dépendances sont en état mineur ou

* Mineur si 20% des entités sont en état majeur ou

* OK si au moins 1 dépendance est en état OK

  

Pour cela, des conditions peuvent être définies pour chaque état final du
service, comme dans l'illustration ci-dessous :

![services-calcul-etat2](../img/services-calcul-etat2.png)




