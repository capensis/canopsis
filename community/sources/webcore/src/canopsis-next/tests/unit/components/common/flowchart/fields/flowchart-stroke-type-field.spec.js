import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';
import { LINE_TYPES, STROKE_TYPES } from '@/constants';

import FlowchartStrokeTypeField from '@/components/common/flowchart/fields/flowchart-stroke-type-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('flowchart-stroke-type-field', () => {
  const factory = generateShallowRenderer(FlowchartStrokeTypeField, { stubs });
  const snapshotFactory = generateRenderer(FlowchartStrokeTypeField);

  test('Value changed after trigger select field', () => {
    const wrapper = factory();

    const selectField = selectSelectField(wrapper);

    selectField.vm.$emit('input', STROKE_TYPES.dashed);

    expect(wrapper).toEmit('input', STROKE_TYPES.dashed);
  });

  test('Renders `flowchart-stroke-type-field` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `flowchart-stroke-type-field` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: LINE_TYPES.line,
        label: 'Custom label',
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
