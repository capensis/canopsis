import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import AlarmsListRow from '@/components/widgets/alarm/partials/alarms-list-row.vue';

const AlarmsExpandPanelBtnStub = {
  props: ['expanded'],
  emits: ['update:expanded'],
  methods: { showExpandPanel() {} },
  template: '<button class="expand-panel-btn" @click="$emit(\'update:expanded\', !expanded)" />',
};

const stubs = {
  'alarms-list-row-instructions-icon': true,
  'alarms-list-row-bookmark-icon': true,
  'alarms-expand-panel-btn': AlarmsExpandPanelBtnStub,
  'alarm-column-value': true,
  'actions-panel': true,
  'c-booted-placeholder-loader': true,
  'c-light-checkbox': {
    props: ['modelValue'],
    emits: ['update:modelValue'],
    template:
      '<button class="light-checkbox" @click="$emit(\'update:modelValue\', !modelValue)" />',
  },
};

const selectExpandButton = wrapper => wrapper.find('button.expand-panel-btn');
const selectTableRow = wrapper => wrapper.find('tr');
const selectCheckbox = wrapper => wrapper.find('button.light-checkbox');
const selectAlarmColumnValue = wrapper => wrapper.find('alarm-column-value-stub');

describe('alarms-list-row', () => {
  const factory = generateShallowRenderer(AlarmsListRow, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
        $intersectionObserver: {
          observe: vi.fn(),
          unobserve: vi.fn(),
        },
      },
    },
  });

  it('Alarm selected after trigger checkbox', async () => {
    const wrapper = factory({
      propsData: {
        alarm: {
          v: {
            status: {},
          },
        },
        booted: true,
        expanded: false,
        widget: {},
        headers: [{ key: 'icons' }],
        rowIconsColumnKey: 'icons',
        selectable: true,
        selectableAlarm: true,
      },
    });

    await selectCheckbox(wrapper).trigger('click');

    expect(wrapper).toEmit('update:selected', true);
  });

  it('Alarm selection is toggled while mouse-selecting', async () => {
    const wrapper = factory({
      propsData: {
        alarm: {
          v: {
            status: {},
          },
        },
        booted: true,
        expanded: false,
        widget: {},
        headers: [{ key: 'icons' }],
        rowIconsColumnKey: 'icons',
        selecting: true,
        selected: false,
      },
    });

    await selectTableRow(wrapper).trigger('mousedown', { button: 0 });

    expect(wrapper).toEmit('update:selected', true);
  });

  it('Row expanded after trigger expand button with hidden groups and without filtered alarms', async () => {
    const alarm = {
      _id: 'alarm-id',
      v: {
        status: {},
      },
    };

    const newExpanded = false;
    const wrapper = factory({
      propsData: {
        alarm,
        booted: true,
        expanded: false,
        widget: {},
        headers: [{ key: 'icons' }, { key: 'value' }],
        rowIconsColumnKey: 'icons',
        expandable: true,
        hideGroups: true,
      },
    });

    const expandButton = selectExpandButton(wrapper);

    expandButton.triggerCustomEvent('update:expanded', newExpanded);

    await flushPromises();

    expect(wrapper).toEmit('update:expanded', newExpanded);
  });

  it('Row closed after trigger expand button with expanded: true', async () => {
    const alarm = {
      _id: 'alarm-id',
      v: {
        status: {},
      },
    };

    const newExpanded = false;
    const wrapper = factory({
      propsData: {
        alarm,
        booted: true,
        expanded: true,
        widget: {},
        headers: [{ key: 'icons' }, { key: 'value' }],
        rowIconsColumnKey: 'icons',
        expandable: true,
      },
    });

    const expandButton = selectExpandButton(wrapper);

    expandButton.triggerCustomEvent('update:expanded', newExpanded);

    await flushPromises();

    expect(wrapper).toEmit('update:expanded', newExpanded);
  });

  it('Click state emitted after trigger click state event', async () => {
    const entity = {
      _id: 'alarm-entity',
    };
    const alarm = {
      _id: 'alarm-id',
      v: {
        status: {},
      },
      entity,
    };

    const wrapper = factory({
      propsData: {
        alarm,
        booted: true,
        expanded: true,
        widget: {},
        headers: [
          { key: 'first', value: 'first' },
          { key: 'second', value: 'second' },
        ],
        rowIconsColumnKey: 'icons',
        expandable: true,
      },
    });

    await flushPromises();

    selectAlarmColumnValue(wrapper).triggerCustomEvent('click:state');

    expect(wrapper).toHaveBeenEmit('click:state');
  });

  it('Expand panel opened after double click on row when openExpandPanelByDoubleClick enabled', () => {
    const showExpandPanel = vi.fn();

    AlarmsListRow.methods.openExpandPanelByDoubleClick.call({
      expandable: true,
      expanded: false,
      widget: { parameters: { openExpandPanelByDoubleClick: true } },
      $refs: { expandPanelBtn: { showExpandPanel } },
    });

    expect(showExpandPanel).toHaveBeenCalledWith(true);
  });
});
