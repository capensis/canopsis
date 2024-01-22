import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createActivatorElementStub } from '@unit/stubs/vuetify';

import { LINE_TYPES } from '@/constants';

import { uid } from '@/helpers/uid';
import {
  arrowLineShapeToForm,
  bidirectionalArrowLineShapeToForm,
  circleShapeToForm,
  documentShapeToForm,
  ellipseShapeToForm,
  imageShapeToForm,
  lineShapeToForm,
  parallelogramShapeToForm,
  processShapeToForm,
  rectShapeToForm,
  rhombusShapeToForm,
  storageShapeToForm,
} from '@/helpers/flowchart/shapes';
import { getFileDataUrlContent } from '@/helpers/file/file-select';
import { getImageProperties } from '@/helpers/file/image';

import FlowchartSidebar from '@/components/common/flowchart/flowchart-sidebar.vue';

jest.mock('@/helpers/uid');
jest.mock('@/helpers/file/file-select', () => ({
  getFileDataUrlContent: jest.fn(),
}));
jest.mock('@/helpers/file/image', () => ({
  getImageProperties: jest.fn(),
}));

const clearFileSelector = jest.fn();
const stubs = {
  'flowchart-color-field': true,
  'file-selector': {
    template: `
      <div class='file-selector'>
        <slot />
      </div>
    `,
    methods: {
      clear: clearFileSelector,
    },
  },
  'image-shape-icon': true,
  'v-tooltip': createActivatorElementStub('v-tooltip'),
};

const snapshotStubs = {
  'flowchart-color-field': true,
  'file-selector': true,
  'image-shape-icon': true,
};

const selectButtons = wrapper => wrapper.findAll('v-btn-stub');
const selectButtonByIndex = (wrapper, index) => selectButtons(wrapper).at(index);
const selectIconButtons = wrapper => wrapper.findAll('v-btn-stub.flowchart-sidebar__button-icon');
const selectIconButtonByIndex = (wrapper, index) => selectIconButtons(wrapper)
  .at(index);
const selectFileSelector = wrapper => wrapper.find('.file-selector');

