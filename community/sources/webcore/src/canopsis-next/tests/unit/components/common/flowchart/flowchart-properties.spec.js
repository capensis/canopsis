import { mount, createVueInstance } from '@unit/utils/vue';
import { lineShapeToForm, rectShapeToForm } from '@/helpers/flowchart/shapes';

import FlowchartProperties from '@/components/common/flowchart/flowchart-properties.vue';

const localVue = createVueInstance();

const stubs = {
  'flowchart-color-field': true,
  'flowchart-number-field': true,
  'flowchart-stroke-type-field': true,
  'flowchart-line-type-field': true,
};

const snapshotFactory = (options = {}) => mount(FlowchartProperties, {
  localVue,
  stubs,

  ...options,
});

describe('flowchart-properties', () => {
  test('Renders `flowchart-properties` with all properties', () => {
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

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `flowchart-properties` with lines', () => {
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
