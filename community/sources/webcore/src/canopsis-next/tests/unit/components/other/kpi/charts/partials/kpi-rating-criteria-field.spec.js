import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { KPI_RATING_CRITERIA } from '@/constants';

import KpiRatingCriteriaField from '@/components/other/kpi/charts/partials/kpi-rating-criteria-field';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(KpiRatingCriteriaField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiRatingCriteriaField, {
  localVue,

  ...options,
});

describe('kpi-rating-criteria-field', () => {
  it('Criteria changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: KPI_RATING_CRITERIA.user,
      },
    });

    const valueElement = wrapper.find('select.v-select');

    valueElement.vm.$emit('input', KPI_RATING_CRITERIA.role);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(KPI_RATING_CRITERIA.role);
  });

  it('Renders `kpi-rating-criteria-field` without props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: KPI_RATING_CRITERIA.user,
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `kpi-rating-criteria-field` with mocked $te', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: KPI_RATING_CRITERIA.user,
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
