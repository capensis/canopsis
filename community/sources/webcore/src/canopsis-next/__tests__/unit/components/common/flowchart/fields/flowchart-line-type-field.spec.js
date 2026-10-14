import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { LINE_TYPES } from '@/constants';

import FlowchartLineTypeField from '@/components/common/flowchart/fields/flowchart-line-type-field.vue';

const stubs = {
  'points-line-path': true,
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('flowchart-line-type-field', () => {
  const factory = generateShallowRenderer(FlowchartLineTypeField, { stubs });

  test('Value changed after trigger select field', () => {
    const wrapper = factory();

    const selectField = selectSelectField(wrapper);

    selectField.triggerCustomEvent('input', LINE_TYPES.verticalCurve);

    expect(wrapper).toEmitInput(LINE_TYPES.verticalCurve);
  });
});
