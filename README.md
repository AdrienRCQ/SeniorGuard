# SeniorGuard - La cybersécurité pour nos aînés

SeniorGuard est un service de cybersécurité orienté EDR pour garantir la cybersécurité de nos ainés.
L'objectif est de garantir un niveau de sécurité informatique sans impliquer l'utilisateur et sans le déranger.

Architecture cible : 

```
                    SeniorGuard
                        │
                Agent Inventory
                        │
       ┌────────────────┴───────────────┐
       │                                │
   Agent A                          Agent B
       │                                │
heartbeat                           heartbeat
       │                                │
       └──────────── API ───────────────┘
                        │
                     last_seen
```