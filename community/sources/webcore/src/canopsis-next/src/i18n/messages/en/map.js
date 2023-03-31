import { MAP_TYPES } from '@/constants';

export default {
  defineEntity: 'Define entity',
  addLink: 'Add link',
  addPoint: 'Add point',
  editPoint: 'Edit point',
  removePoint: 'Remove point',
  latitude: 'Latitude',
  longitude: 'Longitude',
  toggleAddingPointMode: 'Toggle adding point mode',
  usingMap: 'Map is linked',
  showAll: 'Show all ({count})',
  types: {
    [MAP_TYPES.geo]: 'Geo',
    [MAP_TYPES.flowchart]: 'Flowchart',
    [MAP_TYPES.mermaid]: 'Mermaid',
    [MAP_TYPES.treeOfDependencies]: 'Tree of dependencies',
  },
  layers: {
    openStreetMap: 'Open street map',
    points: 'Points',
  },
};
