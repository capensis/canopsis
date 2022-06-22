import { range } from 'lodash';

import uid from '@/helpers/uid';

/**
 * @typedef {Object} ConnectorPoint
 * @property {number} x
 * @property {number} y
 * @property {string} _id
 */

/**
 * Generate connector
 *
 * @param {number} x
 * @param {number} y
 * @returns {ConnectorPoint}
 */
export const generateConnectorPoint = (x, y) => ({
  x,
  y,
  _id: uid(),
});

/**
 * Generate connector points by rectangle
 *
 * @param {Rect} rect
 * @param {number} count
 * @returns {ConnectorPoint[]}
 */
export const generateConnectorsByRect = (rect, count) => {
  const { width, height, rx, ry, x, y } = rect;
  const sideConnectorsCount = Math.round(count / 4);
  const sideSpacesCount = sideConnectorsCount + 1;
  const startX = x + rx;
  const startY = x + ry;

  const widthConnectorOffset = (width - rx * 2) / (sideSpacesCount);
  const heightConnectorOffset = (height - ry * 2) / (sideSpacesCount);

  const sideConnectors = range(1, sideConnectorsCount + 1);

  const topConnectors = sideConnectors.map(index => generateConnectorPoint(
    startX + widthConnectorOffset * index,
    y,
  ));

  const rightConnectors = sideConnectors.map(index => generateConnectorPoint(
    x + width,
    startY + heightConnectorOffset * index,
  ));

  const bottomConnectors = sideConnectors.map(index => generateConnectorPoint(
    startX + widthConnectorOffset * index,
    y + height,
  ));

  const leftConnectors = sideConnectors.map(index => generateConnectorPoint(
    x,
    startY + heightConnectorOffset * index,
  ));

  const connectors = [
    ...topConnectors,
    ...rightConnectors,
    ...bottomConnectors,
    ...leftConnectors,
  ];

  if (!rx && !ry) {
    /* Angles connectors */
    connectors.push(
      generateConnectorPoint(x, y),
      generateConnectorPoint(x + width, y),
      generateConnectorPoint(x + width, y + height),
      generateConnectorPoint(x, y + height),
    );
  }

  return connectors;
};
