# F.A.Q. : Canopsis est-il concerné par la faille Log4j ? (CVE-2021-45046)

Canopsis en lui-même n'utilise pas de mécanisme de journalisation reposant sur Log4j, il n'est donc pas directement concerné par cette faille de sécurité.

Certains composants tiers peuvent cependant avoir été ajoutés dans votre environnement et être concernés par cette faille.

## Logstash

Certaines installations peuvent contenir une brique Logstash, traitant certains évènements en entrée ou en sortie de Canopsis : vérifiez pour cela si un service `logstash` est présent dans votre environnement.

Si vous êtes dans ce cas de figure, l'éditeur de Logstash recommande à ce jour de réaliser une mise à jour vers Logstash 7.16.1 ou 6.8.21.

Suivez pour cela les recommendations officielles de l'éditeur à cet endroit :  
<https://discuss.elastic.co/t/apache-log4j2-remote-code-execution-rce-vulnerability-cve-2021-44228-esa-2021-31/291476>

## Rundeck

Dans le cadre de la mise en place [de la remédiation](../../guide-administration/remediation/index.md), certaines installations peuvent comporter une brique Rundeck.

Si vous êtes dans ce cas de figure, l'éditeur de Rundeck recommande à ce jour de réaliser une mise à jour vers Rundeck 3.4.8 ou 3.3.16.

Suivez pour cela les recommendations officielles de l'éditeur à cet endroit :  
<https://docs.rundeck.com/docs/history/CVEs/#log4shell-cves>
