import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { randomArrayItem } from '@unit/utils/array';
import { createMockedStoreModules, createUserPreferenceModule } from '@unit/utils/store';

import {
  ADVANCED_SEARCH_FIELDS,
  AVAILABILITY_DISPLAY_PARAMETERS,
  AVAILABILITY_SHOW_TYPE,
  AVAILABILITY_VALUE_FILTER_METHODS,
  QUICK_RANGES,
} from '@/constants';

import { formToAdvancedSearch } from '@/helpers/search/advanced-search';
import { uuid } from '@/helpers/uuid';

import AvailabilityWidgetFilters from '@/components/widgets/availability/partials/availability-widget-filters.vue';

const stubs = {
  'c-availability-advanced-search': true,
  'c-quick-date-interval-field': true,
  'filter-selector': true,
  'filters-list-btn': true,
  'availability-display-parameter-field': true,
  'availability-show-type-field': true,
  'c-enabled-field': true,
  'availability-value-filter-field': true,
  'c-action-btn': true,
};

const selectSearch = wrapper => wrapper.find('c-availability-advanced-search-stub');
const selectQuickDateIntervalField = wrapper => wrapper.find('c-quick-date-interval-field-stub');
const selectAvailabilityShowTypeField = wrapper => wrapper.find('availability-show-type-field-stub');
const selectAvailabilityDisplayParameterField = wrapper => wrapper.find('availability-display-parameter-field-stub');
const selectFilterSelector = wrapper => wrapper.find('filter-selector-stub');
const selectTrendField = wrapper => wrapper.find('c-enabled-field-stub');
const selectExportButton = wrapper => wrapper.find('c-action-btn-stub');
const selectAvailabilityValueFilterField = wrapper => wrapper.find('availability-value-filter-field-stub');

describe('availability-widget-filters', () => {
  jest.useFakeTimers({ now: 1386435500000 });

  const widgetId = 'widget-id';

  const { userPreferenceModule } = createUserPreferenceModule();
  const store = createMockedStoreModules([userPreferenceModule]);

  const factory = generateShallowRenderer(AvailabilityWidgetFilters, {
    stubs,
    store,
  });
  const snapshotFactory = generateRenderer(AvailabilityWidgetFilters, {
    stubs,
    store,
  });

  test('Interval changed after trigger quick date interval field', async () => {
    const wrapper = factory({
      propsData: {
        widgetId,
        query: {
          interval: {
            from: QUICK_RANGES.today.start,
            to: QUICK_RANGES.today.stop,
          },
        },
        showInterval: true,
      },
    });

    const newInterval = {
      from: QUICK_RANGES.today.start,
      to: QUICK_RANGES.today.stop,
    };

    await selectQuickDateIntervalField(wrapper).triggerCustomEvent('input', newInterval);

    expect(wrapper).toEmit('update:interval', newInterval);
  });

  test('Search changed after trigger search field', async () => {
    const wrapper = factory({
      propsData: {
        widgetId,
        query: {},
        showInterval: true,
      },
    });

    const newSearch = Faker.lorem.word();
    const advancedSearch = {
      ...formToAdvancedSearch([], ADVANCED_SEARCH_FIELDS.entity),
      search: newSearch,
      _id: uuid(),
      pinned: false,
    };

    await selectSearch(wrapper).triggerCustomEvent('submit', advancedSearch);

    expect(wrapper).toEmit('update:query', expect.objectContaining({ search: newSearch, page: 1 }));
  });

  test('Filters changed after trigger filters field', async () => {
    const wrapper = factory({
      propsData: {
        widgetId,
        showFilter: true,
      },
    });

    const newFilter = Faker.datatype.string();

    await selectFilterSelector(wrapper).triggerCustomEvent('input', newFilter);

    expect(wrapper).toEmit('update:filters', newFilter);
  });

  test('Filters changed after trigger filters field', async () => {
    const wrapper = factory({
      propsData: {
        widgetId,
        showFilter: true,
      },
    });

    const newFilter = Faker.datatype.string();

    await selectFilterSelector(wrapper).triggerCustomEvent('input', newFilter);

    expect(wrapper).toEmit('update:filters', newFilter);
  });

  test('Display parameter changed after trigger availability display parameter field', async () => {
    const wrapper = factory({
      propsData: { widgetId },
    });

    const newDisplayParameter = randomArrayItem(Object.values(AVAILABILITY_DISPLAY_PARAMETERS));

    await selectAvailabilityDisplayParameterField(wrapper).triggerCustomEvent('input', newDisplayParameter);

    expect(wrapper).toEmit('update:display-parameter', newDisplayParameter);
  });

  test('Show type changed after trigger availability show type field', async () => {
    const wrapper = factory({
      propsData: { widgetId },
    });

    const newShowType = randomArrayItem(Object.values(AVAILABILITY_SHOW_TYPE));

    await selectAvailabilityShowTypeField(wrapper).triggerCustomEvent('input', newShowType);

    expect(wrapper).toEmit('update:type', newShowType);
  });

  test('Trend changed after trigger enabled field', async () => {
    const wrapper = factory({
      propsData: { widgetId },
    });

    const newTrend = Faker.datatype.boolean();

    await selectTrendField(wrapper).triggerCustomEvent('input', newTrend);

    expect(wrapper).toEmit('update:trend', newTrend);
  });

  test('Value filter changed after trigger value filter field', async () => {
    const wrapper = factory({
      propsData: { widgetId },
    });

    const newValueFilter = {
      method: randomArrayItem(Object.values(AVAILABILITY_VALUE_FILTER_METHODS)),
      value: Faker.datatype.number({ min: 0, max: 100 }),
    };

    await selectAvailabilityValueFilterField(wrapper).triggerCustomEvent('input', newValueFilter);

    expect(wrapper).not.toHaveBeenEmit('update:value-filter');

    jest.runAllTimers();

    expect(wrapper).toEmit('update:value-filter', newValueFilter);
  });

  test('Local value filter changed after update props', async () => {
    const valueFilter = {
      method: randomArrayItem(Object.values(AVAILABILITY_VALUE_FILTER_METHODS)),
      value: Faker.datatype.number({ min: 0, max: 100 }),
    };

    const wrapper = factory({
      propsData: {
        widgetId,
        query: { valueFilter },
      },
    });

    expect(wrapper.vm.localValueFilter).toEqual(valueFilter);

    const newValueFilter = {
      ...valueFilter,
      value: valueFilter.value + 1,
    };

    await wrapper.setProps({
      query: { valueFilter: newValueFilter },
    });

    expect({ ...wrapper.vm.localValueFilter }).toEqual(newValueFilter);
  });

  test('Export emitted after click on button', async () => {
    const wrapper = factory({
      propsData: {
        widgetId,
        showExport: true,
      },
    });

    await selectExportButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toHaveBeenEmit('export');
  });

  test('Renders `availability-widget-filters` with default props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        widgetId,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-widget-filters` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        widgetId,
        query: {
          search: 'Custom search',
          interval: {
            from: QUICK_RANGES.yesterday.start,
            to: QUICK_RANGES.yesterday.stop,
          },
          filter: 'test-filter',
          lockedFilter: 'locked-filter',
          displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
          showType: AVAILABILITY_SHOW_TYPE.percent,
          showTrend: true,
          valueFilter: {
            method: AVAILABILITY_VALUE_FILTER_METHODS.less,
            value: 80,
          },
        },
        userFilters: [{}],
        widgetFilters: [{}, {}],
        showInterval: true,
        showFilter: true,
        filterDisabled: true,
        exporting: true,
        maxValueFilterSeconds: 123,
        showExport: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
