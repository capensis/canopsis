import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { INSTRUCTION_EXECUTION_ICONS } from '@/constants';

import AlarmsListRowInstructionsIcon from '@/components/widgets/alarm/partials/alarms-list-row-instructions-icon.vue';

const stubs = {
  'c-simple-tooltip': true,
};

const snapshotStubs = {
  'c-simple-tooltip': true,
};

describe('alarms-list-row-instructions-icon', () => {
  const factory = generateShallowRenderer(AlarmsListRowInstructionsIcon, {
    stubs,
  });

  const snapshotFactory = generateRenderer(AlarmsListRowInstructionsIcon, {
    stubs: snapshotStubs,
  });

  describe('Icon computation', () => {
    it('Shows assignment_warning icon when withoutAnyExecution is true', () => {
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.withoutAnyExecution,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [Faker.random.word()],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconName).toBe('$vuetify.icons.assignment_warning');
      expect(wrapper.vm.iconClass).toContain('instruction-icon--warning');
    });

    it('Shows manual_instruction icon when manual instruction is available', () => {
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualAvailable,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [Faker.random.word()],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconName).toBe('$vuetify.icons.manual_instruction');
    });

    it('Shows assignment icon when auto instruction is available', () => {
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoSuccessful,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [Faker.random.word()],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconName).toBe('assignment');
    });

    it('Defaults to manualAvailable when instruction_execution_icon is not provided', () => {
      const alarm = {
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [Faker.random.word()],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconName).toBe('$vuetify.icons.manual_instruction');
    });
  });

  describe('Icon classes', () => {
    it('Applies warning class for withoutAnyExecution', () => {
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.withoutAnyExecution,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [Faker.random.word()],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconClass).toContain('instruction-icon--warning');
    });

    it('Applies blinking and dotted classes for running instructions', () => {
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualInProgress,
        running_manual_instructions: [Faker.random.word()],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconClass).toContain('blinking');
      expect(wrapper.vm.iconClass).toContain('instruction-icon--dotted');
    });

    it('Applies failed class for failed instructions', () => {
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualFailed,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [Faker.random.word()],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconClass).toContain('instruction-icon--failed');
    });

    it('Applies completed class for successful instructions', () => {
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualSuccessful,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [Faker.random.word()],
        successful_auto_instructions: [],
        assigned_instructions: [],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconClass).toContain('instruction-icon--completed');
    });
  });

  describe('Tooltip content', () => {
    it('Shows running manual instruction tooltip', () => {
      const instructionName = Faker.random.word();
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualInProgress,
        running_manual_instructions: [instructionName],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconTooltip).toContain('Manual instruction');
      expect(wrapper.vm.iconTooltip).toContain('in progress');
      expect(wrapper.vm.iconTooltip).toContain(instructionName);
    });

    it('Shows running auto instruction tooltip', () => {
      const instructionName = Faker.random.word();
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoInProgress,
        running_manual_instructions: [],
        running_auto_instructions: [instructionName],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconTooltip).toContain('Automatic instruction');
      expect(wrapper.vm.iconTooltip).toContain('in progress');
      expect(wrapper.vm.iconTooltip).toContain(instructionName);
    });

    it('Shows failed manual instruction tooltip', () => {
      const instructionName = Faker.random.word();
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualFailed,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [instructionName],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconTooltip).toContain('Manual instruction');
      expect(wrapper.vm.iconTooltip).toContain('is failed');
      expect(wrapper.vm.iconTooltip).toContain(instructionName);
    });

    it('Shows failed auto instruction tooltip', () => {
      const instructionName = Faker.random.word();
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoFailed,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [instructionName],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconTooltip).toContain('Automatic instruction');
      expect(wrapper.vm.iconTooltip).toContain('is failed');
      expect(wrapper.vm.iconTooltip).toContain(instructionName);
    });

    it('Shows successful manual instruction tooltip', () => {
      const instructionName = Faker.random.word();
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualSuccessful,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [instructionName],
        successful_auto_instructions: [],
        assigned_instructions: [],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconTooltip).toContain('Manual instruction');
      expect(wrapper.vm.iconTooltip).toContain('is successful');
      expect(wrapper.vm.iconTooltip).toContain(instructionName);
    });

    it('Shows successful auto instruction tooltip', () => {
      const instructionName = Faker.random.word();
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoSuccessful,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [instructionName],
        assigned_instructions: [],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconTooltip).toContain('Automatic instruction');
      expect(wrapper.vm.iconTooltip).toContain('is successful');
      expect(wrapper.vm.iconTooltip).toContain(instructionName);
    });

    it('Shows withoutAnyExecution tooltip', () => {
      const instructionName = Faker.random.word();
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.withoutAnyExecution,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [instructionName],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconTooltip).toContain('Manual instruction wasn\'t executed');
    });

    it('Shows hasManualInstruction tooltip when not withoutAnyExecution', () => {
      const instructionName = Faker.random.word();
      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualAvailable,
        running_manual_instructions: [],
        running_auto_instructions: [],
        failed_manual_instructions: [],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [instructionName],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      expect(wrapper.vm.iconTooltip).toContain('There is a manual instruction');
    });

    it('Combines multiple instruction tooltips', () => {
      const runningInstruction = Faker.random.word();
      const failedInstruction = Faker.random.word();
      const assignedInstruction = Faker.random.word();

      const alarm = {
        instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualFailedWithInProgress,
        running_manual_instructions: [runningInstruction],
        running_auto_instructions: [],
        failed_manual_instructions: [failedInstruction],
        failed_auto_instructions: [],
        successful_manual_instructions: [],
        successful_auto_instructions: [],
        assigned_instructions: [assignedInstruction],
      };

      const wrapper = factory({
        propsData: { alarm },
      });

      const tooltip = wrapper.vm.iconTooltip;
      expect(tooltip).toContain(runningInstruction);
      expect(tooltip).toContain(failedInstruction);
      expect(tooltip).toContain('in progress');
      expect(tooltip).toContain('is failed');
      expect(tooltip).toContain('There is a manual instruction');
    });
  });

  it('Renders `alarms-list-row-instructions-icon` with default props', () => {
    const alarm = {
      running_manual_instructions: [],
      running_auto_instructions: [],
      failed_manual_instructions: [],
      failed_auto_instructions: [],
      successful_manual_instructions: [],
      successful_auto_instructions: [],
      assigned_instructions: [],
    };

    const wrapper = snapshotFactory({
      propsData: { alarm },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-row-instructions-icon` with withoutAnyExecution state', () => {
    const alarm = {
      instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.withoutAnyExecution,
      running_manual_instructions: [],
      running_auto_instructions: [],
      failed_manual_instructions: [],
      failed_auto_instructions: [],
      successful_manual_instructions: [],
      successful_auto_instructions: [],
      assigned_instructions: ['instruction'],
    };

    const wrapper = snapshotFactory({
      propsData: { alarm },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-row-instructions-icon` with running manual instruction', () => {
    const alarm = {
      instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualInProgress,
      running_manual_instructions: ['instruction'],
      running_auto_instructions: [],
      failed_manual_instructions: [],
      failed_auto_instructions: [],
      successful_manual_instructions: [],
      successful_auto_instructions: [],
      assigned_instructions: [],
    };

    const wrapper = snapshotFactory({
      propsData: { alarm },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-row-instructions-icon` with failed instruction', () => {
    const alarm = {
      instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.autoFailed,
      running_manual_instructions: [],
      running_auto_instructions: [],
      failed_manual_instructions: [],
      failed_auto_instructions: ['instruction'],
      successful_manual_instructions: [],
      successful_auto_instructions: [],
      assigned_instructions: [],
    };

    const wrapper = snapshotFactory({
      propsData: { alarm },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-row-instructions-icon` with successful instruction', () => {
    const alarm = {
      instruction_execution_icon: INSTRUCTION_EXECUTION_ICONS.manualSuccessful,
      running_manual_instructions: [],
      running_auto_instructions: [],
      failed_manual_instructions: [],
      failed_auto_instructions: [],
      successful_manual_instructions: ['instruction'],
      successful_auto_instructions: [],
      assigned_instructions: [],
    };

    const wrapper = snapshotFactory({
      propsData: { alarm },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
