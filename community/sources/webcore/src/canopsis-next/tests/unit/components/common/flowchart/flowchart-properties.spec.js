import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { lineShapeToForm, rectShapeToForm } from '@/helpers/flowchart/shapes';

import FlowchartProperties from '@/components/common/flowchart/flowchart-properties.vue';
import { LINE_TYPES, STROKE_TYPES } from '@/constants';

const stubs = {
  'flowchart-color-field': true,
  'flowchart-number-field': true,
  'flowchart-stroke-type-field': true,
  'flowchart-line-type-field': true,
};

const selectFillColorField = wrapper => wrapper.find('flowchart-color-field-stub[label="Fill"]');
const selectStrokeColorField = wrapper => wrapper.find('flowchart-color-field-stub[label="Stroke"]');
const selectStrokeWidthField = wrapper => wrapper.find('flowchart-number-field-stub[label="Stroke width"]');
const selectStrokeTypeField = wrapper => wrapper.find('flowchart-stroke-type-field-stub[label="Stroke type"]');
const selectLineTypeField = wrapper => wrapper.find('flowchart-line-type-field-stub');
const selectFontColorField = wrapper => wrapper.find('flowchart-color-field-stub[label="Font color"]');
const selectFontBackgroundColorField = wrapper => wrapper.find('flowchart-color-field-stub[label="Font background color"]');
const selectFontSizeField = wrapper => wrapper.find('flowchart-number-field-stub[label="Font size"]');

