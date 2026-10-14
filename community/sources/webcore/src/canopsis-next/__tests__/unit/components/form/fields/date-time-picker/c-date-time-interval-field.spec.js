import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CDateTimeIntervalField from '@/components/forms/fields/date-time-picker/c-date-time-interval-field.vue';

const stubs = {
  'date-time-picker-field': true,
};

const selectDateTimePickerField = wrapper => wrapper.findAll('date-time-picker-field-stub');
const selectFromDateTimePickerField = wrapper => selectDateTimePickerField(wrapper).at(0);
const selectToDateTimePickerField = wrapper => selectDateTimePickerField(wrapper).at(1);

describe('c-date-time-interval-field', () => {
  const timestamp = 1386435600000;

  const factory = generateShallowRenderer(CDateTimeIntervalField, { stubs });

  test('From changed after trigger from date time picker field', () => {
    const value = {
      from: 0,
      to: timestamp,
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const fromField = selectFromDateTimePickerField(wrapper);

    const newValue = Faker.datatype.number({
      max: timestamp,
    });

    fromField.triggerCustomEvent('input', newValue);
    expect(wrapper).toEmitInput({ ...value, from: newValue });
  });

  test('To changed after trigger to date time picker field', () => {
    const value = {
      from: 0,
      to: 0,
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const toField = selectToDateTimePickerField(wrapper);

    const newValue = Faker.datatype.number({
      max: timestamp,
    });

    toField.triggerCustomEvent('input', newValue);
    expect(wrapper).toEmitInput({ ...value, to: newValue });
  });
});
