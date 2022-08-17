import { COLORS } from '@/config';

import { MODALS } from './modal';

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

export const CREATE_MAP_MODAL_NAMES_BY_TYPE = {
  [MAP_TYPES.geo]: MODALS.createGeoMap,
  [MAP_TYPES.flowchart]: MODALS.createFlowchartMap,
  [MAP_TYPES.mermaid]: MODALS.createMermaidMap,
  [MAP_TYPES.treeOfDependencies]: MODALS.createTreeOfDependenciesMap,
};

export const MERMAID_THEMES = {
  default: 'default',
  base: 'base',
  forest: 'forest',
  dark: 'dark',
  neutral: 'neutral',
  canopsis: 'canopsis',
};

export const MERMAID_THEME_PROPERTIES_BY_NAME = {
  [MERMAID_THEMES.canopsis]: {
    theme: MERMAID_THEMES.base,
    themeVariables: COLORS.mermaid,
  },
};

export const TREE_OF_DEPENDENCIES_GRAPH_OPTIONS = {
  fitPadding: 20,
  wheelSensitivity: 0.3,
  minZoom: 0.02,
  maxZoom: 1.5,
};

export const TREE_OF_DEPENDENCIES_GRAPH_LAYOUT_OPTIONS = {
  name: 'cise',
  animate: 'end',
  fit: true,
  padding: 50,
  nodeSeparation: 110,
  nodeRepulsion: 120,
  idealInterClusterEdgeLengthCoefficient: 2,
};
