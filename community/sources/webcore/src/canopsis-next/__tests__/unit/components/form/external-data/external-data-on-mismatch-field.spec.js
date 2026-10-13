import { generateShallowRenderer } from '@unit/utils/vue';

import { EXTERNAL_DATA_ON_MISMATCH_TYPES } from '@/constants';

import ExternalDataOnMismatchField from '@/components/forms/external-data/external-data-on-mismatch-field.vue';

const stubs = {
  'c-radio-group': {
    props: ['modelValue'],
    emits: ['update:modelValue'],
    template: '<div class="c-radio-group"><slot /></div>',
  },
  'v-radio': true,
};

const selectRadioGroup = wrapper => wrapper.findComponent('.c-radio-group');

describe('external-data-on-mismatch-field', () => {
  const factory = generateShallowRenderer(ExternalDataOnMismatchField, { stubs });

  it('renders every supported mismatch behavior', () => {
    const wrapper = factory();

    expect(wrapper.findAll('v-radio-stub')).toHaveLength(
      Object.values(EXTERNAL_DATA_ON_MISMATCH_TYPES).length,
    );
  });

  it('emits the selected mismatch behavior', () => {
    const wrapper = factory();

    selectRadioGroup(wrapper).triggerCustomEvent(
      'update:modelValue',
      EXTERNAL_DATA_ON_MISMATCH_TYPES.skipRule,
    );

    expect(wrapper.emitted('update:modelValue')[0][0]).toBe(
      EXTERNAL_DATA_ON_MISMATCH_TYPES.skipRule,
    );
  });
});
