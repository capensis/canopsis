import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import { AVAILABILITY_DISPLAY_PARAMETERS } from '@/constants';

import AvailabilityDisplayParameterField from '@/components/other/availability/form/fields/availability-display-parameter-field.vue';

const stubs = {
  'v-select': createInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('availability-display-parameter-field', () => {
  const factory = generateShallowRenderer(AvailabilityDisplayParameterField, { stubs });
  const snapshotFactory = generateRenderer(AvailabilityDisplayParameterField);

  test('Show type changed after trigger radio field', () => {
    const wrapper = factory({
      propsData: {
        value: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
      },
    });

    selectSelectField(wrapper).triggerCustomEvent('input', AVAILABILITY_DISPLAY_PARAMETERS.downtime);

    expect(wrapper).toEmitInput(AVAILABILITY_DISPLAY_PARAMETERS.downtime);
  });

  test('Renders `availability-display-parameter-field` with required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-display-parameter-field` with custom props', () => {
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
