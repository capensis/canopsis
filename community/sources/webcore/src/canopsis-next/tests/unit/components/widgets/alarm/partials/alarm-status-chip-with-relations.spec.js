import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';
import { createActivatorElementStub } from '@unit/stubs/vuetify';

import { ALARM_STATUSES, MODALS } from '@/constants';

import AlarmStatusChipWithRelations from '@/components/widgets/alarm/partials/alarm-status-chip-with-relations.vue';

const stubs = {
  'c-simple-tooltip': createActivatorElementStub('c-simple-tooltip'),
};

const selectButton = wrapper => wrapper.find('v-btn-stub');

describe('alarm-status-chip-with-relations', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(AlarmStatusChipWithRelations, {
    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(AlarmStatusChipWithRelations, {
    stubs,
    mocks: { $modals },
  });

  test('Shows modal with correct config when button clicked and entity has upstream', async () => {
    const entity = {
      _id: Faker.datatype.uuid(),
      upstream: Faker.datatype.uuid(),
    };
    const status = ALARM_STATUSES.ongoing;
    const wrapper = factory({
      propsData: {
        alarm: {
          entity,
          v: {
            status: {
              val: status,
            },
          },
        },
      },
    });

    await flushPromises();

    const button = selectButton(wrapper);
    button.triggerCustomEvent('click', new Event('click'));

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.entityUpstream,
      config: {
        entity: { ...entity, status },
        title: expect.any(String),
      },
    });
  });

  test('Shows modal with correct config when button clicked and entity has no upstream', async () => {
    const entity = {
      _id: Faker.datatype.uuid(),
    };
    const status = ALARM_STATUSES.ongoing;
    const wrapper = factory({
      propsData: {
        alarm: {
          entity,
          v: {
            status: {
              val: status,
            },
          },
        },
      },
    });

    await flushPromises();

    const button = selectButton(wrapper);
    button.triggerCustomEvent('click', new Event('click'));

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.entityUpstream,
      config: {
        entity: { ...entity, status },
        title: expect.any(String),
      },
    });
  });

  test('Renders `alarm-status-chip-with-relations` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {},
          v: {
            status: {
              val: ALARM_STATUSES.ongoing,
            },
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-status-chip-with-relations` with small prop', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {},
          v: {
            status: {
              val: ALARM_STATUSES.ongoing,
            },
          },
        },
        small: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-status-chip-with-relations` with color prop', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {},
          v: {
            status: {
              val: ALARM_STATUSES.ongoing,
            },
          },
        },
        color: 'primary',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-status-chip-with-relations` with iconColor prop', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {},
          v: {
            status: {
              val: ALARM_STATUSES.ongoing,
            },
          },
        },
        iconColor: 'error',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-status-chip-with-relations` with outlined prop', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {},
          v: {
            status: {
              val: ALARM_STATUSES.ongoing,
            },
          },
        },
        outlined: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-status-chip-with-relations` with entity that has upstream', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {
            _id: Faker.datatype.uuid(),
            upstream: Faker.datatype.uuid(),
          },
          v: {
            status: {
              val: ALARM_STATUSES.ongoing,
            },
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-status-chip-with-relations` with entity that has no upstream', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {
            _id: Faker.datatype.uuid(),
          },
          v: {
            status: {
              val: ALARM_STATUSES.ongoing,
            },
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-status-chip-with-relations` with all custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {
            _id: Faker.datatype.uuid(),
            upstream: Faker.datatype.uuid(),
          },
          v: {
            status: {
              val: ALARM_STATUSES.ongoing,
            },
          },
        },
        small: true,
        color: 'error',
        iconColor: 'white',
        outlined: true,
      },
      slots: {
        default: '<span>Custom content</span>',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
