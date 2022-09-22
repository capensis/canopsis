import flushPromises from 'flush-promises';
import { omit, pick } from 'lodash';
import Faker from 'faker';

import { mount, createVueInstance } from '@unit/utils/vue';
import { CONNECTOR_SIDES, FLOWCHART_KEY_CODES, SHAPES } from '@/constants';
import { shapeToForm } from '@/helpers/flowchart/shapes';
import { readTextFromClipboard, writeTextToClipboard } from '@/helpers/clipboard';
import uid from '@/helpers/uid';

import FlowchartSidebar from '@/components/common/flowchart/flowchart-editor.vue';

jest.mock('@/helpers/uid', () => {
  const originalModule = jest.requireActual('@/helpers/uid');

  return jest.fn(originalModule.default);
});
jest.mock('@/helpers/clipboard', () => ({
  readTextFromClipboard: jest.fn(),
  writeTextToClipboard: jest.fn(),
}));

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(FlowchartSidebar, {
  localVue,
  attachTo: document.body,

  ...options,
});

const selectSvg = wrapper => wrapper.find('svg');
const selectShapeByType = (wrapper, type) => {
  switch (type) {
    case SHAPES.line:
      return wrapper.find(`[data-type="${type}"] ~ path`);
    default:
      return wrapper.find(`[data-type="${type}"]`);
  }
};

const triggerDocumentEvent = (event) => {
  document.dispatchEvent(event);
};

const triggerDocumentMouseEvent = (type, data) => {
  triggerDocumentEvent(new MouseEvent(type, data));
};

const triggerDocumentKeyboardEvent = (type, data) => {
  triggerDocumentEvent(new KeyboardEvent(type, data));
};

