import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPbehaviorReasonModule } from '@unit/utils/store';
import { createSelectInputStub } from '@unit/stubs/input';
import { MAX_LIMIT } from '@/constants';

import CPbehaviorReasonField from '@/components/forms/fields/pbehavior/c-pbehavior-reason-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('c-pbehavior-reason-field', () => {
  const { pbehaviorReasonModule, fetchPbehaviorReasonsListWithoutStore } = createPbehaviorReasonModule();
  const store = createMockedStoreModules([
    pbehaviorReasonModule,
  ]);

  const factory = generateShallowRenderer(CPbehaviorReasonField, { stubs, store });
  const snapshotFactory = generateRenderer(CPbehaviorReasonField, { store });

  test('Reasons fetched after mount', async () => {
    factory({
      store,
      propsData: {},
    });

    await flushPromises();

    expect(fetchPbehaviorReasonsListWithoutStore).toBeCalledWith(
      expect.any(Object),
      { params: { limit: MAX_LIMIT } },
      undefined,
    );
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
        value: '',
      },
    });

    selectSelectField(wrapper).vm.$emit('input', reason._id);

    expect(wrapper).toEmit('input', reason._id);
  });

  test('Renders `c-pbehavior-reason-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `c-pbehavior-reason-field` with custom props', async () => {
    fetchPbehaviorReasonsListWithoutStore.mockResolvedValueOnce({
      data: [
        { _id: 'reason-id-1', name: 'reason-name-1' },
        { _id: 'reason-id-2', name: 'reason-name-2' },
      ],
    });

    const wrapper = snapshotFactory({
      propsData: {
        reason: 'reason-id-1',
        name: 'customName',
        required: true,
        returnObject: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
