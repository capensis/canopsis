import { fakeAlarmDetails } from '@unit/data/alarm';
import { createMockedStoreModule, createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';

import { CANOPSIS_EDITION, ENTITY_TYPES, GRID_SIZES, USER_PERMISSIONS } from '@/constants';

import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import AlarmsExpandPanel from '@/components/widgets/alarm/expand-panel/alarms-expand-panel.vue';

const stubs = {
  'more-infos': true,
  'alarm-timeline': true,
  'alarms-expand-panel-children': true,
  'service-dependencies': true,
  'declared-tickets-list': true,
  'entity-gantt': true,
  'pbehaviors-simple-list': true,
  'alarms-expand-panel-charts': true,
  'alarms-expand-panel-remediation': true,
  'alarms-expand-panel-entity-enrichments': true,
  'entity-availability': true,
  'c-compiled-template': CCompiledTemplate,
  'c-runtime-template': CRuntimeTemplate,
};

describe('alarms-expand-panel', () => {
  // TODO: add tests for children, timeline, query
  const infoModule = {
    name: 'info',
    getters: { edition: CANOPSIS_EDITION.community },
  };

  const authModule = {
    name: 'auth',
    getters: {
      currentUserPermissionsById: {
        [USER_PERMISSIONS.business.alarmsList.actions.pbehaviorAdd]: { actions: [] },
      },
    },
  };

  const fetchAlarmDetails = vi.fn();
  const fetchAlarmsDetailsList = vi.fn();
  const updateAlarmDetailsQuery = vi.fn();
  const removeAlarmDetailsQuery = vi.fn();

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

  const store = createMockedStoreModules([authModule, alarmModule, infoModule]);

  const factory = generateShallowRenderer(AlarmsExpandPanel, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it.each([undefined, null, []])(
    'Uses the full grid width when expandGridRangeSize is %s',
    (expandGridRangeSize) => {
      const range = AlarmsExpandPanel.computed.expandGridRangeSize.call({
        widget: {
          parameters: {
            expandGridRangeSize,
          },
        },
      });

      expect(range).toEqual([GRID_SIZES.min, GRID_SIZES.max]);
    },
  );

  it('Uses the configured expand grid range size', () => {
    const expandGridRangeSize = [2, 10];
    const range = AlarmsExpandPanel.computed.expandGridRangeSize.call({
      widget: {
        parameters: {
          expandGridRangeSize,
        },
      },
    });

    expect(range).toBe(expandGridRangeSize);
  });

  it('More-info template id is derived from the widget id', () => {
    const widget = {
      _id: 'widget-id',
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

    expect(wrapper.vm.moreInfosTemplateId).toBe('widget-id-moreInfoTemplate');
  });
});
