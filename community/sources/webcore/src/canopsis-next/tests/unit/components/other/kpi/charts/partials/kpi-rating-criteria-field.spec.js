import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';
import { createMockedStoreModules } from '@unit/utils/store';

import { KPI_RATING_CRITERIA, MAX_LIMIT } from '@/constants';

import KpiRatingCriteriaField from '@/components/other/kpi/charts/form/fields/kpi-rating-criteria-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

describe('kpi-rating-criteria-field', () => {
  const ratingSettings = [
    { id: 1, label: 'Rating setting 1' },
    { id: 2, label: 'Rating setting 2' },
    { id: 3, label: 'Rating setting 3' },
    { id: 4, label: 'Rating setting 4' },
    { id: 5, label: 'Rating setting 5' },
  ];

  const factory = generateShallowRenderer(KpiRatingCriteriaField, {
    stubs,
    store: createMockedStoreModules([{
      name: 'ratingSettings',
      getters: {
        pending: false,
        items: [],
        updatedAt: null,
      },
      actions: {
        fetchListWithoutStore: jest.fn(),
      },
    }]),
  });
  const snapshotFactory = generateRenderer(KpiRatingCriteriaField);

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
        getters: {
          updatedAt: null,
        },
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
        getters: {
          updatedAt: null,
        },
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);

    expect(wrapper).toEmit('input', ratingSetting);
  });

  it('First rating not set after fetch without value, if data is empty', async () => {
    const fetchRatingSettings = jest.fn(() => ({
      data: [],
    }));
    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        getters: {
          updatedAt: null,
        },
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);

    expect(wrapper).not.toEmit('input');
  });

  it('First rating settled after fetch, if items doesn\'t includes value', async () => {
    const ratingSetting = { id: 321 };
    const fetchRatingSettings = jest.fn(() => ({
      data: [ratingSetting],
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
        getters: {
          updatedAt: null,
        },
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);

    expect(wrapper).toEmit('input', ratingSetting);
  });

  it('First rating not settled after fetch, if items includes value', async () => {
    const fetchRatingSettings = jest.fn(() => ({
      data: ratingSettings,
    }));
    const wrapper = factory({
      propsData: {
        value: ratingSettings[0],
        mandatory: true,
      },
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        getters: {
          updatedAt: null,
        },
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);

    expect(wrapper).not.toEmit('input');
  });

  it('First rating settled after fetch, if value is undefined', async () => {
    const fetchRatingSettings = jest.fn(() => ({
      data: ratingSettings,
    }));
    const wrapper = factory({
      propsData: {
        mandatory: true,
      },
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        getters: {
          updatedAt: null,
        },
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    expect(fetchRatingSettings).toBeCalledTimes(1);

    expect(wrapper).toEmit('input', ratingSettings[0]);
  });

  it('Criteria changed after trigger select field', () => {
    const fetchRatingSettings = jest.fn(() => ({
      data: [],
    }));
    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        getters: {
          updatedAt: null,
        },
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    wrapper.find('select.v-select').triggerCustomEvent('input', KPI_RATING_CRITERIA.role);

    expect(wrapper).toEmit('input', KPI_RATING_CRITERIA.role);
  });

  it('Renders `kpi-rating-criteria-field` without props', async () => {
    const fetchRatingSettings = jest.fn(() => ({
      data: [],
    }));
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'ratingSettings',
        getters: {
          updatedAt: null,
        },
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
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
        getters: {
          updatedAt: null,
        },
        actions: {
          fetchListWithoutStore: fetchRatingSettings,
        },
      }]),
    });

    await flushPromises();

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
