# SOCme

## Introduction

### Un Security Operations Center (SOC)

Un **Security Operations Center (SOC)** est une équipe centralisée de cybersécurité qui surveille, détecte, analyse et répond aux menaces en temps réel. Il est généralement composé d’analystes, d’ingénieurs et parfois d’experts en réponse aux incidents. Un SOC utilise diverses technologies comme les **SIEM (Security Information and Event Management)**, les **EDR (Endpoint Detection and Response)** et d'autres outils de surveillance réseau et système.

Les principales missions d’un SOC sont :

- Surveillance et détection des incidents de sécurité 24/7
- Analyse des menaces
- Gestion des logs et corrélation des événements
- Threat intelligence et huntings pour identifier des menaces avancées
- Amélioration continue via des simulations et analyses post-mortem

Un **SOCme** est une version open-source d'un SOC qui vise à offrir des solutions accessibles pour la détection et la réponse aux menaces. Il permet aux entreprises et aux chercheurs de mettre en place une infrastructure de sécurité sans coûts de licence élevés.

### Le projet

L'idée de ce projet est de créer un SOC open-source qui permet de surveiller les réseaux et les logs des clients de manière décentralisée, mais aussi de fournir des outils pour les entreprises ou particuliers qui souhaitent mettre en place leur propre SOC. Le tout en utilisant des technologies open-source, des bonnes pratiques de sécurité, des configurations reproductibles et une documentation claire.

## Fonctionnement

Récupérer l'intégralité des logs des clients sur un seul et même serveur nous appartenant pour les analyser et les traiter serait trop coûteux en terme de puissance et de stockage. C'est pourquoi SOCme est décentralisé. Chaque client possède un nœud (Artemis) qui va surveiller son réseau et ses logs. Ces nœuds vont ensuite envoyer **les alertes** (***et uniquement les alertes***) à un serveur central (Zeus) qui va les analyser.

### Les composants

SOCme est composé de plusieurs type de machine :

- Zeus:
  - C'est le cerveau de SOCme
  - Il permet de voir les alertes sans avoir à se connecter à chaque artemis, gérer les utilisateurs, les clients, les accès aux noeuds, les notifications, etc...
  - Il permet également de gérer la flotte de machines Artemis
- Artemis:
  - Les nœuds
  - Ils sont disposé chez chaque clients
  - Ils sont les machines qui vont surveiller le réseau, les logs, les fichiers, etc...
- Hermes:
  - Un outil pour visualiser rapidement les alertes des différents clients sur un écran.
  - Fait pour un Raspberry Pi

### Les applications

