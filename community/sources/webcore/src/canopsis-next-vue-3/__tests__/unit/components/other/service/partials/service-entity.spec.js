import { mockModals } from '@unit/utils/mock-hooks';
import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ENTITY_TYPES, MODALS, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import ServiceEntity from '@/components/other/service/partials/service-entity.vue';

const stubs = {
  'v-expansion-panel-content': {
    template: '<div><slot name="header"/><slot/></div>',
  },
  'service-entity-header': true,
  'service-entity-info-tab': true,
  'service-entity-tree-of-dependencies-tab': true,
  'pbehaviors-simple-list': true,
};

const selectServiceEntityHeader = wrapper => wrapper.find('service-entity-header-stub');
const selectServiceEntityInfoTab = wrapper => wrapper.find('service-entity-info-tab-stub');

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

  const updateSelected = vi.fn();
  const removeUnavailable = vi.fn();

  const { authModule } = createAuthModule();

  const store = createMockedStoreModules([authModule]);

  const factory = generateShallowRenderer(ServiceEntity, {
    store,

    stubs,
    mocks: { $modals },
    listeners: {
      'update:selected': updateSelected,
      'remove:unavailable': removeUnavailable,
    },
    propsData: {
      entity,
    },
  });

  test('Selected update emitted after trigger entity header', async () => {
    const wrapper = factory();

    const header = selectServiceEntityHeader(wrapper);

    await header.triggerCustomEvent('update:selected');

    expect(updateSelected).toHaveBeenCalled();
  });

  test('Remove unavailable emitted after trigger entity header', async () => {
    const wrapper = factory();

    const header = selectServiceEntityHeader(wrapper);

    await header.triggerCustomEvent('remove:unavailable');

    expect(removeUnavailable).toHaveBeenCalled();
  });

  test('Instruction executed after trigger entity info tab', async () => {
    const wrapper = factory();

    const info = selectServiceEntityInfoTab(wrapper);
    const assignedInstruction = {};

    await info.triggerCustomEvent('execute', assignedInstruction);

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
