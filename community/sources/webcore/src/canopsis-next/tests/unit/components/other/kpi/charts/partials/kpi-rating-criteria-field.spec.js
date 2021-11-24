import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { createMockedStoreModules } from '@unit/utils/store';
import { KPI_RATING_CRITERIA, MAX_LIMIT } from '@/constants';

import KpiRatingCriteriaField from '@/components/other/kpi/charts/partials/kpi-rating-criteria-field';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(KpiRatingCriteriaField, {
  localVue,
  stubs,
  store: createMockedStoreModules([{
    name: 'ratingSettings',
    getters: {
      pending: false,
      items: [],
    },
    actions: {
      fetchListWithoutStore: jest.fn(),
    },
  }]),

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiRatingCriteriaField, {
  localVue,

  ...options,
});

describe('kpi-rating-criteria-field', () => {
  const ratingSettings = [
    { id: 1, label: 'Rating setting 1' },
    { id: 2, label: 'Rating setting 2' },
    { id: 3, label: 'Rating setting 3' },
    { id: 4, label: 'Rating setting 4' },
    { id: 5, label: 'Rating setting 5' },
  ];

  it('Rating settings fetched after mount', async () => {
    const fetchRatingSettings = jest.fn(() => ({
      data: [],
    }));
    factory({
      propsData: {
        value: null,
        mandatory: true,
      },
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);
    expect(fetchRatingSettings).toBeCalledWith(
      expect.any(Object),
      { params: { limit: MAX_LIMIT, enabled: true } },
      undefined,
    );
  });

  it('First rating setting settled after fetch without value', async () => {
    const ratingSetting = { id: Faker.datatype.number() };
    const fetchRatingSettings = jest.fn(() => ({
      data: [ratingSetting],
    }));
    const wrapper = factory({
      propsData: {
        value: undefined,
        mandatory: true,
      },
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData).toEqual(ratingSetting);
  });

  it('First rating not set after fetch without value, if data is empty', async () => {
    const fetchRatingSettings = jest.fn(() => ({
      data: [],
    }));
    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toBeFalsy();
  });

  it('First rating not set after fetch without value, if value doesn\'t exist', async () => {
    const fetchRatingSettings = jest.fn(() => ({
      data: [{ id: 321 }],
    }));
    const wrapper = factory({
      propsData: {
        value: {
          id: 123,
        },
        mandatory: true,
      },
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData).toBe(undefined);
  });

  it('Criteria changed after trigger select field', () => {
    const fetchRatingSettings = jest.fn(() => ({
      data: [],
    }));
    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    const valueElement = wrapper.find('select.v-select');

    valueElement.vm.$emit('input', KPI_RATING_CRITERIA.role);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(KPI_RATING_CRITERIA.role);
  });

  it('Renders `kpi-rating-criteria-field` without props', async () => {
    const fetchRatingSettings = jest.fn(() => ({
      data: [],
    }));
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `kpi-rating-criteria-field` with custom props', async () => {
    const fetchRatingSettings = jest.fn(() => ({
      data: ratingSettings,
    }));
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          id: 2,
        },
      },
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
