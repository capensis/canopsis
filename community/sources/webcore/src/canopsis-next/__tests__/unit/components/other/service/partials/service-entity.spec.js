import { mockModals } from '@unit/utils/mock-hooks';
import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ENTITY_TYPES, MODALS, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import ServiceEntity from '@/components/other/service/partials/service-entity.vue';

const stubs = {
  'v-expansion-panel-content': {
    template: '<div><slot name="header"/><slot/></div>',
  },
  'service-entity-header': {
    name: 'ServiceEntityHeader',
    emits: ['update:selected', 'remove:unavailable'],
    template: '<div />',
  },
  'service-entity-info-tab': {
    name: 'ServiceEntityInfoTab',
    emits: ['apply', 'execute'],
    template: '<div />',
  },
  'service-entity-tree-of-dependencies-tab': true,
  'pbehaviors-simple-list': true,
};

const selectServiceEntityHeader = wrapper => wrapper.findComponent({ name: 'ServiceEntityHeader' });

describe('service-entity', () => {
  const $modals = mockModals();
  const entity = {
    _id: 'service-id',
    alarm_id: 'service-alarm-id',
    source_type: ENTITY_TYPES.service,
    pbehaviors: [
      {
        type: {
          type: PBEHAVIOR_TYPE_TYPES.pause,
        },
      },
    ],
    ack: null,
  };

  const { authModule } = createAuthModule();

  const store = createMockedStoreModules([authModule]);

  const factory = generateShallowRenderer(ServiceEntity, {
    store,

    stubs,
    mocks: { $modals },
    propsData: {
      entity,
    },
  });

  test('Selected update emitted after trigger entity header', async () => {
    const wrapper = factory();

    const header = selectServiceEntityHeader(wrapper);

    header.vm.$emit('update:selected');

    expect(wrapper).toHaveBeenEmit('update:selected');
  });

  test('Remove unavailable emitted after trigger entity header', async () => {
    const wrapper = factory();

    const header = selectServiceEntityHeader(wrapper);

    header.vm.$emit('remove:unavailable');

    expect(wrapper).toHaveBeenEmit('remove:unavailable');
  });

  test('Instruction execution opens the remediation modal', () => {
    const wrapper = factory();

    const assignedInstruction = {};

    wrapper.vm.executeAlarmInstruction(assignedInstruction);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.executeRemediationSimpleInstruction,
      config: {
        assignedInstruction,
        alarmId: entity.alarm_id,
        onClose: expect.any(Function),
        onComplete: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.onClose();

    expect(wrapper).toHaveBeenEmit('refresh');
  });
});
