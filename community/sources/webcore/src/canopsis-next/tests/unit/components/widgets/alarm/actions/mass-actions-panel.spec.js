import Faker from 'faker';
import { range } from 'lodash';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import {
  createAlarmModule,
  createAuthModule,
  createMetaAlarmModule,
  createMockedStoreModules,
} from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';

import {
  ALARM_LIST_ACTIONS_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  ALARM_STATUSES,
  ENTITY_PATTERN_FIELDS,
  LINK_RULE_ACTIONS,
  META_ALARMS_RULE_TYPES,
  MODALS,
  PATTERN_CONDITIONS,
  TIME_UNITS,
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
  jest.useFakeTimers({ now: timestamp });

  const $modals = mockModals();

  const alarm = {
    _id: 'alarm-id',
    entity: {
      _id: 'alarm-entity-id',
    },
    assigned_declare_ticket_rules: [{}],
    links: {
      Category: [{
        rule_id: 'rule_id',
        label: 'with rule id',
        icon_name: '',
        url: 'url',
        action: LINK_RULE_ACTIONS.open,
      }, {
        label: 'without rule id',
        icon_name: '',
        url: 'url',
        action: LINK_RULE_ACTIONS.open,
      }],
    },
    v: {
      state: {},
      status: {},
      tickets: [],
    },
  };
  const alarmWithAck = {
    ...alarm,
    v: {
      ...alarm.v,
      ack: {},
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
  const fastActionAlarms = range(2).map(index => ({
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
        val: ALARM_STATUSES.ongoing,
      },
      state: {
        val: `state-val-${index}`,
      },
    },
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
  } = createAlarmModule();
  const { metaAlarmModule, addAlarmsIntoMetaAlarm, createMetaAlarm } = createMetaAlarmModule();

  const items = [alarm, metaAlarm];

  const store = createMockedStoreModules([authModuleWithAccess, alarmModule, metaAlarmModule]);

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

  test('Create pbehavior modal showed after trigger pbehavior add action', () => {
    const wrapper = factory({
      store,
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

    expect($modals.show).toHaveBeenCalledWith(
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
          entities: [alarm.entity, metaAlarm.entity],
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    expect(wrapper).toHaveBeenEmit('clear:items');
  });

  test('Ack modal showed after trigger ack action', async () => {
    const isNoteRequired = Faker.datatype.boolean();
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        isAckNoteRequired: isNoteRequired,
      },
    };

    const wrapper = factory({
      store,
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

    expect($modals.show).toHaveBeenCalledWith(
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

    expect(wrapper).toHaveBeenEmit('clear:items');
    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Fast ack event sent after trigger fast ack action', async () => {
    const comment = Faker.datatype.string();
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        fastAckOutput: {
          enabled: true,
          value: comment,
        },
      },
    };

    const wrapper = factory({
      store,
      propsData: {
        items: fastActionAlarms,
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

    expect(bulkCreateAlarmAckEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: fastActionAlarms.map(({ _id: alarmId }) => ({ _id: alarmId, comment })),
      },
    );

    expect(wrapper).toHaveBeenEmit('clear:items');
    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Fast ack event sent after trigger fast ack action without parameters', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store,
      propsData: {
        items: fastActionAlarms,
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

    expect(bulkCreateAlarmAckEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: fastActionAlarms.map(({ _id: alarmId }) => ({ _id: alarmId, comment: '' })),
      },
    );

    expect(wrapper).toHaveBeenEmit('clear:items');
    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Ack remove modal showed after trigger ack remove action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const itemsForAck = [...items, alarmWithAck];
    const wrapper = factory({
      store,
      propsData: {
        items: itemsForAck,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.ackRemove).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createEvent,
        config: {
          title: 'Remove ack',
          items: itemsForAck,
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
        data: itemsForAck.map(({ _id: alarmId }) => ({ _id: alarmId, comment })),
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
    expect(wrapper).toHaveBeenEmit('clear:items');
  });

  test('Cancel modal showed after trigger cancel action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store,
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

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createEvent,
        config: {
          items,
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
        data: items.map(({ _id: alarmId }) => ({ _id: alarmId, ...cancelEvent })),
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
    expect(wrapper).toHaveBeenEmit('clear:items');
  });

  test('Fast cancel event sent after trigger fast cancel action', async () => {
    const comment = Faker.datatype.string();
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        fastCancelOutput: {
          enabled: true,
          value: comment,
        },
      },
    };

    const wrapper = factory({
      store,
      propsData: {
        items: fastActionAlarms,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastCancel).trigger('click');

    await flushPromises();

    expect(bulkCreateAlarmCancelEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: fastActionAlarms.map(({ _id: alarmId }) => ({ _id: alarmId, comment })),
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
    expect(wrapper).toHaveBeenEmit('clear:items');
  });

  test('Fast cancel event sent after trigger fast cancel action without parameters', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store,
      propsData: {
        items: fastActionAlarms,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.fastCancel).trigger('click');

    await flushPromises();

    expect(bulkCreateAlarmCancelEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: fastActionAlarms.map(({ _id: alarmId }) => ({ _id: alarmId, comment: '' })),
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
    expect(wrapper).toHaveBeenEmit('clear:items');
  });

  test('Associate ticket modal showed after trigger associate ticket action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store,
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
      ticket: Faker.datatype.string(),
      ticket_url: Faker.datatype.string(),
      ticket_system_name: Faker.datatype.string(),
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

    expect(wrapper).toHaveBeenEmit('clear:items');
    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
  });

  test('Remove associated ticket modal showed after trigger remove associated ticket action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const alarmsWithTickets = [
      {
        ...alarm,
        _id: 'alarm-with-ticket-1',
        v: {
          ...alarm.v,
          tickets: [
            {
              ticket: 'TICKET-123',
              ticket_system_name: 'Jira',
            },
          ],
        },
      },
      {
        ...alarm,
        _id: 'alarm-with-ticket-2',
        v: {
          ...alarm.v,
          tickets: [
            {
              ticket: 'TICKET-456',
              ticket_system_name: 'ServiceNow',
            },
          ],
        },
      },
    ];

    const wrapper = factory({
      store,
      propsData: {
        items: alarmsWithTickets,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.removeAssociatedTicket).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.removeAssociatedTicketEvent,
        config: {
          items: alarmsWithTickets,
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
        data: alarmsWithTickets.map(({ _id: alarmId }) => ({
          _id: alarmId,
          ...removeTicketEvent,
        })),
      },
    );

    expect(wrapper).toHaveBeenEmit('clear:items');
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
      store,
      propsData: {
        items,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.snooze).trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createSnoozeEvent,
        config: {
          isNoteRequired,
          items,
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
        data: items.map(({ _id: alarmId }) => ({ _id: alarmId, ...snoozeEvent })),
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
    expect(wrapper).toHaveBeenEmit('clear:items');
  });

  test('Manual meta alarm group modal showed after trigger manual meta alarm group action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store,
      propsData: {
        items,
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    const ackRemoveAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.linkToMetaAlarm);

    ackRemoveAction.trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.linkToMetaAlarm,
        config: {
          title: 'Link to a meta alarm',
          items,
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const manualMetaAlarmEventWithId = {
      id: Faker.datatype.string(),
      comment: Faker.datatype.string(),
    };

    await config.action(manualMetaAlarmEventWithId);

    expect(addAlarmsIntoMetaAlarm).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: manualMetaAlarmEventWithId.id,
        data: manualMetaAlarmEventWithId,
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
    expect(wrapper).toHaveBeenEmit('clear:items');

    addAlarmsIntoMetaAlarm.mockClear();
    refreshAlarmsList.mockClear();

    const manualMetaAlarmEventWithoutId = {
      comment: Faker.datatype.string(),
      metaAlarm: Faker.datatype.string(),
    };

    await config.action(manualMetaAlarmEventWithoutId);

    expect(createMetaAlarm).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: manualMetaAlarmEventWithoutId.id,
        data: manualMetaAlarmEventWithoutId,
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
    expect(wrapper).toHaveBeenEmit('clear:items');
  });

  test('Comment modal showed after trigger comment action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store,
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

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createCommentEvent,
        config: {
          items,
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
        data: items.map(({ _id: alarmId }) => ({ _id: alarmId, comment })),
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
    expect(wrapper).toHaveBeenEmit('clear:items');
  });

  test('Change state modal showed after trigger change state action', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store,
      propsData: {
        items: [alarmWithAck],
        refreshAlarmsList,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    const changeStateAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.changeState);

    changeStateAction.trigger('click');

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.createChangeStateEvent,
        config: {
          items: [alarmWithAck],
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const changeStateEvent = {
      state: Faker.datatype.number(),
      comment: Faker.datatype.string(),
    };

    await config.action(changeStateEvent);

    expect(bulkCreateAlarmChangestateEvent).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: [{ _id: alarmWithAck._id, ...changeStateEvent }],
      },
    );

    expect(refreshAlarmsList).toHaveBeenCalledTimes(1);
    expect(wrapper).toHaveBeenEmit('clear:items');
  });

  test('inlineCount reflects quickMassActions count', () => {
    const quickMassActions = [
      ALARM_LIST_ACTIONS_TYPES.fastAck,
      ALARM_LIST_ACTIONS_TYPES.ack,
      ALARM_LIST_ACTIONS_TYPES.cancel,
    ];
    const wrapper = factory({
      store,
      propsData: {
        items,
        widget: { parameters: { quickMassActions } },
      },
    });
    // getActionsInlineCount returns quickMassActions.length + 1 (menu button) if not all actions are quick
    expect(wrapper.vm.inlineCount).toBeGreaterThanOrEqual(quickMassActions.length);
  });

  test('Renders `mass-actions-panel` with non empty items', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        items,
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `mass-actions-panel` with empty items', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        items: [],
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `mass-actions-panel` with one item', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        items: [alarm],
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `mass-actions-panel` with meta alarm', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        items,
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `mass-actions-panel` with meta ack', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        items: [...items, alarmWithAck],
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `mass-actions-panel` with quickMassActions set (custom order)', () => {
    const quickMassActions = [
      ALARM_LIST_ACTIONS_TYPES.fastAck,
      ALARM_LIST_ACTIONS_TYPES.ack,
      ALARM_LIST_ACTIONS_TYPES.cancel,
    ];
    const wrapper = snapshotFactory({
      store,
      propsData: {
        items,
        widget: { parameters: { quickMassActions } },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
