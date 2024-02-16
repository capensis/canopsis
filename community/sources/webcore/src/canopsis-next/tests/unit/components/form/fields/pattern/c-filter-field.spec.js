import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { createSelectInputStub } from '@unit/stubs/input';

import { MAX_LIMIT } from '@/constants';

import CFilterField from '@/components/forms/fields/pattern/c-filter-field.vue';

const stubs = {
  'v-autocomplete': createSelectInputStub('v-autocomplete'),
};

const selectAutocomplete = wrapper => wrapper.find('.v-autocomplete');

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

  const factory = generateShallowRenderer(CFilterField, { stubs });
  const snapshotFactory = generateRenderer(CFilterField);

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

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-filter-field` with custom props', () => {
    filtersGetter.mockReturnValueOnce(filters);
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([filterModule]),
      propsData: {
        value: 'id1',
        label: 'Custom label',
        name: 'customName',
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-filter-field` with old_entity_patterns', () => {
    filtersGetter.mockReturnValueOnce(filters.map(filter => ({ ...filter, old_entity_patterns: true })));
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([filterModule]),
      propsData: {
        label: 'Custom label',
        name: 'customName',
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
