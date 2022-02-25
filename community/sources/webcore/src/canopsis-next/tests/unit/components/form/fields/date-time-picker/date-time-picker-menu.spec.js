import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import DateTimePickerMenu from '@/components/forms/fields/date-time-picker/date-time-picker-menu.vue';

const localVue = createVueInstance();

const stubs = {
  'date-time-picker': true,
};

const factory = (options = {}) => shallowMount(DateTimePickerMenu, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(DateTimePickerMenu, {
  localVue,
  stubs,

  ...options,
});

const selectMenuButton = wrapper => wrapper.find('.v-btn, v-btn-stub');
const selectDateTimePicker = wrapper => wrapper.find('date-time-picker-stub');

describe('date-time-picker-menu', () => {
  const nowTimestamp = 1386435600000;
  mockDateNow(nowTimestamp);

  test('Value changed after trigger input on the date time picker', async () => {
    const date = new Date(1387435600000);
    const input = jest.fn();
    const wrapper = factory({
      propsData: {
        value: date,
      },
      listeners: {
        input,
      },
    });

    const dateTimePicker = selectDateTimePicker(wrapper);

    const newDate = new Date(1387536600000);

    await dateTimePicker.vm.$emit('input', newDate);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(newDate);
  });

  test('Menu opened after trigger button', async () => {
    const wrapper = snapshotFactory();

    const menuButton = selectMenuButton(wrapper);

    await menuButton.trigger('click');

    const dateTimePicker = selectDateTimePicker(wrapper);

    expect(dateTimePicker.element).toBeTruthy();
  });

  test('Menu closed after trigger close event on the date time picker', async () => {
    const wrapper = factory();

    const menuButton = selectMenuButton(wrapper);

    await menuButton.trigger('click');

    const dateTimePicker = selectDateTimePicker(wrapper);

    dateTimePicker.vm.$emit('close');

    await flushPromises();

    expect(wrapper.vm.opened).toBeFalsy();
  });

  test('Renders `date-time-picker-menu` with default props', () => {
    const dateObject = new Date(nowTimestamp);
    const dateSpy = jest.spyOn(global, 'Date')
      .mockReturnValue(dateObject);
    Date.now = jest.fn().mockReturnValue(nowTimestamp);

    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();

    dateSpy.mockClear();
  });

  test('Renders `date-time-picker-menu` with opened menu', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: new Date(1387435600000),
        label: 'label',
        roundHours: true,
      },
    });

    const menuButton = selectMenuButton(wrapper);

    await menuButton.trigger('click');

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
