import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { ALARM_STATES } from '@/constants';

import CAlarmStateField from '@/components/forms/fields/alarm/c-alarm-state-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

describe('c-alarm-state-field', () => {
  const factory = generateShallowRenderer(CAlarmStateField, { stubs });
  const snapshotFactory = generateRenderer(CAlarmStateField);

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

  it('Renders `c-alarm-state-field` with default props', () => {
    const wrapper = snapshotFactory();

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-alarm-state-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ALARM_STATES.major,
        label: 'Custom label',
        name: 'name',
        disabled: true,
        required: true,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-alarm-state-field` with validator error', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        required: true,
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(wrapper).toMatchSnapshot();
  });
});
