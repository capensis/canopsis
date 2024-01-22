import { cloneDeep, isFunction, kebabCase } from 'lodash';
import { toMatchSnapshot } from 'jest-snapshot';
import registerRequireContextHook from 'babel-plugin-require-context-hook/register';
import { toMatchImageSnapshot } from 'jest-image-snapshot';
import ResizeObserver from 'resize-observer-polyfill';
import flatten from 'flat';

registerRequireContextHook();

global.ResizeObserver = ResizeObserver;
global.IntersectionObserver = jest.fn(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
}));

Object.defineProperty(HTMLElement.prototype, 'innerText', {
  set(value) {
    this.textContent = value;
  },
  get() {
    return this.textContent;
  },
});

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
      failureThreshold: 0.1,
      failureThresholdType: 'percent',
      customSnapshotIdentifier: ({ currentTestName, counter }) => (
        kebabCase(`${currentTestName.replace(/(.*\sRenders\s)|(.$)/g, '')}-${counter}`)
      ),

    };

    return toMatchImageSnapshot.call(this, data, newOptions, ...args);
  },
  toMatchTooltipSnapshot(wrapper) {
    const tooltip = wrapper.findTooltip();

    return toMatchSnapshot.call(this, tooltip);
  },
  toMatchMenuSnapshot(wrapper) {
    const menu = wrapper.findMenu();

    return toMatchSnapshot.call(this, menu);
  },
  toEmit(wrapper, event, ...data) {
    // eslint-disable-next-line no-restricted-syntax
    const emittedEvents = wrapper.emitted(event);

    if (this.isNot) {
      try {
        expect(emittedEvents).not.toBeTruthy();
      } catch (err) {
        return err.matcherResult;
      }
    }

    try {
      if (!data.length) {
        expect(emittedEvents).toBeTruthy();

        return { pass: true };
      }

      expect(emittedEvents).toHaveLength(data.length);
    } catch (err) {
      return {
        pass: false,
        message: () => `Event '${event}' not emitted`,
      };
    }

    try {
      expect(
        cloneDeep(emittedEvents.map(events => events[0])),
      ).toEqual(data);
    } catch (err) {
      const { pass, message } = err.matcherResult ?? { pass: false, message: '' };

      return {
        pass,
        message: isFunction(message) ? message : () => message,
      };
    }

    return { pass: true };
  },
  toStructureEqual(received, expected) {
    const flattenReceived = flatten(received);
    const flattenExpected = flatten(expected);

    try {
      expect(flattenReceived).toEqual(Object.keys(flattenExpected).reduce((acc, key) => {
        acc[key] = expect.any(String);

        return acc;
      }, {}));

      return { pass: true };
    } catch (err) {
      return err.matcherResult;
    }
  },
  toBeDispatchedWith(received, expected) {
    try {
      expect(received).toBeCalledWith(
        expect.any(Object),
        expected,
        undefined,
      );

      return { pass: true };
    } catch (err) {
      return err.matcherResult;
    }
  },
});
