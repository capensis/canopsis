import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { ALARM_METRIC_PARAMETERS } from '@/constants';
import CAlarmMetricParametersField from '@/components/forms/fields/c-alarm-metric-parameters-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CAlarmMetricParametersField, {
  localVue,
  stubs,
  ...options,
});

const snapshotFactory = (options = {}) => mount(CAlarmMetricParametersField, {
  localVue,
  ...options,
});

describe('c-alarm-metric-parameters-field', () => {
  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        value: [],
      },
    });
    const newValue = [ALARM_METRIC_PARAMETERS.ratioInstructions];
    const selectElement = wrapper.find('select.v-select');

    selectElement.vm.$emit('input', newValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(newValue);
  });

  it('Renders `c-alarm-metric-parameters-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: [ALARM_METRIC_PARAMETERS.createdAlarms],
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-alarm-metric-parameters-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: [ALARM_METRIC_PARAMETERS.createdAlarms, ALARM_METRIC_PARAMETERS.ratioInstructions],
        min: 2,
        name: 'customName',
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-alarm-metric-parameters-field` with all values', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: Object.values(ALARM_METRIC_PARAMETERS),
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});