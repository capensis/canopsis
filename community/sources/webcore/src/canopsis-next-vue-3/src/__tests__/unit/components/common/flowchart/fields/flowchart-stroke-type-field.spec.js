import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { STROKE_TYPES } from '@/constants';

import FlowchartStrokeTypeField from '@/components/common/flowchart/fields/flowchart-stroke-type-field.vue';

const stubs = {
  'v-select': {
    template: '<button class="v-select" @click="$emit(\'update:modelValue\', { value: selected })" />',
    props: ['selected'],
  },
};

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('flowchart-stroke-type-field', () => {
  const factory = generateShallowRenderer(FlowchartStrokeTypeField, { stubs });

  test('Value changed after trigger select field', () => {
    const wrapper = factory();

    const selectField = selectSelectField(wrapper);

    selectField.triggerCustomEvent('update:modelValue', { value: STROKE_TYPES.dashed });

    expect(wrapper.emitted('update:value')).toEqual([[STROKE_TYPES.dashed]]);
  });
});
