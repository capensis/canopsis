import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import DateTimeSplittedPickerField from '@/components/forms/fields/date-time-picker/date-time-splitted-picker-field.vue';

const stubs = {
  'time-picker-field': true,
  'c-date-picker-field': true,
};

const selectTimePickerField = wrapper => wrapper.find('time-picker-field-stub');
const selectDatePickerField = wrapper => wrapper.find('c-date-picker-field-stub');

describe('date-time-splitted-picker-field', () => {
  const factory = generateShallowRenderer(DateTimeSplittedPickerField, {

    stubs,

    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  const snapshotFactory = generateRenderer(DateTimeSplittedPickerField, {

    stubs,

    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  test('Date changed after trigger date field', () => {
    const wrapper = factory({
      propsData: {
        value: new Date(1998, 10, 1, 15, 55),
      },
    });

    selectDatePickerField(wrapper).vm.$emit('input', '2022-03-11');

    expect(wrapper).toEmit('input', new Date('2022-03-11T14:55:00.000Z'));
  });

  test('Time changed after trigger time field', () => {
    const wrapper = factory({
      propsData: {
        value: new Date(1998, 10, 1, 15, 55),
      },
    });

    selectTimePickerField(wrapper).vm.$emit('input', '12:15');

    expect(wrapper).toEmit('input', new Date('1998-11-01T11:15:00.000Z'));
  });

  test('Renders `date-time-splitted-picker-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `date-time-splitted-picker-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: new Date(123321123),
        fullDay: true,
        label: 'Custom label',
        name: 'custom_name',
        reverse: true,
        disabled: true,
        min: '1970-01-02',
        max: '1970-01-03',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `date-time-splitted-picker-field` with errors', async () => {
    const name = 'validate_name';
    const wrapper = snapshotFactory({
      propsData: {
        name,
      },
    });

    const validator = wrapper.getValidator();

    validator.errors.add({
      field: name,
      msg: 'Error message',
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
