import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';
import { ref } from 'vue';

import { MAX_LIMIT } from '@/constants';

import CFilterField from '@/components/forms/fields/pattern/c-filter-field.vue';

const { fetchFiltersList } = vi.hoisted(() => ({
  fetchFiltersList: vi.fn(),
}));

const filters = ref([]);
const filtersPending = ref(false);

vi.mock('@/hooks/store/modules/filter', () => ({
  useFilter: () => ({
    filters,
    filtersPending,
    fetchFiltersList,
  }),
}));

const stubs = {
  'v-autocomplete': createSelectInputStub('v-autocomplete'),
};

const selectAutocomplete = wrapper => wrapper.find('.v-autocomplete');

describe('c-filter-field', () => {
  const factory = generateShallowRenderer(CFilterField, { stubs });

  afterEach(() => {
    filters.value = [];
    filtersPending.value = false;
    fetchFiltersList.mockClear();
  });

  it('Filters fetched after mount if pending is false', async () => {
    factory({
      propsData: {
        value: null,
      },
    });

    await flushPromises();

    expect(fetchFiltersList).toHaveBeenCalledTimes(1);
    expect(fetchFiltersList).toHaveBeenCalledWith({ params: { limit: MAX_LIMIT } });
  });

  it('Filters not fetched after mount if pending is true', async () => {
    filtersPending.value = true;

    factory({
      propsData: {
        value: null,
      },
    });

    await flushPromises();

    expect(fetchFiltersList).not.toHaveBeenCalled();
  });

  it('Filter changed after trigger select field', () => {
    const newValue = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    selectAutocomplete(wrapper).triggerCustomEvent('update:modelValue', newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
