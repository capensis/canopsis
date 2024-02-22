import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import { DATETIME_FORMATS } from '@/constants';

import DateTimePicker from '@/components/forms/fields/date-time-picker/date-time-picker.vue';

const stubs = {
  'time-picker-field': true,
};

const selectButtons = wrapper => wrapper.findAll('v-btn-stub');
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectApplyButton = wrapper => selectButtons(wrapper).at(1);
const selectTimePickerField = wrapper => wrapper.find('time-picker-field-stub');
const selectDatePicker = wrapper => wrapper.find('v-date-picker-stub');

describe('date-time-picker', () => {
  const nowTimestamp = 1386435600000;
  const nowDate = new Date(nowTimestamp);

  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(DateTimePicker, {
    stubs,
    listeners: {
      close: jest.fn(),
    },
  });
  const snapshotFactory = generateRenderer(DateTimePicker, {
    stubs,
    listeners: {
      close: jest.fn(),
    },
  });

  test('Time updated after trigger time picker field', async () => {
    const wrapper = factory({
      propsData: {
        value: nowDate,
      },
    });

    const timePickerField = selectTimePickerField(wrapper);

    const newTime = '12:45';

    await timePickerField.triggerCustomEvent('input', newTime);

    expect(timePickerField.attributes('value')).toBe(newTime);
  });

  test('Date updated after trigger date picker', async () => {
    const wrapper = factory({
      propsData: {
        value: nowDate,
      },
    });

    const datePicker = selectDatePicker(wrapper);

    const newDate = '2015-12-12';

    await datePicker.triggerCustomEvent('input', newDate);

    expect(datePicker.attributes('value')).toBe(newDate);
  });

  test('Value changed after trigger apply button without changes', () => {
    const wrapper = factory({
      propsData: {
        value: nowTimestamp,
      },
    });

    selectApplyButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toHaveBeenEmit('close');
    expect(wrapper).toEmitInput(new Date(nowTimestamp));
  });

  test('Value changed after trigger apply button with changes', async () => {
    const wrapper = factory({
      propsData: {
        value: null,
      },
    });
    const timePickerField = selectTimePickerField(wrapper);
    const datePicker = selectDatePicker(wrapper);

    const newTime = '12:45';
    const newDate = '2015-12-12';
    const resultDate = new Date(`${newDate} ${newTime}`);

    await timePickerField.triggerCustomEvent('input', newTime);
    await datePicker.triggerCustomEvent('input', newDate);

    selectApplyButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toEmitInput(resultDate);
  });

  test('Close emitted after trigger cancel button', () => {
    const close = jest.fn();
    const wrapper = factory({
      propsData: {
        value: null,
      },
      listeners: {
        close,
      },
    });

    selectCancelButton(wrapper).triggerCustomEvent('click');

    expect(close).toBeCalled();
  });

  test('Renders `date-time-picker` with default props', () => {
    const dateObject = new Date(nowTimestamp);
    const dateSpy = jest.spyOn(global, 'Date')
      .mockReturnValue(dateObject);

    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();

    dateSpy.mockClear();
  });

  test('Renders `date-time-picker` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: new Date(1387435600000),
        label: 'label',
        roundHours: true,
        dateFormat: DATETIME_FORMATS.long,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `date-time-picker` without value', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: null,
      },
      listeners: {
        close: jest.fn(),
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
