import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockDateNow, mockModals } from '@unit/utils/mock-hooks';
import {
  ALARM_LIST_ACTIONS_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  ENTITIES_STATUSES,
  ENTITY_PATTERN_FIELDS,
  EVENT_DEFAULT_ORIGIN,
  EVENT_ENTITY_TYPES,
  EVENT_INITIATORS,
  INSTRUCTION_EXECUTION_ICONS,
  META_ALARMS_RULE_TYPES,
  MODALS,
  PATTERN_CONDITIONS,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES,
  REMEDIATION_INSTRUCTION_TYPES,
} from '@/constants';

import featuresService from '@/services/features';

import { generateDefaultAlarmListWidget } from '@/helpers/entities';

import ActionsPanel from '@/components/widgets/alarm/actions/actions-panel.vue';

const localVue = createVueInstance();

const stubs = {
  'shared-actions-panel': {
    props: ['actions', 'dropDownActions'],
    template: `
      <div class="shared-actions-panel">
        <button
          v-for="action in actions"
          :class="'action-' + action.type"
          :disabled="action.disabled"
          @click="action.method"
        >{{ action.title }}|{{ action.icon }}|{{ action.type }}</button>
        <button
          v-for="action in dropDownActions"
          :class="'drop-down-action-' + action.type"
          :disabled="action.disabled"
          @click="action.method"
        >{{ action.title }}|{{ action.icon }}|{{ action.type }}</button>
      </div>
    `,
  },
};

