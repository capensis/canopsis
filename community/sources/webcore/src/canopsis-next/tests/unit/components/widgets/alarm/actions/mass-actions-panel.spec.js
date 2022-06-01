import Faker from 'faker';
import { range } from 'lodash';
import flushPromises from 'flush-promises';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockDateNow, mockModals } from '@unit/utils/mock-hooks';
import {
  ALARM_LIST_ACTIONS_TYPES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  ENTITIES_STATUSES,
  ENTITIES_TYPES,
  ENTITY_PATTERN_FIELDS,
  EVENT_DEFAULT_ORIGIN,
  EVENT_ENTITY_TYPES,
  EVENT_INITIATORS,
  META_ALARMS_RULE_TYPES,
  MODALS,
  PATTERN_CONDITIONS,
} from '@/constants';

import MassActionsPanel from '@/components/widgets/alarm/actions/mass-actions-panel.vue';

const localVue = createVueInstance();

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

const factory = (options = {}) => shallowMount(MassActionsPanel, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(MassActionsPanel, {
  localVue,
  stubs,

  ...options,
});

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
  };

  const metaAlarm = {
    _id: 'meta-alarm-id',
    metaalarm: true,
    entity: {
      _id: 'meta-alarm-entity-id',
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
  const entitiesModule = {
    name: 'entities',
    getters: {
      getList: () => () => [alarm, metaAlarm],
    },
  };
  const fetchAlarmsListWithPreviousParams = jest.fn();
  const alarmModule = {
    name: 'alarm',
    actions: {
      fetchListWithPreviousParams: fetchAlarmsListWithPreviousParams,
    },
  };
  const createEvent = jest.fn();
  const eventModule = {
    name: 'event',
    actions: {
      create: createEvent,
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

  it('Create pbehavior modal showed after trigger pbehavior add action', () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        entitiesModule,
      ]),
      propsData: {
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

  it('Ack modal showed after trigger ack action', () => {
    const isNoteRequired = Faker.datatype.boolean();
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {
        isAckNoteRequired: isNoteRequired,
      },
    };

    const itemsIds = [Faker.datatype.string(), Faker.datatype.string()];
    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        entitiesModule,
        alarmModule,
      ]),
      propsData: {
        itemsIds,
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
          itemsIds,
          itemsType: ENTITIES_TYPES.alarm,
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);

    expect(fetchAlarmsListWithPreviousParams).toBeCalledWith(
      expect.any(Object),
      { widgetId: widgetData._id },
      undefined,
    );
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
        alarmModule,
        eventModule,
        {
          ...entitiesModule,
          getters: {
            getList: () => () => fastAckAlarms,
          },
        },
      ]),
      propsData: {
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

    expect(fetchAlarmsListWithPreviousParams).toBeCalledWith(
      expect.any(Object),
      { widgetId: widgetData._id },
      undefined,
    );
  });

  it('Fast ack event sent after trigger fast ack action without parameters', async () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        eventModule,
        {
          ...entitiesModule,
          getters: {
            getList: () => () => fastAckAlarms,
          },
        },
      ]),
      propsData: {
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

    expect(fetchAlarmsListWithPreviousParams).toBeCalledWith(
      expect.any(Object),
      { widgetId: widgetData._id },
      undefined,
    );
  });

  it('Ack remove modal showed after trigger ack remove action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const itemsIds = [Faker.datatype.string(), Faker.datatype.string()];
    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        entitiesModule,
      ]),
      propsData: {
        itemsIds,
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
          itemsIds,
          itemsType: ENTITIES_TYPES.alarm,
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);

    expect(fetchAlarmsListWithPreviousParams).toBeCalledWith(
      expect.any(Object),
      { widgetId: widgetData._id },
      undefined,
    );
  });

  it('Cancel modal showed after trigger cancel action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const itemsIds = [Faker.datatype.string(), Faker.datatype.string()];
    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        entitiesModule,
        alarmModule,
      ]),
      propsData: {
        itemsIds,
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
          itemsIds,
          itemsType: ENTITIES_TYPES.alarm,
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

    expect(fetchAlarmsListWithPreviousParams).toBeCalledWith(
      expect.any(Object),
      { widgetId: widgetData._id },
      undefined,
    );
  });

  it('Associate ticket modal showed after trigger associate ticket action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const itemsIds = [Faker.datatype.string(), Faker.datatype.string()];
    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        entitiesModule,
        alarmModule,
      ]),
      propsData: {
        itemsIds,
        widget: widgetData,
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
          itemsIds,
          itemsType: ENTITIES_TYPES.alarm,
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);

    expect(fetchAlarmsListWithPreviousParams).toBeCalledWith(
      expect.any(Object),
      { widgetId: widgetData._id },
      undefined,
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

    const itemsIds = [Faker.datatype.string(), Faker.datatype.string()];
    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        entitiesModule,
        alarmModule,
      ]),
      propsData: {
        itemsIds,
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
          itemsIds,
          itemsType: ENTITIES_TYPES.alarm,
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);

    expect(fetchAlarmsListWithPreviousParams).toBeCalledWith(
      expect.any(Object),
      { widgetId: widgetData._id },
      undefined,
    );
  });

  it('Group alarm modal showed after trigger group alarm action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const itemsIds = [Faker.datatype.string(), Faker.datatype.string()];
    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        {
          ...entitiesModule,
          getters: {
            getList: () => () => [alarm],
          },
        },
      ]),
      propsData: {
        itemsIds,
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
          itemsIds,
          itemsType: ENTITIES_TYPES.alarm,
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);

    expect(fetchAlarmsListWithPreviousParams).toBeCalledWith(
      expect.any(Object),
      { widgetId: widgetData._id },
      undefined,
    );
  });

  it('Manual meta alarm group modal showed after trigger manual meta alarm group action', () => {
    const widgetData = {
      _id: Faker.datatype.string(),
      parameters: {},
    };

    const itemsIds = [Faker.datatype.string(), Faker.datatype.string()];
    const wrapper = factory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        alarmModule,
        {
          ...entitiesModule,
          getters: {
            getList: () => () => [alarm],
          },
        },
      ]),
      propsData: {
        itemsIds,
        widget: widgetData,
      },
      mocks: {
        $modals,
      },
    });

    const ackRemoveAction = selectActionByType(wrapper, ALARM_LIST_ACTIONS_TYPES.manualMetaAlarmGroup);

    ackRemoveAction.trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createManualMetaAlarm,
        config: {
          title: 'Manual meta alarm management',
          itemsIds,
          itemsType: ENTITIES_TYPES.alarm,
          afterSubmit: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    config.afterSubmit();

    const clearItemsEvent = wrapper.emitted('clear:items');

    expect(clearItemsEvent).toHaveLength(1);

    expect(fetchAlarmsListWithPreviousParams).toBeCalledWith(
      expect.any(Object),
      { widgetId: widgetData._id },
      undefined,
    );
  });

  it('Renders `mass-actions-panel` with empty items', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        {
          ...entitiesModule,
          getters: {
            getList: () => () => [],
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `mass-actions-panel` with meta alarm', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModuleWithAccess,
        {
          ...entitiesModule,
          getters: {
            getList: () => () => [alarm, metaAlarm],
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
