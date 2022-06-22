import { CONNECTOR_SIDES } from '@/plugins/flowchart/constants';

export const calculateConnectorPointBySide = (rect, side, { x: percentX, y: percentY }) => {
  const { x, y, width, height } = rect;
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
