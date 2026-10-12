import { generateShallowRenderer } from '@unit/utils/vue';

import KpiRatingCriteriaField from '@/components/other/kpi/charts/form/fields/kpi-rating-criteria-field.vue';

const stubs = {
  'c-select': true,
};

describe('kpi-rating-criteria-field', () => {
  const factory = generateShallowRenderer(KpiRatingCriteriaField, { stubs });

  it('selects the first available rating setting when mandatory', () => {
    const wrapper = factory({ props: { mandatory: true } });
    const ratingSetting = { id: 1, label: 'Rating setting' };

    wrapper.vm.items = [ratingSetting];
    wrapper.vm.setAvailableValue();

    expect(wrapper.emitted('update:value')).toEqual([[ratingSetting]]);
  });

  it('clears a mandatory value when no rating settings are available', () => {
    const wrapper = factory({ props: { mandatory: true } });

    wrapper.vm.items = [];
    wrapper.vm.setAvailableValue();

    expect(wrapper.emitted('update:value')).toEqual([[undefined]]);
  });

  it('keeps the current value when it is still available', () => {
    const ratingSetting = { id: 1, label: 'Rating setting' };
    const wrapper = factory({ props: { mandatory: true, value: ratingSetting } });

    wrapper.vm.items = [ratingSetting];
    wrapper.vm.setAvailableValue();

    expect(wrapper.emitted('update:value')).toBeUndefined();
  });

  it('emits a changed criteria value', () => {
    const wrapper = factory();
    const value = { id: 2, label: 'Another setting' };

    wrapper.vm.$emit('update:value', value);

    expect(wrapper.emitted('update:value')).toEqual([[value]]);
  });
});
