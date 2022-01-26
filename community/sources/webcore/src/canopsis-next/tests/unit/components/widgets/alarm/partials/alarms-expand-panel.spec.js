import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createMockedStoreModules } from '@unit/utils/store';
import AlarmsExpandPanel from '@/components/widgets/alarm/partials/alarms-expand-panel.vue';
import { CANOPSIS_EDITION, ENTITY_TYPES, JUNIT_ALARM_CONNECTOR } from '@/constants';

const localVue = createVueInstance();

const stubs = {
  'more-infos': true,
  'time-line': true,
  'group-alarms-list': true,
  'service-dependencies': true,
  'entity-gantt': true,
};

const factory = (options = {}) => shallowMount(AlarmsExpandPanel, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(AlarmsExpandPanel, {
  localVue,
  stubs,

  ...options,
});

const selectTabs = wrapper => wrapper.find('v-tabs-stub');

describe('alarms-expand-panel', () => {
  const infoModule = {
    name: 'info',
    getters: { edition: CANOPSIS_EDITION.core },
  };
  const catInfoModule = {
    name: 'info',
    getters: { edition: CANOPSIS_EDITION.cat },
  };

  it('Tab key updated after change tour enabled', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        infoModule,
      ]),
      propsData: {
        alarm: {
          _id: 'alarm-id',
          entity: {
            type: ENTITY_TYPES.connector,
            impact: [],
          },
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

    // eslint-disable-next-line no-underscore-dangle
    const prevKey = selectTabs(wrapper).vm._vnode.key;

    await wrapper.setProps({
      isTourEnabled: true,
    });

    // eslint-disable-next-line no-underscore-dangle
    expect(prevKey !== selectTabs(wrapper).vm._vnode.key).toBe(true);
  });

  it('Tab key updated after change moreInfoTemplate', async () => {
    const widget = {
      parameters: {
        moreInfoTemplate: 'template',
      },
    };
    const wrapper = factory({
      store: createMockedStoreModules([
        infoModule,
      ]),
      propsData: {
        alarm: {
          _id: 'alarm-id',
          entity: {
            type: ENTITY_TYPES.connector,
            impact: [],
          },
        },
        widget,
      },
    });

    // eslint-disable-next-line no-underscore-dangle
    const prevKey = selectTabs(wrapper).vm._vnode.key;

    await wrapper.setProps({
      widget: {
        parameters: {
          moreInfoTemplate: 'template2',
        },
      },
    });

    // eslint-disable-next-line no-underscore-dangle
    expect(prevKey !== selectTabs(wrapper).vm._vnode.key).toBe(true);
  });

  it('Renders `alarms-expand-panel` with required props', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        infoModule,
      ]),
      propsData: {
        alarm: {
          _id: 'alarm-id',
          entity: {
            type: ENTITY_TYPES.connector,
            impact: [],
          },
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
      store: createMockedStoreModules([
        infoModule,
      ]),
      propsData: {
        alarm: {
          _id: 'alarm-id',
          causes: {},
          consequences: {},
          entity: {
            type: ENTITY_TYPES.connector,
            impact: [],
          },
        },
        widget: {
          parameters: {
            moreInfoTemplate: 'template',
            isHtmlEnabledOnTimeLine: false,
            serviceDependenciesColumns: [],
          },
        },
        isEditingMode: true,
        hideGroups: true,
        isTourEnabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel` with full alarm', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        infoModule,
      ]),
      propsData: {
        alarm: {
          _id: 'alarm-id',
          causes: {},
          consequences: {},
          entity: {
            type: ENTITY_TYPES.service,
            impact: ['test'],
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
        catInfoModule,
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
