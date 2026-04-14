import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import {
  createAlarmModule,
  createAlarmDetailsModule,
  createAuthModule,
  createDeclareTicketModule,
  createMetaAlarmModule,
  createMockedStoreModules,
  createPbehaviorModule,
  createPbehaviorTypesModule,
  createPbehaviorReasonModule,
} from '@unit/utils/store';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';

import {
  ALARM_LIST_ACTIONS_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  ALARM_STATES,
  ALARM_STATUSES,
  ENTITY_PATTERN_FIELDS,
  INSTRUCTION_EXECUTION_ICONS,
  META_ALARMS_RULE_TYPES,
  MODALS,
  PATTERN_CONDITIONS,
  PBEHAVIOR_ORIGINS,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES,
  REMEDIATION_INSTRUCTION_TYPES,
  TIME_UNITS,
} from '@/constants';

import { featuresService } from '@/services/features';

import { generateDefaultAlarmListWidget } from '@/helpers/entities/widget/form';
import { prepareAlarmListWidget } from '@/helpers/entities/widget/forms/alarm';
import { exportAlarmToPdf } from '@/helpers/file/pdf';

import ActionsPanel from '@/components/widgets/alarm/actions/actions-panel.vue';

jest.mock('@/helpers/file/pdf', () => {
  const original = jest.requireActual('@/helpers/file/pdf');
  return {
    ...original,

    exportAlarmToPdf: jest.fn(),
  };
});

jest.mock('@/helpers/async', () => ({
  promisedTimeout: callback => (callback ? Promise.resolve(callback()) : Promise.resolve()),
}));

const stubs = {
  'shared-actions-panel': {
    props: ['actions'],
    template: `
      <div class="shared-actions-panel">
        <button
          v-for="action in actions"
          :class="'action-' + action.type"
          :disabled="action.disabled"
          @click="action.method"
        >{{ action.title }}|{{ action.icon }}|{{ action.type }}</button>
      </div>
    `,
  },
};

const selectActionByType = (wrapper, type) => wrapper.find(`.action-${type}`);

