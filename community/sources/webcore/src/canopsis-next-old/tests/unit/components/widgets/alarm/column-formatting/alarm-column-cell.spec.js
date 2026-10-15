import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { ALARM_FIELDS, DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';
import {
  getAlarmsListWidgetColumnComponentGetter,
  getAlarmsListWidgetColumnValueFilter,
} from '@/helpers/entities/widget/forms/alarm';

import AlarmColumnCell from '@/components/widgets/alarm/columns-formatting/alarm-column-cell.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';

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

const selectOpenButton = wrapper => wrapper.find('.v-btn');
const selectEllipsis = wrapper => wrapper.find('c-ellipsis-stub');
const selectAlarmColumnPopupBody = wrapper => wrapper.find('alarm-column-cell-popup-body-stub');

describe('alarm-column-cell', () => {
  const timestamp = 1641768553245;
  const duration = 164176;
  const widget = {
    parameters: {},
  };

  const factory = generateShallowRenderer(AlarmColumnCell, {

    stubs,
    attachTo: document.body,
    provide: {
      $selectAdvancedSearchField: jest.fn(),
    },
  });
  const snapshotFactory = generateRenderer(AlarmColumnCell, {

    stubs,
    attachTo: document.body,
    provide: {
      $selectAdvancedSearchField: jest.fn(),
    },
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

    expect(ellipsis.attributes('text')).toBe('09/01/2022 23:49:13');
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

    expect(ellipsis.attributes('text')).toBe('02/01/1970');
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
    const mockSelectAdvancedSearchField = jest.fn();
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

    expect(mockSelectAdvancedSearchField).toHaveBeenCalledWith(ALARM_FIELDS.displayName, 'Test Name');
  });

  it('Calls selectAdvancedSearchField with sanitized HTML value when column is HTML and filter', async () => {
    const mockSelectAdvancedSearchField = jest.fn();
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

    expect(mockSelectAdvancedSearchField).toHaveBeenCalledWith(ALARM_FIELDS.entityName, '&lt;script&gt;alert("xss")Safe Name');
  });

  it('Calls selectAdvancedSearchField with filtered value when column has custom filter', async () => {
    const mockSelectAdvancedSearchField = jest.fn();
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

    expect(mockSelectAdvancedSearchField).toHaveBeenCalledWith('custom_field', 'filtered_original_value');
  });

  it('Does not call selectAdvancedSearchField when column is not a filter', async () => {
    const mockSelectAdvancedSearchField = jest.fn();
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
    const mockSelectAdvancedSearchField = jest.fn();
    const mockCustomClickHandler = jest.fn();
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

    expect(mockSelectAdvancedSearchField).toHaveBeenCalledWith(ALARM_FIELDS.displayName, 'Test Name');
    expect(mockCustomClickHandler).toHaveBeenCalled();
  });

  it('Handles empty or null values when filter column is clicked', async () => {
    const mockSelectAdvancedSearchField = jest.fn();
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

  it('Renders `alarm-column-cell` with column state', async () => {
    const column = {
      value: ALARM_FIELDS.state,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.state),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.state }),
    };
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {},
        widget,
        column,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-cell` with column status', async () => {
    const column = {
      value: ALARM_FIELDS.status,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.status),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.status }),
    };

    const wrapper = snapshotFactory({
      propsData: {
        alarm: {},
        widget,
        column,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-cell` with column impact state', async () => {
    const column = {
      value: ALARM_FIELDS.impactState,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.impactState),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.impactState }),
    };

    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {},
        },
        widget,
        column,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-cell` with column links', async () => {
    const column = {
      value: ALARM_FIELDS.links,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.links),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.links }),
    };

    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          links: {},
        },
        widget,
        column,
      },
      listeners: {
        activate: jest.fn(),
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-cell` with column links only icon', async () => {
    const column = {
      value: ALARM_FIELDS.links,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.links),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.links, onlyIcon: true }),
    };

    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          links: {},
        },
        widget,
        column,
      },
      listeners: {
        activate: jest.fn(),
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-cell` with column links with inline links count', async () => {
    const originalColumn = {
      value: ALARM_FIELDS.links,
      inlineLinksCount: 2,
    };

    const column = {
      ...originalColumn,

      filter: getAlarmsListWidgetColumnValueFilter(originalColumn.value),
      getComponent: getAlarmsListWidgetColumnComponentGetter(originalColumn),
    };

    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          links: {},
        },
        widget,
        column,
      },
      listeners: {
        activate: jest.fn(),
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-cell` with column links category', async () => {
    const columnValue = 'links.test';
    const column = {
      value: 'links.test',
      filter: getAlarmsListWidgetColumnValueFilter(columnValue),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: columnValue }),
    };
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          links: {
            test: [],
          },
        },
        widget,
        column,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-cell` with column extra details', async () => {
    const column = {
      value: ALARM_FIELDS.extraDetails,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.extraDetails),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.extraDetails }),
    };

    const wrapper = snapshotFactory({
      propsData: {
        alarm: {},
        widget,
        column,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-cell` with invalid html', async () => {
    const columnValue = 'entity.test';
    const column = {
      value: columnValue,
      filter: getAlarmsListWidgetColumnValueFilter(columnValue),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: columnValue }),
      isHtml: true,
    };

    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {
            test: '<div Name',
          },
        },
        widget,
        column,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-cell` with popups', async () => {
    const column = {
      value: ALARM_FIELDS.displayName,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.displayName),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.displayName }),
    };

    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          [ALARM_FIELDS.displayName]: 'Name',
        },
        widget,
        column: {
          ...column,

          popupTemplate: 'template',
        },
      },
    });

    await flushPromises();

    const openButton = selectOpenButton(wrapper);

    openButton.trigger('click');

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `alarm-column-cell` with popups after hide', async () => {
    const column = {
      value: ALARM_FIELDS.displayName,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.displayName),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.displayName }),
    };

    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          [ALARM_FIELDS.displayName]: 'Name',
        },
        widget,
        column: {
          ...column,

          popupTemplate: 'template',
        },
      },
    });

    await flushPromises();

    const openButton = selectOpenButton(wrapper);

    openButton.trigger('click');

    await flushPromises();

    const popupBody = selectAlarmColumnPopupBody(wrapper);

    popupBody.triggerCustomEvent('close');

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `alarm-column-cell` with popups and html', async () => {
    const column = {
      value: ALARM_FIELDS.entityName,
      filter: getAlarmsListWidgetColumnValueFilter(ALARM_FIELDS.entityName),
      getComponent: getAlarmsListWidgetColumnComponentGetter({ value: ALARM_FIELDS.entityName }),
      isHtml: true,
    };

    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {
            name: '<div class="custom-html-value" data-test="123">Name</div>',
          },
        },
        widget,
        column: {
          ...column,

          popupTemplate: 'template',
        },
      },
    });

    await flushPromises();

    const openButton = selectOpenButton(wrapper);

    openButton.trigger('click');

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
