import { mockModals } from '@unit/utils/mock-hooks';
import {
  createAlarmModule,
  createAuthModule,
  createDeclareTicketModule,
  createEntityModule,
  createMockedStoreModules,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import {
  ALARM_STATES,
  ENTITY_TYPES,
  MODALS,
  PBEHAVIOR_TYPE_TYPES,
  USER_PERMISSIONS,
  WEATHER_ACK_EVENT_OUTPUT,
  WEATHER_ACTIONS_TYPES,
} from '@/constants';

import ServiceEntitiesList from '@/components/other/service/partials/service-entities-list.vue';

const stubs = {
  'service-entity-actions': true,
  'service-entity': true,
  'c-table-pagination': {
    name: 'CTablePagination',
    emits: ['update:page', 'update:items-per-page', 'update:pagination'],
    template: '<div class="c-table-pagination-stub" />',
  },
  'v-checkbox': {
    name: 'VCheckbox',
    emits: ['update:modelValue'],
    template: '<div class="v-checkbox-stub" />',
  },
};

const selectServiceEntityByIndex = (wrapper, index) =>
  wrapper.findAllComponents({ name: 'ServiceEntity' }).at(index);
const selectTablePagination = wrapper => wrapper.findComponent({ name: 'CTablePagination' });

const applyEntitiesAction = async (wrapper, type) => {
  wrapper.vm.applyActionForSelected({ type });

  await flushPromises();
};

const selectAllEntities = async (wrapper) => {
  wrapper.vm.isAllSelected = true;
  await flushPromises();
};

const unselectAllEntities = async (wrapper) => {
  wrapper.vm.isAllSelected = false;
  await flushPromises();
};

describe('service-entities-list', () => {
  const $modals = mockModals();
  const service = {
    _id: 'service-id',
  };
  const serviceEntities = [
    {
      _id: 'service-entity-1-id',
      pbehaviors: [],
      state: {
        val: ALARM_STATES.major,
      },
      alarm_id: 'alarm-id',
    },
    {
      _id: 'service-entity-2-id',
      pbehaviors: [
        {
          type: {
            type: PBEHAVIOR_TYPE_TYPES.pause,
          },
        },
      ],
      alarm_id: 'alarm-id',
    },
  ];
  const { authModule, currentUserPermissionsById } = createAuthModule();
  const {
    alarmModule,
    bulkCreateAlarmAckEvent,
    bulkCreateAlarmAckremoveEvent,
    bulkCreateAlarmAssocticketEvent,
    bulkCreateAlarmCommentEvent,
    bulkCreateAlarmCancelEvent,
    bulkCreateAlarmChangestateEvent,
  } = createAlarmModule();
  const {
    declareTicketRuleModule,
    bulkCreateDeclareTicketExecution,
    fetchAssignedDeclareTicketsWithoutStore,
  } = createDeclareTicketModule();
  const { entityModule, bulkPause, bulkResume } = createEntityModule();

  const store = createMockedStoreModules([
    alarmModule,
    authModule,
    declareTicketRuleModule,
    entityModule,
  ]);

  const refresh = vi.fn();

  const factory = generateShallowRenderer(ServiceEntitiesList, {
    store,
    stubs,
    mocks: {
      $modals,
    },
    propsData: {
      service,
      options: {},
    },
    listeners: {
      refresh,
    },
  });

  test('Entities selected after trigger select event on service entity', async () => {
    const wrapper = factory({
      propsData: {
        serviceEntities,
      },
    });

    selectServiceEntityByIndex(wrapper, 1).vm.$emit('update:selected', true);
    await flushPromises();

    expect(wrapper.vm.selectedEntities).toEqual([serviceEntities[1]]);

    selectServiceEntityByIndex(wrapper, 1).vm.$emit('update:selected', false);
    await flushPromises();

    expect(wrapper.vm.selectedEntities).toEqual([]);
  });

  test('All entities selected and unselected after trigger checkbox', async () => {
    const wrapper = factory({
      propsData: {
        serviceEntities,
      },
    });

    await selectAllEntities(wrapper);
    expect(wrapper.vm.selectedEntities).toEqual(serviceEntities);

    await unselectAllEntities(wrapper);
    expect(wrapper.vm.selectedEntities).toEqual([]);
  });

  test('Ack remove action applied after trigger mass ackremove  action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      alarm_id: Faker.datatype.string(),
      ack: {},
      pbehaviors: [],
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectAllEntities(wrapper);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityAckRemove);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.textFieldEditor,
      config: {
        title: 'modals.createAckRemove.title',
        field: {
          name: 'comment',
          label: 'common.note',
          validationRules: 'required',
        },
        action: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];
    const comment = Faker.datatype.string();

    await config.action({ comment });

    expect(bulkCreateAlarmAckremoveEvent).toHaveBeenCalledWith(expect.any(Object), {
      data: [
        {
          _id: entity.alarm_id,
          comment,
        },
      ],
    });
  });

  test('Associate ticket action applied after trigger mass associate ticket action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      ack: {},
      pbehaviors: [],
      alarm_id: 'alarm-id',
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectAllEntities(wrapper);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityAssocTicket);

    await flushPromises();

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createAssociateTicketEvent,
      config: {
        items: [{}],
        action: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];
    const event = {
      ticket: Faker.datatype.string(),
      ticket_url: Faker.datatype.string(),
    };

    await config.action(event);

    expect(bulkCreateAlarmAssocticketEvent).toHaveBeenCalledWith(expect.any(Object), {
      data: [
        {
          _id: entity.alarm_id,
          ...event,
        },
      ],
    });
  });

  test('Validate action applied after trigger mass validate action', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USER_PERMISSIONS.business.serviceWeather.actions.entityValidate]: {
        actions: [],
      },
    });

    const entity = {
      _id: Faker.datatype.string(),
      state: {
        val: ALARM_STATES.major,
      },
      pbehaviors: [],
      alarm_id: 'alarm-id',
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectAllEntities(wrapper);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityValidate);
    await flushPromises();

    expect(bulkCreateAlarmAckEvent).toHaveBeenCalledWith(expect.any(Object), {
      data: [
        {
          _id: entity.alarm_id,
          comment: WEATHER_ACK_EVENT_OUTPUT.validateOk,
        },
      ],
    });

    expect(bulkCreateAlarmChangestateEvent).toHaveBeenCalledWith(expect.any(Object), {
      data: [
        {
          _id: entity.alarm_id,
          comment: WEATHER_ACK_EVENT_OUTPUT.validateOk,
          state: ALARM_STATES.critical,
        },
      ],
    });
  });

  test('Invalidate action applied after trigger mass invalidate action', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USER_PERMISSIONS.business.serviceWeather.actions.entityInvalidate]: {
        actions: [],
      },
    });

    const entity = {
      _id: Faker.datatype.string(),
      state: {
        val: ALARM_STATES.major,
      },
      pbehaviors: [],
      alarm_id: 'alarm-id',
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectAllEntities(wrapper);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityInvalidate);

    await flushPromises();

    expect(bulkCreateAlarmAckEvent).toHaveBeenCalledWith(expect.any(Object), {
      data: [
        {
          _id: entity.alarm_id,
          comment: WEATHER_ACK_EVENT_OUTPUT.ack,
        },
      ],
    });

    expect(bulkCreateAlarmCancelEvent).toHaveBeenCalledWith(expect.any(Object), {
      data: [
        {
          _id: entity.alarm_id,
          comment: WEATHER_ACK_EVENT_OUTPUT.validateCancel,
        },
      ],
    });
  });

  test('Pause action applied after trigger mass pause action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
      type: ENTITY_TYPES.service,
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectAllEntities(wrapper);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityPause);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createEntityPauseEvent,
      config: {
        withInherited: true,
        action: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];
    const comment = Faker.datatype.string();
    const type = Faker.datatype.string();

    await config.action({ comment, type, inherited: true });

    expect(bulkPause).toHaveBeenCalledWith(expect.any(Object), {
      data: [
        {
          _id: entity._id,
          comment,
          type,
          inherited: true,
        },
      ],
    });
  });

  test('Play action applied after trigger mass play action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      pause: { timestamp: 123 },
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectAllEntities(wrapper);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityPlay);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createEntityResumeEvent,
      config: { action: expect.any(Function) },
    });

    const [{ config }] = $modals.show.mock.calls[0];
    const comment = Faker.datatype.string();

    await config.action({ comment });

    expect(bulkResume).toHaveBeenCalledWith(expect.any(Object), {
      data: [{ _id: entity._id, comment }],
    });
  });

  test('Cancel action applied after trigger mass cancel action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      alarm_display_name: 'alarm_display_name',
      status: {},
      pbehaviors: [],
      alarm_id: 'alarm-id',
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectAllEntities(wrapper);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityCancel);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.textFieldEditor,
      config: {
        title: 'common.note',
        action: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];

    const comment = Faker.datatype.string();

    await config.action(comment);

    expect(bulkCreateAlarmCancelEvent).toHaveBeenCalledWith(expect.any(Object), {
      data: [{ _id: entity.alarm_id, comment }],
    });
  });

  test('Comment action applied after trigger mass comment action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      pbehaviors: [],
      alarm_id: 'alarm-id',
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectAllEntities(wrapper);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityComment);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createCommentEvent,
      config: {
        action: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];

    const comment = Faker.datatype.string();

    await config.action({ comment });

    expect(bulkCreateAlarmCommentEvent).toHaveBeenCalledWith(expect.any(Object), {
      data: [{ _id: entity.alarm_id, comment }],
    });
  });

  test('Declare ticket action applied after trigger mass declare ticket action', async () => {
    bulkCreateDeclareTicketExecution.mockResolvedValueOnce([
      {
        id: 'execution-id',
        item: { _id: 'ticket-id', alarms: [] },
        status: 200,
      },
    ]);
    const rule = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
    };

    const byRules = {
      [rule._id]: {
        name: rule.name,
        alarms: [],
      },
    };
    const byAlarms = {};
    fetchAssignedDeclareTicketsWithoutStore.mockResolvedValueOnce({
      by_rules: byRules,
      by_alarms: byAlarms,
    });
    const entity = {
      _id: Faker.datatype.string(),
      pbehaviors: [],
    };

    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectAllEntities(wrapper);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityDeclareTicket);

    await flushPromises();

    expect(fetchAssignedDeclareTicketsWithoutStore).toBeCalled();

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createDeclareTicketEvent,
      config: {
        alarmsByTickets: byRules,
        ticketsByAlarms: {},
        items: [{}],
        action: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];

    const events = [{ _id: rule._id }];

    $modals.show.mockReset();
    await config.action(events);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.executeDeclareTickets,
      config: {
        executions: events,
        alarms: [{}],
        tickets: [rule],
        onExecute: expect.any(Function),
      },
    });
  });

  test('Unavailable action not applied', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      pbehaviors: [],
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectAllEntities(wrapper);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityAck);

    expect(bulkCreateAlarmAckEvent).toHaveBeenCalledWith(expect.any(Object), { data: [] });
  });

  test('Unavailable action flag removed after trigger', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      pbehaviors: [],
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectServiceEntityByIndex(wrapper, 0).triggerCustomEvent('remove:unavailable');

    expect(wrapper.vm.unavailableEntitiesAction).toEqual({
      [entity._id]: false,
    });
  });

  test('Entities list refreshed after trigger refresh on entity element', async () => {
    const wrapper = factory({
      propsData: {
        serviceEntities,
      },
    });

    selectServiceEntityByIndex(wrapper, 0).vm.$emit('refresh');
    await flushPromises();

    expect(wrapper).toHaveBeenEmit('refresh');
  });

  test('Page updated after trigger pagination', async () => {
    const options = {
      itemsPerPage: 10,
      page: 1,
    };
    const wrapper = factory({
      propsData: {
        serviceEntities,
        totalItems: 20,
        options,
      },
    });

    const newPage = 2;
    selectTablePagination(wrapper).vm.$emit('update:page', newPage);
    await flushPromises();

    expect(wrapper).toEmit('update:options', {
      ...options,
      page: newPage,
    });
  });

  test('Records per page updated after trigger pagination', async () => {
    const options = {
      itemsPerPage: 10,
      page: 1,
    };
    const wrapper = factory({
      propsData: {
        serviceEntities,
        totalItems: 20,
        options,
      },
    });

    const newItemsPerPage = 11;
    selectTablePagination(wrapper).vm.$emit('update:items-per-page', newItemsPerPage);
    await flushPromises();

    expect(wrapper).toEmit('update:options', {
      ...options,
      itemsPerPage: newItemsPerPage,
    });
  });

  test('Selected entities cleared after service entities updated', async () => {
    const wrapper = factory({
      propsData: {
        serviceEntities,
      },
    });

    selectServiceEntityByIndex(wrapper, 1).vm.$emit('update:selected', true);
    await flushPromises();

    expect(wrapper.vm.selectedEntities).toEqual([serviceEntities[1]]);

    await wrapper.setProps({
      serviceEntities: [...serviceEntities],
    });

    expect(wrapper.vm.selectedEntities).toEqual([]);
  });
});
