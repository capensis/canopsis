import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { LINE_TYPES, STROKE_TYPES } from '@/constants';

import { lineShapeToForm, rectShapeToForm } from '@/helpers/flowchart/shapes';

import FlowchartProperties from '@/components/common/flowchart/flowchart-properties.vue';

const stubs = {
  'flowchart-color-field': true,
  'flowchart-number-field': true,
  'flowchart-stroke-type-field': true,
  'flowchart-line-type-field': true,
};

describe('flowchart-properties', () => {
  const factory = generateShallowRenderer(FlowchartProperties, { stubs });

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

    const newFill = Faker.internet.color();

    wrapper.vm.updateFill(newFill);

    expect(wrapper).toEmitInput({
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

    const newStroke = Faker.internet.color();

    wrapper.vm.updateStroke(newStroke);

    expect(wrapper).toEmitInput({
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

    const newStrokeWidth = Faker.datatype.number();

    wrapper.vm.updateStrokeWidth(newStrokeWidth);

    expect(wrapper).toEmitInput({
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

    wrapper.vm.updateStrokeType(STROKE_TYPES.dashed);

    expect(wrapper).toEmitInput({
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

    wrapper.vm.updateLineType(LINE_TYPES.rightElbow);

    expect(wrapper).toEmitInput({
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

    const newFontColor = Faker.internet.color();

    wrapper.vm.updateTextColor(newFontColor);

    expect(wrapper).toEmitInput({
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

    const newFontBackgroundColor = Faker.internet.color();

    wrapper.vm.updateTextBackgroundColor(newFontBackgroundColor);

    expect(wrapper).toEmitInput({
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

    const newFontSize = Faker.datatype.number();

    wrapper.vm.updateFontSize(newFontSize);

    expect(wrapper).toEmitInput({
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
});
