import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPbehaviorTypesModule } from '@unit/utils/store';
import { createSelectInputStub } from '@unit/stubs/input';

import CPbehaviorTypeField from '@/components/forms/fields/pbehavior/c-pbehavior-type-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('c-pbehavior-type-field', () => {
  const { pbehaviorTypesModule, fieldPbehaviorTypes } = createPbehaviorTypesModule();
  const store = createMockedStoreModules([
    pbehaviorTypesModule,
  ]);

  const factory = generateShallowRenderer(CPbehaviorTypeField, { stubs, store });
  const snapshotFactory = generateRenderer(CPbehaviorTypeField, { store });

  test('Type changed after trigger select field', () => {
    const type = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
    };
    fieldPbehaviorTypes.mockReturnValueOnce([type]);

    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    selectSelectField(wrapper).triggerCustomEvent('input', type._id);

    expect(wrapper).toEmit('input', type._id);
  });

  test('Renders `c-pbehavior-type-field` with default props', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        pbehaviorTypesModule,
      ]),
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `c-pbehavior-type-field` with custom props', async () => {
    fieldPbehaviorTypes.mockReturnValueOnce([
      { _id: 'type-id-1', name: 'type-name-1' },
      { _id: 'type-id-2', name: 'type-name-2' },
      { _id: 'type-id-3', name: 'type-name-3' },
    ]);

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
      store: createMockedStoreModules([
        pbehaviorTypesModule,
      ]),
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `c-pbehavior-type-field` with max prop', async () => {
    fieldPbehaviorTypes.mockReturnValueOnce([
      { _id: 'type-id-1', name: 'type-name-1' },
      { _id: 'type-id-2', name: 'type-name-2' },
      { _id: 'type-id-3', name: 'type-name-3' },
    ]);

    const wrapper = snapshotFactory({
      propsData: {
        value: 'type-id-1',
        max: 1,
      },
      store: createMockedStoreModules([
        pbehaviorTypesModule,
      ]),
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
