import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_FIELDS, DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';
import {
  getAlarmsListWidgetColumnComponentGetter,
  getAlarmsListWidgetColumnValueFilter,
} from '@/helpers/entities/widget/forms/alarm';

import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import AlarmColumnCell from '@/components/widgets/alarm/columns-formatting/alarm-column-cell.vue';

const stubs = {
  'alarm-column-cell-popup-body': true,
  'alarm-column-value-status': true,
  'color-indicator-wrapper': true,
  'alarm-column-value-categories': true,
  'alarm-column-value-extra-details': true,
  'c-alarm-state-chip': true,
  'c-alarm-links-chips': true,
  'c-runtime-template': CRuntimeTemplate,
  'c-compiled-template': CCompiledTemplate,
  'c-ellipsis': true,
};

const selectEllipsis = wrapper => wrapper.find('c-ellipsis-stub');

describe('alarm-column-cell', () => {
  const timestamp = 1641768553245;
  const duration = 164176;
  const widget = {
    parameters: {},
  };

  const factory = generateShallowRenderer(AlarmColumnCell, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $selectAdvancedSearchField: (...args) =>
          globalThis.__TEST_SELECT_ADVANCED_SEARCH_FIELD__?.(...args),
      },
    },
  });

  afterEach(() => {
    delete globalThis.__TEST_SELECT_ADVANCED_SEARCH_FIELD__;
  });

  it.each([
    ALARM_FIELDS.lastUpdateDate,
    ALARM_FIELDS.creationDate,
    ALARM_FIELDS.lastEventDate,
    ALARM_FIELDS.activationDate,
    ALARM_FIELDS.stateAt,
    ALARM_FIELDS.statusAt,
    ALARM_FIELDS.resolved,
    ALARM_FIELDS.timestamp,
  ])('Default filter for date field: "%s" converted value to time', async (field) => {
    const column = {
      value: field,
      filter: getAlarmsListWidgetColumnValueFilter(field),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: field }),
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          t: timestamp,
          v: {
            last_update_date: timestamp,
            creation_date: timestamp,
            last_event_date: timestamp,
            activation_date: timestamp,
            resolved: timestamp,
            state: {
              t: timestamp,
            },
            status: {
              t: timestamp,
            },
          },
        },
        widget,
        column,
      },
    });

    const ellipsis = selectEllipsis(wrapper);

    expect(ellipsis.attributes('text')).toBe('10/01/2022 04:49:13');
  });

  it.each([
    ALARM_FIELDS.duration,
    ALARM_FIELDS.currentStateDuration,
    ALARM_FIELDS.activeDuration,
    ALARM_FIELDS.snoozeDuration,
    ALARM_FIELDS.pbhInactiveDuration,
  ])('Default filter for duration field: "%s" converted value to duration', async (field) => {
    const column = {
      value: field,
      filter: getAlarmsListWidgetColumnValueFilter(field),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: field }),
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          v: {
            duration,
            current_state_duration: duration,
            active_duration: duration,
            snooze_duration: duration,
            pbh_inactive_duration: duration,
          },
        },
        widget,
        column,
      },
    });

    const ellipsis = selectEllipsis(wrapper);

    expect(ellipsis.attributes('text')).toBe('1 day 21 hrs 36 mins 16 secs');
  });

  it('Custom filter for field converted value correctly', async () => {
    const filter = value => convertDateToString(value, DATETIME_FORMATS.short);
    const columnValue = 'custom_field';
    const column = {
      value: columnValue,
      filter,
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: columnValue }),
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          custom_field: duration,
        },
        widget,
        column,
      },
    });

    const ellipsis = selectEllipsis(wrapper);

    expect(ellipsis.attributes('text')).toBe('03/01/1970');
  });

  it('Applies filter class when column has isFilter set to true', async () => {
    const column = {
      value: ALARM_FIELDS.displayName,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.displayName),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.displayName }),
      isFilter: true,
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          [ALARM_FIELDS.displayName]: 'Test Name',
        },
        widget,
        column,
      },
    });

    const ellipsis = selectEllipsis(wrapper);

    expect(ellipsis.classes('alarms-column-cell__filter')).toBe(true);
  });

  it('Does not apply filter class when column has isFilter set to false', async () => {
    const column = {
      value: ALARM_FIELDS.displayName,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.displayName),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.displayName }),
      isFilter: false,
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          [ALARM_FIELDS.displayName]: 'Test Name',
        },
        widget,
        column,
      },
    });

    const ellipsis = selectEllipsis(wrapper);

    expect(ellipsis.classes('alarms-column-cell__filter')).toBe(false);
  });

  it('Calls selectAdvancedSearchField when filter column is clicked', async () => {
    const mockSelectAdvancedSearchField = vi.fn();
    globalThis.__TEST_SELECT_ADVANCED_SEARCH_FIELD__ = mockSelectAdvancedSearchField;
    const column = {
      value: ALARM_FIELDS.displayName,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.displayName),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.displayName }),
      isFilter: true,
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          [ALARM_FIELDS.displayName]: 'Test Name',
        },
        widget,
        column,
      },
      provide: {
        $selectAdvancedSearchField: mockSelectAdvancedSearchField,
      },
    });

    // The component with isFilter=true should render a clickable ellipsis
    const ellipsis = wrapper.find('c-ellipsis-stub');
    expect(ellipsis.exists()).toBe(true);
    expect(ellipsis.classes()).toContain('alarms-column-cell__filter');

    // Manually trigger the click event by calling the vm method
    await wrapper.vm.$nextTick();

    // Simulate click by manually calling the handler
    const clickHandler = wrapper.vm.componentOn.click;

    if (clickHandler) {
      await clickHandler();
    }

    expect(mockSelectAdvancedSearchField).toHaveBeenCalledWith(
      ALARM_FIELDS.displayName,
      'Test Name',
    );
  });

  it('Calls selectAdvancedSearchField with sanitized HTML value when column is HTML and filter', async () => {
    const mockSelectAdvancedSearchField = vi.fn();
    globalThis.__TEST_SELECT_ADVANCED_SEARCH_FIELD__ = mockSelectAdvancedSearchField;
    const column = {
      value: ALARM_FIELDS.entityName,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.entityName),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.entityName }),
      isFilter: true,
      isHtml: true,
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          entity: {
            name: '<script>alert("xss")</script>Safe Name',
          },
        },
        widget,
        column,
      },
      provide: {
        $selectAdvancedSearchField: mockSelectAdvancedSearchField,
      },
    });

    // For HTML columns, use manual approach to trigger click
    await wrapper.vm.$nextTick();

    // Simulate click by manually calling the handler
    const clickHandler = wrapper.vm.componentOn.click;

    if (clickHandler) {
      await clickHandler();
    }

    expect(mockSelectAdvancedSearchField).toHaveBeenCalledWith(
      ALARM_FIELDS.entityName,
      'Safe Name',
    );
  });

  it('Calls selectAdvancedSearchField with filtered value when column has custom filter', async () => {
    const mockSelectAdvancedSearchField = vi.fn();
    globalThis.__TEST_SELECT_ADVANCED_SEARCH_FIELD__ = mockSelectAdvancedSearchField;
    const customFilter = value => `filtered_${value}`;
    const column = {
      value: 'custom_field',
      filter: customFilter,
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: 'custom_field' }),
      isFilter: true,
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          custom_field: 'original_value',
        },
        widget,
        column,
      },
      provide: {
        $selectAdvancedSearchField: mockSelectAdvancedSearchField,
      },
    });

    // Manually trigger the click event by calling the vm method
    await wrapper.vm.$nextTick();

    // Simulate click by manually calling the handler
    const clickHandler = wrapper.vm.componentOn.click;

    if (clickHandler) {
      await clickHandler();
    }

    expect(mockSelectAdvancedSearchField).toHaveBeenCalledWith(
      'custom_field',
      'filtered_original_value',
    );
  });

  it('Does not call selectAdvancedSearchField when column is clicked without filter styling', async () => {
    const mockSelectAdvancedSearchField = vi.fn();
    globalThis.__TEST_SELECT_ADVANCED_SEARCH_FIELD__ = mockSelectAdvancedSearchField;
    const column = {
      value: ALARM_FIELDS.displayName,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.displayName),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.displayName }),
      isFilter: false,
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          [ALARM_FIELDS.displayName]: 'Test Name',
        },
        widget,
        column,
      },
      provide: {
        $selectAdvancedSearchField: mockSelectAdvancedSearchField,
      },
    });

    const ellipsis = selectEllipsis(wrapper);

    await ellipsis.trigger('click');

    expect(mockSelectAdvancedSearchField).not.toHaveBeenCalled();
  });

  it('Calls both custom click handler and selectAdvancedSearchField when column has both', async () => {
    const mockSelectAdvancedSearchField = vi.fn();
    globalThis.__TEST_SELECT_ADVANCED_SEARCH_FIELD__ = mockSelectAdvancedSearchField;
    const mockCustomClickHandler = vi.fn();
    const column = {
      value: ALARM_FIELDS.displayName,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.displayName),
      getComponent: props => ({
        ...getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.displayName })(props),
        on: { click: mockCustomClickHandler },
      }),
      isFilter: true,
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          [ALARM_FIELDS.displayName]: 'Test Name',
        },
        widget,
        column,
      },
      provide: {
        $selectAdvancedSearchField: mockSelectAdvancedSearchField,
      },
    });

    // Manually trigger the click event by calling the vm method
    await wrapper.vm.$nextTick();

    // Simulate click by manually calling the handler
    const clickHandler = wrapper.vm.componentOn.click;

    if (clickHandler) {
      await clickHandler();
    }

    expect(mockSelectAdvancedSearchField).toHaveBeenCalledWith(
      ALARM_FIELDS.displayName,
      'Test Name',
    );
    expect(mockCustomClickHandler).toHaveBeenCalled();
  });

  it('Handles empty or null values when filter column is clicked', async () => {
    const mockSelectAdvancedSearchField = vi.fn();
    globalThis.__TEST_SELECT_ADVANCED_SEARCH_FIELD__ = mockSelectAdvancedSearchField;
    const column = {
      value: ALARM_FIELDS.displayName,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.displayName),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.displayName }),
      isFilter: true,
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          [ALARM_FIELDS.displayName]: null,
        },
        widget,
        column,
      },
      provide: {
        $selectAdvancedSearchField: mockSelectAdvancedSearchField,
      },
    });

    // Manually trigger the click event by calling the vm method
    await wrapper.vm.$nextTick();

    // Simulate click by manually calling the handler
    const clickHandler = wrapper.vm.componentOn.click;

    if (clickHandler) {
      await clickHandler();
    }

    expect(mockSelectAdvancedSearchField).toHaveBeenCalledWith(ALARM_FIELDS.displayName, null);
  });

  it('Handles undefined selectAdvancedSearchField gracefully', async () => {
    const column = {
      value: ALARM_FIELDS.displayName,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.displayName),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.displayName }),
      isFilter: true,
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          [ALARM_FIELDS.displayName]: 'Test Name',
        },
        widget,
        column,
      },
      provide: {
        $selectAdvancedSearchField: undefined,
      },
    });

    const ellipsis = selectEllipsis(wrapper);

    // Should not throw error when selectAdvancedSearchField is undefined
    expect(() => ellipsis.trigger('click')).not.toThrow();
  });
});
