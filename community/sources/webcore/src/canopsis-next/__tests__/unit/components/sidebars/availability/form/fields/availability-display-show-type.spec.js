import { generateShallowRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityDisplayShowType from '@/components/sidebars/availability/form/fields/availability-display-show-type.vue';

const stubs = {
  'widget-settings-item': true,
  'availability-show-type-radio-field': true,
};

const selectAvailabilityShowTypeRadioField = wrapper =>
  wrapper.find('availability-show-type-radio-field-stub');

describe('availability-display-show-type', () => {
  const factory = generateShallowRenderer(AvailabilityDisplayShowType, { stubs });

  test('Show type changed after trigger radio field', () => {
    const wrapper = factory({
      propsData: {
        value: AVAILABILITY_SHOW_TYPE.percent,
      },
    });

    selectAvailabilityShowTypeRadioField(wrapper).triggerCustomEvent(
      'input',
      AVAILABILITY_SHOW_TYPE.duration,
    );

    expect(wrapper).toEmitInput(AVAILABILITY_SHOW_TYPE.duration);
  });
});
