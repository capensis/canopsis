import { createActivatorElementStub } from '@unit/stubs/vuetify';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { LINE_TYPES } from '@/constants';

import { getFileDataUrlContent } from '@/helpers/file/file-select';
import { getImageProperties } from '@/helpers/file/image';
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
import { uid } from '@/helpers/uid';

import FlowchartSidebar from '@/components/common/flowchart/flowchart-sidebar.vue';

vi.mock('@/helpers/uid');
vi.mock('@/helpers/file/file-select', () => ({
  getFileDataUrlContent: vi.fn(),
}));
vi.mock('@/helpers/file/image', () => ({
  getImageProperties: vi.fn(),
}));

const clearFileSelector = vi.fn();
const FileSelectorStub = {
  name: 'file-selector',
  emits: ['change'],
  methods: {
    clear: clearFileSelector,
  },
  template: `
    <div class='file-selector'>
      <slot />
    </div>
  `,
};
const stubs = {
  FileSelector: FileSelectorStub,
  'flowchart-color-field': true,
  'file-selector': FileSelectorStub,
  'image-shape-icon': true,
  'c-tooltip': createActivatorElementStub('c-tooltip'),
};

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

  test('Sidebar uses the fixed editor toolbar width', () => {
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    expect(wrapper.findComponent({ name: 'VNavigationDrawer' }).props('width')).toBe(320);
  });

  test('Rect shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);
    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    wrapper.vm.buttons[0].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
      [id]: rectShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'flowchart.shapes.rectangle',
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

    wrapper.vm.buttons[1].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
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
        text: 'flowchart.shapes.roundedRectangle',
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

    wrapper.vm.buttons[2].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
      [id]: rectShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        aspectRatio: true,
        properties,
        text: 'flowchart.shapes.square',
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

    wrapper.vm.buttons[3].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
      [id]: rhombusShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'flowchart.shapes.rhombus',
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

    wrapper.vm.buttons[4].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
      [id]: circleShapeToForm({
        _id: id,
        diameter: 150,
        x,
        y,
        properties,
        text: 'flowchart.shapes.circle',
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

    wrapper.vm.buttons[5].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
      [id]: ellipseShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'flowchart.shapes.ellipse',
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

    wrapper.vm.buttons[6].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
      [id]: parallelogramShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'flowchart.shapes.parallelogram',
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

    wrapper.vm.buttons[7].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
      [id]: processShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'flowchart.shapes.process',
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

    wrapper.vm.buttons[8].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
      [id]: documentShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'flowchart.shapes.document',
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

    wrapper.vm.buttons[9].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
      [id]: storageShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties,
        text: 'flowchart.shapes.storage',
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

    wrapper.vm.buttons[10].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
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

    wrapper.vm.buttons[11].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
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

    wrapper.vm.buttons[12].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
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

    wrapper.vm.buttons[13].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
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

    wrapper.vm.buttons[14].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
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

    wrapper.vm.buttons[15].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
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

    wrapper.vm.buttons[16].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
      [id]: rectShapeToForm({
        _id: id,
        height,
        width,
        x,
        y,
        properties: {
          fill: 'transparent',
        },
        text: 'flowchart.shapes.text',
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

    wrapper.vm.buttons[17].action();
    await flushPromises();

    expect(wrapper).toEmitInput({
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

    wrapper.vm.$.refs = {
      fileSelector: {
        clear: clearFileSelector,
      },
    };

    await wrapper.vm.addImage([file]);

    await flushPromises();

    expect(wrapper).toEmitInput({
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

    wrapper.vm.$.refs = {
      fileSelector: {
        clear: clearFileSelector,
      },
    };

    await wrapper.vm.addImage([file]);

    await flushPromises();

    expect(wrapper).toEmitInput({
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

  test('Icon asset shape added after trigger button', async () => {
    const id = Faker.datatype.string();
    uid.mockReturnValueOnce(id);

    const wrapper = factory({
      propsData: {
        shapes: {},
        viewBox,
      },
    });

    await wrapper.vm.addIconAsset('<svg><text>special-asset</text></svg>');

    await flushPromises();

    expect(wrapper).toEmitInput({
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
});
