import { isNumber } from 'lodash';

export const getControlPosition = (evt) => {
  const offsetParent = evt.target.offsetParent || document.body;
  const offsetParentRect = evt.offsetParent === document.body
    ? { left: 0, top: 0 } : offsetParent.getBoundingClientRect();

  const x = evt.clientX + offsetParent.scrollLeft - offsetParentRect.left;
  const y = evt.clientY + offsetParent.scrollTop - offsetParentRect.top;

  return { x, y };
};

export const createCoreData = (lastX, lastY, x, y) => (
  !isNumber(lastX)
    ? {
      deltaX: 0,
      deltaY: 0,
      lastX: x,
      lastY: y,
      x,
      y,
    }
    : {
      deltaX: x - lastX,
      deltaY: y - lastY,
      lastX,
      lastY,
      x,
      y,
    }
);