describe('flowchart-sidebar', () => {
  const viewBox = {
    x: 10,
    y: 10,
    width: 1000,
    height: 1000,
  };
  const properties = {
    fill: 'white',
    stroke: 'black',
    'stroke-width': 2,
  };
  const textProperties = {
    alignCenter: true,
    fontColor: 'black',
    fontSize: 12,
    justifyCenter: true,
  };
  const lineProperties = {
    stroke: 'black',
    'stroke-width': 2,
  };
  const lineTextProperties = {
    fontColor: 'black',
    fontSize: 12,
  };
  const x = 435;
  const y = 435;
  const width = 150;
  const height = 150;

  const factory = generateShallowRenderer(FlowchartSidebar, { stubs });
  const snapshotFactory = generateRenderer(FlowchartSidebar, { stubs: snapshotStubs });

  test('Rect shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 0);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: rectShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'Rectangle',
        textProperties,
      }),
    });
  });

  test('Rounded rect shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 1);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: rectShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties: {
          rx: 20,
          ry: 20,
          fill: 'white',
          stroke: 'black',
          'stroke-width': 2,
        },
        text: 'Rounded rectangle',
        textProperties,
      }),
    });
  });

  test('Square shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 2);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: rectShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        aspectRatio: true,
        properties,
        text: 'Square',
        textProperties,
      }),
    });
  });

  test('Rhombus shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 3);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: rhombusShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'Rhombus',
        textProperties,
      }),
    });
  });

  test('Circle shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 4);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: circleShapeToForm({
        _id: id,
        diameter: 150,
        x,
        y,
        properties,
        text: 'Circle',
        textProperties,
      }),
    });
  });

  test('Ellipse shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 5);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: ellipseShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'Ellipse',
        textProperties,
      }),
    });
  });

  test('Parallelogram shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 6);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: parallelogramShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'Parallelogram',
        textProperties,
      }),
    });
  });

  test('Process shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 7);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: processShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'Process',
        textProperties,
      }),
    });
  });

  test('Document shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 8);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: documentShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'Document',
        textProperties,
      }),
    });
  });

  test('Storage shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 9);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: storageShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'Storage',
        textProperties,
      }),
    });
  });

  test('Curve shape added after trigger button', async () => {
    const firstPointId = Faker.datatype.string();
    const secondPointId = Faker.datatype.string();
    const id = Faker.datatype.string();
    uid
      .mockReturnValueOnce(firstPointId)
      .mockReturnValueOnce(secondPointId)
      .mockReturnValueOnce(id);

    const points = [
      {
        _id: firstPointId,
        x,
        y: 585,
      },
      {
        _id: secondPointId,
        x: 585,
        y,
      },
    ];
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 10);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: lineShapeToForm({
        _id: id,
        points,
        properties: lineProperties,
        lineType: LINE_TYPES.horizontalCurve,
        text: '',
        textProperties: lineTextProperties,
      }),
    });
  });

  test('Curve arrow line shape added after trigger button', async () => {
    const firstPointId = Faker.datatype.string();
    const secondPointId = Faker.datatype.string();
    const id = Faker.datatype.string();
    uid
      .mockReturnValueOnce(firstPointId)
      .mockReturnValueOnce(secondPointId)
      .mockReturnValueOnce(id);

    const points = [
      {
        _id: firstPointId,
        x,
        y: 585,
      },
      {
        _id: secondPointId,
        x: 585,
        y,
      },
    ];
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 11);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: arrowLineShapeToForm({
        _id: id,
        points,
        properties: lineProperties,
        lineType: LINE_TYPES.horizontalCurve,
        text: '',
        textProperties: lineTextProperties,
      }),
    });
  });

  test('Bidirectional arrow curve line shape added after trigger button', async () => {
    const firstPointId = Faker.datatype.string();
    const secondPointId = Faker.datatype.string();
    const id = Faker.datatype.string();
    uid
      .mockReturnValueOnce(firstPointId)
      .mockReturnValueOnce(secondPointId)
      .mockReturnValueOnce(id);

    const points = [
      {
        _id: firstPointId,
        x,
        y: 585,
      },
      {
        _id: secondPointId,
        x: 585,
        y,
      },
    ];
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 12);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: bidirectionalArrowLineShapeToForm({
        _id: id,
        points,
        properties: lineProperties,
        lineType: LINE_TYPES.horizontalCurve,
        text: '',
        textProperties: lineTextProperties,
      }),
    });
  });

  test('Line shape added after trigger button', async () => {
    const firstPointId = Faker.datatype.string();
    const secondPointId = Faker.datatype.string();
    const id = Faker.datatype.string();
    uid
      .mockReturnValueOnce(firstPointId)
      .mockReturnValueOnce(secondPointId)
      .mockReturnValueOnce(id);

    const points = [
      {
        _id: firstPointId,
        x,
        y: 585,
      },
      {
        _id: secondPointId,
        x: 585,
        y,
      },
    ];
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 13);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: lineShapeToForm({
        _id: id,
        points,
        properties: lineProperties,
        lineType: LINE_TYPES.line,
        text: '',
        textProperties: lineTextProperties,
      }),
    });
  });

  test('Arrow line shape added after trigger button', async () => {
    const firstPointId = Faker.datatype.string();
    const secondPointId = Faker.datatype.string();
    const id = Faker.datatype.string();
    uid
      .mockReturnValueOnce(firstPointId)
      .mockReturnValueOnce(secondPointId)
      .mockReturnValueOnce(id);

    const points = [
      {
        _id: firstPointId,
        x,
        y: 585,
      },
      {
        _id: secondPointId,
        x: 585,
        y,
      },
    ];
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 14);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: arrowLineShapeToForm({
        _id: id,
        points,
        properties: lineProperties,
        lineType: LINE_TYPES.line,
        text: '',
        textProperties: lineTextProperties,
      }),
    });
  });

  test('Bidirectional arrow line shape added after trigger button', async () => {
    const firstPointId = Faker.datatype.string();
    const secondPointId = Faker.datatype.string();
    const id = Faker.datatype.string();
    uid
      .mockReturnValueOnce(firstPointId)
      .mockReturnValueOnce(secondPointId)
      .mockReturnValueOnce(id);

    const points = [
      {
        _id: firstPointId,
        x,
        y: 585,
      },
      {
        _id: secondPointId,
        x: 585,
        y,
      },
    ];
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 15);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: bidirectionalArrowLineShapeToForm({
        _id: id,
        points,
        properties: lineProperties,
        lineType: LINE_TYPES.line,
        text: '',
        textProperties: lineTextProperties,
      }),
    });
  });

  test('Text shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 16);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: rectShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties: {
          fill: 'transparent',
        },
        text: 'Text',
        textProperties,
      }),
    });
  });

  test('Textbox shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const button = selectButtonByIndex(wrapper, 17);
    await button.triggerCustomEvent('click');

    expect(wrapper).toEmit('input', {
      [id]: rectShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties: {
          fill: 'transparent',
        },
        text: '<h2>Heading</h2><p>Paragraph</p>',
        textProperties: {
          hidden: true,
        },
      }),
    });
  });

  test('Image shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    const filePath = `/${Faker.datatype.string()}.png`;
    const fileName = 'file.png';
    getFileDataUrlContent.mockReturnValueOnce(filePath);
    const imageProperties = {
      width: 100,
      height: 200,
    };
    getImageProperties.mockReturnValueOnce(imageProperties);
    uid.mockReturnValueOnce(id);

    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const file = new File([new ArrayBuffer(1)], fileName);

    const fileSelector = selectFileSelector(wrapper);
    await fileSelector.triggerCustomEvent('change', [file]);

    await flushPromises();

    expect(wrapper).toEmit('input', {
      [id]: imageShapeToForm({
        ...imageProperties,
        _id: id,
        x: 460,
        y: 410,
        src: filePath,
        properties: {
          fill: 'transparent',
          stroke: 'transparent',
        },
        text: fileName,
        aspectRatio: true,
      }),
    });
  });

  test('Image shape added and cropped after trigger button', async () => {
    const id = Faker.datatype.string();
    const filePath = `/${Faker.datatype.string()}.png`;
    const fileName = 'file.png';
    getFileDataUrlContent.mockReturnValueOnce(filePath);
    const imageProperties = {
      width: 1500,
      height: 2000,
    };
    getImageProperties.mockReturnValueOnce(imageProperties);
    uid.mockReturnValueOnce(id);

    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const file = new File([new ArrayBuffer(1)], fileName);

    const fileSelector = selectFileSelector(wrapper);
    await fileSelector.triggerCustomEvent('change', [file]);

    await flushPromises();

    expect(wrapper).toEmit('input', {
      [id]: imageShapeToForm({
        _id: id,
        width: 562.5,
        height: 750,
        x: 228.75,
        y: 135,
        src: filePath,
        properties: {
          fill: 'transparent',
          stroke: 'transparent',
        },
        text: fileName,
        aspectRatio: true,
      }),
    });
  });

  test('Image shape added and cropped after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);

    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    const iconButton = selectIconButtonByIndex(wrapper, 1);
    await iconButton.triggerCustomEvent('click');

    await flushPromises();

    expect(wrapper).toEmit('input', {
      [id]: imageShapeToForm({
        _id: id,
        x: 435,
        y: 435,
        height: 150,
        width: 150,
        src: undefined,
        svg: '<svg><text>special-asset</text></svg>',
        properties: {
          fill: 'black',
        },
        aspectRatio: true,
      }),
    });
  });

  test('Renders `flowchart-sidebar` with form', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    await wrapper.openAllExpansionPanels();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `flowchart-sidebar` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        shapes: {
          rect: rectShapeToForm({ _id: 'rect' }),
          circle: circleShapeToForm({ _id: 'circle' }),
        },
        cursorStyle: 'pointer',
        backgroundColor: '#000',
        viewBox,
      },
    });

    await wrapper.openAllExpansionPanels();

    expect(wrapper).toMatchSnapshot();
  });
});
