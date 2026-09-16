import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import CPbehaviorReasonField from '@/components/forms/fields/pbehavior/c-pbehavior-reason-field.vue';

const { fetchPbehaviorReasonsListWithoutStore } = vi.hoisted(() => ({
  fetchPbehaviorReasonsListWithoutStore: vi.fn().mockResolvedValue({
    data: [],
    meta: {},
  }),
}));

vi.mock('@/hooks/store/modules/pbehavior-reason', () => ({
  usePbehaviorReason: () => ({
    fetchPbehaviorReasonsListWithoutStore,
  }),
}));

const stubs = {
  'c-autocomplete': createSelectInputStub('c-autocomplete'),
};

const selectSelectField = wrapper => wrapper.find('select.c-autocomplete');

describe('c-pbehavior-reason-field', () => {
  const factory = generateShallowRenderer(CPbehaviorReasonField, { stubs });

  afterEach(() => {
    fetchPbehaviorReasonsListWithoutStore.mockClear();
  });

  test('Reasons fetched after mount', async () => {
    factory();

    await flushPromises();

    expect(fetchPbehaviorReasonsListWithoutStore).toHaveBeenCalledWith({
      params: {
        paginate: false,
        with_hidden: true,
      },
    });
  });

  test('Reason changed after trigger select field', () => {
    const reason = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
    };
    fetchPbehaviorReasonsListWithoutStore.mockResolvedValueOnce({
      data: [reason],
    });

    const wrapper = factory({
      propsData: {
        modelValue: '',
      },
    });

    selectSelectField(wrapper).triggerCustomEvent('update:modelValue', reason._id);

    expect(wrapper).toEmitInput(reason._id);
  });
});
