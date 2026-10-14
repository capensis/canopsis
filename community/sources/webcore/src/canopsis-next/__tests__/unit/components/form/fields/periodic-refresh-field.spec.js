import { generateShallowRenderer } from '@unit/utils/vue';

import { TIME_UNITS } from '@/constants';

import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';

const stubs = {
  'c-duration-field': {
    props: ['name', 'duration'],
    emits: ['update:duration'],
    template: `
      <input
        class="c-duration-field"
        :value="duration?.value"
      />
    `,
  },
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectDurationField = wrapper => wrapper.find('.c-duration-field');

describe('periodic-refresh-field', () => {
  const factory = generateShallowRenderer(PeriodicRefreshField, { stubs });

  it('Enabled changed after trigger the enabled field', () => {
    const periodicRefresh = {
      enabled: true,
      unit: TIME_UNITS.second,
      value: 2,
    };
    const wrapper = factory({
      propsData: {
        periodicRefresh,
      },
    });

    selectEnabledField(wrapper).triggerCustomEvent('input', false);

    expect(wrapper).toEmitInput({
      ...periodicRefresh,
      enabled: false,
    });
  });

  it('Duration changed after trigger the duration field', () => {
    const periodicRefresh = {
      enabled: true,
      unit: TIME_UNITS.second,
      value: 2,
    };
    const wrapper = factory({
      propsData: {
        periodicRefresh,
      },
    });

    const newDuration = {
      enabled: true,
      unit: TIME_UNITS.week,
      value: 5,
    };

    selectDurationField(wrapper).triggerCustomEvent('update:duration', newDuration);

    expect(wrapper).toEmitInput(newDuration);
  });
});
