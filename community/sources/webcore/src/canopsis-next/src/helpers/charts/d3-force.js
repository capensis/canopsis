import { forceCollide, forceLink, forceManyBody, forceSimulation, forceX, forceY } from 'd3-force';
import { keyBy, cloneDeep } from 'lodash';

const DEFAULT_FORCE_BODY_STRENGTH = -2300;
const DEFAULT_FORCE_X_STRENGTH = 0.6;
const DEFAULT_FORCE_Y_STRENGTH = 3;

/**
 * Calculate dependencies count
 *
 * @param {Object} nodeByName
 * @param {Object} node
 * @param {number} length
 * @returns {number}
 */
const calculateNodeDepth = (node, nodeByName, length = 1) => {
  if (!node || !node.dependencies || node.dependencies.length < 1) {
    return length;
  }

  const childrenNodesDepth = node.dependencies.map(id => calculateNodeDepth(
    nodeByName[id],
    nodeByName,
    length + 1,
  ));

  return Math.max.apply(null, childrenNodesDepth);
};

/**
 * Mix depth in each of node
 *
 * @param {Array} nodes
 * @returns {Array}
 */
const prepareNodesForSimulation = (nodes) => {
  const nodeByName = keyBy(nodes, 'name');

  return cloneDeep(nodes).map(node => ({
    ...node,
    depth: calculateNodeDepth(node, nodeByName),
  }));
};

/**
 * Simulate network graph
 *
 * @param {Array} nodes
 * @param {Array} links
 * @param {number} width
 * @param {number} nodeRadius
 * @param {number} linkDistance
 * @param {number} iterations
 * @returns {{
 *  nodes: array
 *  links: array
 *  height: number
 * }}
 */
export const simulateNetworkGraph = ({
  nodes,
  links,
  nodeRadius,
  width,
  linkDistance,
  iterations = 300,
}) => {
  const diameter = nodeRadius * 2;

  const preparedNodes = prepareNodesForSimulation(nodes);
  const maxDepth = Math.max.apply(null, preparedNodes.map(({ depth }) => depth));
  const copiedLinks = cloneDeep(links);

  const forceLinks = forceLink(copiedLinks)
    .id(data => data.name)
    .links(copiedLinks);

  const forceXAxis = forceX(width / 2).strength(DEFAULT_FORCE_X_STRENGTH);

  const forceYAxis = forceY()
    .y(node => ((node.depth - 1) * linkDistance) + diameter)
    .strength(DEFAULT_FORCE_Y_STRENGTH);

  const forceCharge = forceManyBody().strength(DEFAULT_FORCE_BODY_STRENGTH);

  const simulation = forceSimulation(preparedNodes)
    .force('link', forceLinks)
    .force('x', forceXAxis)
    .force('charge', forceCharge)
    .force('y', forceYAxis)
    .force('collide', forceCollide(diameter));

  simulation.tick(iterations);

  return {
    nodes: preparedNodes,
    links: copiedLinks.map(link => ({
      x1: link.source.x,
      y1: link.source.y + nodeRadius,
      x2: link.target.x,
      y2: link.target.y - nodeRadius,
    })),
    height: ((maxDepth - 1) * linkDistance) + (diameter * 2),
  };
};
