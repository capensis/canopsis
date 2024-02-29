import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';
import { getWidgetRequestWithNewParametersProperty, getWidgetRequestWithNewProperty } from '@unit/utils/settings';

import { AVAILABILITY_DISPLAY_PARAMETERS, QUICK_RANGES, WIDGET_TYPES } from '@/constants';

import { widgetToForm } from '@/helpers/entities/widget/form';

import AvailabilityForm from '@/components/sidebars/availability/form/availability-form.vue';

const stubs = {
  'widget-settings-group': true,
  'field-title': createInputStub('field-title'),
  'field-periodic-refresh': createInputStub('field-periodic-refresh'),
  'field-columns': createInputStub('field-columns'),
  'field-quick-date-interval-type': createInputStub('field-quick-date-interval-type'),
  'field-availability-display-parameter': createInputStub('field-availability-display-parameter'),
  'field-availability-display-show-type': createInputStub('field-availability-display-show-type'),
  'field-filters': createInputStub('field-filters'),
};
const snapshotStubs = {
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-columns': true,
  'field-quick-date-interval-type': true,
  'field-availability-display-parameter': true,
  'field-availability-display-show-type': true,
  'field-filters': true,
};

const selectFieldTitle = wrapper => wrapper.find('input.field-title');
const selectFieldPeriodicRefresh = wrapper => wrapper.find('input.field-periodic-refresh');
const selectFieldColumnsByLabel = (wrapper, label) => wrapper.find(`input.field-columns[label="${label}"]`);
const selectWidgetColumnsField = wrapper => selectFieldColumnsByLabel(wrapper, 'Column names');
const selectActiveAlarmsColumnsField = wrapper => selectFieldColumnsByLabel(
  wrapper,
  'Column names for active alarms',
);
const selectFieldQuickDateIntervalType = wrapper => wrapper.find('input.field-quick-date-interval-type');
const selectFieldAvailabilityDisplayParameter = wrapper => wrapper.find('input.field-availability-display-parameter');
const selectFieldFilters = wrapper => wrapper.find('input.field-filters');

describe('availability-form', () => {
  const factory = generateShallowRenderer(AvailabilityForm, { stubs });
  const snapshotFactory = generateRenderer(AvailabilityForm, { stubs: snapshotStubs });

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

    expect(wrapper).toEmitInput(
      getWidgetRequestWithNewProperty(widget, 'title', newTitle),
    );
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

    selectFieldPeriodicRefresh(wrapper).triggerCustomEvent('input', {
      ...widget.parameters,
      periodic_refresh: newPeriodicRefresh,
    });

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

    const newWidgetColumns = [{
      value: Faker.lorem.word(),
      label: Faker.lorem.word(),
    }];

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

    const newWidgetColumns = [{
      value: Faker.lorem.word(),
      label: Faker.lorem.word(),
    }];

    selectActiveAlarmsColumnsField(wrapper).triggerCustomEvent('input', newWidgetColumns);

    expect(wrapper).toEmitInput(
      getWidgetRequestWithNewParametersProperty(widget, 'active_alarms_columns', newWidgetColumns),
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

  test('Renders `availability-form` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: widgetToForm({ type: WIDGET_TYPES.availability }),
        widgetId: 'widget-id',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-form` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: widgetToForm({ type: WIDGET_TYPES.availability }),
        widgetId: 'widget-id',
        entityColumnsWidgetTemplates: [{}],
        entityColumnsWidgetTemplatesPending: true,
        alarmColumnsWidgetTemplates: [{}, {}],
        alarmColumnsWidgetTemplatesPending: true,
        filterAddable: true,
        filterEditable: true,
        showInterval: true,
        showFilter: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
