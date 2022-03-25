import flushPromises from 'flush-promises';
import Faker from 'faker';

import { createVueInstance, mount, shallowMount } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import { MAX_LIMIT } from '@/constants';

import CFiltersField from '@/components/forms/fields/pattern/c-filter-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-autocomplete': {
    props: ['value', 'items'],
    template: `
      <select class="v-autocomplete" :value="value" @change="$listeners.input($event.target.value)">
        <option v-for="item in items" :value="item._id" :key="item._id">
          {{ item.name }}
        </option>
      </select>
  `,
  },
};

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

const factory = (options = {}) => shallowMount(CFiltersField, {
  localVue,
  stubs,
  store: createMockedStoreModules([{
    name: 'filter',
    getters: {
      pending: false,
      items: filters,
    },
    actions: {
      fetchList: jest.fn(),
    },
  }]),

  ...options,
});

const snapshotFactory = (options = {}) => mount(CFiltersField, {
  localVue,
  store: createMockedStoreModules([{
    name: 'filter',
    getters: {
      pending: false,
      items: filters,
    },
    actions: {
      fetchList: jest.fn(),
    },
  }]),

  ...options,
});

describe('c-filter-field', () => {
  it('Filters fetched after mount if pending is false', async () => {
    const fetchFilters = jest.fn();
    factory({
      propsData: {
        value: null,
      },
      store: createMockedStoreModules([{
        name: 'filter',
        getters: {
          pending: false,
          items: [],
        },
        actions: {
          fetchList: fetchFilters,
        },
      }]),
    });

    await flushPromises();

    expect(fetchFilters).toBeCalledTimes(1);
    expect(fetchFilters).toBeCalledWith(
      expect.any(Object),
      { params: { limit: MAX_LIMIT } },
      undefined,
    );
  });

  it('Filters fetched after mount if pending is true', async () => {
    const fetchFilters = jest.fn();
    factory({
      propsData: {
        value: null,
      },
      store: createMockedStoreModules([{
        name: 'filter',
        getters: {
          pending: true,
          items: [],
        },
        actions: {
          fetchList: fetchFilters,
        },
      }]),
    });

    await flushPromises();

    expect(fetchFilters).toBeCalledTimes(0);
  });

  it('Filter changed after trigger select field', () => {
    const newValue = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        value: '',
      },
      store: createMockedStoreModules([{
        name: 'filter',
        getters: {
          pending: true,
          items: [{
            name: '',
            _id: newValue,
          }],
        },
        actions: {
          fetchList: jest.fn(),
        },
      }]),
    });

    const valueElement = wrapper.find('select.v-autocomplete');

    valueElement.setValue(newValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(newValue);
  });

  it('Renders `c-filters-field` with default props', () => {
    const wrapper = snapshotFactory();

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-filters-field` with default custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'id1',
        label: 'Custom label',
        name: 'customName',
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
