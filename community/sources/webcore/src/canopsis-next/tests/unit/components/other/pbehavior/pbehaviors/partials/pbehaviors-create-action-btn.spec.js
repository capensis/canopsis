import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';

import { CRUD_ACTIONS, MODALS, USERS_PERMISSIONS } from '@/constants';

import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';

import PbehaviorsCreateActionBtn from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-create-action-btn.vue';

const stubs = {
  'c-action-btn': true,
};

const selectActionButtonNode = wrapper => wrapper.vm.$children[0];

describe('pbehaviors-create-action-btn', () => {
  const $modals = mockModals();

  const { authModule, currentUserPermissionsById } = createAuthModule();
  const store = createMockedStoreModules([authModule]);

  const factory = generateShallowRenderer(PbehaviorsCreateActionBtn, {

    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(PbehaviorsCreateActionBtn, {

    stubs,
    mocks: { $modals },
  });

  test('Pbehavior planning modal opened after trigger button', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.technical.exploitation.pbehavior]: { actions: [CRUD_ACTIONS.create] },
    });
    const entity = { _id: Faker.datatype.string() };

    const wrapper = factory({
      store: createMockedStoreModules([authModule]),
      propsData: {
        entity,
      },
    });

    await selectActionButtonNode(wrapper).$emit('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(entity._id),
          entities: [entity],
        },
      },
    );
  });

  test('Renders `pbehaviors-create-action-btn` without access', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        entity: { _id: 'entityId' },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pbehaviors-create-action-btn` with access', () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.technical.exploitation.pbehavior]: { actions: [CRUD_ACTIONS.create] },
    });

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([authModule]),
      propsData: {
        entity: { _id: 'entity-id' },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
