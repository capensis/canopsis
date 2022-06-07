# Guide de migration vers Canopsis 4.6.0

Ce guide donne des instructions vous permettant de mettre à jour Canopsis 4.5 vers [la version 4.6.0](../4.6.0.md).

## Prérequis

L'ensemble de cette procédure doit être lu avant son exécution.

Ce document ne prend en compte que Canopsis Community et Canopsis Pro : tout développement personnalisé dont vous pourriez bénéficier ne fait pas partie du cadre de ce Guide de migration.

## Procédure de mise à jour

### Réalisation d'une sauvegarde

Des sauvegardes sont toujours recommandées, qu'elles soient régulières ou lors de modifications importantes.

La restructuration apportée dans les environnements docker-compose pour cette version de Canopsis nous amène à insister d'autant plus sur ce point. Il est donc fortement recommandé de réaliser une **sauvegarde complète** des VM hébergeant vos services Canopsis, avant cette mise à jour.

### Migration vers les nouveaux environnements docker-compose

1. Arrêter Canopsis :
```
docker-compose -f 00-data.docker-compose.yml -f 01-prov.docker-compose.yml -f 02-app.docker-compose.yml down
```
2. Mettre à jour les environnements
```
git checkout 4.6.0
```
3. Démarrer Canopsis :

=== "Canopsis Pro"

	```
	CPS_EDITION=pro docker-compose up -d
	```

=== "Canopsis Community"

	```
	CPS_EDITION=community docker-compose up -d
	```

Docker-compose doit démarrer les conteneurs par ordre de dépendance. Vous pouvez
surveiller cela en lançant la commande suivante dans un autre terminal pendant 
le démarrage:

=== "Canopsis Pro"

	```
	watch -n1 CPS_EDITION=pro docker-compose ps
	```

=== "Canopsis Community"

	```
	watch -n1 CPS_EDITION=community docker-compose ps
	```


