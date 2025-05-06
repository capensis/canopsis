# Supervision de Canopsis

La supervision de la plateforme Canopsis est un élément clé pour assurer la disponibilité, la performance et la stabilité de votre solution d'hypervision.

Cette page présente les outils et mécanismes mis à disposition pour superviser l'état de santé de Canopsis, détecter les dysfonctionnements, et intégrer les métriques dans vos solutions de monitoring existantes.

## Bilan de santé de Canopsis

Canopsis propose un module dédié permettant de consulter en temps réel l'état des différents composants du système : les moteurs, les files RabbitMQ, les bases de données, etc.

[Consulter la documentation du Bilan de santé](../../../guide-utilisation/menu-administration/bilan-de-sante/)

## Exporter Prometheus

Un **exporter Prometheus** est disponible pour exposer des métriques internes au format compatible avec Prometheus.

[Voir la documentation détaillée de l’exporter Prometheus](./exporter-prometheus.md)

## À venir

Nous proposerons prochainement des **templates de supervision** pour les solutions Centreon et Zabbix.

