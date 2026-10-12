import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_STATES } from '@/constants';

import CAlarmStateField from '@/components/forms/fields/alarm/c-alarm-state-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

describe('c-alarm-state-field', () => {
  const factory = generateShallowRenderer(CAlarmStateField, { stubs });

  it('State type changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: ALARM_STATES.ok,
      },
    });

    const valueElement = wrapper.find('select.v-select');

    valueElement.triggerCustomEvent('input', ALARM_STATES.critical);

    expect(wrapper).toEmitInput(ALARM_STATES.critical);
  });
});
