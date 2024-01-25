import Faker from 'faker';
import { range } from 'lodash';
import flushPromises from 'flush-promises';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import {
  createAlarmModule,
  createAuthModule,
  createManualMetaAlarmModule,
  createMockedStoreModules,
} from '@unit/utils/store';
import { mockDateNow, mockModals } from '@unit/utils/mock-hooks';
import {
  ALARM_LIST_ACTIONS_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  ENTITIES_STATUSES,
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
  mockDateNow(timestamp);

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
        val: ENTITIES_STATUSES.ongoing,
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
    bulkCreateAlarmCommentEvent,
    bulkCreateAlarmCancelEvent,
  } = createAlarmModule();
  const { manualMetaAlarmModule, addAlarmsIntoManualMetaAlarm, createManualMetaAlarm } = createManualMetaAlarmModule();

  const items = [alarm, metaAlarm];

  const store = createMockedStoreModules([authModuleWithAccess, alarmModule, manualMetaAlarmModule]);

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

  const factory = generateShallowRenderer(MassActionsPanel, { stubs });

  const snapshotFactory = generateRenderer(MassActionsPanel, { stubs });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('Create pbehavior modal showed after trigger pbehavior add action', () => {
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

    expect(wrapper).toEmit('clear:items');
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

    expect(bulkCreateAlarmAckEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: fastActionAlarms.map(({ _id: alarmId }) => ({ _id: alarmId, comment })),
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

    expect(bulkCreateAlarmAckEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: fastActionAlarms.map(({ _id: alarmId }) => ({ _id: alarmId, comment: '' })),
      },
      undefined,
    );

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);
    expect(refreshAlarmsList).toBeCalledTimes(1);
  });

  it('Ack remove modal showed after trigger ack remove action', async () => {
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

    selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.ackRemove).trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createEvent,
        config: {
          title: 'Remove ack',
          items,
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
        data: items.map(({ _id: alarmId }) => ({ _id: alarmId, comment })),
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
    expect(wrapper).toEmit('clear:items');
  });

  it('Cancel modal showed after trigger cancel action', async () => {
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

    expect($modals.show).toBeCalledWith(
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

    expect(bulkCreateAlarmCancelEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: items.map(({ _id: alarmId }) => ({ _id: alarmId, ...cancelEvent })),
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
    expect(wrapper).toEmit('clear:items');
  });

  it('Fast cancel event sent after trigger fast cancel action', async () => {
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

    expect(bulkCreateAlarmCancelEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: fastActionAlarms.map(({ _id: alarmId }) => ({ _id: alarmId, comment })),
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
    expect(wrapper).toEmit('clear:items');
  });

  it('Fast cancel event sent after trigger fast cancel action without parameters', async () => {
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

    expect(bulkCreateAlarmCancelEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: fastActionAlarms.map(({ _id: alarmId }) => ({ _id: alarmId, comment: '' })),
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
    expect(wrapper).toEmit('clear:items');
  });

  it('Associate ticket modal showed after trigger associate ticket action', async () => {
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
      ticket: Faker.datatype.string(),
      ticket_url: Faker.datatype.string(),
      ticket_system_name: Faker.datatype.string(),
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

    expect(wrapper).toEmit('clear:items');
    expect(refreshAlarmsList).toBeCalledTimes(1);
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

    expect($modals.show).toBeCalledWith(
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

    expect(bulkCreateAlarmSnoozeEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: items.map(({ _id: alarmId }) => ({ _id: alarmId, ...snoozeEvent })),
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
    expect(wrapper).toEmit('clear:items');
  });

  it('Manual meta alarm group modal showed after trigger manual meta alarm group action', async () => {
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

    const ackRemoveAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.createManualMetaAlarm);

    ackRemoveAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createManualMetaAlarm,
        config: {
          title: 'Manual meta alarm management',
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

    expect(addAlarmsIntoManualMetaAlarm).toBeCalledWith(
      expect.any(Object),
      {
        id: manualMetaAlarmEventWithId.id,
        data: manualMetaAlarmEventWithId,
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
    expect(wrapper).toEmit('clear:items');

    addAlarmsIntoManualMetaAlarm.mockClear();
    refreshAlarmsList.mockClear();

    const manualMetaAlarmEventWithoutId = {
      comment: Faker.datatype.string(),
      metaAlarm: Faker.datatype.string(),
    };

    await config.action(manualMetaAlarmEventWithoutId);

    expect(createManualMetaAlarm).toBeCalledWith(
      expect.any(Object),
      {
        id: manualMetaAlarmEventWithoutId.id,
        data: manualMetaAlarmEventWithoutId,
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
    expect(wrapper).toEmit('clear:items');
  });

  it('Comment modal showed after trigger comment action', async () => {
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

    expect($modals.show).toBeCalledWith(
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

    expect(bulkCreateAlarmCommentEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: items.map(({ _id: alarmId }) => ({ _id: alarmId, comment })),
      },
      undefined,
    );

    expect(refreshAlarmsList).toBeCalledTimes(1);
    expect(wrapper).toEmit('clear:items');
  });

  it('Renders `mass-actions-panel` with non empty items', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        items,
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `mass-actions-panel` with empty items', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        items: [],
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `mass-actions-panel` with one item', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        items: [alarm],
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `mass-actions-panel` with meta alarm', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        items,
        widget,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
