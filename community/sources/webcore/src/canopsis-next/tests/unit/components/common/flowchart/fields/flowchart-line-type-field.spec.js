import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';
import { LINE_TYPES } from '@/constants';

import FlowchartLineTypeField from '@/components/common/flowchart/fields/flowchart-line-type-field.vue';

const stubs = {
  'points-line-path': true,
  'v-select': createSelectInputStub('v-select'),
};

const snapshotStubs = {
  'points-line-path': true,
};

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('flowchart-line-type-field', () => {
  const factory = generateShallowRenderer(FlowchartLineTypeField, { stubs });
  const snapshotFactory = generateRenderer(FlowchartLineTypeField, { stubs: snapshotStubs });

  test('Value changed after trigger select field', () => {
    const wrapper = factory();

    const selectField = selectSelectField(wrapper);

    selectField.vm.$emit('input', LINE_TYPES.verticalCurve);

    expect(wrapper).toEmit('input', LINE_TYPES.verticalCurve);
  });

  test('Renders `flowchart-line-type-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `flowchart-line-type-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: LINE_TYPES.line,
        label: 'Custom label',
        averagePoints: [{ x: 10, y: 10 }, { x: 1, y: 1 }],
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `flowchart-line-type-field` with default averagePoints', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: LINE_TYPES.line,
        label: 'Custom label',
        averagePoints: [{ x: 1, y: 1 }, { x: 10, y: 10 }],
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
