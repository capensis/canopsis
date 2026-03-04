# Nettoyage et rétention des bases de données

### Nettoyage

!!! warning
    Le nettoyage des bases de données `MongoDB` et `TimescaleDB` ne doit jamais être réalisé directement en base.  
    Il faut pour cela utiliser la fonctionnalité [`Paramètres de stockage`](../../../guide-utilisation/menu-administration/parametres-de-stockage.md) qui centralise toutes les politiques de rétention des données de Canopsis.