const factory = (options = {}) => shallowMount(ActionsPanel, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(ActionsPanel, {
  localVue,
  stubs,

  ...options,
});

const selectActionByType = (wrapper, type) => wrapper.find(`.action-${type}`);
const selectDropDownActionByType = (wrapper, type) => wrapper.find(`.drop-down-action-${type}`);

describe('actions-panel', () => {
  const timestamp = 1386435600000;
  mockDateNow(timestamp);

  const $modals = mockModals();
  const authModule = {
    name: 'auth',
    getters: {
      currentUserPermissionsById: {},
    },
  };
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
  const fetchAlarmItem = jest.fn();
  const alarmModule = {
    name: 'alarm',
    actions: {
      fetchItem: fetchAlarmItem,
    },
  };
  const createEvent = jest.fn();
  const eventModule = {
    name: 'event',
    actions: {
      create: createEvent,
    },
  };

  const store = createMockedStoreModules([
    authModule,
    alarmModule,
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

  const alarm = {
    _id: 'alarm-id',
    assigned_instructions: assignedInstructions,
    entity: {},
    v: {
      ack: {},
      status: {
        val: ENTITIES_STATUSES.flapping,
      },
    },
  };

  const widget = {
    parameters: {
      isMultiAckEnabled: true,
    },
  };

  const parentAlarm = {
    rule: {
      type: META_ALARMS_RULE_TYPES.manualgroup,
    },
    d: 'parent-d',
  };

  const refreshAlarmsList = jest.fn();

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('Ack modal showed after trigger ack action', () => {
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
      mocks: {
        $modals,
      },
    });

    const ackAction = selectDropDownActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.ack);

    ackAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createAckEvent,
        config: {
          isNoteRequired,
          items: [alarm],
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Fast ack event sent after trigger fast ack action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        fastAckOutput: {
          enabled: true,
        },
      },
    };
    const fastAckAlarm = {
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
        eventModule,
      ]),
      propsData: {
        item: fastAckAlarm,
        widget: widgetData,
        parentAlarm,
      },
      mocks: {
        $modals,
      },
    });

    const fastAckAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastAck);

    fastAckAction.trigger('click');

    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: {
          timestamp: timestamp / 1000,
          component: fastAckAlarm.v.component,
          connector: fastAckAlarm.v.connector,
          connector_name: fastAckAlarm.v.connector_name,
          resource: fastAckAlarm.v.resource,
          state: fastAckAlarm.v.state.val,
          state_type: fastAckAlarm.v.status.val,
          source_type: fastAckAlarm.entity.type,
          crecord_type: 'ack',
          event_type: 'ack',
          id: fastAckAlarm._id,
          initiator: EVENT_INITIATORS.user,
          origin: EVENT_DEFAULT_ORIGIN,
          ref_rk: `${fastAckAlarm.v.resource}/${fastAckAlarm.v.component}`,
        },
      },
      undefined,
    );
  });

  it('Ack remove modal showed after trigger ack remove action', () => {
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
      mocks: {
        $modals,
      },
    });

    const ackRemoveAction = selectDropDownActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.ackRemove);

    ackRemoveAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createEvent,
        config: {
          title: 'Remove ack',
          eventType: EVENT_ENTITY_TYPES.ackRemove,
          items: [alarm],
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

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
      ]),
      propsData: {
        item: { ...alarm, entity },
        widget,
        parentAlarm,
      },
      mocks: {
        $modals,
      },
    });

    const pbehaviorAddAction = selectDropDownActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd);

    pbehaviorAddAction.trigger('click');

    expect($modals.show).toBeCalledWith(
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
        },
      },
    );
  });

  it('Snooze modal showed after trigger snooze action', () => {
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
      mocks: {
        $modals,
      },
    });

    const snoozeAction = selectDropDownActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.snooze);

    snoozeAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createSnoozeEvent,
        config: {
          isNoteRequired,
          items: [alarm],
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Declare ticket modal showed after trigger declare action', () => {
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
      mocks: {
        $modals,
      },
    });

    const declareTicketAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.declareTicket);

    declareTicketAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createDeclareTicketEvent,
        config: {
          items: [alarm],
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Associate ticket modal showed after trigger associate ticket action', () => {
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
      mocks: {
        $modals,
      },
    });

    const associateTicketAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.associateTicket);

    associateTicketAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createAssociateTicketEvent,
        config: {
          items: [alarm],
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Change state modal showed after trigger change state action', () => {
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
      mocks: {
        $modals,
      },
    });

    const changeStateAction = selectDropDownActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.changeState);

    changeStateAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createChangeStateEvent,
        config: {
          items: [alarm],
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Cancel modal showed after trigger cancel action', () => {
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
      mocks: {
        $modals,
      },
    });

    const cancelAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.cancel);

    cancelAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createEvent,
        config: {
          items: [alarm],
          afterSubmit: expect.any(Function),
          title: 'Cancel',
          eventType: EVENT_ENTITY_TYPES.cancel,
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    expect(refreshAlarmsList).toBeCalledTimes(1);
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
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
      ]),
      propsData: {
        item: alarmData,
        widget: widgetData,
        parentAlarm,
        isResolvedAlarm: true,
      },
      mocks: {
        $modals,
      },
    });

    const variablesHelpAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.variablesHelp);

    variablesHelpAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.variablesHelp,
        config: {
          items: [alarmData],
          afterSubmit: expect.any(Function),
          variables: [
            {
              name: 'alarm',
              children: [{ name: '_id', path: 'alarm._id', value: alarmData._id }],
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
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        widgetColumns: [
          {
            value: Faker.datatype.string(),
            label: Faker.datatype.string(),
          },
        ],
      },
    };
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
      mocks: {
        $modals,
      },
    });

    const historyAction = selectDropDownActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.history);

    historyAction.trigger('click');

    const defaultWidget = generateDefaultAlarmListWidget();

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
            },
          },
        },
      },
    );
  });

  it('Comment modal showed after trigger comment action', () => {
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
        eventModule,
      ]),
      propsData: {
        item: commentAlarm,
        widget: widgetData,
        parentAlarm,
        refreshAlarmsList,
      },
      mocks: {
        $modals,
      },
    });

    const commentAction = selectDropDownActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.comment);

    commentAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createCommentEvent,
        config: {
          items: [commentAlarm],
          afterSubmit: expect.any(Function),
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.action();
    config.afterSubmit();

    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: {
          timestamp: timestamp / 1000,
          component: commentAlarm.v.component,
          connector: commentAlarm.v.connector,
          connector_name: commentAlarm.v.connector_name,
          resource: commentAlarm.v.resource,
          state: commentAlarm.v.state.val,
          state_type: commentAlarm.v.status.val,
          source_type: undefined,
          crecord_type: 'comment',
          event_type: 'comment',
          id: commentAlarm._id,
          initiator: EVENT_INITIATORS.user,
          origin: EVENT_DEFAULT_ORIGIN,
          ref_rk: `${commentAlarm.v.resource}/${commentAlarm.v.component}`,
        },
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Manual meta alarm modal showed after trigger manual meta alarm ungroup action', () => {
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
      mocks: {
        $modals,
      },
    });

    const manualMetaAlarmUngroupAction = selectDropDownActionByType(
      wrapper,
      ALARM_LIST_ACTIONS_TYPES.manualMetaAlarmUngroup,
    );

    manualMetaAlarmUngroupAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createEvent,
        config: {
          items: [alarm],
          afterSubmit: expect.any(Function),
          title: 'Unlink alarm from manual meta alarm',
          eventType: EVENT_ENTITY_TYPES.manualMetaAlarmUngroup,
          parentsIds: [parentAlarm.d],
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

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
      mocks: {
        $modals,
      },
    });

    const executeInstructionAction = selectDropDownActionByType(
      wrapper,
      ALARM_LIST_ACTIONS_TYPES.executeInstruction,
    );

    executeInstructionAction.trigger('click');

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

  it('Custom action called after trigger button', () => {
    const customAction = {
      type: 'custom-type',
      icon: 'custom-icon',
      title: 'custom-title',
      method: jest.fn(),
    };
    const featureHasSpy = jest.spyOn(featuresService, 'has')
      .mockReturnValueOnce(true);
    const featureGetSpy = jest.spyOn(featuresService, 'get')
      .mockReturnValueOnce((
      ) => ({
        inline: [customAction],
        dropDown: [],
      }));

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

    const executeInstructionAction = selectActionByType(wrapper, customAction.type);

    executeInstructionAction.trigger('click');

    expect(customAction.method).toBeCalled();

    featureHasSpy.mockClear();
    featureGetSpy.mockClear();
  });

  // TODO: put tests for: no active instructions, one active instruction
  it('Renders `actions-panel` with manual instruction in running', () => {
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

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with manual instruction in waiting result', () => {
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

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with auto instruction in running', () => {
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

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with paused manual instruction', () => {
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

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with unresolved alarm and flapping status', () => {
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

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with unresolved alarm and ongoing status', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
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

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with resolved alarm', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        item: alarm,
        widget,
        parentAlarm,
        isResolvedAlarm: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` without entity, instructions, but with status stealthy', () => {
    const wrapper = snapshotFactory({
      store,
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
