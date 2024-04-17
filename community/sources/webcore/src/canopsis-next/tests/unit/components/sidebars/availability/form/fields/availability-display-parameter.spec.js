import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { AVAILABILITY_DISPLAY_PARAMETERS } from '@/constants';

import AvailabilityDisplayParameter from '@/components/sidebars/availability/form/fields/availability-display-parameter.vue';

const stubs = {
  'widget-settings-item': true,
  'availability-display-parameter-radio-field': true,
};

const selectAvailabilityDisplayParameterRadioField = wrapper => wrapper.find('availability-display-parameter-radio-field-stub');

describe('availability-display-parameter', () => {
  const factory = generateShallowRenderer(AvailabilityDisplayParameter, { stubs });
  const snapshotFactory = generateRenderer(AvailabilityDisplayParameter, { stubs });

  test('Display parameter type changed after trigger radio field', () => {
    const wrapper = factory({
      propsData: {
        value: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
      },
    });

    selectAvailabilityDisplayParameterRadioField(wrapper).triggerCustomEvent('input', AVAILABILITY_DISPLAY_PARAMETERS.downtime);

    expect(wrapper).toEmitInput(AVAILABILITY_DISPLAY_PARAMETERS.downtime);
  });

  test('Renders `availability-display-parameter` with required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-display-parameter` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        name: 'customName',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
