import { generateRenderer } from '@unit/utils/vue';

import { INSTRUCTION_EXECUTION_ICONS } from '@/constants';

import AlarmsListRowInstructionsIcon from '@/components/widgets/alarm/partials/alarms-list-row-instructions-icon.vue';

const snapshotStubs = {
  'c-simple-tooltip': true,
};

describe('alarms-list-row-instructions-icon', () => {
  const instructions = ['Instruction 1', 'Instruction 2'];

  const snapshotFactory = generateRenderer(AlarmsListRowInstructionsIcon, {
    attachTo: document.body,
    stubs: snapshotStubs,
  });

  test('Renders `alarms-list-row-instructions-icon` without instruction', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-row-instructions-icon` with all auto instructions successful', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoSuccessful,
          successful_auto_instructions: instructions,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-row-instructions-icon` with auto instruction failed', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoFailed,
          failed_auto_instructions: instructions,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-row-instructions-icon` with auto instruction running', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoFailed,
          running_auto_instructions: instructions,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-row-instructions-icon` with manual instruction failed', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualFailed,
          failed_manual_instructions: instructions,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-row-instructions-icon` with manual instruction running', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualInProgress,
          running_manual_instructions: instructions,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
