import { getMaxZIndex } from '@/helpers/vuetify';

/**
 * Get the dimensions for a tooltip relative to a target element.
 *
 * @param {Object} options - The options object.
 * @param {HTMLElement} options.targetElement - The target element for the tooltip.
 * @param {HTMLElement} options.contentElement - The content element of the tooltip.
 * @param {boolean} [options.top=false] - Position the tooltip on top of the target element.
 * @param {boolean} [options.bottom=false] - Position the tooltip at the bottom of the target element.
 * @param {boolean} [options.left=false] - Position the tooltip to the left of the target element.
 * @param {boolean} [options.right=false] - Position the tooltip to the right of the target element.
 * @param {number} [options.nudgeTop=10] - Nudge value for top positioning.
 * @param {number} [options.nudgeLeft=10] - Nudge value for left positioning.
 * @param {number} [options.nudgeBottom=10] - Nudge value for bottom positioning.
 * @param {number} [options.nudgeRight=10] - Nudge value for right positioning.
 * @param {number} [options.contentMargin=12] - Margin around the content of the tooltip.
 * @returns {Object} The dimensions object for positioning the tooltip.
 */
export const getTooltipDimensions = ({
  targetElement,
  contentElement,
  top = false,
  bottom = false,
  left = false,
  right = false,
  nudgeTop = 10,
  nudgeLeft = 10,
  nudgeBottom = 10,
  nudgeRight = 10,
  contentMargin = 12,
}) => {
  const dimensions = {
    zIndex: getMaxZIndex(contentElement),
  };

  if (!targetElement || !contentElement) {
    return dimensions;
  }

  const { x, y, width, height } = targetElement.getBoundingClientRect();
  const { width: contentWidth, height: contentHeight } = contentElement.getBoundingClientRect();
  const targetLeft = window.scrollX + x;
  const targetTop = window.scrollY + y;
  const windowWidth = window.innerWidth;
  const windowHeight = window.innerHeight;

  const nudge = {
    [top]: nudgeTop,
    [right]: nudgeLeft,
    [bottom]: nudgeBottom,
    [left]: nudgeRight,
  }.true ?? nudgeBottom;

  if (top) {
    dimensions.left = targetLeft + (width / 2) - (contentWidth / 2);
    dimensions.top = targetTop - contentHeight;
    dimensions.transform = `translate(0, ${-nudge}px)`;
  } else if (bottom) {
    dimensions.left = targetLeft + (width / 2) - (contentWidth / 2);
    dimensions.top = targetTop + height;
    dimensions.transform = `translate(0, ${nudge}px)`;
  } else if (left) {
    dimensions.left = targetLeft - contentWidth;
    dimensions.top = targetTop + (height / 2) - (contentHeight / 2);
    dimensions.transform = `translate(${-nudge}px, 0)`;
  } else if (right) {
    dimensions.left = targetLeft + width;
    dimensions.top = targetTop + (height / 2) - (contentHeight / 2);
    dimensions.transform = `translate(${nudge}px, 0)`;
  }

  const minTop = contentMargin + nudge + window.scrollY;
  const minLeft = contentMargin + nudge + window.scrollX;
  const maxTop = windowHeight - contentMargin - nudge + window.scrollY;
  const maxLeft = windowWidth - contentMargin - nudge + window.scrollX;

  if (dimensions.top < minTop) {
    dimensions.top = minTop;
  }

  if (dimensions.top > maxTop) {
    dimensions.top = maxTop;
  }

  if (dimensions.left < minLeft) {
    dimensions.left = minLeft;
  }

  if (dimensions.left > maxLeft) {
    dimensions.left = maxLeft;
  }

  return dimensions;
};
