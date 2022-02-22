import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { ENTITIES_STATUSES } from '@/constants';

import CAlarmStatusField from '@/components/forms/fields/c-alarm-status-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CAlarmStatusField, {
  localVue,
  stubs,
  ...options,
});

const snapshotFactory = (options = {}) => mount(CAlarmStatusField, {
  localVue,
  ...options,
});

describe('c-alarm-status-field', () => {
  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        value: ENTITIES_STATUSES.closed,
      },
    });
    const selectElement = wrapper.find('select.v-select');

    selectElement.vm.$emit('input', ENTITIES_STATUSES.cancelled);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(ENTITIES_STATUSES.cancelled);
  });

  it('Renders `c-alarm-status-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ENTITIES_STATUSES.stealthy,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-alarm-status-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ENTITIES_STATUSES.flapping,
        label: 'Custom label',
        name: 'customAlarmStatusName',
        disabled: true,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
