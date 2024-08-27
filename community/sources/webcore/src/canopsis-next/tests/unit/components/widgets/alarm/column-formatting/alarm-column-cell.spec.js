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
  'alarm-column-value-state': true,
  'alarm-column-value-status': true,
  'color-indicator-wrapper': true,
  'alarm-column-value-categories': true,
  'alarm-column-value-extra-details': true,
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
  });
  const snapshotFactory = generateRenderer(AlarmColumnCell, {

    stubs,
    attachTo: document.body,
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    const menu = wrapper.findMenu();

    expect(document.body.innerHTML).toMatchSnapshot();
    expect(menu.element).toMatchSnapshot();
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

    popupBody.vm.$emit('close');

    await flushPromises();

    const menu = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menu.element).toMatchSnapshot();
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

    const menu = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menu.element).toMatchSnapshot();
  });
});
