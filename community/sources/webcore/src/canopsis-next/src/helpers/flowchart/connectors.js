import { CONNECTOR_SIDES } from '@/constants';

/**
 * @typedef { 'top' | 'right' | 'bottom' | 'left' } ConnectorSide
 */

/**
 * Calculate connect point by connected side
 *
 * @param {Shape} shape
 * @param {ConnectorSide} side
 * @param {number} percentX
 * @param {number} percentY
 * @returns {Point}
 */
export const calculateConnectorPointBySide = (
  shape,
  side,
  { x: percentX = 0.5, y: percentY = 0.5 } = {},
) => {
  const { x, y } = shape;
  const width = shape.width ?? shape.diameter;
  const height = shape.height ?? shape.diameter;
  const isTop = side === CONNECTOR_SIDES.top;

  if (isTop || side === CONNECTOR_SIDES.bottom) {
    const resultX = x + width * percentX;

    if (isTop) {
      return { x: resultX, y };
    }

    return { x: resultX, y: y + height };
  }

  const resultY = y + height * percentY;

  if (side === CONNECTOR_SIDES.left) {
    return { x, y: resultY };
  }

  return { x: x + width, y: resultY };
};
