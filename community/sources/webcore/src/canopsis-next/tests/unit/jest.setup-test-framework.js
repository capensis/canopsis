import registerRequireContextHook from 'babel-plugin-require-context-hook/register';
import { toMatchImageSnapshot } from 'jest-image-snapshot';
import ResizeObserver from 'resize-observer-polyfill';

registerRequireContextHook();

global.ResizeObserver = ResizeObserver;

expect.extend({
  toMatchImageSnapshot,
  toMatchCanvasSnapshot(canvas, ...args) {
    const img = canvas.toDataURL();
    const data = img.replace(/^data:image\/(png|jpg);base64,/, '');

    return toMatchImageSnapshot.call(this, data, ...args);
  },
});
