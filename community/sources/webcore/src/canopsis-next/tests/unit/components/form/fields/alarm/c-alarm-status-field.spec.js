import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { ALARM_STATUSES } from '@/constants';

import CAlarmStatusField from '@/components/forms/fields/alarm/c-alarm-status-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

describe('c-alarm-status-field', () => {
  const factory = generateShallowRenderer(CAlarmStatusField, { stubs });
  const snapshotFactory = generateRenderer(CAlarmStatusField);

  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        value: ALARM_STATUSES.closed,
      },
    });
    const selectElement = wrapper.find('select.v-select');

    selectElement.triggerCustomEvent('input', ALARM_STATUSES.cancelled);

    expect(wrapper).toEmitInput(ALARM_STATUSES.cancelled);
  });

  it('Renders `c-alarm-status-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ALARM_STATUSES.stealthy,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-alarm-status-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ALARM_STATUSES.flapping,
        label: 'Custom label',
        name: 'customAlarmStatusName',
        disabled: true,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
