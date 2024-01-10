import Faker from 'faker';

import { createVueInstance, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import {
  createAlarmModule,
  createAuthModule, createDeclareTicketModule,
  createMockedStoreModules,
  createPbehaviorTypesModule,
} from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';
import flushPromises from 'flush-promises';
import {
  ENTITIES_STATES,
  MODALS,
  PBEHAVIOR_TYPE_TYPES,
  WEATHER_ACK_EVENT_OUTPUT,
  WEATHER_ACTIONS_TYPES,
} from '@/constants';

import ServiceEntitiesList from '@/components/other/service/partials/service-entities-list.vue';

const localVue = createVueInstance();

const stubs = {
  'service-entity-actions': true,
  'service-entity': true,
  'c-table-pagination': true,
};

const selectEntityActions = wrapper => wrapper.find('service-entity-actions-stub');
const selectServiceEntityByIndex = (wrapper, index) => wrapper
  .findAll('service-entity-stub')
  .at(index);
const selectTablePagination = wrapper => wrapper.find('c-table-pagination-stub');
const selectCheckboxFunctional = wrapper => wrapper.find('v-checkbox-functional-stub');

const applyEntitiesAction = async (wrapper, type) => {
  const entityActions = selectEntityActions(wrapper);

  await entityActions.vm.$emit('apply', { type });
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
  const { authModule } = createAuthModule();
  const { alarmModule } = createAlarmModule();
  const {
    declareTicketRuleModule,
    bulkCreateDeclareTicketExecution,
    fetchAssignedDeclareTicketsWithoutStore,
  } = createDeclareTicketModule();
  const { pbehaviorTypesModule, fetchPbehaviorTypesListWithoutStore } = createPbehaviorTypesModule();

  const store = createMockedStoreModules([alarmModule, authModule, declareTicketRuleModule, pbehaviorTypesModule]);

  const applyAction = jest.fn();
  const refresh = jest.fn();

  const factory = generateShallowRenderer(ServiceEntitiesList, {
    store,
    localVue,
    stubs,
    mocks: {
      $modals,
    },
    propsData: {
      service,
      pagination: {},
    },
    listeners: {
      'apply:action': applyAction,
      refresh,
    },
  });

  const snapshotFactory = generateRenderer(ServiceEntitiesList, {
    store,
    localVue,
    stubs,
    propsData: {
      service,
      pagination: {},
    },
    listeners: {
      'apply:action': applyAction,
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

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);
    expect(wrapper.vm.selectedEntities).toEqual(serviceEntities);

    await selectCheckboxFunctional(wrapper).vm.$emit('change', false);
    expect(wrapper.vm.selectedEntities).toEqual([]);
  });

  test('Ack action applied after trigger mass ack action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      ack: null,
      pbehaviors: [],
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityAck);

    expect(wrapper).toEmit('apply:action', {
      actionType: WEATHER_ACTIONS_TYPES.entityAck,
      entities: [entity],
    });
  });

  test('Ack remove action applied after trigger mass ackremove  action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      ack: {},
      pbehaviors: [],
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityAckRemove);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.textFieldEditor,
        config: {
          title: 'Remove ack',
          field: {
            name: 'output',
            label: 'Note',
            validationRules: 'required',
          },
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];
    const output = Faker.datatype.string();

    modalArguments.config.action({ output });

    expect(wrapper).toEmit('apply:action', {
      actionType: WEATHER_ACTIONS_TYPES.entityAckRemove,
      payload: { output },
      entities: [entity],
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
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);

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

    const [modalArguments] = $modals.show.mock.calls[0];
    const event = {
      ticket: Faker.datatype.string(),
      ticket_url: Faker.datatype.string(),
    };

    modalArguments.config.action(event);

    expect(wrapper).toEmit('apply:action', {
      actionType: WEATHER_ACTIONS_TYPES.entityAssocTicket,
      payload: event,
      entities: [entity],
    });
  });

  test('Validate action applied after trigger mass validate action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      state: {
        val: ENTITIES_STATES.major,
      },
      pbehaviors: [],
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityValidate);

    expect(wrapper).toEmit('apply:action', {
      actionType: WEATHER_ACTIONS_TYPES.entityValidate,
      entities: [entity],
    });
  });

  test('Invalidate action applied after trigger mass invalidate action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      state: {
        val: ENTITIES_STATES.major,
      },
      pbehaviors: [],
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityInvalidate);

    expect(wrapper).toEmit('apply:action', {
      actionType: WEATHER_ACTIONS_TYPES.entityInvalidate,
      entities: [entity],
      payload: {
        output: WEATHER_ACK_EVENT_OUTPUT.ack,
      },
    });
  });

  test('Pause action applied after trigger mass pause action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      pbh_origin_icon: '',
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityPause);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createServicePauseEvent,
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];
    const comment = Faker.datatype.string();
    const reason = Faker.datatype.string();

    const pauseType = {
      type: PBEHAVIOR_TYPE_TYPES.pause,
    };

    fetchPbehaviorTypesListWithoutStore.mockResolvedValue({
      data: [pauseType],
    });

    await modalArguments.config.action({ comment, reason });

    expect(wrapper).toEmit('apply:action', {
      actionType: WEATHER_ACTIONS_TYPES.entityPause,
      entities: [entity],
      payload: {
        comment,
        reason,
        type: pauseType,
      },
    });
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
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityPlay);

    expect(wrapper).toEmit('apply:action', {
      actionType: WEATHER_ACTIONS_TYPES.entityPlay,
      entities: [entity],
    });
  });

  test('Cancel action applied after trigger mass cancel action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      alarm_display_name: 'alarm_display_name',
      status: {},
      pbehaviors: [],
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);

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

    const [modalArguments] = $modals.show.mock.calls[0];

    const output = Faker.datatype.string();

    await modalArguments.config.action(output);

    expect(wrapper).toEmit('apply:action', {
      actionType: WEATHER_ACTIONS_TYPES.entityCancel,
      entities: [entity],
      payload: { output },
    });
  });

  test('Comment action applied after trigger mass comment action', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      pbehaviors: [],
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityComment);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createCommentEvent,
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const output = Faker.datatype.string();

    await modalArguments.config.action({ output });

    expect(wrapper).toEmit('apply:action', {
      actionType: WEATHER_ACTIONS_TYPES.entityComment,
      entities: [entity],
      payload: { output },
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
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);

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
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectCheckboxFunctional(wrapper).vm.$emit('change', true);

    await applyEntitiesAction(wrapper, WEATHER_ACTIONS_TYPES.entityAck);

    expect(wrapper).toEmit('apply:action', {
      actionType: WEATHER_ACTIONS_TYPES.entityAck,
      entities: [],
    });
  });

  test('Unavailable action flag removed after trigger', async () => {
    const entity = {
      _id: Faker.datatype.string(),
      pbehaviors: [],
    };
    const wrapper = factory({
      propsData: {
        serviceEntities: [
          entity,
        ],
      },
    });

    await selectServiceEntityByIndex(wrapper, 0).vm.$emit('remove:unavailable');

    expect(wrapper.vm.unavailableEntitiesAction).toEqual({
      [entity._id]: false,
    });
  });

  test('Action applied after trigger apply on service entity', async () => {
    const wrapper = factory({
      propsData: {
        serviceEntities,
      },
    });

    await selectServiceEntityByIndex(wrapper, 0).vm.$emit('apply:action');

    expect(applyAction).toHaveBeenCalled();
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
    const pagination = {
      rowsPerPage: 10,
      page: 1,
    };
    const wrapper = factory({
      propsData: {
        serviceEntities,
        totalItems: 20,
        pagination,
      },
    });

    const newPage = 2;
    await selectTablePagination(wrapper).vm.$emit('update:page', newPage);

    expect(wrapper).toEmit('update:pagination', {
      ...pagination,
      page: newPage,
    });
  });

  test('Records per page updated after trigger pagination', async () => {
    const pagination = {
      rowsPerPage: 10,
      page: 1,
    };
    const wrapper = factory({
      propsData: {
        serviceEntities,
        totalItems: 20,
        pagination,
      },
    });

    const newRowsPerPage = 11;
    await selectTablePagination(wrapper).vm.$emit('update:rows-per-page', newRowsPerPage);

    expect(wrapper).toEmit('update:pagination', {
      ...pagination,
      rowsPerPage: newRowsPerPage,
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
