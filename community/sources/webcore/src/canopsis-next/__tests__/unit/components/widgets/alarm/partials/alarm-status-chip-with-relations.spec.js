import { createActivatorElementStub } from '@unit/stubs/vuetify';
import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { ALARM_STATES, ALARM_STATUSES, MODALS } from '@/constants';

import AlarmStatusChipWithRelations from '@/components/widgets/alarm/partials/alarm-status-chip-with-relations.vue';

const stubs = {
  'c-tooltip': createActivatorElementStub('c-tooltip'),
};

describe('alarm-status-chip-with-relations', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(AlarmStatusChipWithRelations, {
    stubs,
    mocks: { $modals },
  });

  test('Shows modal with correct config when button clicked and entity has upstream', async () => {
    const entity = {
      _id: Faker.datatype.uuid(),
      upstream: Faker.datatype.uuid(),
    };
    const status = ALARM_STATUSES.ongoing;
    const state = ALARM_STATES.ok;
    const wrapper = factory({
      propsData: {
        alarm: {
          entity,
          v: {
            status: {
              val: status,
            },
            state: {
              val: state,
            },
          },
        },
      },
    });

    await flushPromises();

    wrapper.vm.showRelationsModal();

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.entityUpstream,
      config: {
        entity: { ...entity, status, state },
        title: expect.any(String),
      },
    });
  });

  test('Shows modal with correct config when button clicked and entity has no upstream', async () => {
    const entity = {
      _id: Faker.datatype.uuid(),
    };
    const status = ALARM_STATUSES.ongoing;
    const state = ALARM_STATES.ok;
    const wrapper = factory({
      propsData: {
        alarm: {
          entity,
          v: {
            status: {
              val: status,
            },
            state: {
              val: state,
            },
          },
        },
      },
    });

    await flushPromises();

    wrapper.vm.showRelationsModal();

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.entityUpstream,
      config: {
        entity: { ...entity, status, state },
        title: expect.any(String),
      },
    });
  });
});
