import { kebabCase } from 'lodash';
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
      failureThreshold: 0.02,
      failureThresholdType: 'percent',
      customSnapshotIdentifier: ({ currentTestName, counter }) => (
        kebabCase(`${currentTestName.replace(/(.*\sRenders\s)|(.$)/g, '')}-${counter}`)
      ),

      ...options,
    };

    return toMatchImageSnapshot.call(this, data, newOptions, ...args);
  },
});
