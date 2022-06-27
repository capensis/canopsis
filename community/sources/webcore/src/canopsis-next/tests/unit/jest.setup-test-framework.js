import { kebabCase } from 'lodash';
import { toMatchSnapshot } from 'jest-snapshot';
import registerRequireContextHook from 'babel-plugin-require-context-hook/register';
import { toMatchImageSnapshot } from 'jest-image-snapshot';
import ResizeObserver from 'resize-observer-polyfill';

registerRequireContextHook();

global.ResizeObserver = ResizeObserver;

expect.extend({
  toMatchImageSnapshot,
  toMatchCanvasSnapshot(canvas, options, ...args) {
    const img = canvas.toDataURL();
    const data = img.replace(/^data:image\/(png|jpg);base64,/, '');
    const newOptions = {
      comparisonMethod: 'ssim',
      diffDirection: 'vertical',
      customDiffConfig: {
        ssim: 'fast',
      },
      failureThreshold: 0.05,
      failureThresholdType: 'percent',
      customSnapshotIdentifier: ({ currentTestName, counter }) => (
        kebabCase(`${currentTestName.replace(/(.*\sRenders\s)|(.$)/g, '')}-${counter}`)
      ),

      ...options,
    };

    return toMatchImageSnapshot.call(this, data, newOptions, ...args);
  },
  toMatchTooltipSnapshot(wrapper) {
    const tooltip = wrapper.findTooltip();

    return toMatchSnapshot.call(this, tooltip.element);
  },
  toMatchMenuSnapshot(wrapper) {
    const menu = wrapper.findMenu();

    return toMatchSnapshot.call(this, menu.element);
  },
});
