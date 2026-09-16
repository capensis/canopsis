import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { STROKE_TYPES } from '@/constants';

import FlowchartStrokeTypeField from '@/components/common/flowchart/fields/flowchart-stroke-type-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('flowchart-stroke-type-field', () => {
  const factory = generateShallowRenderer(FlowchartStrokeTypeField, { stubs });

  test('Value changed after trigger select field', () => {
    const wrapper = factory();

    const selectField = selectSelectField(wrapper);

    selectField.triggerCustomEvent('input', STROKE_TYPES.dashed);

    expect(wrapper).toEmitInput(STROKE_TYPES.dashed);
  });
});
