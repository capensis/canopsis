import { SNMP_TEMPLATE_STATE_STATES } from '@/constants';

export default {
  oid: 'OID',
  module: 'Sélectionnez un module MIB',
  output: 'Message',
  resource: 'Ressource',
  component: 'Composant',
  connectorName: 'Nom du connecteur',
  state: 'Criticité',
  toCustom: 'Personnaliser',
  writeTemplate: 'Écrire un modèle',
  defineVar: 'Définir la variable SNMP correspondante',
  moduleMibObjects: 'Champ d\'association des variables SNMP',
  regex: 'Expression régulière',
  formatter: 'Format (groupe de capture avec \\x)',
  uploadMib: 'Envoyer un fichier MIB',
  addSnmpRule: 'Ajouter une règle SNMP',
  uploadedMibPopup:
    'Le fichier a été téléchargé.\nAvis: {notification}\nObjets: {object}'
    + '|Les fichiers ont été téléchargés.\nAvis: {notification}\nObjets: {object}',
  states: {
    [SNMP_TEMPLATE_STATE_STATES.info]: 'Info',
    [SNMP_TEMPLATE_STATE_STATES.minor]: 'Mineur',
    [SNMP_TEMPLATE_STATE_STATES.major]: 'Majeur',
    [SNMP_TEMPLATE_STATE_STATES.critical]: 'Critique',
  },
};