describe('flowchart-editor', () => {
  let wrapper;
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

  afterEach(() => {
    wrapper.destroy();
  });

  test('Shape selected after mouse events triggered', async () => {
    wrapper = snapshotFactory({
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
    wrapper = snapshotFactory({
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
    wrapper = snapshotFactory({
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
    wrapper = snapshotFactory({
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
    wrapper = snapshotFactory({
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

    wrapper = snapshotFactory({
      propsData: {
        shapes: { rect, circle, line },
        viewBox,
      },
    });

    await flushPromises();

    const svg = selectSvg(wrapper);

    triggerDocumentMouseEvent('mousemove', {
      clientX: rect.x - 10,
      clientY: rect.y - 10,
    });

    await svg.trigger('mousedown');

    triggerDocumentMouseEvent('mousemove', {
      clientX: rect.x + rect.width + 10,
      clientY: rect.y + rect.height + 10,
    });

    await flushPromises();

    triggerDocumentMouseEvent('mouseup');

    await flushPromises();

    expect(wrapper).toEmit('update:selected', [SHAPES.rect, SHAPES.line]);
  });

  test('Shapes selected after mouse events triggered on container with shift', async () => {
    const { rhombus, ellipse, storage } = shapes;

    wrapper = snapshotFactory({
      propsData: {
        shapes: { rhombus, ellipse, storage },
        viewBox,
        selected: [SHAPES.rhombus, SHAPES.ellipse],
      },
    });

    await flushPromises();

    const svg = selectSvg(wrapper);

    triggerDocumentMouseEvent('mousemove', {
      clientX: rhombus.x - 10,
      clientY: rhombus.y - 10,
    });

    await svg.trigger('mousedown');

    triggerDocumentMouseEvent('mousemove', {
      clientX: rhombus.x + rhombus.width + 10,
      clientY: rhombus.y + rhombus.height + 10,
    });

    await flushPromises();

    triggerDocumentMouseEvent('mouseup', { shiftKey: true });

    await flushPromises();

    expect(wrapper).toEmit('update:selected', [SHAPES.storage]);
  });

  test('Shapes removed after keyboard event triggered', async () => {
    wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
        selected: [SHAPES.rhombus, SHAPES.ellipse],
      },
    });

    await flushPromises();

    triggerDocumentKeyboardEvent('keydown', {
      keyCode: FLOWCHART_KEY_CODES.delete,
    });

    await flushPromises();

    expect(wrapper).toEmit('input', omit(shapes, [SHAPES.ellipse, SHAPES.rhombus]));
  });

  test('Shapes moved up after keyboard event triggered', async () => {
    wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
        selected: [SHAPES.line, SHAPES.document],
      },
    });

    await flushPromises();

    triggerDocumentKeyboardEvent('keydown', {
      keyCode: FLOWCHART_KEY_CODES.arrowUp,
    });

    await flushPromises();

    expect(wrapper).toEmit('input', {
      ...shapes,
      line: {
        ...shapes.line,
        points: shapes.line.points.map(point => ({
          ...point,
          y: point.y - 5,
        })),
      },
      document: {
        ...shapes.document,
        y: shapes.document.y - 5,
      },
    });
  });

  test('Shapes moved right after keyboard event triggered', async () => {
    wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
        selected: [SHAPES.process, SHAPES.arrowLine],
      },
    });

    await flushPromises();

    triggerDocumentKeyboardEvent('keydown', {
      keyCode: FLOWCHART_KEY_CODES.arrowRight,
    });

    await flushPromises();

    const arrowLine = shapes[SHAPES.arrowLine];

    expect(wrapper).toEmit('input', {
      ...shapes,
      [SHAPES.arrowLine]: {
        ...arrowLine,
        points: arrowLine.points.map(point => ({
          ...point,
          x: point.x + 5,
        })),
      },
      process: {
        ...shapes.process,
        x: shapes.process.x + 5,
      },
    });
  });

  test('Shapes moved up after keyboard event triggered', async () => {
    wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
        selected: [SHAPES.line, SHAPES.ellipse],
      },
    });

    await flushPromises();

    triggerDocumentKeyboardEvent('keydown', {
      keyCode: FLOWCHART_KEY_CODES.arrowDown,
    });

    await flushPromises();

    expect(wrapper).toEmit('input', {
      ...shapes,
      line: {
        ...shapes.line,
        points: shapes.line.points.map(point => ({
          ...point,
          y: point.y + 5,
        })),
      },
      ellipse: {
        ...shapes.ellipse,
        y: shapes.ellipse.y + 5,
      },
    });
  });

  test('Shapes moved left after keyboard event triggered', async () => {
    wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
        selected: [SHAPES.process, SHAPES.arrowLine],
      },
    });

    await flushPromises();

    triggerDocumentKeyboardEvent('keydown', {
      keyCode: FLOWCHART_KEY_CODES.arrowLeft,
    });

    await flushPromises();

    const arrowLine = shapes[SHAPES.arrowLine];

    expect(wrapper).toEmit('input', {
      ...shapes,
      [SHAPES.arrowLine]: {
        ...arrowLine,
        points: arrowLine.points.map(point => ({
          ...point,
          x: point.x - 5,
        })),
      },
      process: {
        ...shapes.process,
        x: shapes.process.x - 5,
      },
    });
  });

  test('Shapes copied and pasted after keyboard event triggered', async () => {
    const copiedRectId = Faker.datatype.string();
    const copiedCircleId = Faker.datatype.string();
    uid
      .mockReturnValueOnce(copiedRectId)
      .mockReturnValueOnce(copiedCircleId);

    wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
        selected: [SHAPES.rect, SHAPES.circle],
      },
    });

    await flushPromises();

    triggerDocumentKeyboardEvent('keydown', {
      keyCode: FLOWCHART_KEY_CODES.keyC,
      ctrlKey: true,
    });

    await flushPromises();

    const copiedData = JSON.stringify(pick(shapes, [SHAPES.rect, SHAPES.circle]));

    expect(writeTextToClipboard).toBeCalledWith(copiedData);

    readTextFromClipboard.mockReturnValueOnce(copiedData);

    await flushPromises();

    triggerDocumentKeyboardEvent('keydown', {
      keyCode: FLOWCHART_KEY_CODES.keyV,
      ctrlKey: true,
    });

    await flushPromises();

    expect(wrapper).toEmit('input', {
      ...shapes,
      [copiedRectId]: {
        ...shapes.rect,
        _id: copiedRectId,
      },
      [copiedCircleId]: {
        ...shapes.circle,
        _id: copiedCircleId,
      },
    });
  });

  test('Shape moved after mouse event triggered', async () => {
    const startX = Faker.datatype.number({ min: 0, precision: 5 });
    const startY = Faker.datatype.number({ min: 0, precision: 5 });
    const diffX = Faker.datatype.number({ min: 0, precision: 5 });
    const diffY = Faker.datatype.number({ min: 0, precision: 5 });

    wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
      },
    });

    await flushPromises();

    triggerDocumentMouseEvent('mousemove', {
      clientX: startX,
      clientY: startY,
    });

    await selectShapeByType(wrapper, SHAPES.circle).trigger('mousedown');

    triggerDocumentMouseEvent('mousemove', {
      clientX: startX + diffX,
      clientY: startY + diffY,
    });

    await triggerDocumentMouseEvent('mouseup');

    const offsetX = diffX + startX;
    const offsetY = diffY + startY;

    expect(wrapper).toEmit('input', {
      ...shapes,
      circle: {
        ...shapes.circle,
        x: shapes.circle.x + offsetX,
        y: shapes.circle.y + offsetY,
      },
    });
  });

  test('Shapes moved after mouse event triggered', async () => {
    const startX = Faker.datatype.number({ min: 0, precision: 5 });
    const startY = Faker.datatype.number({ min: 0, precision: 5 });
    const diffX = Faker.datatype.number({ min: 0, precision: 5 });
    const diffY = Faker.datatype.number({ min: 0, precision: 5 });

    wrapper = snapshotFactory({
      propsData: {
        shapes,
        viewBox,
        selected: [SHAPES.rect, SHAPES.line, SHAPES.circle],
      },
    });

    await flushPromises();

    triggerDocumentMouseEvent('mousemove', {
      clientX: startX,
      clientY: startY,
    });

    await selectShapeByType(wrapper, SHAPES.circle).trigger('mousedown');

    triggerDocumentMouseEvent('mousemove', {
      clientX: startX + diffX,
      clientY: startY + diffY,
    });

    await triggerDocumentMouseEvent('mouseup');

    const offsetX = diffX + startX;
    const offsetY = diffY + startY;

    expect(wrapper).toEmit('input', {
      ...shapes,
      circle: {
        ...shapes.circle,
        x: shapes.circle.x + offsetX,
        y: shapes.circle.y + offsetY,
      },
      rect: {
        ...shapes.rect,
        x: shapes.rect.x + offsetX,
        y: shapes.rect.y + offsetY,
      },
      line: {
        ...shapes.line,
        points: shapes.line.points.map(point => ({
          ...point,
          x: point.x + offsetX,
          y: point.y + offsetY,
        })),
      },
    });
  });

  test('Connected line moved after move connected shape', async () => {
    const clientX = 50;
    const clientY = 50;

    const lineShape = {
      ...shapes.line,
      connectedTo: [shapes.rect._id],
    };
    const [connectedPoint] = lineShape.points;
    const rectShape = {
      ...shapes.rect,
      connections: [{
        shapeId: lineShape._id,
        pointId: connectedPoint._id,
        side: CONNECTOR_SIDES.top,
      }],
    };

    wrapper = snapshotFactory({
      propsData: {
        shapes: {
          line: lineShape,
          rect: rectShape,
        },
        viewBox,
        selected: [SHAPES.rect],
      },
    });

    await flushPromises();

    triggerDocumentMouseEvent('mousemove', {
      clientX: 0,
      clientY: 0,
    });

    await selectShapeByType(wrapper, SHAPES.rect).trigger('mousedown');

    triggerDocumentMouseEvent('mousemove', {
      clientX,
      clientY,
    });

    await triggerDocumentMouseEvent('mouseup');

    expect(wrapper).toEmit('input', {
      line: {
        ...lineShape,
        points: [
          {
            ...connectedPoint,
            x: clientX + rectShape.height / 2,
            y: clientY,
          },
          lineShape.points[1],
        ],
      },
      rect: {
        ...rectShape,
        x: clientX,
        y: clientY,
      },
    });
  });

  test('Line connected after mouseover with point', async () => {
    const clientX = 50;
    const clientY = 50;

    wrapper = snapshotFactory({
      propsData: {
        shapes: pick(shapes, [SHAPES.rect, SHAPES.line]),
        viewBox,
      },
    });

    await flushPromises();

    const line = selectShapeByType(wrapper, SHAPES.line);

    await line.trigger('mousedown');
    await line.trigger('mouseup');

    const linePointCircle = wrapper.find('circle[cursor="crosshair"]');
    await linePointCircle.trigger('mousedown');

    const rectTopConnector = wrapper.findAll('g')
      .at(2)
      .find('path');

    await rectTopConnector.trigger('mouseenter');
    await rectTopConnector.trigger('mouseup', {
      clientX,
      clientY,
    });

    await flushPromises();

    expect(wrapper).toEmit('input', {
      line: {
        ...shapes.line,
        points: [
          {
            ...shapes.line.points[0],
            x: shapes.rect.width / 2,
            y: 0,
          },
          shapes.line.points[1],
        ],
        connectedTo: [SHAPES.rect],
      },
      rect: {
        ...shapes.rect,
        connections: [
          {
            pointId: shapes.line.points[0]._id,
            side: CONNECTOR_SIDES.top,
            shapeId: shapes.line._id,
          },
        ],
      },
    });
  });

  test('Line unconnected after mouseleave', async () => {
    const lineShape = {
      ...shapes.line,
      connectedTo: [shapes.rect._id],
    };
    const [connectedPoint] = lineShape.points;
    const rectShape = {
      ...shapes.rect,
      connections: [{
        shapeId: lineShape._id,
        pointId: connectedPoint._id,
        side: CONNECTOR_SIDES.top,
      }],
    };

    wrapper = snapshotFactory({
      propsData: {
        shapes: {
          rect: rectShape,
          line: lineShape,
        },
        viewBox,
      },
    });

    await flushPromises();

    const line = selectShapeByType(wrapper, SHAPES.line);

    await line.trigger('mousedown');
    await line.trigger('mouseup');

    const linePointCircle = wrapper.find('circle[cursor="crosshair"]');
    await linePointCircle.trigger('mousedown');

    const rectTopConnector = wrapper.findAll('g')
      .at(2)
      .find('path');

    await rectTopConnector.trigger('mouseleave');

    await selectSvg(wrapper).trigger('mouseup');

    await flushPromises();

    expect(wrapper).toEmit('input', {
      line: shapes.line,
      rect: shapes.rect,
    });
  });

  test('Line unconnected after move connected shape', async () => {
    const clientX = 50;
    const clientY = 50;

    const lineShape = {
      ...shapes.line,
      connectedTo: [shapes.rect._id],
    };
    const [connectedPoint] = lineShape.points;
    const rectShape = {
      ...shapes.rect,
      connections: [{
        shapeId: lineShape._id,
        pointId: connectedPoint._id,
        side: CONNECTOR_SIDES.top,
      }],
    };

    wrapper = snapshotFactory({
      propsData: {
        shapes: {
          line: lineShape,
          rect: rectShape,
        },
        viewBox,
      },
    });

    await flushPromises();

    triggerDocumentMouseEvent('mousemove', {
      clientX: 0,
      clientY: 0,
    });

    await selectShapeByType(wrapper, SHAPES.line).trigger('mousedown');

    triggerDocumentMouseEvent('mousemove', {
      clientX,
      clientY,
    });

    await flushPromises();

    await selectShapeByType(wrapper, SHAPES.line).trigger('mouseup', {
      clientX,
      clientY,
    });

    expect(wrapper).toEmit('input', {
      line: {
        ...lineShape,
        connectedTo: [],
        points: lineShape.points.map(point => ({
          ...point,
          x: point.x + clientX,
          y: point.y + clientY,
        })),
      },
      rect: {
        ...rectShape,
        connections: [],
      },
    });
  });

  test('Renders `flowchart-editor` with all shapes', async () => {
    getPointAtLength
      .mockReturnValueOnce({ x: 1, y: 2 })
      .mockReturnValueOnce({ x: 3, y: 4 })
      .mockReturnValueOnce({ x: 5, y: 6 })
      .mockReturnValueOnce({ x: 7, y: 8 })
      .mockReturnValueOnce({ x: 9, y: 10 })
      .mockReturnValueOnce({ x: 11, y: 12 });

    wrapper = snapshotFactory({
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

    wrapper = snapshotFactory({
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
