import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { QUICK_RANGES } from '@/constants';

import CQuickDateIntervalTypeField from '@/components/forms/fields/c-quick-date-interval-type-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CQuickDateIntervalTypeField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CQuickDateIntervalTypeField, {
  localVue,

  ...options,
});

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('c-quick-date-interval-type-field', () => {
  test('Value changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: QUICK_RANGES.last2Days,
      },
    });

    const selectField = selectSelectField(wrapper);

    selectField.vm.$emit('input', QUICK_RANGES.yesterday);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(QUICK_RANGES.yesterday);
  });

  test('Renders `c-quick-date-interval-type-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: QUICK_RANGES.last2Days,
      },
    });

    const dropdownContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  test('Renders `c-quick-date-interval-type-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: QUICK_RANGES.last2Days,
        hideDetails: true,
        disabled: true,
      },
    });

    const dropdownContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  test('Renders `c-quick-date-interval-type-field` with custom ranges', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {},
        ranges: [QUICK_RANGES.last2Days, QUICK_RANGES.custom],
      },
    });

    const dropdownContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });
});
