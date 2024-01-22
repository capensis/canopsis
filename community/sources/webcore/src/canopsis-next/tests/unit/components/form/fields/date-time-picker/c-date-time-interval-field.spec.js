import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CDateTimeIntervalField from '@/components/forms/fields/date-time-picker/c-date-time-interval-field.vue';

const stubs = {
  'date-time-picker-field': true,
};

const selectDateTimePickerField = wrapper => wrapper.findAll('date-time-picker-field-stub');
const selectFromDateTimePickerField = wrapper => selectDateTimePickerField(wrapper)
  .at(0);
const selectToDateTimePickerField = wrapper => selectDateTimePickerField(wrapper)
  .at(1);

describe('c-date-time-interval-field', () => {
  const timestamp = 1386435600000;

  const factory = generateShallowRenderer(CDateTimeIntervalField, { stubs });
  const snapshotFactory = generateRenderer(CDateTimeIntervalField, { stubs });

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
    expect(wrapper).toEmit('input', { ...value, from: newValue });
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
    expect(wrapper).toEmit('input', { ...value, to: newValue });
  });

  test('Renders `c-date-time-interval-field` with default props', () => {
    const dateObject = new Date(timestamp);
    const dateSpy = jest.spyOn(global, 'Date')
      .mockReturnValue(dateObject);
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();

    dateSpy.mockClear();
  });

  test('Renders `c-date-time-interval-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          from: timestamp,
          to: timestamp - 1000,
        },
        name: 'custom_name',
        disabled: true,
        hideDetails: true,
        required: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
