import { CONNECTOR_SIDES } from '@/constants';

export const calculateConnectorPointBySide = (
  shape,
  side,
  { x: percentX = 0.5, y: percentY = 0.5 } = {},
) => {
  const { x, y } = shape;
  const width = shape.width ?? shape.diameter ?? shape.size;
  const height = shape.height ?? shape.diameter ?? shape.size;
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
