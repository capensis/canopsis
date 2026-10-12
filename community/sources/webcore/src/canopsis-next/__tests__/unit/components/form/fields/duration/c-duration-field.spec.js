import { createNumberInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { TIME_UNITS } from '@/constants';

import CDurationField from '@/components/forms/fields/compat/c-duration-field.vue';

const stubs = {
  'c-number-field': createNumberInputStub('c-number-field'),
};

const selectNumberField = wrapper => wrapper.find('input.c-number-field');
const selectSelectField = wrapper => wrapper.find('select.c-select');

describe('c-duration-field', () => {
  const factory = generateShallowRenderer(CDurationField, {
    stubs,
  });

  it('Value changed after trigger text field', () => {
    const duration = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.week,
    };
    const wrapper = factory({
      propsData: {
        duration,
      },
    });
    const newValue = Faker.datatype.number();

    selectNumberField(wrapper).triggerCustomEvent('update:value', newValue);

    expect(wrapper).toEmitInput({
      unit: duration.unit,
      value: newValue,
    });
  });

  it('Unit changed after trigger select field', () => {
    const duration = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.week,
    };
    const wrapper = factory({
      propsData: {
        duration,
        long: true,
      },
    });

    selectSelectField(wrapper).triggerCustomEvent('update:modelValue', TIME_UNITS.month);

    expect(wrapper).toEmitInput({
      unit: TIME_UNITS.month,
      value: duration.value,
    });
  });

  it('Value cleared after trigger select field without value', () => {
    const duration = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.week,
    };
    const wrapper = factory({
      propsData: {
        duration,
        long: true,
      },
    });

    selectSelectField(wrapper).triggerCustomEvent('update:modelValue');

    expect(wrapper).toEmitInput({
      unit: undefined,
      value: undefined,
    });
  });
});
