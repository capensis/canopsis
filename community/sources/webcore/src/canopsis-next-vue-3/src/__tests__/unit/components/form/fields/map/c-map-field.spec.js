import { createSelectInputStub } from '@unit/stubs/input';
import { createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { MAX_LIMIT } from '@/constants';

import CMapField from '@/components/forms/fields/map/c-map-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('c-map-field', () => {
  const fetchMapsListWithoutStore = vi.fn().mockReturnValue({
    meta: {},
    data: [],
  });
  const mapModule = {
    name: 'map',
    actions: {
      fetchListWithoutStore: fetchMapsListWithoutStore,
    },
  };
  const store = createMockedStoreModules([mapModule]);

  const factory = generateShallowRenderer(CMapField, {
    stubs,
    attachTo: document.body,
  });

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
    expect(fetchMapsListWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      params: { limit: MAX_LIMIT },
    });
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
});
