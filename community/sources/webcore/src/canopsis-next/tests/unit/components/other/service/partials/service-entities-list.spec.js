import Faker from 'faker';
import flushPromises from 'flush-promises';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import {
  createAlarmModule,
  createAuthModule,
  createDeclareTicketModule,
  createMockedStoreModules,
  createPbehaviorModule,
  createPbehaviorTypesModule,
} from '@unit/utils/store';
import { uid } from '@/helpers/uid';
import { mockDateNow, mockModals } from '@unit/utils/mock-hooks';
import {
  ENTITIES_STATES,
  MODALS,
  PBEHAVIOR_ORIGINS,
  PBEHAVIOR_TYPE_TYPES,
  USERS_PERMISSIONS,
  WEATHER_ACK_EVENT_OUTPUT,
  WEATHER_ACTIONS_TYPES,
  WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE,
} from '@/constants';
import { COLORS } from '@/config';

import ServiceEntitiesList from '@/components/other/service/partials/service-entities-list.vue';
import { createCheckboxInputStub } from '@unit/stubs/input';

jest.mock('@/helpers/uid');

const stubs = {
  'service-entity-actions': true,
  'service-entity': true,
  'c-table-pagination': true,
  'v-simple-checkbox': createCheckboxInputStub('v-simple-checkbox'),
};

const selectEntityActions = wrapper => wrapper.find('service-entity-actions-stub');
const selectServiceEntityByIndex = (wrapper, index) => wrapper
  .findAll('service-entity-stub')
  .at(index);
const selectTablePagination = wrapper => wrapper.find('c-table-pagination-stub');
const selectCheckbox = wrapper => wrapper.find('.v-simple-checkbox');

const applyEntitiesAction = async (wrapper, type) => {
  const entityActions = selectEntityActions(wrapper);

  await entityActions.vm.$emit('apply', { type });
};