- Zeus:
  Zeus est un serveur virtualisé sous NixOS (qui permet la reproductibilité de l'infrastructure) qui va gérer l'ensemble des alertes des Artemis. Il est composé de plusieurs applications :
  - SOCme:
    Application web pour visualiser les alertes, gérer les utilisateurs, les clients, etc...
    C'est cette application qui va communiquer avec les autres applications (Wazuh, TheHive, FETCHme, ...)
  - Headscale (alternative open-source & self hosted à Tailscale server):
    Serveur VPN pour que les Artemis puissent communiquer avec Zeus
    Avec des ACL, nous pourrons contrôler les accès des Artemis à Zeus. Les artemis ne pourront pas accéder aux autres artemis.
  - +Hardening, CI/CD
- Artemis:
  - Wazuh:
    Wazuh est un XDR (Extended Detection and Response) open-source qui permet de surveiller les logs, les fichiers, les processus, etc... Il est composé de *Wazuh Manager*, *Wazuh Dashboard*, *Wazuh Indexer* et de plusieurs *Wazuh Agents*.
    Il va analyser les logs des clients et générer des alertes en cas d'activité suspecte.
    Fonctionne avec des règles de détection, donc mettre à jour régulièrement les règles via la CI/CD.
  - Suricata:
    Suricata est un IDS open-source qui permet de surveiller le réseau et générer des alertes en cas d'activité suspecte.
    Sur un switch, on configure un port mirroring pour envoyer une copie de tous les paquets à Suricata (ou via un tap) via une interface dédiée.
    Fonctionne avec des règles de détection de signatures, donc mettre à jour régulièrement les règles via la CI/CD.
    Suricata est configuré pour envoyer les alertes à Wazuh.
  - Tailscale (client):
    Tailscale est un VPN open-source qui permet de connecter les Artemis à Zeus de manière sécurisée sans avoir à ouvrir des ports sur le routeur du client. (Pratique pour une mise en place rapide et facile)
  - TheHive:
    TheHive est un gestionnaire de cas d'incidents open-source qui permet de gérer les alertes, les incidents, les preuves, ...
    Il va servir à gérer la communication entre les différents acteurs de la sécurité.
    Les alertes récolté par l'SOCme qui nous semblent suspectes seront envoyées à TheHive pour être traitées par les clients.
    Permet la création de sortes de tickets pour chaque alerte.
    A voir dans quelle mesure TheHive peut également s'occuper de la partie notification.
  - FETCHme (custom):
    Application pour afficher les signaux vitaux des Artemis, basé sur fastfetch.
  - +Hardening, CI/CD
- Hermes:
  - Will see

#### SOCme

L'application principale.
En sveltekit avec un backend golang.

Vas utiliser l'API de Wazuh pour récupérer les alertes.
Vas utiliser des commandes système pour gérer les Artemis (headscale, ssh, ...): Vas falloir bien géré les permissions pour que le backend est le droit d'utiliser tout ça

**Headscale:**

Par exemple, pour ajouter un client:

```bash
user="acme-inc"
headscale users create $user
headscale preauthkeys create --user $user
```

Un user headscale doit également être créé pour chaque analyste membre de l'équipe afin de pouvoir se connecter au Zeus depuis son ordinateur personnel.
Pour se connecter à un Artemis, l'analyste doit passer par le Zeus (session SSH dans une session SSH).

Le retour de la commande est l'authkey qui doit être écrit dans le fichier de configuration de l'Artemis.

#### FETCHme

L'application pour afficher les signaux vitaux des Artemis.

- [x] Flake Version
- [x] Wazuh status for dashboard, indexer & manager
- [ ] Suricata status and number of logs in 24h
- [ ] TheHive status
- [ ] tailscale name
- [ ] tailscale ip
- [x] Uptime
- [x] hostname
- [x] storage
- [x] CPU & Memory usage
- [x] ip

#### OPSme

C'est un outil présent sur Zeus pour gérer la flotte des artémis, nous allons créer un outil en golang qui va:

- Récupérer les informations des Artemis (ip, hostname, password) depuis pocketbase
- Se connecter en SSH aux Artemis pour exécuter des commandes (status, mise à jour, redémarrage, etc...)
- Une version TUI de OpenOps pourra être créée pour la création de TTY sur zeus

#### SOCme Website

Un site pour présenter le projet, les différentes applications, les technologies utilisées, les documentations, etc...
Peux aussi contenir un blog pour les mises à jour, les nouvelles fonctionnalités, des statistiques anonymiser
Un calculateur de stockage/puissance nécessaire pour les clients en fonction de leur nombre de postes et de serveurs

## En pratique

### Build

Lors de l'arrivé chez un client:

- Création du client sur Zeus pour obtenir un token (headscale)
- Installation de la machine sur le réseau
- Modification de la configuration NixOS pour match les interfaces réseaux, les mot de passes, l'hostname
- Déploiement de l'agent Wazuh sur le SI (via GPO par exemple)
- Déploiement de Sysmon sur les postes clients pour augmenter la visibilité
- Port Mirroring sur le/les switch pour envoyer les logs à l'Artemis (interface suricata)

C'est une procédure simple et rapide, qui ne devrait pas prendre très longtemps. A voir au niveau des SI cibles si on peut automatiser tout ça ou pas. Selon l'architecture réseau du client, il faudra peut-être faire des ajustements.

Ensuite, nous établissons avec le client une liste de règles à mettre en place pour les alertes. Nous allons également configurer les notifications pour que le client soit informé en cas d'alerte.

Par exemple:

- Si bruteforce, envoyer un mail à l'administrateur
- Si compromission de compte, envoyer un mail à l'administrateur et téléphoner
- Si telle alerte, ajouter sur TheHive
- ...

### Run

Durant le RUN, nous allons surveiller les alertes, les incidents, les logs et les performances.
Nous allons également mettre à jour les règles de détection, les logiciels et les configurations.
Lorsqu'une alerte est détectée, nous allons la traiter en fonction de sa criticité et des procédures établies avec le client. Si elle ne correspond à aucune règle, nous allons l'analyser et la documenter pour l'ajouter à une règle et notifier le client si ça nous semble nécessaire.
