### Migration vers les nouveaux environnements docker-compose

1. Arrêter Canopsis :
```
docker-compose -f 00-data.docker-compose.yml -f 01-prov.docker-compose.yml -f 02-app.docker-compose.yml down
```
2. Mettre à jour les environnements
```
git checkout 22.10
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


