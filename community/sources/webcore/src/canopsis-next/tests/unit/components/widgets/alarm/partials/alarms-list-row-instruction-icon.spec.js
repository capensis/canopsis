import { generateRenderer } from '@unit/utils/vue';
import { INSTRUCTION_EXECUTION_ICONS } from '@/constants';

import AlarmsListRowInstructionsIcon from '@/components/widgets/alarm/partials/alarms-list-row-instructions-icon.vue';

describe('alarms-list-row-instructions-icon', () => {
  const snapshotFactory = generateRenderer(AlarmsListRowInstructionsIcon, {
    attachTo: document.body,
  });

  it('Renders `alarms-list-row-instructions-icon` with instruction', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `alarms-list-row-instructions-icon` with all auto instructions', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoSuccessful,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `alarms-list-row-instructions-icon` with auto instruction failed', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoFailed,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `alarms-list-row-instructions-icon` with auto instruction running', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoFailed,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `alarms-list-row-instructions-icon` with manual instruction waiting result', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualFailed,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `alarms-list-row-instructions-icon` with manual instruction running', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualInProgress,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});