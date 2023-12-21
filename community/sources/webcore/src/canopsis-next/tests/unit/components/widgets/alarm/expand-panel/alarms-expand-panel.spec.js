import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createMockedStoreModule, createMockedStoreModules } from '@unit/utils/store';

import { fakeAlarmDetails } from '@unit/data/alarm';

import { CANOPSIS_EDITION, ENTITY_TYPES, JUNIT_ALARM_CONNECTOR, USERS_PERMISSIONS } from '@/constants';

import AlarmsExpandPanel from '@/components/widgets/alarm/expand-panel/alarms-expand-panel.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';

const stubs = {
  'more-infos': true,
  'alarms-time-line': true,
  'alarms-expand-panel-children': true,
  'service-dependencies': true,
  'declared-tickets-list': true,
  'entity-gantt': true,
  'pbehaviors-simple-list': true,
  'alarms-expand-panel-charts': true,
  'c-compiled-template': CCompiledTemplate,
  'c-runtime-template': CRuntimeTemplate,
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

  const factory = generateShallowRenderer(AlarmsExpandPanel, { stubs });
  const snapshotFactory = generateRenderer(AlarmsExpandPanel, { stubs });

  afterEach(() => {
    jest.clearAllMocks();
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
            impacts_count: 0,
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

  it('Renders `alarms-expand-panel` with required props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm: {
          _id: 'alarm-id',
          entity: {
            type: ENTITY_TYPES.connector,
            impacts_count: 0,
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

    await wrapper.activateAllTabs();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm: {
          _id: 'alarm-id',
          causes: {},
          consequences: {},
          entity: {
            type: ENTITY_TYPES.connector,
            impacts_count: 0,
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
      },
    });

    await wrapper.activateAllTabs();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel` with full alarm', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm: {
          _id: 'alarm-id',
          causes: {},
          consequences: {},
          entity: {
            type: ENTITY_TYPES.service,
            impacts_count: 1,
          },
          v: {
            tickets: [{}],
          },
        },
        widget: {
          parameters: {
            charts: [{
              metric: 'cpu',
            }],
            moreInfoTemplate: 'template',
            isHtmlEnabledOnTimeLine: false,
            serviceDependenciesColumns: [{}],
          },
        },
      },
    });

    await wrapper.activateAllTabs();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel` with gantt', async () => {
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
            impacts_count: 0,
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

    await wrapper.activateAllTabs();

    expect(wrapper).toMatchSnapshot();
  });
});
