import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { mockDateNow } from '@unit/utils/mock-hooks';
import { createInputStub } from '@unit/stubs/input';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

const stubs = {
  'v-text-field': createInputStub('v-text-field'),
  'date-time-picker-menu': true,
};

const snapshotStubs = {
  'date-time-picker-menu': true,
};

const selectTextField = wrapper => wrapper.find('.v-text-field');
const selectDateTimePickerButton = wrapper => wrapper.find('date-time-picker-menu-stub');

describe('date-time-picker-text-field', () => {
  const nowTimestamp = 1386435600000;
  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(DateTimePickerTextField, { stubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });
  const snapshotFactory = generateRenderer(DateTimePickerTextField, {
    stubs: snapshotStubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  test('String value changed after trigger text field', () => {
    const wrapper = factory();

    const textField = selectTextField(wrapper);

    const newValue = 'now-1d';

    textField.vm.$emit('input', newValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(newValue);
  });

  test('Date object not updated if text field is focused', async () => {
    const value = '11/12/2021 21:00';
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const textField = selectTextField(wrapper);
    const dateTimePickerButton = selectDateTimePickerButton(wrapper);

    textField.vm.$emit('focus');

    await wrapper.setProps({
      value: '11/12/2022 21:00',
    });

    expect(dateTimePickerButton.vm.value.getTime()).toBe(1639252800000);
  });

  test('Date object updated if text field is blurred', async () => {
    const value = '11/12/2021 21:00';
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const textField = selectTextField(wrapper);
    const dateTimePickerButton = selectDateTimePickerButton(wrapper);

    textField.vm.$emit('focus');

    await wrapper.setProps({
      value: '11/12/2022 21:00',
    });

    textField.vm.$emit('blur');

    await flushPromises();

    expect(dateTimePickerButton.vm.value.getTime()).toBe(1670788800000);
  });

  test('Errors added after validate date time text field', async () => {
    const dateObjectPreparer = jest.fn().mockReturnValue(null);

    const wrapper = factory({
      propsData: {
        value: 'now',
        name: 'name',
        dateObjectPreparer,
      },
    });

    const validator = wrapper.getValidator();

    const isFormValid = await validator.validateAll();

    expect(isFormValid).toBeFalsy();

    expect(validator.errors.items.map(({ msg }) => msg)).toEqual([
      'End date should be after start date',
    ]);

    wrapper.destroy();
  });

  test('Error added after validate date time text field with failed preparer', async () => {
    const dateObjectPreparer = jest.fn()
      .mockImplementation(() => { throw Error(); });

    const wrapper = factory({
      propsData: {
        value: ' ',
        name: 'name',
        dateObjectPreparer,
      },
    });

    await flushPromises();

    const validator = wrapper.getValidator();

    const isFormValid = await validator.validateAll();

    expect(isFormValid).toBeFalsy();

    expect(validator.errors.items).toHaveLength(1);
  });

  test('Object field updated after trigger date time picker button', () => {
    const wrapper = factory();

    const dateTimePickerButton = selectDateTimePickerButton(wrapper);

    dateTimePickerButton.vm.$emit('input', nowTimestamp);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe('07/12/2013 18:00');
  });

  test('Renders `date-time-picker-text-field` with default props', () => {
    const dateObject = new Date(nowTimestamp);
    const dateSpy = jest.spyOn(global, 'Date')
      .mockReturnValue(dateObject);

    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();

    dateSpy.mockClear();
  });

  test('Renders `date-time-picker-text-field` with custom props', () => {
    const dateObject = new Date(nowTimestamp);
    const dateSpy = jest.spyOn(global, 'Date')
      .mockReturnValue(dateObject);

    const wrapper = snapshotFactory({
      propsData: {
        value: 'now',
        label: 'label',
        name: 'name',
        roundHours: true,
        dateObjectPreparer: jest.fn().mockImplementation(date => new Date(date)),
      },
    });

    expect(wrapper.element).toMatchSnapshot();

    dateSpy.mockClear();
  });
});
