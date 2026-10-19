import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_STATUSES } from '@/constants';

import CAlarmStatusField from '@/components/forms/fields/alarm/c-alarm-status-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

describe('c-alarm-status-field', () => {
  const factory = generateShallowRenderer(CAlarmStatusField, { stubs });

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
});
