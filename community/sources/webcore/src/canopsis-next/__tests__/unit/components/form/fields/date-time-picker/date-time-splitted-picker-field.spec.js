import { generateShallowRenderer } from '@unit/utils/vue';

import DateTimeSplittedPickerField from '@/components/forms/fields/date-time-picker/date-time-splitted-picker-field.vue';

const stubs = {
  'v-input': {
    inheritAttrs: false,
    template: '<div class="v-input"><slot /></div>',
  },
  'time-picker-field': true,
  'c-date-picker-field': true,
};

const selectTimePickerField = wrapper => wrapper.find('time-picker-field-stub');
const selectDatePickerField = wrapper => wrapper.find('c-date-picker-field-stub');

describe('date-time-splitted-picker-field', () => {
  const factory = generateShallowRenderer(DateTimeSplittedPickerField, { stubs });

  test('Date changed after trigger date field', () => {
    const value = new Date(1998, 10, 1, 15, 55);
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    selectDatePickerField(wrapper).triggerCustomEvent('update:value', '2022-03-11');

    const expected = new Date(value);
    expected.setFullYear(2022, 2, 11);
    expected.setSeconds(0, 0);

    expect(wrapper).toEmitInput(expected);
  });

  test('Time changed after trigger time field', () => {
    const value = new Date(1998, 10, 1, 15, 55);
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    selectTimePickerField(wrapper).triggerCustomEvent('update:value', '12:15');

    const expected = new Date(value);
    expected.setHours(12, 15, 0, 0);

    expect(wrapper).toEmitInput(expected);
  });
});
