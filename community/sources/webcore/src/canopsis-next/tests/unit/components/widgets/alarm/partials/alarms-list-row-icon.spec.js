import { mount, createVueInstance } from '@unit/utils/vue';

import AlarmsListRowIcon from '@/components/widgets/alarm/partials/alarms-list-row-icon.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(AlarmsListRowIcon, {
  localVue,

  ...options,
});

describe('alarms-list-row-icon', () => {
  it('Renders `alarms-list-row-icon` with instruction', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `alarms-list-row-icon` with all auto instructions', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          is_all_auto_instructions_completed: true,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `alarms-list-row-icon` with auto instruction running', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          is_auto_instruction_running: true,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `alarms-list-row-icon` with manual instruction waiting result', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          is_manual_instruction_waiting_result: true,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `alarms-list-row-icon` with manual instruction running', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          is_manual_instruction_running: true,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});
