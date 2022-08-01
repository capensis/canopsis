import { MODALS } from '@/constants/modal';

export const MAP_TYPES = {
  geo: 'geo',
  flowchart: 'flowchart',
  mermaid: 'mermaid',
  treeOfDependencies: 'treeOfDependencies',
};

export const MAP_CREATE_MODAL_NAMES_BY_TYPE = {
  [MAP_TYPES.geo]: MODALS.createGeoMap,
  [MAP_TYPES.flowchart]: MODALS.createFlowchartMap,
  [MAP_TYPES.mermaid]: MODALS.createMermaidMap,
  [MAP_TYPES.treeOfDependencies]: MODALS.createTreeOfDependenciesMap,
};
