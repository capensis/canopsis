import Faker from 'faker';
import { range } from 'lodash';
import flushPromises from 'flush-promises';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createAuthModule, createEventModule, createMockedStoreModules } from '@unit/utils/store';
import { mockDateNow, mockModals } from '@unit/utils/mock-hooks';
import {
  ALARM_LIST_ACTIONS_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  ENTITIES_STATUSES,
  ENTITY_PATTERN_FIELDS,
  EVENT_DEFAULT_ORIGIN,
  EVENT_ENTITY_TYPES,
  EVENT_INITIATORS,
  META_ALARMS_RULE_TYPES,
  MODALS,
  PATTERN_CONDITIONS,
} from '@/constants';

import MassActionsPanel from '@/components/widgets/alarm/actions/mass-actions-panel.vue';

const stubs = {
  'shared-mass-actions-panel': {
    props: ['actions', 'dropDownActions'],
    template: `
      <div class="shared-actions-panel">
        <button
          v-for="action in actions"
          :class="'action-' + action.type"
          @click="action.method"
        >{{ action.title }}|{{ action.icon }}|{{ action.type }}</button>
      </div>
    `,
  },
};

const selectActionByType = (wrapper, type) => wrapper.find(`.action-${type}`);

describe('mass-actions-panel', () => {
  const timestamp = 1386435600000;
  mockDateNow(timestamp);

  const $modals = mockModals();

  const alarm = {
    _id: 'alarm-id',
    entity: {
      _id: 'alarm-entity-id',
    },
    assigned_declare_ticket_rules: [{}],
    v: {
      state: {},
      status: {},
      tickets: [],
    },
  };

  const metaAlarm = {
    _id: 'meta-alarm-id',
    metaalarm: true,
    entity: {
      _id: 'meta-alarm-entity-id',
    },
    assigned_declare_ticket_rules: [{}],
    v: {
      state: {},
      status: {},
      tickets: [],
    },
  };
  const fastAckAlarms = range(2).map(index => ({
    _id: `alarm-id-${index}`,
    entity: {
      type: `entity-type-${index}`,
    },
    v: {
      connector: `alarm-connector-${index}`,
      connector_name: `alarm-connector-name-${index}`,
      component: `alarm-component-${index}`,
      resource: `alarm-resource-${index}`,
      status: {
        val: ENTITIES_STATUSES.ongoing,
      },
      state: {
        val: `state-val-${index}`,
      },
    },
  }));
  const fastAckEvents = fastAckAlarms.map(fastAckAlarm => ({
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
  }));

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
  const items = [alarm, metaAlarm];
  const { eventModule, createEvent } = createEventModule();

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

  const factory = generateShallowRenderer(MassActionsPanel, {
    stubs,
  });

  const snapshotFactory = generateRenderer(MassActionsPanel, {
    stubs,
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('Create pbehavior modal showed after trigger pbehavior add action', () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        eventModule,
      ]),
      propsData: {
        items,
        widget,
        parentAlarm,
      },
      mocks: {
        $modals,
      },
    });

    const pbehaviorAddAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd);

    pbehaviorAddAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: [[{
            field: ENTITY_PATTERN_FIELDS.id,
            cond: {
              type: PATTERN_CONDITIONS.isOneOf,
              value: [alarm.entity._id, metaAlarm.entity._id],
            },
          }]],
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);
  });

  it('Ack modal showed after trigger ack action', async () => {
    const isNoteRequired = Faker.datatype.boolean();
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        isAckNoteRequired: isNoteRequired,
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        eventModule,
      ]),
      propsData: {
        items,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    const ackAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.ack);

    ackAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createAckEvent,
        config: {
          isNoteRequired,
          items,
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.action({ output: 'OUTPUT', ack_resources: false }, {});

    await flushPromises();

    expect(wrapper).toEmit('clear:items');
    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Fast ack event sent after trigger fast ack action', async () => {
    const output = Faker.datatype.string();
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        fastAckOutput: {
          enabled: true,
          value: output,
        },
      },
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        eventModule,
      ]),
      propsData: {
        items: fastAckAlarms,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    const fastAckAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastAck);

    fastAckAction.trigger('click');

    await flushPromises();

    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: fastAckEvents.map(fastAckEvent => ({
          ...fastAckEvent,
          output,
        })),
      },
      undefined,
    );

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);
    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Fast ack event sent after trigger fast ack action without parameters', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        eventModule,
      ]),
      propsData: {
        items: fastAckAlarms,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    const fastAckAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastAck);

    fastAckAction.trigger('click');

    await flushPromises();

    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: fastAckEvents,
      },
      undefined,
    );

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);
    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Ack remove modal showed after trigger ack remove action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        eventModule,
      ]),
      propsData: {
        items,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    const ackRemoveAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.ackRemove);

    ackRemoveAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createEvent,
        config: {
          title: 'Remove ack',
          eventType: EVENT_ENTITY_TYPES.ackRemove,
          items,
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);
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
        eventModule,
      ]),
      propsData: {
        items,
        refreshAlarmsList,
        widget: widgetData,
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
          items,
          afterSubmit: expect.any(Function),
          title: 'Cancel',
          eventType: EVENT_ENTITY_TYPES.cancel,
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);
    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Associate ticket modal showed after trigger associate ticket action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        eventModule,
        authModuleWithAccess,
      ]),
      propsData: {
        items: [alarm],
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
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

    const event = {
      ticket: Faker.datatype.string(),
      ticket_url: Faker.datatype.string(),
      ticket_system_name: Faker.datatype.string(),
    };

    config.action(event);

    await flushPromises();

    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          component: undefined,
          connector: undefined,
          connector_name: undefined,
          crecord_type: EVENT_ENTITY_TYPES.ack,
          event_type: EVENT_ENTITY_TYPES.ack,
          id: alarm._id,
          initiator: 'user',
          origin: 'canopsis',
          output: '',
          ref_rk: 'undefined/undefined',
          resource: undefined,
          source_type: undefined,
          state: undefined,
          state_type: undefined,
          timestamp: 1386435600,
        }, {
          component: undefined,
          connector: undefined,
          connector_name: undefined,
          crecord_type: EVENT_ENTITY_TYPES.assocTicket,
          event_type: EVENT_ENTITY_TYPES.assocTicket,
          id: alarm._id,
          initiator: 'user',
          origin: 'canopsis',
          ref_rk: 'undefined/undefined',
          resource: undefined,
          source_type: undefined,
          state: undefined,
          state_type: undefined,
          timestamp: 1386435600,
          ...event,
        }],
      },
      undefined,
    );

    expect(wrapper).toEmit('clear:items');
    expect(refreshAlarmsList).toBeCalledTimes(1);
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
        eventModule,
      ]),
      propsData: {
        items,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    const snoozeAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.snooze);

    snoozeAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createSnoozeEvent,
        config: {
          isNoteRequired,
          items,
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);
    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Group alarm modal showed after trigger group alarm action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        eventModule,
      ]),
      propsData: {
        items,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    const ackRemoveAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.groupRequest);

    ackRemoveAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createEvent,
        config: {
          title: 'Suggest group request for meta alarm',
          eventType: EVENT_ENTITY_TYPES.groupRequest,
          items,
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);
    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Manual meta alarm group modal showed after trigger manual meta alarm group action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        eventModule,
      ]),
      propsData: {
        items,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    const ackRemoveAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.createManualMetaAlarm);

    ackRemoveAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createManualMetaAlarm,
        config: {
          title: 'Manual meta alarm management',
          items,
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);
    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Comment modal showed after trigger snooze action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        eventModule,
      ]),
      propsData: {
        items,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    const commentAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.comment);

    commentAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createCommentEvent,
        config: {
          items,
          afterSubmit: expect.any(Function),
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);
    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Renders `mass-actions-panel` with empty items', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        eventModule,
      ]),
      propsData: {
        items,
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `mass-actions-panel` with meta alarm', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
      ]),
      propsData: {
        items,
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
