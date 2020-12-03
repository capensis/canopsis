import { forceCollide, forceLink, forceManyBody, forceSimulation, forceX, forceY } from 'd3-force';
import { keyBy, cloneDeep } from 'lodash';

const DEFAULT_FORCE_BODY_STRENGTH = -2300;

/**
 * Calculate dependencies count
 *
 * @param {Object} nodeById
 * @param {Object} node
 * @param {Number} length
 * @returns {number}
 */
const calculateNodeDepth = (node, nodeById, length = 1) => {
  if (!node || !node.dependencies || node.dependencies.length < 1) {
    return length;
  }

  const childrenNodesDepth = node.dependencies.map(id => calculateNodeDepth(nodeById[id], nodeById, length + 1));

  return Math.max.apply(null, childrenNodesDepth);
};

const fixCoordinate = (value, step) => {
  const inv = 1.0 / step;

  return Math.round(value * inv) / inv;
};

const prepareNodesForSimulation = (nodes) => {
  const nodeById = keyBy(nodes, 'id');

  return cloneDeep(nodes).map(node => ({
    ...node,
    depth: calculateNodeDepth(node, nodeById),
  }));
};

/**
 * Simulate network graph with multi root
 *
 * @param {Array} nodes
 * @param {Array} links
 * @param {number} width
 * @param {number} nodeRadius
 * @param {number} linkDistance
 * @returns {Promise<unknown>}
 */
export const simulateNetworkGraph = ({
  nodes,
  links,
  nodeRadius,
  width = 1000,
  linkDistance = 100,
}) => {
  const diameter = nodeRadius * 2;

  const preparedNodes = prepareNodesForSimulation(nodes);
  const maxDepth = Math.max.apply(null, preparedNodes.map(({ depth }) => depth));
  const copiedLinks = cloneDeep(links);

  const forceLinks = forceLink(copiedLinks)
    .id(data => data.id)
    .links(copiedLinks);

  const forceXAxis = forceX(width / 2).strength(0.5);

  const forceYAxis = forceY()
    .y(node => ((node.depth - 1) * linkDistance) + diameter)
    .strength(3);

  const forceCharge = forceManyBody().distanceMin(100).distanceMax(150).strength(DEFAULT_FORCE_BODY_STRENGTH);

  const simulation = forceSimulation(preparedNodes)
    .force('link', forceLinks)
    .force('x', forceXAxis)
    .force('charge', forceCharge)
    .force('y', forceYAxis)
    .force('collide', forceCollide(diameter));

  return new Promise((resolve) => {
    simulation.on('end', () => {
      resolve({
        nodes: preparedNodes.map(({ x, ...node }) => ({
          ...node,
          x: fixCoordinate(x, 50),
        })),
        links: copiedLinks.map(link => ({
          x1: fixCoordinate(link.source.x, 50),
          y1: link.source.y + nodeRadius,
          x2: fixCoordinate(link.target.x, 50),
          y2: link.target.y - nodeRadius,
        })),
        height: ((maxDepth - 1) * linkDistance) + (diameter * 2),
      });

      simulation.stop();
    });
  });
};
