import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';

import { CRUD_ACTIONS, MODALS, USERS_PERMISSIONS } from '@/constants';

import PbehaviorsListActionBtn from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-list-action-btn.vue';

const stubs = {
  'c-action-btn': true,
};

const selectActionButtonNode = wrapper => wrapper.vm.$children[0];

describe('pbehaviors-list-action-btn', () => {
  const $modals = mockModals();

  const { authModule, currentUserPermissionsById } = createAuthModule();
  const store = createMockedStoreModules([authModule]);

  const factory = generateShallowRenderer(PbehaviorsListActionBtn, {

    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(PbehaviorsListActionBtn, {

    stubs,
    mocks: { $modals },
  });

  test('Pbehavior planning modal opened after trigger button', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.technical.exploitation.pbehavior]: { actions: [CRUD_ACTIONS.read] },
    });
    const entityId = Faker.datatype.string();

    const wrapper = factory({
      store: createMockedStoreModules([authModule]),
      propsData: {
        entityId,
      },
    });

    await selectActionButtonNode(wrapper).$emit('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.pbehaviorList,
        config: {
          entityId,
          availableActions: [CRUD_ACTIONS.delete, CRUD_ACTIONS.update],
        },
      },
    );
  });

  test('Renders `pbehaviors-list-action-btn` without access', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        entityId: 'entityId',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pbehaviors-list-action-btn` with access', () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.technical.exploitation.pbehavior]: { actions: [CRUD_ACTIONS.read] },
    });

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([authModule]),
      propsData: {
        entityId: 'entity-id',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
