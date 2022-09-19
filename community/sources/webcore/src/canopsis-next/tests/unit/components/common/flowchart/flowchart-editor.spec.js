import flushPromises from 'flush-promises';

import { mount, createVueInstance } from '@unit/utils/vue';
import { SHAPES } from '@/constants';
import { shapeToForm } from '@/helpers/flowchart/shapes';

import FlowchartSidebar from '@/components/common/flowchart/flowchart-editor.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(FlowchartSidebar, {
  localVue,

  ...options,
});

describe('flowchart-editor', () => {
  const viewBox = {
    x: 10,
    y: 10,
    width: 1000,
    height: 1000,
  };
  const shapes = Object.values(SHAPES).reduce((acc, type) => {
    acc[type] = shapeToForm({ type, _id: type });

    return acc;
  }, {});

  const getTotalLength = jest.fn();
  const getPointAtLength = jest.fn();

  SVGElement.prototype.getTotalLength = getTotalLength;
  SVGElement.prototype.getPointAtLength = getPointAtLength;

  jest.spyOn(window, 'getComputedStyle').mockReturnValue(viewBox);

  beforeEach(() => {
    getPointAtLength.mockClear();
    getTotalLength.mockClear();
  });

  test('Renders `flowchart-editor` with all shapes', async () => {
    getPointAtLength
      .mockReturnValueOnce({ x: 1, y: 2 })
      .mockReturnValueOnce({ x: 3, y: 4 })
      .mockReturnValueOnce({ x: 5, y: 6 })
      .mockReturnValueOnce({ x: 7, y: 8 })
      .mockReturnValueOnce({ x: 9, y: 10 })
      .mockReturnValueOnce({ x: 11, y: 12 });

    const wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `flowchart-editor` with custom props', async () => {
    getPointAtLength
      .mockReturnValueOnce({ x: 1, y: 2 })
      .mockReturnValueOnce({ x: 3, y: 4 })
      .mockReturnValueOnce({ x: 5, y: 6 })
      .mockReturnValueOnce({ x: 7, y: 8 })
      .mockReturnValueOnce({ x: 9, y: 10 })
      .mockReturnValueOnce({ x: 11, y: 12 });

    const wrapper = snapshotFactory({
      propsData: {
        shapes,
        gridSize: 10,
        readonly: true,
        backgroundColor: '#000',
        pointSize: 32,
        selectionColor: '#000',
        selectionPadding: 9,
        cursorStyle: 'pointer',
        viewBox,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
