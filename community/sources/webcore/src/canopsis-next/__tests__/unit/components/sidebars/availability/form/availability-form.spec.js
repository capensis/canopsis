import {
  getWidgetRequestWithNewParametersProperty,
  getWidgetRequestWithNewProperty,
} from '@unit/utils/settings';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { AVAILABILITY_DISPLAY_PARAMETERS, QUICK_RANGES, WIDGET_TYPES } from '@/constants';

import { widgetToForm } from '@/helpers/entities/widget/form';

import AvailabilityForm from '@/components/sidebars/availability/form/availability-form.vue';

const stubs = {
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-columns': true,
  'field-quick-date-interval-type': true,
  'field-availability-display-parameter': true,
  'field-availability-display-show-type': true,
  'export-csv-form': true,
  'field-filters': true,
};

const selectFieldTitle = wrapper => wrapper.find('field-title-stub');
const selectFieldPeriodicRefresh = wrapper => wrapper.find('field-periodic-refresh-stub');
const selectFieldColumns = wrapper => wrapper.findAll('field-columns-stub');
const selectWidgetColumnsField = wrapper => selectFieldColumns(wrapper).at(0);
const selectActiveAlarmsColumnsField = wrapper => selectFieldColumns(wrapper).at(1);
const selectResolvedAlarmsColumnsField = wrapper => selectFieldColumns(wrapper).at(2);
const selectFieldQuickDateIntervalType = wrapper =>
  wrapper.find('field-quick-date-interval-type-stub');
const selectFieldAvailabilityDisplayParameter = wrapper =>
  wrapper.find('field-availability-display-parameter-stub');
const selectFieldFilters = wrapper => wrapper.find('field-filters-stub');

describe('availability-form', () => {
  const factory = generateShallowRenderer(AvailabilityForm, { stubs });

  test('Title changed after trigger title field', () => {
    const widget = widgetToForm({ type: WIDGET_TYPES.availability });

    const wrapper = factory({
      propsData: {
        form: widget,
        widgetId: 'widget-id',
      },
    });

    const newTitle = Faker.datatype.string();

    selectFieldTitle(wrapper).triggerCustomEvent('input', newTitle);

    expect(wrapper).toEmitInput(getWidgetRequestWithNewProperty(widget, 'title', newTitle));
  });

  test('Periodic refresh changed after trigger periodic refresh field', () => {
    const widget = widgetToForm({ type: WIDGET_TYPES.availability });

    const wrapper = factory({
      propsData: {
        form: widget,
        widgetId: 'widget-id',
      },
    });

    const newPeriodicRefresh = {
      enabled: Faker.datatype.boolean(),
      value: Faker.datatype.number(),
      unit: Faker.datatype.string(),
    };

    selectFieldPeriodicRefresh(wrapper).triggerCustomEvent('input', newPeriodicRefresh);

    expect(wrapper).toEmitInput(
      getWidgetRequestWithNewParametersProperty(widget, 'periodic_refresh', newPeriodicRefresh),
    );
  });

  test('Widget columns changed after trigger widget columns field', () => {
    const widget = widgetToForm({ type: WIDGET_TYPES.availability });

    const wrapper = factory({
      propsData: {
        form: widget,
        widgetId: 'widget-id',
      },
    });

    const newWidgetColumns = [
      {
        value: Faker.lorem.word(),
        label: Faker.lorem.word(),
      },
    ];

    selectWidgetColumnsField(wrapper).triggerCustomEvent('input', newWidgetColumns);

    expect(wrapper).toEmitInput(
      getWidgetRequestWithNewParametersProperty(widget, 'widget_columns', newWidgetColumns),
    );
  });

  test('Active alarms columns changed after trigger columns field', () => {
    const widget = widgetToForm({ type: WIDGET_TYPES.availability });

    const wrapper = factory({
      propsData: {
        form: widget,
        widgetId: 'widget-id',
      },
    });

    const newWidgetColumns = [
      {
        value: Faker.lorem.word(),
        label: Faker.lorem.word(),
      },
    ];

    selectActiveAlarmsColumnsField(wrapper).triggerCustomEvent('input', newWidgetColumns);

    expect(wrapper).toEmitInput(
      getWidgetRequestWithNewParametersProperty(widget, 'active_alarms_columns', newWidgetColumns),
    );
  });

  test('REsolved alarms columns changed after trigger columns field', () => {
    const widget = widgetToForm({ type: WIDGET_TYPES.availability });

    const wrapper = factory({
      propsData: {
        form: widget,
        widgetId: 'widget-id',
      },
    });

    const newWidgetColumns = [
      {
        value: Faker.lorem.word(),
        label: Faker.lorem.word(),
      },
    ];

    selectResolvedAlarmsColumnsField(wrapper).triggerCustomEvent('input', newWidgetColumns);

    expect(wrapper).toEmitInput(
      getWidgetRequestWithNewParametersProperty(
        widget,
        'resolved_alarms_columns',
        newWidgetColumns,
      ),
    );
  });

  test('Default time range changed after trigger quick interval type field', () => {
    const widget = widgetToForm({ type: WIDGET_TYPES.availability });

    const wrapper = factory({
      propsData: {
        form: widget,
        widgetId: 'widget-id',
        showInterval: true,
      },
    });

    const newValue = QUICK_RANGES.thisMonth.value;

    selectFieldQuickDateIntervalType(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(
      getWidgetRequestWithNewParametersProperty(widget, 'default_time_range', newValue),
    );
  });

  test('Default display parameter changed after trigger default display parameter field', () => {
    const widget = widgetToForm({ type: WIDGET_TYPES.availability });

    const wrapper = factory({
      propsData: {
        form: widget,
        widgetId: 'widget-id',
      },
    });

    const newValue = AVAILABILITY_DISPLAY_PARAMETERS.downtime;

    selectFieldAvailabilityDisplayParameter(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(
      getWidgetRequestWithNewParametersProperty(widget, 'default_display_parameter', newValue),
    );
  });

  test('Filters changed after trigger filters field', () => {
    const widget = widgetToForm({ type: WIDGET_TYPES.availability });

    const wrapper = factory({
      propsData: {
        form: widget,
        widgetId: 'widget-id',
        showFilter: true,
      },
    });

    const newValue = {
      title: Faker.datatype.string(),
    };

    selectFieldFilters(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(
      getWidgetRequestWithNewParametersProperty(widget, 'mainFilter', newValue),
    );
  });
});
