import flushPromises from 'flush-promises';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { DATETIME_FORMATS } from '@/constants';

import AlarmColumnCell from '@/components/widgets/alarm/columns-formatting/alarm-column-cell.vue';

const localVue = createVueInstance();

const stubs = {
  'alarm-column-cell-popup-body': true,
  'alarm-column-value-state': true,
  'alarm-column-value-status': true,
  'color-indicator-wrapper': true,
  'alarm-column-value-categories': true,
  'alarm-column-value-extra-details': true,
  'alarm-column-value-links': true,
  'c-ellipsis': true,
};

const factory = (options = {}) => shallowMount(AlarmColumnCell, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(AlarmColumnCell, {
  localVue,
  stubs,

  ...options,
});

const selectOpenButton = wrapper => wrapper.find('.v-btn');
const selectEllipsis = wrapper => wrapper.find('c-ellipsis-stub');
const selectAlarmColumnPopupBody = wrapper => wrapper.find('alarm-column-cell-popup-body-stub');

describe('alarm-column-cell', () => {
  const timestamp = 1641768553245;
  const duration = 164176;
  const widget = {
    parameters: {},
  };

  it.each([
    'v.last_update_date',
    'v.creation_date',
    'v.last_event_date',
    'v.activation_date',
    'v.state.t',
    'v.status.t',
    'v.resolved',
    't',
  ])('Default filter for date field: "%s" converted value to time', async (field) => {
    const column = {
      value: field,
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
    'v.duration',
    'v.current_state_duration',
    'v.active_duration',
    'v.snooze_duration',
    'v.pbh_inactive_duration',
  ])('Default filter for duration field: "%s" converted value to duration', async (field) => {
    const column = {
      value: field,
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
    const column = {
      value: 'custom_field',
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          custom_field: duration,
        },
        columnsFilters: [{
          column: column.value,
          filter: 'date',
          attributes: [DATETIME_FORMATS.short],
        }],
        widget,
        column,
      },
    });

    const ellipsis = selectEllipsis(wrapper);

    expect(ellipsis.attributes('text')).toBe('02/01/1970');
  });

  it('Default filter for creation date field converted value to time', async () => {
    const column = {
      value: 'v.creation_date',
    };

    const wrapper = factory({
      propsData: {
        alarm: {
          v: {
            creation_date: timestamp,
          },
        },
        widget,
        column,
      },
    });

    const ellipsis = selectEllipsis(wrapper);

    expect(ellipsis.attributes('text')).toBe('09/01/2022 23:49:13');
  });

  it('Renders `alarm-column-cell` with column state', async () => {
    const column = {
      value: 'v.state.val',
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
      value: 'v.status.val',
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

  it('Renders `alarm-column-cell` with column priority', async () => {
    const column = {
      value: 'priority',
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

  it('Renders `alarm-column-cell` with column impact state', async () => {
    const column = {
      value: 'impact_state',
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
      value: 'links',
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

  it('Renders `alarm-column-cell` with column links as list', async () => {
    const column = {
      value: 'links',
    };
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          links: {},
        },
        widget: {
          parameters: {
            linksCategoriesAsList: {
              enabled: true,
              limit: 2,
            },
          },
        },
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
    const column = {
      value: 'links.test',
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
      value: 'extra_details',
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
    const column = {
      value: 'entity.test',
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
      value: 'name',
    };
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          name: 'Name',
        },
        widget,
        column: {
          ...column,

          popupTemplate: 'template',
        },
      },
    });

    const openButton = selectOpenButton(wrapper);

    openButton.trigger('click');

    await flushPromises();

    const menu = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menu.element).toMatchSnapshot();
  });

  it('Renders `alarm-column-cell` with popups after hide', async () => {
    const column = {
      value: 'name',
    };
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          name: 'Name',
        },
        widget,
        column: {
          ...column,

          popupTemplate: 'template',
        },
      },
    });

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
      value: 'entity.name',
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
