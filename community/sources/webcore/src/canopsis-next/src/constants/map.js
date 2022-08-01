import { MODALS } from '@/constants/modal';

export const MAP_TYPES = {
  flowchart: 'flowchart',
  mermaid: 'mermaid',
  geo: 'geo',
  treeOfDependencies: 'treeOfDependencies',
};

export const MAP_ICON_BY_TYPES = {
  [MAP_TYPES.geo]: 'place',
  [MAP_TYPES.flowchart]: 'category',
  [MAP_TYPES.mermaid]: 'code',
  [MAP_TYPES.treeOfDependencies]: 'scatter_plot',
};

export const MAP_CREATE_MODAL_NAMES_BY_TYPE = {
  [MAP_TYPES.geo]: MODALS.createGeoMap,
  [MAP_TYPES.flowchart]: MODALS.createFlowchartMap,
  [MAP_TYPES.mermaid]: MODALS.createMermaidMap,
  [MAP_TYPES.treeOfDependencies]: MODALS.createTreeOfDependenciesMap,
};
