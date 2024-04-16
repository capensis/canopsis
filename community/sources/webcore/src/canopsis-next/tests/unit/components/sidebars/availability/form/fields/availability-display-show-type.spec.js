import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityDisplayShowType from '@/components/sidebars/availability/form/fields/availability-display-show-type.vue';

const stubs = {
  'widget-settings-item': true,
  'availability-show-type-radio-field': true,
};

const selectAvailabilityShowTypeRadioField = wrapper => wrapper.find('availability-show-type-radio-field-stub');

describe('availability-display-show-type', () => {
  const factory = generateShallowRenderer(AvailabilityDisplayShowType, { stubs });
  const snapshotFactory = generateRenderer(AvailabilityDisplayShowType, { stubs });

  test('Show type changed after trigger radio field', () => {
    const wrapper = factory({
      propsData: {
        value: AVAILABILITY_SHOW_TYPE.percent,
      },
    });

    selectAvailabilityShowTypeRadioField(wrapper).triggerCustomEvent('input', AVAILABILITY_SHOW_TYPE.duration);

    expect(wrapper).toEmitInput(AVAILABILITY_SHOW_TYPE.duration);
  });

  test('Renders `availability-display-show-type` with required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-display-show-type` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: AVAILABILITY_SHOW_TYPE.duration,
        name: 'customName',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
