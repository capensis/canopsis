import Faker from 'faker';
import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import {
  createAlarmModule,
  createAlarmDetailsModule,
  createAuthModule,
  createDeclareTicketModule,
  createManualMetaAlarmModule,
  createMetaAlarmModule,
  createMockedStoreModules,
} from '@unit/utils/store';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import {
  ALARM_LIST_ACTIONS_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  ENTITY_PATTERN_FIELDS,
  INSTRUCTION_EXECUTION_ICONS,
  META_ALARMS_RULE_TYPES,
  MODALS,
  PATTERN_CONDITIONS,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES,
  REMEDIATION_INSTRUCTION_TYPES,
  TIME_UNITS,
} from '@/constants';

import featuresService from '@/services/features';

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
  mockDateNow(timestamp);
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
    bulkCreateAlarmCommentEvent,
    bulkCreateAlarmCancelEvent,
    bulkCreateAlarmChangestateEvent,
    addBookmarkToAlarm,
    removeBookmarkFromAlarm,
  } = createAlarmModule();
  const { manualMetaAlarmModule, removeAlarmsFromManualMetaAlarm } = createManualMetaAlarmModule();
  const { metaAlarmModule, removeAlarmsFromMetaAlarm } = createMetaAlarmModule();
  const { alarmDetailsModule, fetchAlarmDetailsWithoutStore } = createAlarmDetailsModule();

  const {
    declareTicketRuleModule,
    fetchAssignedDeclareTicketsWithoutStore,
  } = createDeclareTicketModule();

  const store = createMockedStoreModules([
    manualMetaAlarmModule,
    metaAlarmModule,
    authModule,
    alarmModule,
    alarmDetailsModule,
    declareTicketRuleModule,
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
        val: ENTITIES_STATUSES.flapping,
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
    mocks: { $modals },
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

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('Ack modal showed after trigger ack action', async () => {
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
        manualMetaAlarmModule,
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

  it('Fast ack event sent after trigger fast ack action', () => {
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
          val: ENTITIES_STATUSES.ongoing,
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
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: fastActionAlarm,
        widget: widgetData,
        parentAlarm,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastAck).trigger('click');

    expect(bulkCreateAlarmAckEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          comment: widgetData.parameters.fastAckOutput.value,
        }],
      },
      undefined,
    );
  });

  it('Ack remove modal showed after trigger ack remove action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.ackRemove).trigger('click');

    expect($modals.show).toBeCalledWith(
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

    expect(bulkCreateAlarmAckremoveEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          comment,
        }],
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Create pbehavior modal showed after trigger pbehavior add action', () => {
    const entity = {
      _id: Faker.datatype.string(),
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: { ...alarm, entity },
        widget,
        parentAlarm,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd).trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.pbehaviorPlanning,
        config: {
          afterSubmit: expect.any(Function),
          entityPattern: [[{
            field: ENTITY_PATTERN_FIELDS.id,
            cond: {
              type: PATTERN_CONDITIONS.equal,
              value: entity._id,
            },
          }]],
        },
      },
    );
  });

  it('Snooze modal showed after trigger snooze action', async () => {
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
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.snooze).trigger('click');

    expect($modals.show).toBeCalledWith(
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

    expect(bulkCreateAlarmSnoozeEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          ...snoozeEvent,
        }],
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Declare ticket modal showed after trigger declare action', async () => {
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
        manualMetaAlarmModule,
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

    expect($modals.show).toBeCalledWith(
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

    expect($modals.show).toBeCalledWith({
      name: MODALS.executeDeclareTickets,
      config: {
        executions: events,
        alarms: [alarm],
        tickets: [rule],
        onExecute: expect.any(Function),
      },
    });
  });

  it('Associate ticket modal showed after trigger associate ticket action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        manualMetaAlarmModule,
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

    expect($modals.show).toBeCalledWith(
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

    expect(bulkCreateAlarmAssocticketEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          ...ticketEvent,
        }],
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Change state modal showed after trigger change state action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.changeState).trigger('click');

    expect($modals.show).toBeCalledWith(
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
      state: ENTITIES_STATES.critical,
      comment: Faker.datatype.string(),
    };

    await config.action(changeStateEvent);

    expect(bulkCreateAlarmChangestateEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          ...changeStateEvent,
        }],
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Cancel modal showed after trigger cancel action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.cancel).trigger('click');

    expect($modals.show).toBeCalledWith(
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

    expect(bulkCreateAlarmCancelEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          ...cancelEvent,
        }],
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Fast cancel event sent after trigger fast cancel action', async () => {
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
          val: ENTITIES_STATUSES.ongoing,
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
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: fastActionAlarm,
        widget: widgetData,
        parentAlarm,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastCancel).trigger('click');

    await flushPromises();

    expect(bulkCreateAlarmCancelEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          _id: alarm._id,
          comment: 'Output',
        }],
      },
      undefined,
    );
  });

  it('Fast cancel and cancel action available when status is flapping and state is ok', async () => {
    const flappingAlarm = {
      ...alarm,
      v: {
        connector: 'alarm-connector',
        connector_name: 'alarm-connector-name',
        component: 'alarm-component',
        resource: 'alarm-resource',
        status: {
          val: ENTITIES_STATUSES.flapping,
        },
        state: {
          val: ENTITIES_STATES.ok,
        },
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        manualMetaAlarmModule,
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

  it('Fast cancel and cancel action available when alarm is not resolved and state is ok', async () => {
    const flappingAlarm = {
      ...alarm,
      v: {
        connector: 'alarm-connector',
        connector_name: 'alarm-connector-name',
        component: 'alarm-component',
        resource: 'alarm-resource',
        resolved: null,
        status: {},
        state: {
          val: ENTITIES_STATES.ok,
        },
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        manualMetaAlarmModule,
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

  it('Variables modal showed after trigger variables help action', () => {
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
          val: ENTITIES_STATUSES.closed,
        },
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: {
          ...alarmData,
          v: {
            status: {
              val: ENTITIES_STATUSES.closed,
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

  it('History modal showed after trigger history action', () => {
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
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: { ...alarm, entity },
        widget: widgetData,
        parentAlarm,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.history).trigger('click');

    const defaultWidget = prepareAlarmListWidget(generateDefaultAlarmListWidget());

    expect($modals.show).toBeCalledWith(
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

  it('Comment modal showed after trigger comment action', async () => {
    const commentAlarm = {
      ...alarm,
      v: {
        connector: 'alarm-connector',
        connector_name: 'alarm-connector-name',
        component: 'alarm-component',
        resource: 'alarm-resource',
        status: {
          val: ENTITIES_STATUSES.ongoing,
        },
        state: {
          val: 'state-val',
        },
      },
    };
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: commentAlarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.comment).trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createCommentEvent,
        config: {
          items: [commentAlarm],
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
        data: [{
          _id: alarm._id,
          comment,
        }],
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Remove alarms from manual meta alarm modal showed after trigger remove alarms from manual meta alarm action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.removeAlarmsFromManualMetaAlarm).trigger('click');

    expect($modals.show).toBeCalledWith(
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

    expect(removeAlarmsFromManualMetaAlarm).toBeCalledWith(
      expect.any(Object),
      {
        id: parentAlarm?._id,
        data: newRemoveAlarmsEvent,
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Remove alarms from manual meta alarm modal showed after trigger remove alarms from manual meta alarm action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.removeAlarmsFromManualMetaAlarm).trigger('click');

    expect($modals.show).toBeCalledWith(
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

    expect(removeAlarmsFromManualMetaAlarm).toBeCalledWith(
      expect.any(Object),
      {
        id: parentAlarm?._id,
        data: newRemoveAlarmsEvent,
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Remove alarms from auto meta alarm modal showed after trigger remove alarms from auto meta alarm action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        manualMetaAlarmModule,
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

    expect($modals.show).toBeCalledWith(
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

    expect(removeAlarmsFromMetaAlarm).toBeCalledWith(
      expect.any(Object),
      {
        id: parentAlarm?._id,
        data: newRemoveAlarmsEvent,
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Execute instruction alarm modal showed after trigger execute instruction action', () => {
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
        manualMetaAlarmModule,
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

    expect($modals.show).toBeCalledWith(
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

    expect(refreshAlarmsList).toBeCalledTimes(3);
  });

  it('Export PDF action', async () => {
    const wrapper = factory({
      store,
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

    expect(fetchAlarmDetailsWithoutStore).toBeCalled();
    expect(exportAlarmToPdf).toBeCalled();
  });

  it('Custom action called after trigger button', () => {
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

    expect(customAction.method).toBeCalled();

    featureHasSpy.mockClear();
    featureGetSpy.mockClear();
  });

  it('Add bookmark request was sent after trigger add bookmark', async () => {
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
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: alarm,
        widget: widgetData,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.addBookmark).trigger('click');

    await flushPromises();

    expect(addBookmarkToAlarm).toBeCalledWith(
      expect.any(Object),
      { id: alarm._id },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Remove bookmark request was sent after trigger remove bookmark', async () => {
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
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: { ...alarm, bookmark: true },
        widget: widgetData,
        refreshAlarmsList,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.removeBookmark).trigger('click');

    await flushPromises();

    expect(removeBookmarkFromAlarm).toBeCalledWith(
      expect.any(Object),
      { id: alarm._id },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Renders `actions-panel` with manual instruction in running', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        manualMetaAlarmModule,
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

  it('Renders `actions-panel` with simple manual instruction in running', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        manualMetaAlarmModule,
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

  it('Renders `actions-panel` with manual instruction in waiting result', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        manualMetaAlarmModule,
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

  it('Renders `actions-panel` with auto instruction in running', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        manualMetaAlarmModule,
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

  it('Renders `actions-panel` with paused manual instruction', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        manualMetaAlarmModule,
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

  it('Renders `actions-panel` with unresolved alarm and flapping status', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: alarm,
        widget,
        parentAlarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `actions-panel` with unresolved alarm and ongoing status', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: {
          ...alarm,
          v: {
            status: {
              val: ENTITIES_STATUSES.ongoing,
            },
          },
        },
        widget,
        parentAlarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `actions-panel` with resolved alarm', () => {
    const resolvedAlarmData = {
      ...alarm,
      v: {
        status: {
          val: ENTITIES_STATUSES.closed,
        },
      },
    };
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        manualMetaAlarmModule,
      ]),
      propsData: {
        item: resolvedAlarmData,
        widget,
        parentAlarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `actions-panel` with full unresolved alarm, but without access', () => {
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

  it('Renders `actions-panel` without entity, instructions, but with status stealthy', () => {
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
              val: ENTITIES_STATUSES.stealthy,
            },
          },
        },
        widget,
        parentAlarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `actions-panel` without assigned_declare_ticket_rules', () => {
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

  it('Renders `actions-panel` with links in alarm', () => {
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

  it('Renders `actions-panel` with links in resolved alarm', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: {
          ...alarm,

          v: {
            status: {
              val: ENTITIES_STATUSES.closed,
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

  it('Renders `actions-panel` with parentAlarm with auto meta alarm', () => {
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
});