describe('service-entities-list', () => {
  const nowTimestamp = 1386435600000;
  mockDateNow(nowTimestamp);

  const $modals = mockModals();
  const service = {
    _id: 'service-id',
  };
  const serviceEntities = [
    {
      _id: 'service-entity-1-id',
      pbehaviors: [],
      state: {
        val: ENTITIES_STATES.major,
      },
      alarm_id: 'alarm-id',
    },
    {
      _id: 'service-entity-2-id',
      pbehaviors: [{
        type: {
          type: PBEHAVIOR_TYPE_TYPES.pause,
        },
      }],
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
  const { pbehaviorTypesModule, fetchPbehaviorTypesListWithoutStore } = createPbehaviorTypesModule();
  const { pbehaviorModule, createEntityPbehaviors, removeEntityPbehaviors } = createPbehaviorModule();

  const store = createMockedStoreModules([
    alarmModule,
    authModule,
    declareTicketRuleModule,
    pbehaviorTypesModule,
    pbehaviorModule,
  ]);

  const refresh = jest.fn();

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

  const snapshotFactory = generateRenderer(ServiceEntitiesList, {
    store,

    stubs,
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

    await selectServiceEntityByIndex(wrapper, 1).vm.$emit('update:selected', true);

    expect(wrapper.vm.selectedEntities).toEqual([serviceEntities[1]]);

    await selectServiceEntityByIndex(wrapper, 1).vm.$emit('update:selected', false);

    expect(wrapper.vm.selectedEntities).toEqual([]);
  });

  test('All entities selected and unselected after trigger checkbox', async () => {
    const wrapper = factory({
      propsData: {
        serviceEntities,
      },
    });

    await selectCheckbox(wrapper).vm.$emit('change', true);
    expect(wrapper.vm.selectedEntities).toEqual(serviceEntities);

    await selectCheckbox(wrapper).vm.$emit('change', false);
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

    await selectCheckbox(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityAckRemove);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.textFieldEditor,
        config: {
          title: 'Remove ack',
          field: {
            name: 'comment',
            label: 'Note',
            validationRules: 'required',
          },
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];
    const comment = Faker.datatype.string();

    await config.action({ comment });

    expect(bulkCreateAlarmAckremoveEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: entity.alarm_id,
          comment,
        }],
      },
      undefined,
    );
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

    await selectCheckbox(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityAssocTicket);

    await flushPromises();

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createAssociateTicketEvent,
        config: {
          items: [{}],
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];
    const event = {
      ticket: Faker.datatype.string(),
      ticket_url: Faker.datatype.string(),
    };

    await config.action(event);

    expect(bulkCreateAlarmAssocticketEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: entity.alarm_id,
          ...event,
        }],
      },
      undefined,
    );
  });

  test('Validate action applied after trigger mass validate action', async () => {
    currentUserPermissionsById.mockReturnValueOnce(({
      [USERS_PERMISSIONS.business.serviceWeather.actions.entityValidate]: {
        actions: [],
      },
    }));

    const entity = {
      _id: Faker.datatype.string(),
      state: {
        val: ENTITIES_STATES.major,
      },
      pbehaviors: [],
      alarm_id: 'alarm-id',
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectCheckbox(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityValidate);
    await flushPromises();

    expect(bulkCreateAlarmAckEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: entity.alarm_id,
          comment: WEATHER_ACK_EVENT_OUTPUT.validateOk,
        }],
      },
      undefined,
    );

    expect(bulkCreateAlarmChangestateEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: entity.alarm_id,
          comment: WEATHER_ACK_EVENT_OUTPUT.validateOk,
          state: ENTITIES_STATES.critical,
        }],
      },
      undefined,
    );
  });

  test('Invalidate action applied after trigger mass invalidate action', async () => {
    currentUserPermissionsById.mockReturnValueOnce(({
      [USERS_PERMISSIONS.business.serviceWeather.actions.entityInvalidate]: {
        actions: [],
      },
    }));

    const entity = {
      _id: Faker.datatype.string(),
      state: {
        val: ENTITIES_STATES.major,
      },
      pbehaviors: [],
      alarm_id: 'alarm-id',
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectCheckbox(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityInvalidate);

    await flushPromises();

    expect(bulkCreateAlarmAckEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: entity.alarm_id,
          comment: WEATHER_ACK_EVENT_OUTPUT.ack,
        }],
      },
      undefined,
    );

    expect(bulkCreateAlarmCancelEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: entity.alarm_id,
          comment: WEATHER_ACK_EVENT_OUTPUT.validateCancel,
        }],
      },
      undefined,
    );
  });

  test('Pause action applied after trigger mass pause action', async () => {
    const uniqId = Faker.datatype.string();

    uid.mockReturnValueOnce(uniqId);
    const entity = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
      pbh_origin_icon: '',
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectCheckbox(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityPause);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createServicePauseEvent,
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];
    const comment = Faker.datatype.string();
    const reason = Faker.datatype.string();

    const pauseType = {
      type: PBEHAVIOR_TYPE_TYPES.pause,
    };

    fetchPbehaviorTypesListWithoutStore.mockResolvedValue({
      data: [pauseType],
    });

    await config.action({ comment, reason });

    expect(createEntityPbehaviors).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          color: COLORS.secondary,
          comment,
          reason,
          comments: [],
          enabled: true,
          entity: entity._id,
          exceptions: [],
          exdates: [],
          name: `${WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE}-${entity.name}-${uniqId}`,
          origin: PBEHAVIOR_ORIGINS.serviceWeather,
          tstart: nowTimestamp / 1000,
          tstop: null,
          type: undefined,
        }],
      },
      undefined,
    );
  });

  test('Play action applied after trigger mass play action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      pbehaviors: [{
        type: {
          type: PBEHAVIOR_TYPE_TYPES.pause,
        },
      }],
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [entity],
      },
    });

    await selectCheckbox(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityPlay);

    expect(removeEntityPbehaviors).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          entity: entity._id,
          origin: PBEHAVIOR_ORIGINS.serviceWeather,
        }],
      },
      undefined,
    );
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

    await selectCheckbox(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityCancel);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.textFieldEditor,
        config: {
          title: 'Note',
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const comment = Faker.datatype.string();

    await config.action(comment);

    expect(bulkCreateAlarmCancelEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{ _id: entity.alarm_id, comment }],
      },
      undefined,
    );
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

    await selectCheckbox(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityComment);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createCommentEvent,
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const comment = Faker.datatype.string();

    await config.action({ comment });

    expect(bulkCreateAlarmCommentEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{ _id: entity.alarm_id, comment }],
      },
      undefined,
    );
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

    await selectCheckbox(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.declareTicket);

    await flushPromises();

    expect(fetchAssignedDeclareTicketsWithoutStore).toBeCalled();

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createDeclareTicketEvent,
        config: {
          alarmsByTickets: byRules,
          ticketsByAlarms: {},
          items: [{}],
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const events = [{ _id: rule._id }];

    $modals.show.mockReset();
    await config.action(events);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.executeDeclareTickets,
        config: {
          executions: events,
          alarms: [{}],
          tickets: [rule],
          onExecute: expect.any(Function),
        },
      },
    );
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

    await selectCheckbox(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityAck);

    expect(bulkCreateAlarmAckEvent).toBeCalledWith(
      expect.any(Object),
      { data: [] },
      undefined,
    );
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

    await selectServiceEntityByIndex(wrapper, 0).vm.$emit('remove:unavailable');

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

    await selectServiceEntityByIndex(wrapper, 0).vm.$emit('refresh');

    expect(refresh).toHaveBeenCalled();
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
    await selectTablePagination(wrapper).vm.$emit('update:page', newPage);

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
    await selectTablePagination(wrapper).vm.$emit('update:items-per-page', newItemsPerPage);

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

    await selectServiceEntityByIndex(wrapper, 1).vm.$emit('update:selected', true);

    expect(wrapper.vm.selectedEntities).toEqual([serviceEntities[1]]);

    await wrapper.setProps({
      serviceEntities: [...serviceEntities],
    });

    expect(wrapper.vm.selectedEntities).toEqual([]);
  });

  test('Renders `service-entities-list` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-entities-list` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        serviceEntities,
        entityNameField: 'custom-entity-name',
        widgetParameters: {},
        totalItems: serviceEntities.length,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-entities-list` with selected entities', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        serviceEntities,
        entityNameField: 'custom-entity-name',
        widgetParameters: {},
        totalItems: serviceEntities.length,
      },
    });

    const firstEntity = selectServiceEntityByIndex(wrapper, 0);

    await firstEntity.vm.$emit('update:selected', true);

    expect(wrapper).toMatchSnapshot();
  });
});
