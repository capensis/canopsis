import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import DateTimeSplittedRangePickerField from '@/components/forms/fields/date-time-splitted-range-picker-field.vue';

const stubs = {
  'date-time-splitted-picker-field': true,
  'c-date-picker-field': true,
};

const selectDateTimeSplittedPickerFieldsByIndex = (wrapper, index) => wrapper
  .findAll('date-time-splitted-picker-field-stub')
  .at(index);
const selectStartDateTimeSplittedPickerFields = wrapper => selectDateTimeSplittedPickerFieldsByIndex(wrapper, 0);
const selectEndDateTimeSplittedPickerFields = wrapper => selectDateTimeSplittedPickerFieldsByIndex(wrapper, 1);

describe('date-time-splitted-range-picker-field', () => {
  const factory = generateShallowRenderer(DateTimeSplittedRangePickerField, {

    stubs,

    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  const snapshotFactory = generateRenderer(DateTimeSplittedRangePickerField, {

    stubs,

    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  test('Start changed after trigger date field', () => {
    const wrapper = factory({
      propsData: {
        value: new Date(1998, 10, 1, 15, 55),
        endRules: {},
        startRules: {},
        startLabel: ' ',
        endLabel: ' ',
      },
    });

    const newStart = new Date(1998, 10, 3, 15);

    selectStartDateTimeSplittedPickerFields(wrapper).vm.$emit('input', newStart);

    expect(wrapper).toEmit('update:start', newStart);
  });

  test('End changed after trigger time field', () => {
    const wrapper = factory({
      propsData: {
        end: new Date(1998, 10, 1, 15, 55),
        endRules: {},
        startRules: {},
        startLabel: ' ',
        endLabel: ' ',
      },
    });

    const newEnd = new Date(1998, 10, 3, 15);

    selectEndDateTimeSplittedPickerFields(wrapper).vm.$emit('input', newEnd);

    expect(wrapper).toEmit('update:end', newEnd);
  });

  test('Renders `date-time-splitted-range-picker-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        endRules: { required: true },
        startRules: { required: true },
        startLabel: 'Required start label',
        endLabel: 'Required end label',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `date-time-splitted-range-picker-field` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        start: new Date(2012, 2, 3),
        end: new Date(2012, 2, 4),
        endRules: { required: true },
        startRules: { required: true },
        startLabel: 'Custom start label',
        endLabel: 'Custom end label',
        name: 'custom_name',
        fullDay: true,
        disabled: true,
        endMin: '1970-01-02',
        endMax: '1970-01-03',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
