import { MAP_TYPES } from '@/constants';

export default {
  defineEntity: 'Définir l\'entité',
  addLink: 'Ajouter un lien',
  addPoint: 'Ajouter un point',
  editPoint: 'Modifier le point',
  removePoint: 'Supprimer le point',
  latitude: 'Latitude',
  longitude: 'Longitude',
  toggleAddingPointMode: 'Basculer le mode d\'ajout de point',
  usingMap: 'La carte est liée',
  showAll: 'Afficher tout ({count})',
  types: {
    [MAP_TYPES.geo]: 'Géographique',
    [MAP_TYPES.flowchart]: 'Flowchart',
    [MAP_TYPES.mermaid]: 'Mermaid',
    [MAP_TYPES.treeOfDependencies]: 'Arbre des dépendances',
  },
  layers: {
    openStreetMap: 'Open street map',
    points: 'Points',
  },
};
