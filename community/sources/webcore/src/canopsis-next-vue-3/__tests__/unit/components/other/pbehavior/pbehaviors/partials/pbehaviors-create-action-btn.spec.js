import { mockModals } from '@unit/utils/mock-hooks';
import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { CRUD_ACTIONS, MODALS, USER_PERMISSIONS } from '@/constants';

import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';

import PbehaviorsCreateActionBtn from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-create-action-btn.vue';

const stubs = {
  'c-action-btn': true,
};

const selectActionButtonNode = wrapper => wrapper.vm.$children[0];

describe('pbehaviors-create-action-btn', () => {
  const $modals = mockModals();

  const { authModule, currentUserPermissionsById } = createAuthModule();

  const factory = generateShallowRenderer(PbehaviorsCreateActionBtn, {
    stubs,
    mocks: { $modals },
  });

  test('Pbehavior planning modal opened after trigger button', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USER_PERMISSIONS.technical.exploitation.pbehavior]: { actions: [CRUD_ACTIONS.create] },
    });
    const entity = { _id: Faker.datatype.string() };

    const wrapper = factory({
      store: createMockedStoreModules([authModule]),
      propsData: {
        entity,
      },
    });

    await selectActionButtonNode(wrapper).$emit('click');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.pbehaviorPlanning,
      config: {
        entityPattern: createEntityIdPatternByValue(entity._id),
        entities: [entity],
      },
    });
  });
});
