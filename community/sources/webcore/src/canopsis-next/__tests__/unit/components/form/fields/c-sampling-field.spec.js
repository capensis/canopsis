import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { SAMPLINGS } from '@/constants';

import CSamplingField from '@/components/forms/fields/c-sampling-field.vue';

const stubs = {
  'c-select': createSelectInputStub('c-select'),
};

const selectSelectField = wrapper => wrapper.find('select.c-select');

describe('c-sampling-field', () => {
  const factory = generateShallowRenderer(CSamplingField, { stubs });

  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        value: SAMPLINGS.day,
      },
    });

    selectSelectField(wrapper).triggerCustomEvent('update:modelValue', SAMPLINGS.hour);

    expect(wrapper).toEmitInput(SAMPLINGS.hour);
  });
});
