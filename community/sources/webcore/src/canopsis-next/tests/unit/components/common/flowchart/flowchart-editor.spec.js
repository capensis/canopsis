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

const selectSvg = wrapper => wrapper.find('svg');
const selectShapeByType = (wrapper, type) => wrapper.find(`[data-type="${type}"]`);

const triggerDocumentEvent = (type, data) => {
  document.dispatchEvent(new MouseEvent(type, data));
};

describe('flowchart-editor', () => {
  const viewBox = {
    x: 10,
    y: 10,
    width: 1000,
    height: 1000,
  };
  const shapes = Object.values(SHAPES).reduce((acc, type) => {
    acc[type] = shapeToForm({
      type,
      _id: type,
      properties: {
        'data-type': type,
      },
    });

    return acc;
  }, {});

  const getTotalLength = jest.fn();
  const getPointAtLength = jest.fn();
  const createSVGPoint = jest.fn()
    .mockImplementation(() => ({
      x: 0,
      y: 0,
      matrixTransform() {
        return this;
      },
    }));
  const getScreenCTM = jest.fn().mockReturnValue({
    inverse: jest.fn(),
  });

  SVGElement.prototype.createSVGPoint = createSVGPoint;
  SVGElement.prototype.getTotalLength = getTotalLength;
  SVGElement.prototype.getPointAtLength = getPointAtLength;
  SVGElement.prototype.getScreenCTM = getScreenCTM;

  jest.spyOn(window, 'getComputedStyle').mockReturnValue(viewBox);

  beforeEach(() => {
    getPointAtLength.mockClear();
    getTotalLength.mockClear();
  });

  test('Shape selected after mouse events triggered', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
      },
    });

    await flushPromises();

    const rect = selectShapeByType(wrapper, SHAPES.rect);

    await rect.trigger('mousedown');
    await rect.trigger('mouseup');

    expect(wrapper).toEmit('update:selected', [SHAPES.rect]);
  });

  test('Second shape selected after mouse events triggered with ctrl', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
        selected: [SHAPES.rect],
      },
    });

    await flushPromises();

    await selectShapeByType(wrapper, SHAPES.circle)
      .trigger('mousedown', { ctrlKey: true });

    expect(wrapper).not.toEmit('update:selected');

    await selectShapeByType(wrapper, SHAPES.circle)
      .trigger('mouseup', { ctrlKey: true });

    expect(wrapper).toEmit('update:selected', [SHAPES.rect, SHAPES.circle]);
  });

  test('Shape unselected after mouse events triggered on already selected shape with ctrl', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
        selected: [SHAPES.rect, SHAPES.circle],
      },
    });

    await flushPromises();

    await selectShapeByType(wrapper, SHAPES.circle)
      .trigger('mouseup', { ctrlKey: true });

    expect(wrapper).toEmit('update:selected', [SHAPES.rect]);
  });

  test('Shapes unselected after mouse events triggered on already selected shape without ctrl', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
        selected: [SHAPES.rect, SHAPES.storage, SHAPES.circle],
      },
    });

    await flushPromises();

    await selectShapeByType(wrapper, SHAPES.storage)
      .trigger('mouseup', { ctrlKey: false });

    expect(wrapper).toEmit('update:selected', [SHAPES.storage]);
  });

  test('Selected shapes cleared and new selected after mouse events triggered without ctrl', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
        selected: [SHAPES.rect, SHAPES.circle],
      },
    });

    await flushPromises();

    await selectShapeByType(wrapper, SHAPES.rhombus)
      .trigger('mousedown', { ctrlKey: false });

    await flushPromises();

    expect(wrapper).toEmit('update:selected', [SHAPES.rhombus]);
  });

  test('Shapes selected after mouse events triggered on container', async () => {
    const { rect, circle, line } = shapes;

    const wrapper = snapshotFactory({
      propsData: {
        shapes: { rect, circle, line },
        viewBox,
      },
    });

    await flushPromises();

    const svg = selectSvg(wrapper);

    triggerDocumentEvent('mousemove', {
      clientX: rect.x - 10,
      clientY: rect.y - 10,
    });

    await svg.trigger('mousedown');

    triggerDocumentEvent('mousemove', {
      clientX: rect.x + rect.width + 10,
      clientY: rect.y + rect.height + 10,
    });

    await flushPromises();

    triggerDocumentEvent('mouseup');

    await flushPromises();

    expect(wrapper).toEmit('update:selected', [SHAPES.rect, SHAPES.line]);
  });

  test('Shapes selected after mouse events triggered on container with shift', async () => {
    const { rhombus, ellipse, storage } = shapes;

    const wrapper = snapshotFactory({
      propsData: {
        shapes: { rhombus, ellipse, storage },
        viewBox,
        selected: [SHAPES.rhombus, SHAPES.ellipse],
      },
    });

    await flushPromises();

    const svg = selectSvg(wrapper);

    triggerDocumentEvent('mousemove', {
      clientX: rhombus.x - 10,
      clientY: rhombus.y - 10,
    });

    await svg.trigger('mousedown');

    triggerDocumentEvent('mousemove', {
      clientX: rhombus.x + rhombus.width + 10,
      clientY: rhombus.y + rhombus.height + 10,
    });

    await flushPromises();

    triggerDocumentEvent('mouseup', { shiftKey: true });

    await flushPromises();

    expect(wrapper).toEmit('update:selected', [SHAPES.storage]);
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
