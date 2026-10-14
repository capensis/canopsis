import { createInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import CDatePickerField from '@/components/forms/fields/date-picker/c-date-picker-field.vue';

const stubs = {
  'v-menu': {
    inheritAttrs: false,
    template: '<div class="v-menu"><slot name="activator" :props="{}" /><slot /></div>',
  },
  'v-date-picker': createInputStub('v-date-picker'),
  'c-text-field': {
    props: ['modelValue', 'value'],
    emits: ['click:appendInner', 'update:modelValue'],
    template: '<input class="c-text-field" :value="modelValue ?? value" />',
  },
};

const selectDatePicker = wrapper => wrapper.find('.v-date-picker');
const selectTextField = wrapper => wrapper.find('.c-text-field');

describe('c-date-picker-field', () => {
  const factory = generateShallowRenderer(CDatePickerField, { stubs });

  test('Value changed after trigger date picker', () => {
    const wrapper = factory({
      propsData: {
        value: '2022-10-03',
      },
    });

    const newValue = '2022-10-04';

    selectDatePicker(wrapper).triggerCustomEvent('update:modelValue', newValue);

    expect(wrapper).toEmitInput(newValue);
  });

  test('Value cleared after trigger text field', () => {
    const wrapper = factory({
      propsData: {
        value: '2022-10-03',
        clearable: true,
      },
    });

    selectTextField(wrapper).triggerCustomEvent('click:appendInner');

    expect(wrapper).toEmitInput(null);
  });

  test('Change event emitted after trigger date picker', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const newValue = '2022-10-04';

    selectDatePicker(wrapper).triggerCustomEvent('update:modelValue', newValue);

    expect(wrapper).toEmit('change', newValue);
  });
});