describe('actions-panel', () => {
  const timestamp = 1386435600000;
  jest.useFakeTimers({ now: timestamp });

  const $modals = mockModals();
  const $popups = mockPopups();

  const { authModule } = createAuthModule();
  const authModuleWithAccess = {
    ...authModule,
    getters: {
      currentUserPermissionsById: Object.values(ALARM_LIST_ACTIONS_TYPES)
        .reduce((acc, type) => ({
          ...acc,
          [BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[type]]: { actions: [] },
        }), {}),
    },
  };

  const {
    alarmModule,
    bulkCreateAlarmAckEvent,
    bulkCreateAlarmAckremoveEvent,
    bulkCreateAlarmSnoozeEvent,
    bulkCreateAlarmAssocticketEvent,
    bulkCreateAlarmTicketremoveEvent,
    bulkCreateAlarmCommentEvent,
    bulkCreateAlarmCancelEvent,
    bulkCreateAlarmChangestateEvent,
    addBookmarkToAlarm,
    removeBookmarkFromAlarm,
  } = createAlarmModule();
  const { metaAlarmModule, removeAlarmsFromMetaAlarm } = createMetaAlarmModule();
  const { alarmDetailsModule, fetchAlarmDetailsWithoutStore } = createAlarmDetailsModule();

  const {
    declareTicketRuleModule,
    fetchAssignedDeclareTicketsWithoutStore,
  } = createDeclareTicketModule();
  const {
    pbehaviorModule,
    createEntityPbehaviors,
    removeEntityPbehaviors,
  } = createPbehaviorModule();
  const { pbehaviorTypesModule } = createPbehaviorTypesModule();
  const { pbehaviorReasonModule } = createPbehaviorReasonModule();

  const store = createMockedStoreModules([
    metaAlarmModule,
    authModule,
    alarmModule,
    alarmDetailsModule,
    declareTicketRuleModule,
    pbehaviorModule,
    pbehaviorTypesModule,
    pbehaviorReasonModule,
  ]);

  const assignedInstructions = [
    {
      _id: 1,
      name: 'Running instruction',
      execution: {
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running,
      },
    },
    {
      _id: 2,
      name: 'New instruction',
      execution: null,
    },
    {
      _id: 3,
      name: 'Paused instruction',
      execution: {
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused,
      },
    },
  ];

  const assignedInstructionsWithPaused = [
    {
      _id: 1,
      name: 'New instruction',
      type: REMEDIATION_INSTRUCTION_TYPES.manual,
      execution: null,
    },
    {
      _id: 2,
      name: 'Paused instruction',
      type: REMEDIATION_INSTRUCTION_TYPES.manual,
      execution: {
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused,
      },
    },
  ];

  const assignedDeclareTicketRules = [
    {
      _id: 1,
      name: 'Name 1',
    },
    {
      _id: 2,
      name: 'Name 2',

    },
    {
      _id: 3,
      name: 'Name 3',
    },
  ];

  const alarm = {
    _id: 'alarm-id',
    assigned_instructions: assignedInstructions,
    assigned_declare_ticket_rules: assignedDeclareTicketRules,
    entity: {},
    v: {
      ack: {},
      status: {
        val: ALARM_STATUSES.flapping,
      },
      state: {},
    },
  };

  const widget = {
    parameters: {
      isMultiAckEnabled: true,
    },
  };

  const parentAlarm = {
    meta_alarm_rule: {
      type: META_ALARMS_RULE_TYPES.manualgroup,
    },
    d: 'parent-d',
  };

  const refreshAlarmsList = jest.fn();

  const factory = generateShallowRenderer(ActionsPanel, {
    stubs,
    mocks: { $modals, $popups },
    provide: {
      $system: {},
    },
  });
  const snapshotFactory = generateRenderer(ActionsPanel, {
    stubs,
    provide: {
      $system: {},
    },
  });

  test('quickActions returns [] if widget.parameters.quickActions is undefined', () => {
    const wrapper = factory({
      store,
      propsData: {
        item: alarm,
        widget: { parameters: {} },
        parentAlarm,
      },
    });
    expect(wrapper.vm.quickActions).toEqual([]);
  });

  test('quickActions returns array if widget.parameters.quickActions is set', () => {
    const quickActions = [
      ALARM_LIST_ACTIONS_TYPES.ack,
      ALARM_LIST_ACTIONS_TYPES.fastAck,
      ALARM_LIST_ACTIONS_TYPES.cancel,
    ];
    const wrapper = factory({
      store,
      propsData: {
        item: alarm,
        widget: { parameters: { quickActions } },
        parentAlarm,
      },
    });
    expect(wrapper.vm.quickActions).toEqual(quickActions);
  });

  test('preparedActions sorts actions according to quickActions order', () => {
    const quickActions = [
      ALARM_LIST_ACTIONS_TYPES.cancel,
      ALARM_LIST_ACTIONS_TYPES.ack,
      ALARM_LIST_ACTIONS_TYPES.fastAck,
    ];
    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: { parameters: { quickActions, isMultiAckEnabled: true } },
        parentAlarm,
      },
    });
    const types = wrapper.vm.preparedActions.map(a => a.type);
    // The first three actions should be in quickActions order
    expect(types.slice(0, 3)).toEqual(quickActions);
  });

  test('additionalProps.inlineCount reflects quickActions count', () => {
    const quickActions = [
      ALARM_LIST_ACTIONS_TYPES.ack,
      ALARM_LIST_ACTIONS_TYPES.fastAck,
      ALARM_LIST_ACTIONS_TYPES.cancel,
    ];
    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: { parameters: { quickActions, isMultiAckEnabled: true } },
        parentAlarm,
      },
    });
    // getActionsInlineCount returns quickActions.length + 1 (menu button) if not all actions are quick
    expect(wrapper.vm.additionalProps.inlineCount).toBeGreaterThanOrEqual(quickActions.length);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test('Ack modal showed after trigger ack action', async () => {
    const isNoteRequired = Faker.datatype.boolean();
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        isMultiAckEnabled: true,
        isAckNoteRequired: isNoteRequired,
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.ack).trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createAckEvent,
        config: {
          isNoteRequired,
          items: [alarm],
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.action({ output: 'OUTPUT', ack_resources: true }, { needDeclareTicket: false, needAssociateTicket: false });

    await flushPromises();

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  test('Fast ack event sent after trigger fast ack action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        fastAckOutput: {
          enabled: true,
          value: 'Output',
        },
      },
    };
    const fastActionAlarm = {
      ...alarm,
      entity: {
        type: 'entity-type',
      },
      v: {
        connector: 'alarm-connector',
        connector_name: 'alarm-connector-name',
        component: 'alarm-component',
        resource: 'alarm-resource',
        status: {
          val: ALARM_STATUSES.ongoing,
        },
        state: {
          val: 'state-val',
        },
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: fastActionAlarm,
        widget: widgetData,
        parentAlarm,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastAck).trigger('click');

    expect(bulkCreateAlarmAckEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          comment: widgetData.parameters.fastAckOutput.value,
        }],
      },
    );
  });

  test('Ack remove modal showed after trigger ack remove action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.ackRemove).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createEvent,
        config: {
          title: 'Remove ack',
          items: [alarm],
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const comment = Faker.datatype.string();

    await config.action({ comment });

    expect(bulkCreateAlarmAckremoveEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          comment,
        }],
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Create pbehavior modal showed after trigger pbehavior add action', () => {
    const entity = {
      _id: Faker.datatype.string(),
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: { ...alarm, entity },
        widget,
        parentAlarm,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: [[{
            field: ENTITY_PATTERN_FIELDS.id,
            cond: {
              type: PATTERN_CONDITIONS.equal,
              value: entity._id,
            },
          }]],
          entities: [entity],
          afterSubmit: expect.any(Function),
        },
      },
    );
  });

  test('Fast pbehavior add creates downtime pbehavior and calls refreshAlarmsList', async () => {
    const entity = {
      _id: Faker.datatype.string(),
    };
    const typeId = Faker.datatype.string();
    const reasonId = Faker.datatype.string();
    const alarmWithEntity = { ...alarm, entity };
    const widgetWithFastPbehavior = {
      parameters: {
        isMultiAckEnabled: true,
        fast_pbehaviors: [{ type: typeId, reason: reasonId, name_prefix: 'Test' }],
      },
    };

    createEntityPbehaviors.mockResolvedValue([{}]);

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        pbehaviorModule,
        pbehaviorTypesModule,
        pbehaviorReasonModule,
      ]),
      propsData: {
        item: alarmWithEntity,
        widget: widgetWithFastPbehavior,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd).trigger('click');

    await flushPromises();

    expect(createEntityPbehaviors).toHaveBeenCalledWith(
      expect.any(Object),
      expect.objectContaining({
        data: expect.arrayContaining([
          expect.objectContaining({
            entity: entity._id,
            type: typeId,
            reason: reasonId,
            origin: PBEHAVIOR_ORIGINS.alarmList,
          }),
        ]),
      }),
    );
    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Fast pbehavior remove removes downtime pbehavior and calls refreshAlarmsList', async () => {
    const entity = {
      _id: Faker.datatype.string(),
    };
    const alarmWithFastPbehavior = {
      ...alarm,
      entity,
      pbh_origin_icon: true,
    };

    removeEntityPbehaviors.mockResolvedValue([{}]);

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        pbehaviorModule,
        pbehaviorTypesModule,
        pbehaviorReasonModule,
      ]),
      propsData: {
        item: alarmWithFastPbehavior,
        widget,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastPbehaviorRemove).trigger('click');

    await flushPromises();

    expect(removeEntityPbehaviors).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: [
          { origin: PBEHAVIOR_ORIGINS.alarmList, entity: entity._id },
        ],
      },
    );
    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Snooze modal showed after trigger snooze action', async () => {
    const isNoteRequired = Faker.datatype.boolean();
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        isSnoozeNoteRequired: isNoteRequired,
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.snooze).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createSnoozeEvent,
        config: {
          isNoteRequired,
          items: [alarm],
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const snoozeEvent = {
      duration: {
        unit: TIME_UNITS.minute,
        value: Faker.datatype.number(),
      },
      comment: Faker.datatype.string(),
    };

    await config.action(snoozeEvent);

    expect(bulkCreateAlarmSnoozeEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          ...snoozeEvent,
        }],
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Declare ticket modal showed after trigger declare action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };
    const rule = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
    };

    const byRules = {
      [rule._id]: {
        name: rule.name,
        alarms: [alarm._id],
      },
    };
    const byAlarms = {
      [alarm._id]: [rule._id],
    };

    fetchAssignedDeclareTicketsWithoutStore.mockResolvedValueOnce({
      by_rules: byRules,
      by_alarms: byAlarms,
    });

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        declareTicketRuleModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.declareTicket).trigger('click');

    await flushPromises();

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createDeclareTicketEvent,
        config: {
          items: [alarm],
          alarmsByTickets: byRules,
          ticketsByAlarms: byAlarms,
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const events = [{ _id: rule._id, alarms: [Faker.datatype.string()] }];

    $modals.show.mockReset();
    config.action(events);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.executeDeclareTickets,
      config: {
        executions: events,
        alarms: [alarm],
        tickets: [rule],
        onExecute: expect.any(Function),
      },
    });
  });

  test('Associate ticket modal showed after trigger associate ticket action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.associateTicket).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createAssociateTicketEvent,
        config: {
          items: [alarm],
          ignoreAck: false,
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const ticketEvent = {
      comment: Faker.datatype.string(),
    };

    await config.action(ticketEvent);

    expect(bulkCreateAlarmAssocticketEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          ...ticketEvent,
        }],
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Remove associated ticket modal showed after trigger remove associated ticket action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const alarmWithTickets = {
      ...alarm,
      v: {
        ...alarm.v,
        tickets: [
          {
            ticket: 'TICKET-123',
            ticket_system_name: 'Jira',
          },
        ],
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarmWithTickets,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.removeAssociatedTicket).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.removeAssociatedTicketEvent,
        config: {
          items: [alarmWithTickets],
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const removeTicketEvent = {
      ticket: 'TICKET-123',
      reason: Faker.datatype.string(),
    };

    await config.action(removeTicketEvent);

    expect(bulkCreateAlarmTicketremoveEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarmWithTickets._id,
          ...removeTicketEvent,
        }],
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Change state modal showed after trigger change state action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.changeState).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createChangeStateEvent,
        config: {
          items: [alarm],
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const changeStateEvent = {
      state: ALARM_STATES.critical,
      comment: Faker.datatype.string(),
    };

    await config.action(changeStateEvent);

    expect(bulkCreateAlarmChangestateEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          ...changeStateEvent,
        }],
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Cancel modal showed after trigger cancel action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.cancel).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createEvent,
        config: {
          items: [alarm],
          action: expect.any(Function),
          title: 'Cancel',
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const cancelEvent = {
      comment: Faker.datatype.string(),
    };

    await config.action(cancelEvent);

    expect(bulkCreateAlarmCancelEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          ...cancelEvent,
        }],
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Fast cancel event sent after trigger fast cancel action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        fastCancelOutput: {
          enabled: true,
          value: 'Output',
        },
      },
    };
    const fastActionAlarm = {
      ...alarm,
      entity: {
        type: 'entity-type',
      },
      v: {
        connector: 'alarm-connector',
        connector_name: 'alarm-connector-name',
        component: 'alarm-component',
        resource: 'alarm-resource',
        status: {
          val: ALARM_STATUSES.ongoing,
        },
        state: {
          val: 'state-val',
        },
        ack: {},
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: fastActionAlarm,
        widget: widgetData,
        parentAlarm,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastCancel).trigger('click');

    await flushPromises();

    expect(bulkCreateAlarmCancelEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          comment: 'Output',
        }],
      },
    );
  });

  test('Fast cancel and cancel action available when status is flapping and state is ok', async () => {
    const flappingAlarm = {
      ...alarm,
      v: {
        connector: 'alarm-connector',
        connector_name: 'alarm-connector-name',
        component: 'alarm-component',
        resource: 'alarm-resource',
        status: {
          val: ALARM_STATUSES.flapping,
        },
        state: {
          val: ALARM_STATES.ok,
        },
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: flappingAlarm,
        widget,
        parentAlarm,
      },
    });

    expect(selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastCancel).exists()).toBeTruthy();
    expect(selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.cancel).exists()).toBeTruthy();
  });

  test('Fast cancel and cancel action available when status is closed and state is ok', async () => {
    const flappingAlarm = {
      ...alarm,
      v: {
        connector: 'alarm-connector',
        connector_name: 'alarm-connector-name',
        component: 'alarm-component',
        resource: 'alarm-resource',
        resolved: null,
        status: {
          val: ALARM_STATUSES.closed,
        },
        state: {
          val: ALARM_STATES.ok,
        },
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: flappingAlarm,
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,
            isActionsAllowWithOkState: true,
          },
        },
        parentAlarm,
      },
    });

    expect(selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastCancel).exists()).toBeTruthy();
    expect(selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.cancel).exists()).toBeTruthy();
  });

  test('Fast cancel and cancel action doesn\'t available when status is closed and state is ok and isActionsAllowWithOkState is disabled', async () => {
    const flappingAlarm = {
      ...alarm,
      v: {
        connector: 'alarm-connector',
        connector_name: 'alarm-connector-name',
        component: 'alarm-component',
        resource: 'alarm-resource',
        resolved: null,
        status: {
          val: ALARM_STATUSES.closed,
        },
        state: {
          val: ALARM_STATES.ok,
        },
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: flappingAlarm,
        widget,
        parentAlarm,
      },
    });

    expect(selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastCancel).exists()).toBeFalsy();
    expect(selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.cancel).exists()).toBeFalsy();
  });

  test('Variables modal showed after trigger variables help action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };
    const entity = {
      _id: Faker.datatype.string(),
    };
    const pbehavior = {
      _id: Faker.datatype.string(),
    };
    const alarmData = {
      _id: Faker.datatype.string(),
      entity,
      pbehavior,
      v: {
        status: {
          val: ALARM_STATUSES.closed,
        },
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: {
          ...alarmData,
          v: {
            status: {
              val: ALARM_STATUSES.closed,
            },
          },
        },
        widget: widgetData,
        parentAlarm,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.variablesHelp).trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.variablesHelp,
        config: {
          exportEntity: alarmData,
          exportEntityName: `alarm-${alarmData._id}`,
          variables: [
            {
              name: 'alarm',
              original: alarmData,
              children: [
                { name: '_id', path: 'alarm._id', value: alarmData._id },
                {
                  name: 'v',
                  children: [{
                    name: 'status',
                    children: [
                      {
                        name: 'val',
                        path: 'alarm.v.status.val',
                        value: 0,
                      },
                    ],
                  }],
                },
              ],
            },
            {
              name: 'entity',
              children: [{ name: '_id', path: 'entity._id', value: entity._id }],
            },
            {
              name: 'pbehavior',
              children: [{ name: '_id', path: 'pbehavior._id', value: pbehavior._id }],
            },
          ],
        },
      },
    );
  });

  test('History modal showed after trigger history action', () => {
    const widgetData = prepareAlarmListWidget({
      _id: Faker.datatype.string(),
      parameters: {
        widgetColumns: [
          {
            value: Faker.datatype.string(),
            label: Faker.datatype.string(),
          },
        ],
      },
    });

    const entity = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: { ...alarm, entity },
        widget: widgetData,
        parentAlarm,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.history).trigger('click');

    const defaultWidget = prepareAlarmListWidget(generateDefaultAlarmListWidget());

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.alarmsList,
        config: {
          title: `${entity._id} - alarm list`,
          fetchList: expect.any(Function),
          widget: {
            ...defaultWidget,

            _id: expect.any(String),
            parameters: {
              ...defaultWidget.parameters,

              widgetColumns: widgetData.parameters.widgetColumns,
              widgetGroupColumns: widgetData.parameters.widgetGroupColumns,
              serviceDependenciesColumns: widgetData.parameters.serviceDependenciesColumns,
            },
          },
        },
      },
    );
  });

  test('Comment modal showed after trigger comment action', async () => {
    const commentAlarm = {
      ...alarm,
      v: {
        connector: 'alarm-connector',
        connector_name: 'alarm-connector-name',
        component: 'alarm-component',
        resource: 'alarm-resource',
        status: {
          val: ALARM_STATUSES.ongoing,
        },
        state: {
          val: 'state-val',
        },
      },
    };
    const widgetData = {
      _id: Faker.datatype.string(),
      comment_templates: [],
      parameters: {
        comment_templates: [],
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: commentAlarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.comment).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createCommentEvent,
        config: {
          items: [commentAlarm],
          templates: [],
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const comment = Faker.datatype.string();

    await config.action({ comment });

    expect(bulkCreateAlarmCommentEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          comment,
        }],
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Remove alarms from manual meta alarm modal showed after trigger remove alarms from manual meta alarm action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        metaAlarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.removeAlarmsFromManualMetaAlarm).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.removeAlarmsFromMetaAlarm,
        config: {
          items: [alarm],
          action: expect.any(Function),
          title: 'Unlink alarm from manual meta alarm',
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const newRemoveAlarmsEvent = {
      comment: Faker.datatype.string(),
      alarms: [Faker.datatype.string()],
    };

    await config.action(newRemoveAlarmsEvent);

    expect(removeAlarmsFromMetaAlarm).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: parentAlarm?._id,
        data: newRemoveAlarmsEvent,
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Remove alarms from auto meta alarm modal showed after trigger remove alarms from auto meta alarm action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        metaAlarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm: {
          meta_alarm_rule: {
            type: META_ALARMS_RULE_TYPES.attribute,
          },
          d: 'parent-d',
        },
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.removeAlarmsFromAutoMetaAlarm).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.removeAlarmsFromMetaAlarm,
        config: {
          items: [alarm],
          action: expect.any(Function),
          title: 'Unlink alarm from meta alarm',
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const newRemoveAlarmsEvent = {
      comment: Faker.datatype.string(),
      alarms: [Faker.datatype.string()],
    };

    await config.action(newRemoveAlarmsEvent);

    expect(removeAlarmsFromMetaAlarm).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: parentAlarm?._id,
        data: newRemoveAlarmsEvent,
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Execute instruction alarm modal showed after trigger execute instruction action', () => {
    const assignedInstruction = assignedInstructionsWithPaused[1];
    const alarmData = {
      ...alarm,
      _id: Faker.datatype.string(),
      assigned_instructions: [assignedInstruction],
    };
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        {
          ...alarmModule,
          getters: {
            getItem: () => () => alarmData,
          },
        },
      ]),
      propsData: {
        item: alarmData,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.executeInstruction).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        id: `${alarmData._id}${assignedInstruction._id}`,
        name: MODALS.executeRemediationInstruction,
        config: {
          alarmId: alarmData._id,
          assignedInstruction,
          onExecute: expect.any(Function),
          onClose: expect.any(Function),
          onComplete: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.onExecute();
    config.onClose();
    config.onComplete();

    expect(refreshAlarmsList).toHaveBeenCalledTimes(3);
  });

  test('Export PDF action', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        metaAlarmModule,
        authModuleWithAccess,
        alarmModule,
        alarmDetailsModule,
        declareTicketRuleModule,
      ]),
      propsData: {
        item: alarm,
        widget,
        parentAlarm,
      },
    });

    const exportPdfAction = selectActionByType(
      wrapper,
      ALARM_LIST_ACTIONS_TYPES.exportPdf,
    );

    exportPdfAction.trigger('click');

    await flushPromises();

    expect(fetchAlarmDetailsWithoutStore).toHaveBeenCalled();
    expect(exportAlarmToPdf).toHaveBeenCalled();
  });

  test('Custom action called after trigger button', () => {
    const customAction = {
      type: 'custom-type',
      icon: 'custom-icon',
      title: 'custom-title',
      method: jest.fn(),
    };
    const featureHasSpy = jest.spyOn(featuresService, 'has').mockReturnValueOnce(true);
    const featureGetSpy = jest.spyOn(featuresService, 'get').mockReturnValueOnce(() => [customAction]);

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarm,
        widget,
        parentAlarm,
      },
    });

    selectActionByType(wrapper, customAction.type).trigger('click');

    expect(customAction.method).toHaveBeenCalled();

    featureHasSpy.mockClear();
    featureGetSpy.mockClear();
  });

  test('Add bookmark request was sent after trigger add bookmark', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      mocks: {
        $popups,
      },
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.addBookmark).trigger('click');

    await flushPromises();

    expect(addBookmarkToAlarm).toHaveBeenCalledWith(
      expect.any(Object),
      { id: alarm._id },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Remove bookmark request was sent after trigger remove bookmark', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      mocks: {
        $popups,
      },
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: { ...alarm, bookmark: true },
        widget: widgetData,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.removeBookmark).trigger('click');

    await flushPromises();

    expect(removeBookmarkFromAlarm).toHaveBeenCalledWith(
      expect.any(Object),
      { id: alarm._id },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Renders `actions-panel` with manual instruction in running', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: {
          ...alarm,

          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualInProgress,
        },
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with simple manual instruction in running', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: {
          ...alarm,

          assigned_instructions: [
            ...assignedInstructions,
            {
              _id: 3,
              name: 'Manual simple instruction',
              type: REMEDIATION_INSTRUCTION_TYPES.simpleManual,
              execution: null,
            },
          ],
          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualInProgress,
        },
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with manual instruction in waiting result', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: {
          ...alarm,

          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualInProgress,
        },
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with auto instruction in running', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: {
          ...alarm,

          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoInProgress,
        },
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with paused manual instruction', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: {
          ...alarm,

          assigned_instructions: assignedInstructionsWithPaused,
        },
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with unresolved alarm and flapping status', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: alarm,
        widget,
        parentAlarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with unresolved alarm and ongoing status', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: {
          ...alarm,
          v: {
            status: {
              val: ALARM_STATUSES.ongoing,
            },
          },
        },
        widget,
        parentAlarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with resolved alarm', () => {
    const resolvedAlarmData = {
      ...alarm,
      v: {
        status: {
          val: ALARM_STATUSES.closed,
        },
      },
    };
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: resolvedAlarmData,
        widget,
        parentAlarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with full unresolved alarm, but without access', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        item: alarm,
        widget,
        parentAlarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` without entity, instructions, but with status stealthy', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: {
          ...alarm,
          assigned_instructions: undefined,
          entity: undefined,
          v: {
            status: {
              val: ALARM_STATUSES.stealthy,
            },
          },
        },
        widget,
        parentAlarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` without assigned_declare_ticket_rules', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: {
          ...alarm,

          assigned_declare_ticket_rules: undefined,
        },
        widget,
        parentAlarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with links in alarm', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: {
          ...alarm,

          links: {
            cat: [
              {
                hide_in_menu: true,
                icon_name: 'hidden_link_icon',
                label: 'Hidden link label',
                url: 'Hidden link URL',
                rule_id: 'Hidden link RuleId',
              },
              {
                icon_name: 'icon',
                label: 'Label',
                url: 'URL',
                rule_id: 'RuleId',
              },
            ],
          },
        },
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with links in resolved alarm', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: {
          ...alarm,

          v: {
            status: {
              val: ALARM_STATUSES.closed,
            },
          },
          links: {
            cat: [
              {
                icon_name: 'icon',
                label: 'Label',
                url: 'URL',
                rule_id: 'RuleId',
              },
            ],
          },
        },
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with parentAlarm with auto meta alarm', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: alarm,
        widget,
        parentAlarm: {
          meta_alarm_rule: {
            type: META_ALARMS_RULE_TYPES.attribute,
          },
          d: 'parent-d',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `actions-panel` with custom quickActions order (snapshot)', () => {
    const quickActions = [
      ALARM_LIST_ACTIONS_TYPES.ack,
      ALARM_LIST_ACTIONS_TYPES.fastAck,
      ALARM_LIST_ACTIONS_TYPES.cancel,
    ];
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: { parameters: { quickActions, isMultiAckEnabled: true } },
        parentAlarm,
      },
    });
    expect(wrapper).toMatchSnapshot();
  });
});
