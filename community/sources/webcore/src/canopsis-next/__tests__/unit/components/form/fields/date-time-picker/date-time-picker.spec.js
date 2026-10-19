import { createButtonStub } from '@unit/stubs/button';
import { mockDateNow } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';

import { DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import DateTimePicker from '@/components/forms/fields/date-time-picker/date-time-picker.vue';

const stubs = {
  'v-sheet': {
    inheritAttrs: false,
    template: '<div class="v-sheet"><slot /></div>',
  },
  'v-divider': true,
  'time-picker-field': {
    name: 'TimePickerFieldStub',
    props: ['value', 'roundHours', 'allowedTime'],
    emits: ['update:value'],
    template: '<input class="time-picker-field" />',
  },
  'v-date-picker': {
    name: 'VDatePickerStub',
    props: ['modelValue', 'allowedDates'],
    emits: ['update:modelValue'],
    template: '<div class="v-date-picker" />',
  },
  'v-btn': createButtonStub('v-btn'),
};

const selectButtons = wrapper => wrapper.findAll('.v-btn');
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectApplyButton = wrapper => selectButtons(wrapper).at(1);
const selectTimePickerField = wrapper => wrapper.findComponent({ name: 'TimePickerFieldStub' });
const selectDatePicker = wrapper => wrapper.findComponent({ name: 'VDatePickerStub' });

describe('date-time-picker', () => {
  const nowTimestamp = 1386435600000;
  const nowDate = new Date(nowTimestamp);

  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(DateTimePicker, { stubs });

  test('Time updated after trigger time picker field', async () => {
    const wrapper = factory({
      propsData: {
        value: nowDate,
      },
    });

    const newTime = '12:45';

    await selectTimePickerField(wrapper).triggerCustomEvent('update:value', newTime);

    expect(wrapper.vm.timeString).toBe(newTime);
  });

  test('Date updated after trigger date picker', async () => {
    const wrapper = factory({
      propsData: {
        value: nowDate,
      },
    });

    const newDate = '2015-12-12';

    await selectDatePicker(wrapper).triggerCustomEvent('update:modelValue', newDate);

    expect(convertDateToString(wrapper.vm.localValue, DATETIME_FORMATS.datePicker)).toBe(newDate);
  });

  test('Value changed after trigger apply button without changes', async () => {
    const wrapper = factory({
      propsData: {
        value: nowTimestamp,
      },
    });

    await selectApplyButton(wrapper).trigger('click');

    expect(wrapper).toHaveBeenEmit('close');
    expect(wrapper).toEmitInput(new Date(nowTimestamp));
  });

  test('Value changed after trigger apply button with changes', async () => {
    const wrapper = factory({
      propsData: {
        value: null,
      },
    });

    const newTime = '12:45';
    const newDate = '2015-12-12';

    await selectTimePickerField(wrapper).triggerCustomEvent('update:value', newTime);
    await selectDatePicker(wrapper).triggerCustomEvent('update:modelValue', newDate);

    const expected = new Date(wrapper.vm.localValue);

    await selectApplyButton(wrapper).trigger('click');

    expect(wrapper).toEmitInput(expected);
    expect(wrapper).toHaveBeenEmit('close');
  });

  test('Close emitted after trigger cancel button', async () => {
    const wrapper = factory({
      propsData: {
        value: null,
      },
    });

    await selectCancelButton(wrapper).trigger('click');

    expect(wrapper).toHaveBeenEmit('close');
  });
});
