import { kebabCase } from 'lodash';
import { toMatchSnapshot } from 'jest-snapshot';
import registerRequireContextHook from 'babel-plugin-require-context-hook/register';
import { toMatchImageSnapshot } from 'jest-image-snapshot';
import ResizeObserver from 'resize-observer-polyfill';

registerRequireContextHook();

global.ResizeObserver = ResizeObserver;
global.IntersectionObserver = jest.fn(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
}));

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
  toMatchTooltipSnapshot(wrapper) {
    const tooltip = wrapper.findTooltip();

    return toMatchSnapshot.call(this, tooltip.element);
  },
  toMatchMenuSnapshot(wrapper) {
    const menu = wrapper.findMenu();

    return toMatchSnapshot.call(this, menu.element);
  },
  toEmit(wrapper, event, data) {
    const emittedEvents = wrapper.emitted(event);

    if (this.isNot) {
      try {
        expect(emittedEvents).not.toBeTruthy();
      } catch (err) {
        return err.matcherResult;
      }
    }

    try {
      expect(emittedEvents).toHaveLength(1);
    } catch (err) {
      return {
        pass: false,
        message: () => `Event '${event}' not emitted`,
      };
    }

    const [eventData] = emittedEvents[0];

    try {
      expect(eventData).toEqual(data);
    } catch (err) {
      return err.matcherResult;
    }

    return { pass: true };
  },
});
