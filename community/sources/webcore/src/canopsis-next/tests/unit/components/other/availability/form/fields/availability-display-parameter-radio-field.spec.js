import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import { AVAILABILITY_DISPLAY_PARAMETERS } from '@/constants';

import AvailabilityDisplayParameterRadioField from '@/components/other/availability/form/fields/availability-display-parameter-radio-field.vue';

const stubs = {
  'v-radio-group': createInputStub('v-radio-group'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('availability-display-parameter-radio-field', () => {
  const factory = generateShallowRenderer(AvailabilityDisplayParameterRadioField, { stubs });
  const snapshotFactory = generateRenderer(AvailabilityDisplayParameterRadioField);

  test('Show type changed after trigger radio field', () => {
    const wrapper = factory({
      propsData: {
        value: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
      },
    });

    selectRadioGroup(wrapper).triggerCustomEvent('input', AVAILABILITY_DISPLAY_PARAMETERS.downtime);

    expect(wrapper).toEmitInput(AVAILABILITY_DISPLAY_PARAMETERS.downtime);
  });

  test('Renders `availability-display-parameter-radio-field` with required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-display-parameter-radio-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        label: 'Custom label',
        name: 'customName',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
