import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createMockedStoreModule, createMockedStoreModules } from '@unit/utils/store';

import { fakeAlarmDetails } from '@unit/data/alarm';

import {
  CANOPSIS_EDITION,
  ENTITY_TYPES,
  JUNIT_ALARM_CONNECTOR,
  USERS_PERMISSIONS,
} from '@/constants';

import AlarmsExpandPanel from '@/components/widgets/alarm/expand-panel/alarms-expand-panel.vue';

const stubs = {
  'more-infos': true,
  'alarms-time-line': true,
  'alarms-expand-panel-children': true,
  'service-dependencies': true,
  'declared-tickets-list': true,
  'entity-gantt': true,
  'pbehaviors-simple-list': true,
};

const selectTabs = wrapper => wrapper.vm.$children[0];

describe('alarms-expand-panel', () => { // TODO: add tests for children, timeline, query
  const infoModule = {
    name: 'info',
    getters: { edition: CANOPSIS_EDITION.community },
  };

  const catInfoModule = {
    name: 'info',
    getters: { edition: CANOPSIS_EDITION.pro },
  };
  const authModule = {
    name: 'auth',
    getters: {
      currentUserPermissionsById: {
        [USERS_PERMISSIONS.technical.exploitation.pbehavior]: { actions: [] },
      },
    },
  };

  const fetchAlarmDetails = jest.fn();
  const fetchAlarmsDetailsList = jest.fn();
  const updateAlarmDetailsQuery = jest.fn();
  const removeAlarmDetailsQuery = jest.fn();

  const alarmDetailsModule = createMockedStoreModule({
    name: 'details',
    getters: {
      getItem: () => () => fakeAlarmDetails(),
      getPending: () => () => false,
      getQuery: () => () => ({ page: 1, limit: 10 }),
      getQueries: () => () => [
        { page: 2, limit: 5 },
        { page: 1, limit: 10 },
      ],
    },
    actions: {
      fetchItem: fetchAlarmDetails,
      fetchList: fetchAlarmsDetailsList,
      updateQuery: updateAlarmDetailsQuery,
      removeQuery: removeAlarmDetailsQuery,
    },
  });

  const alarmModule = {
    name: 'alarm',
    modules: {
      details: alarmDetailsModule,
    },
  };

  const store = createMockedStoreModules([
    authModule,
    alarmModule,
    infoModule,
  ]);

  const factory = generateShallowRenderer(AlarmsExpandPanel, {
    stubs,
  });
  const snapshotFactory = generateRenderer(AlarmsExpandPanel, {
    stubs,
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('Tabs key updated after change tour enabled', async () => {
    const wrapper = factory({
      store,
      propsData: {
        alarm: {
          _id: 'alarm-id',
          entity: {
            type: ENTITY_TYPES.connector,
            impact: [],
          },
          v: {},
        },
        widget: {
          parameters: {
            moreInfoTemplate: 'template',
            isHtmlEnabledOnTimeLine: false,
            serviceDependenciesColumns: [],
          },
        },
      },
    });

    const prevKey = selectTabs(wrapper).$vnode.key;

    await wrapper.setProps({
      isTourEnabled: true,
    });

    expect(prevKey !== selectTabs(wrapper).$vnode.key).toBe(true);
  });

  it('Tab key updated after change moreInfoTemplate', async () => {
    const widget = {
      parameters: {
        moreInfoTemplate: 'template',
      },
    };
    const wrapper = factory({
      store,
      propsData: {
        alarm: {
          _id: 'alarm-id',
          entity: {
            type: ENTITY_TYPES.connector,
            impact: [],
          },
          v: {},
        },
        widget,
      },
    });

    const prevKey = selectTabs(wrapper).$vnode.key;

    await wrapper.setProps({
      widget: {
        parameters: {
          moreInfoTemplate: 'template2',
        },
      },
    });

    expect(prevKey !== selectTabs(wrapper).$vnode.key).toBe(true);
  });

  it('Renders `alarms-expand-panel` with required props', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm: {
          _id: 'alarm-id',
          entity: {
            type: ENTITY_TYPES.connector,
            impact: [],
          },
          v: {},
        },
        widget: {
          parameters: {
            moreInfoTemplate: 'template',
            isHtmlEnabledOnTimeLine: false,
            serviceDependenciesColumns: [],
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel` with custom props', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm: {
          _id: 'alarm-id',
          causes: {},
          consequences: {},
          entity: {
            type: ENTITY_TYPES.connector,
            impact: [],
          },
          v: {},
        },
        widget: {
          parameters: {
            moreInfoTemplate: 'template',
            isHtmlEnabledOnTimeLine: false,
            serviceDependenciesColumns: [],
          },
        },
        hideChildren: true,
        editing: true,
        isTourEnabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel` with full alarm', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm: {
          _id: 'alarm-id',
          causes: {},
          consequences: {},
          entity: {
            type: ENTITY_TYPES.service,
            impact: ['test'],
          },
          v: {
            tickets: [{}],
          },
        },
        widget: {
          parameters: {
            moreInfoTemplate: 'template',
            isHtmlEnabledOnTimeLine: false,
            serviceDependenciesColumns: [{}],
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel` with gantt', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        alarmModule,
        catInfoModule,
        authModule,
      ]),
      propsData: {
        alarm: {
          _id: 'alarm-id',
          v: {
            connector: JUNIT_ALARM_CONNECTOR,
          },
          entity: {
            type: ENTITY_TYPES.resource,
            impact: [],
          },
        },
        widget: {
          parameters: {
            moreInfoTemplate: 'template',
            isHtmlEnabledOnTimeLine: false,
            serviceDependenciesColumns: [{}],
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
