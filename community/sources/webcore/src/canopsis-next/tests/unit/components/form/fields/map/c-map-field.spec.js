import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { createSelectInputStub } from '@unit/stubs/input';

import { MAX_LIMIT } from '@/constants';

import CMapField from '@/components/forms/fields/map/c-map-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

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

  const factory = generateShallowRenderer(CMapField, {
    stubs,
    attachTo: document.body,
  });
  const snapshotFactory = generateRenderer(CMapField, { attachTo: document.body });

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

    selectField.triggerCustomEvent('input', map._id);

    expect(wrapper).toEmitInput(map._id);
  });

  it('Renders `c-map-field` with default props', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([mapModule]),
    });

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
