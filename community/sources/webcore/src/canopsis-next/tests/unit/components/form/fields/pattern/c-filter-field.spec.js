import flushPromises from 'flush-promises';
import Faker from 'faker';

import { createVueInstance, mount, shallowMount } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import { MAX_LIMIT } from '@/constants';

import CFilterField from '@/components/forms/fields/pattern/c-filter-field.vue';

const localVue = createVueInstance();

const factory = (options = {}) => shallowMount(CFilterField, {
  localVue,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CFilterField, {
  localVue,

  ...options,
});

const selectAutocomplete = wrapper => wrapper.find('v-autocomplete-stub');

describe('c-filter-field', () => {
  const filters = [
    {
      _id: 'id1',
      name: 'Fake 1',
    },
    {
      _id: 'id2',
      name: 'Fake 2',
    },
    {
      _id: 'id3',
      name: 'Fake 3',
    },
  ];

  const fetchFiltersList = jest.fn();
  const filtersGetter = jest.fn().mockReturnValue([]);
  const pendingGetter = jest.fn().mockReturnValue(false);
  const filterModule = {
    name: 'filter',
    getters: {
      pending: pendingGetter,
      items: filtersGetter,
    },
    actions: {
      fetchList: fetchFiltersList,
    },
  };
  const store = createMockedStoreModules([
    filterModule,
  ]);

  afterEach(() => {
    filtersGetter.mockClear();
    pendingGetter.mockClear();
    fetchFiltersList.mockClear();
  });

  it('Filters fetched after mount if pending is false', async () => {
    factory({
      store,
      propsData: {
        value: null,
      },
    });

    await flushPromises();

    expect(fetchFiltersList).toBeCalledTimes(1);
    expect(fetchFiltersList).toBeCalledWith(
      expect.any(Object),
      { params: { limit: MAX_LIMIT } },
      undefined,
    );
  });

  it('Filters fetched after mount if pending is true', async () => {
    const fetchFilters = jest.fn();
    factory({
      store,
      propsData: {
        value: null,
      },
    });

    await flushPromises();

    expect(fetchFilters).toBeCalledTimes(0);
  });

  it('Filter changed after trigger select field', () => {
    const newValue = Faker.datatype.string();
    filtersGetter.mockReturnValueOnce([{
      name: '',
      _id: newValue,
    }]);
    const wrapper = factory({
      propsData: {
        value: '',
      },
      store: createMockedStoreModules([filterModule]),
    });

    const valueElement = selectAutocomplete(wrapper);

    valueElement.vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  it('Renders `c-filter-field` with default props', () => {
    filtersGetter.mockReturnValueOnce(filters);
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([filterModule]),
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-filter-field` with default custom props', () => {
    filtersGetter.mockReturnValueOnce(filters);
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([filterModule]),
      propsData: {
        value: 'id1',
        label: 'Custom label',
        name: 'customName',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
