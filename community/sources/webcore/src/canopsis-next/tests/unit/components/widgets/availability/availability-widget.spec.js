import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import {
  createActiveViewModule,
  createAlarmModule,
  createAuthModule,
  createAvailabilityModule,
  createMockedStoreModules,
  createQueryModule,
  createUserPreferenceModule,
} from '@unit/utils/store';

import {
  AVAILABILITY_DISPLAY_PARAMETERS,
  AVAILABILITY_SHOW_TYPE,
  AVAILABILITY_VALUE_FILTER_METHODS,
  QUICK_RANGES,
  USERS_PERMISSIONS,
  WIDGET_TYPES,
} from '@/constants';

import { formToWidget, widgetToForm } from '@/helpers/entities/widget/form';

import AvailabilityWidget from '@/components/widgets/availability/availability-widget.vue';

const stubs = {
  'availability-widget-filters': true,
  'availability-list': true,
};

const selectAvailabilityWidgetFilters = wrapper => wrapper.find('availability-widget-filters-stub');

const generateDefaultAvailabilityWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.availability })),

  _id: `${WIDGET_TYPES.availability}_id`,
});
describe('availability-widget', () => {
  const tabId = Faker.datatype.string();

  const widget = generateDefaultAvailabilityWidget();

  const defaultQuery = {
    displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
    filter: undefined,
    interval: {
      from: QUICK_RANGES.today.start,
      to: QUICK_RANGES.today.stop,
    },
    search: '',
    showTrend: false,
    showType: AVAILABILITY_SHOW_TYPE.percent,
  };

  const { authModule, currentUserPermissionsById } = createAuthModule();
  const { activeViewModule } = createActiveViewModule();
  const { alarmModule } = createAlarmModule();

  const { userPreferenceModule, updateUserPreference } = createUserPreferenceModule();
  const { queryModule, updateQuery } = createQueryModule();
  const { availabilityModule, fetchAvailabilityList } = createAvailabilityModule();

  currentUserPermissionsById.mockReturnValue({
    [USERS_PERMISSIONS.business.availability.actions.userFilter]: { actions: [] },
  });

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    availabilityModule,
    activeViewModule,
    alarmModule,
    queryModule,
  ]);

  const factory = generateShallowRenderer(AvailabilityWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(AvailabilityWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Query changed after component mount', async () => {
    factory({
      store,
      propsData: {
        tabId,
        widget,
      },
    });

    await flushPromises();

    expect(updateQuery).toBeDispatchedWith(
      {
        id: widget._id,
        query: defaultQuery,
      },
    );
  });

  test('Data fetched when correct query', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        widget,
      },
    });

    await flushPromises();

    await wrapper.vm.fetchList();

    expect(fetchAvailabilityList).toBeDispatchedWith(
      {
        widgetId: widget._id,
        params: {
          value_filter: undefined,
          widget_filters: [],
        },
      },
    );
  });

  test('Filters changed after trigger availability filters', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const newFilter = Faker.datatype.string();

    await selectAvailabilityWidgetFilters(wrapper).triggerCustomEvent(
      'update:filters',
      newFilter,
    );

    await flushPromises();

    expect(updateUserPreference).toBeDispatchedWith({
      data: {
        content: {
          mainFilter: newFilter,
        },
      },
    });
    expect(updateQuery).toBeDispatchedWith({
      id: widget._id,
      query: {
        filter: newFilter,
        page: 1,
      },
    });
  });

  test('Interval changed after trigger availability filters', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const newInterval = {
      from: QUICK_RANGES.yesterday.start,
      to: QUICK_RANGES.yesterday.start,
    };

    await selectAvailabilityWidgetFilters(wrapper).triggerCustomEvent(
      'update:interval',
      newInterval,
    );

    await flushPromises();

    expect(updateQuery).toBeDispatchedWith({
      id: widget._id,
      query: {
        interval: newInterval,
      },
    });
  });

  test('Trend changed after trigger availability filters', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const newTrend = true;

    await selectAvailabilityWidgetFilters(wrapper).triggerCustomEvent(
      'update:trend',
      newTrend,
    );

    await flushPromises();

    expect(updateUserPreference).toBeDispatchedWith({
      data: {
        content: {
          show_trend: newTrend,
        },
      },
    });
    expect(updateQuery).toBeDispatchedWith({
      id: widget._id,
      query: {
        showTrend: newTrend,
      },
    });
  });

  test('Type changed after trigger availability filters', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const newType = AVAILABILITY_SHOW_TYPE.duration;

    await selectAvailabilityWidgetFilters(wrapper).triggerCustomEvent(
      'update:type',
      newType,
    );

    await flushPromises();

    expect(updateUserPreference).toBeDispatchedWith({
      data: {
        content: {
          show_type: newType,
        },
      },
    });
    expect(updateQuery).toBeDispatchedWith({
      id: widget._id,
      query: {
        showType: newType,
      },
    });
  });

  test('Display parameter changed after trigger availability filters', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const newDisplayParameter = AVAILABILITY_SHOW_TYPE.duration;

    await selectAvailabilityWidgetFilters(wrapper).triggerCustomEvent(
      'update:display-parameter',
      newDisplayParameter,
    );

    await flushPromises();

    expect(updateUserPreference).toBeDispatchedWith({
      data: {
        content: {
          display_parameter: newDisplayParameter,
        },
      },
    });
    expect(updateQuery).toBeDispatchedWith({
      id: widget._id,
      query: {
        displayParameter: newDisplayParameter,
      },
    });
  });

  test('Value filter changed after trigger availability filters', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const newValueFilter = {
      value: 100,
      method: AVAILABILITY_VALUE_FILTER_METHODS.greater,
    };

    await selectAvailabilityWidgetFilters(wrapper).triggerCustomEvent(
      'update:value-filter',
      newValueFilter,
    );

    await flushPromises();

    expect(updateQuery).toBeDispatchedWith({
      id: widget._id,
      query: {
        valueFilter: newValueFilter,
      },
    });
  });

  test('Renders `availability-widget` with default props', async () => {
    fetchAvailabilityList.mockReturnValue({
      data: [],
    });

    const wrapper = snapshotFactory({
      propsData: {
        tabId: 'tab-id',
        widget,
      },
      store,
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
