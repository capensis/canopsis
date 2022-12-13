import flushPromises from 'flush-promises';
import Faker from 'faker';

import { createVueInstance, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPbehaviorTypesModule } from '@unit/utils/store';
import { createSelectInputStub } from '@unit/stubs/input';
import { MAX_LIMIT } from '@/constants';

import CPbehaviorTypeField from '@/components/forms/fields/pbehavior/c-pbehavior-type-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('c-pbehavior-type-field', () => {
  const { pbehaviorTypesModule, fetchPbehaviorTypesListWithoutStore } = createPbehaviorTypesModule();
  const store = createMockedStoreModules([
    pbehaviorTypesModule,
  ]);

  const factory = generateShallowRenderer(CPbehaviorTypeField, {
    localVue,
    stubs,
    store,
  });

  const snapshotFactory = generateRenderer(CPbehaviorTypeField, {
    localVue,
    store,
  });

  test('Types fetched after mount', async () => {
    factory({
      store,
      propsData: {},
    });

    await flushPromises();

    expect(fetchPbehaviorTypesListWithoutStore).toBeCalledWith(
      expect.any(Object),
      { params: { limit: MAX_LIMIT } },
      undefined,
    );
  });

  test('Type changed after trigger select field', () => {
    const type = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
    };
    fetchPbehaviorTypesListWithoutStore.mockResolvedValueOnce({
      data: [type],
    });

    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    selectSelectField(wrapper).vm.$emit('input', type._id);

    expect(wrapper).toEmit('input', type._id);
  });

  test('Renders `c-pbehavior-type-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `c-pbehavior-type-field` with custom props', async () => {
    fetchPbehaviorTypesListWithoutStore.mockResolvedValueOnce({
      data: [
        { _id: 'type-id-1', name: 'type-name-1' },
        { _id: 'type-id-2', name: 'type-name-2' },
        { _id: 'type-id-3', name: 'type-name-3' },
      ],
    });

    const wrapper = snapshotFactory({
      propsData: {
        value: ['type-id-1', 'type-id-3'],
        name: 'customName',
        label: 'Custom label',
        disabled: true,
        multiple: true,
        chips: true,
        returnObject: true,
        required: true,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `c-pbehavior-type-field` with max prop', async () => {
    fetchPbehaviorTypesListWithoutStore.mockResolvedValueOnce({
      data: [
        { _id: 'type-id-1', name: 'type-name-1' },
        { _id: 'type-id-2', name: 'type-name-2' },
        { _id: 'type-id-3', name: 'type-name-3' },
      ],
    });

    const wrapper = snapshotFactory({
      propsData: {
        value: 'type-id-1',
        max: 1,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
