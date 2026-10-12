import { mockDateNow } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

const stubs = {
  'c-text-field': {
    props: ['modelValue', 'validationRules'],
    emits: ['update:modelValue', 'focus', 'blur'],
    template: `
      <input
        class="v-text-field"
        :value="modelValue"
        @input="$emit('update:modelValue', $event.target.value)"
        @focus="$emit('focus')"
        @blur="$emit('blur')"
      />
    `,
  },
  'date-time-picker-menu': {
    props: ['value'],
    emits: ['update:value'],
    template: '<button class="date-time-picker-menu" type="button" />',
  },
};

const selectTextField = wrapper => wrapper.find('.v-text-field');
const selectDateTimePickerMenu = wrapper => wrapper.find('.date-time-picker-menu');

describe('date-time-picker-text-field', () => {
  const nowTimestamp = 1386435600000;
  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(DateTimePickerTextField, { stubs });

  test('String value changed after trigger text field', async () => {
    const wrapper = factory();

    const newValue = 'now-1d';

    await selectTextField(wrapper).triggerCustomEvent('update:modelValue', newValue);

    expect(wrapper).toEmitInput(newValue);
  });

  test('Date object not updated if text field is focused', async () => {
    const value = '11/12/2021 21:00';
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const originalTime = wrapper.vm.objectValue.getTime();

    await selectTextField(wrapper).triggerCustomEvent('focus');

    await wrapper.setProps({
      value: '11/12/2022 21:00',
    });

    expect(wrapper.vm.objectValue.getTime()).toBe(originalTime);
  });

  test('Date object updated if text field is blurred', async () => {
    const value = '11/12/2021 21:00';
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const originalTime = wrapper.vm.objectValue.getTime();

    await selectTextField(wrapper).triggerCustomEvent('focus');

    await wrapper.setProps({
      value: '11/12/2022 21:00',
    });

    await selectTextField(wrapper).triggerCustomEvent('blur');
    await flushPromises();

    expect(wrapper.vm.objectValue.getTime()).not.toBe(originalTime);
  });

  test('Errors added after validate date time text field', () => {
    const wrapper = factory({
      propsData: {
        value: '11/12/2021 21:00',
        name: 'name',
        rules: {
          after: '11/12/2022 21:00',
        },
      },
    });

    const results = wrapper.vm.validationRules.map(rule => rule(wrapper.vm.value));

    expect(results).toContain('Date must be after start date');
  });

  test('Error added after validate date time text field with failed preparer', () => {
    const dateObjectPreparer = vi.fn().mockImplementation(() => {
      throw Error();
    });

    const wrapper = factory({
      propsData: {
        value: '11/12/2021 21:00',
        name: 'name',
        dateObjectPreparer,
        rules: {
          date_format: true,
        },
      },
    });

    const results = wrapper.vm.validationRules.map(rule => rule(wrapper.vm.value));

    expect(results).toContain('Invalid date format');
  });

  test('Object field updated after trigger date time picker button', async () => {
    const wrapper = factory();

    await selectDateTimePickerMenu(wrapper).triggerCustomEvent('update:value', nowTimestamp);

    expect(wrapper).toEmitInput(convertDateToString(nowTimestamp, DATETIME_FORMATS.dateTimePicker));
  });
});
