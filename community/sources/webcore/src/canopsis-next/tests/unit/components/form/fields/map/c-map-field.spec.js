import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { createSelectInputStub } from '@unit/stubs/input';
import { MAX_LIMIT } from '@/constants';

import CMapField from '@/components/forms/fields/map/c-map-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CMapField, {
  localVue,
  stubs,
  attachTo: document.body,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CMapField, {
  localVue,
  attachTo: document.body,

  ...options,
});

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('c-map-field', () => {
  const fetchMapsListWithoutStore = jest.fn().mockReturnValue({
    meta: {},
    data: [],
  });
  const mapModule = {
    name: 'map',
    actions: {
      fetchListWithoutStore: fetchMapsListWithoutStore,
    },
  };
  const store = createMockedStoreModules([
    mapModule,
  ]);

  afterEach(() => {
    fetchMapsListWithoutStore.mockClear();
  });

  it('Filters fetched after mount', async () => {
    factory({
      store,
      propsData: {},
    });

    await flushPromises();

    expect(fetchMapsListWithoutStore).toBeCalledTimes(1);
    expect(fetchMapsListWithoutStore).toBeCalledWith(
      expect.any(Object),
      { params: { limit: MAX_LIMIT } },
      undefined,
    );
  });

  it('Map changed after trigger select field', () => {
    const map = { _id: Faker.datatype.string() };
    fetchMapsListWithoutStore.mockReturnValueOnce({
      data: [map],
    });

    const wrapper = factory({
      propsData: {
        value: '',
      },
      store: createMockedStoreModules([mapModule]),
    });

    const selectField = selectSelectField(wrapper);

    selectField.vm.$emit('input', map._id);

    expect(wrapper).toEmit('input', map._id);
  });

  it('Renders `c-map-field` with default props', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([mapModule]),
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-map-field` with custom props', async () => {
    fetchMapsListWithoutStore.mockReturnValueOnce({
      data: [
        { _id: 'map-id-1' },
        { _id: 'map-id-2' },
      ],
    });

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([mapModule]),
      propsData: {
        value: 'map-id-1',
        label: 'Custom label',
        name: 'customName',
        disabled: true,
        required: true,
        hideDetails: true,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