describe('flowchart-properties', () => {
  const factory = generateShallowRenderer(FlowchartProperties, { stubs });
  const snapshotFactory = generateRenderer(FlowchartProperties, { stubs });

  test('Fill changed after trigger color field', () => {
    const firstShape = rectShapeToForm({
      _id: 'first',
      properties: { fill: 'red' },
    });
    const secondShape = rectShapeToForm({
      _id: 'second',
      properties: { fill: 'transparent' },
    });
    const shapes = {
      first: firstShape,
      second: secondShape,
    };
    const wrapper = factory({
      propsData: {
        shapes,
        selected: Object.keys(shapes),
      },
    });

    const fillColorField = selectFillColorField(wrapper);

    const newFill = Faker.internet.color();

    fillColorField.vm.$emit('input', newFill);

    expect(wrapper).toEmit('input', {
      first: {
        ...firstShape,
        properties: {
          fill: newFill,
        },
      },
      second: {
        ...secondShape,
        properties: {
          fill: newFill,
        },
      },
    });
  });

  test('Stroke changed after trigger color field', () => {
    const firstShape = rectShapeToForm({
      _id: 'first',
      properties: { stroke: 'red' },
    });
    const secondShape = rectShapeToForm({
      _id: 'second',
      properties: { stroke: 'transparent' },
    });
    const shapes = {
      first: firstShape,
      second: secondShape,
    };
    const wrapper = factory({
      propsData: {
        shapes,
        selected: Object.keys(shapes),
      },
    });

    const strokeColorField = selectStrokeColorField(wrapper);

    const newStroke = Faker.internet.color();

    strokeColorField.vm.$emit('input', newStroke);

    expect(wrapper).toEmit('input', {
      first: {
        ...firstShape,
        properties: {
          stroke: newStroke,
        },
      },
      second: {
        ...secondShape,
        properties: {
          stroke: newStroke,
        },
      },
    });
  });

  test('Stroke width changed after trigger number field', () => {
    const firstShape = rectShapeToForm({
      _id: 'first',
      properties: { stroke: 'red' },
    });
    const secondShape = rectShapeToForm({
      _id: 'second',
      properties: { stroke: 'red' },
    });
    const shapes = {
      first: firstShape,
      second: secondShape,
    };
    const wrapper = factory({
      propsData: {
        shapes,
        selected: Object.keys(shapes),
      },
    });

    const strokeWidthField = selectStrokeWidthField(wrapper);

    const newStrokeWidth = Faker.datatype.number();

    strokeWidthField.vm.$emit('input', newStrokeWidth);

    expect(wrapper).toEmit('input', {
      first: {
        ...firstShape,
        properties: {
          ...firstShape.properties,
          'stroke-width': newStrokeWidth,
        },
      },
      second: {
        ...secondShape,
        properties: {
          ...secondShape.properties,
          'stroke-width': newStrokeWidth,
        },
      },
    });
  });

  test('Stroke type changed after trigger stroke type field', () => {
    const firstShape = rectShapeToForm({
      _id: 'first',
      properties: { stroke: 'red' },
    });
    const secondShape = rectShapeToForm({
      _id: 'second',
      properties: { stroke: 'red' },
    });
    const shapes = {
      first: firstShape,
      second: secondShape,
    };
    const wrapper = factory({
      propsData: {
        shapes,
        selected: Object.keys(shapes),
      },
    });

    const strokeTypeField = selectStrokeTypeField(wrapper);

    strokeTypeField.vm.$emit('input', STROKE_TYPES.dashed);

    expect(wrapper).toEmit('input', {
      first: {
        ...firstShape,
        properties: {
          ...firstShape.properties,
          'stroke-dasharray': '4 4',
        },
      },
      second: {
        ...secondShape,
        properties: {
          ...secondShape.properties,
          'stroke-dasharray': '4 4',
        },
      },
    });
  });

  test('Stroke type changed after trigger stroke type field', () => {
    const firstShape = lineShapeToForm({
      _id: 'first',
    });
    const secondShape = lineShapeToForm({
      _id: 'second',
    });
    const shapes = {
      first: firstShape,
      second: secondShape,
    };
    const wrapper = factory({
      propsData: {
        shapes,
        selected: Object.keys(shapes),
      },
    });

    const lineTypeField = selectLineTypeField(wrapper);

    lineTypeField.vm.$emit('input', LINE_TYPES.rightElbow);

    expect(wrapper).toEmit('input', {
      first: {
        ...firstShape,
        lineType: LINE_TYPES.rightElbow,
      },
      second: {
        ...secondShape,
        lineType: LINE_TYPES.rightElbow,
      },
    });
  });

  test('Text color changed after trigger color field', () => {
    const firstShape = rectShapeToForm({
      _id: 'first',
    });
    const secondShape = rectShapeToForm({
      _id: 'second',
    });
    const shapes = {
      first: firstShape,
      second: secondShape,
    };
    const wrapper = factory({
      propsData: {
        shapes,
        selected: Object.keys(shapes),
      },
    });

    const fontColorField = selectFontColorField(wrapper);

    const newFontColor = Faker.internet.color();

    fontColorField.vm.$emit('input', newFontColor);

    expect(wrapper).toEmit('input', {
      first: {
        ...firstShape,
        textProperties: {
          ...firstShape.textProperties,
          color: newFontColor,
        },
      },
      second: {
        ...secondShape,
        textProperties: {
          ...secondShape.textProperties,
          color: newFontColor,
        },
      },
    });
  });

  test('Text background color changed after trigger color field', () => {
    const firstShape = rectShapeToForm({
      _id: 'first',
    });
    const secondShape = rectShapeToForm({
      _id: 'second',
    });
    const shapes = {
      first: firstShape,
      second: secondShape,
    };
    const wrapper = factory({
      propsData: {
        shapes,
        selected: Object.keys(shapes),
      },
    });

    const fontBackgroundColorField = selectFontBackgroundColorField(wrapper);

    const newFontBackgroundColor = Faker.internet.color();

    fontBackgroundColorField.vm.$emit('input', newFontBackgroundColor);

    expect(wrapper).toEmit('input', {
      first: {
        ...firstShape,
        textProperties: {
          ...firstShape.textProperties,
          backgroundColor: newFontBackgroundColor,
        },
      },
      second: {
        ...secondShape,
        textProperties: {
          ...secondShape.textProperties,
          backgroundColor: newFontBackgroundColor,
        },
      },
    });
  });

  test('text size changed after trigger number field', () => {
    const firstShape = rectShapeToForm({
      _id: 'first',
    });
    const secondShape = rectShapeToForm({
      _id: 'second',
    });
    const shapes = {
      first: firstShape,
      second: secondShape,
    };
    const wrapper = factory({
      propsData: {
        shapes,
        selected: Object.keys(shapes),
      },
    });

    const fontSizeField = selectFontSizeField(wrapper);

    const newFontSize = Faker.datatype.number();

    fontSizeField.vm.$emit('input', newFontSize);

    expect(wrapper).toEmit('input', {
      first: {
        ...firstShape,
        textProperties: {
          ...firstShape.textProperties,
          fontSize: newFontSize,
        },
      },
      second: {
        ...secondShape,
        textProperties: {
          ...secondShape.textProperties,
          fontSize: newFontSize,
        },
      },
    });
  });

  test('Renders `flowchart-properties` with all properties', async () => {
    const fill = 'red';
    const stroke = 'orange';
    const shapes = {
      first: rectShapeToForm({
        _id: 'first',
        properties: { fill, stroke },
      }),
      second: rectShapeToForm({
        _id: 'second',
        properties: { fill, stroke },
      }),
    };
    const wrapper = snapshotFactory({
      propsData: {
        shapes,
        selected: Object.keys(shapes),
      },
    });

    await wrapper.openAllExpansionPanels();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `flowchart-properties` with lines', async () => {
    const shapes = {
      first: lineShapeToForm({
        _id: 'first',
      }),
      second: lineShapeToForm({
        _id: 'second',
      }),
    };
    const wrapper = snapshotFactory({
      propsData: {
        shapes,
        selected: Object.keys(shapes),
      },
    });

    await wrapper.openAllExpansionPanels();

    expect(wrapper).toMatchSnapshot();
  });
});
