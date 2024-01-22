import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createMockedStoreModules } from '@unit/utils/store';
import { createCheckboxInputStub, createInputStub, createSelectInputStub } from '@unit/stubs/input';
import { MAX_LIMIT, REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import RemediationInstructionsFilterForm from '@/components/other/remediation/instructions-filter/form/remediation-instructions-filter-form.vue';

const stubs = {
  'v-radio-group': createInputStub('v-radio-group'),
  'v-select': createSelectInputStub('v-select'),
  'v-switch': createCheckboxInputStub('v-switch'),
  'v-checkbox': createCheckboxInputStub('v-checkbox'),
  'c-help-icon': true,
};

const snapshotStubs = {
  'c-help-icon': true,
};

const selectRadioGroups = wrapper => wrapper.findAll('.v-radio-group');
const selectRunningField = wrapper => selectRadioGroups(wrapper).at(1);
const selectWithField = wrapper => selectRadioGroups(wrapper).at(0);
const selectAllField = wrapper => wrapper.find('.v-switch');
const selectCheckboxFields = wrapper => wrapper.findAll('.v-checkbox');
const selectAutoField = wrapper => selectCheckboxFields(wrapper).at(0);
const selectManualField = wrapper => selectCheckboxFields(wrapper).at(1);
const selectInstructionsField = wrapper => wrapper.find('.v-select');

describe('remediation-instructions-filter-form', () => {
  const fetchList = jest.fn();
  const getInstructions = jest.fn().mockReturnValue([]);
  const remediationInstructionModule = {
    name: 'remediationInstruction',
    getters: {
      items: getInstructions,
      pending: false,
    },
    actions: {
      fetchList,
    },
  };

  const store = createMockedStoreModules([
    remediationInstructionModule,
  ]);

  const factory = generateShallowRenderer(RemediationInstructionsFilterForm, { stubs });
  const snapshotFactory = generateRenderer(RemediationInstructionsFilterForm, { stubs: snapshotStubs });

  test('Instructions fetched after mount', () => {
    factory({
      store,
      propsData: {
        form: {},
      },
    });

    expect(fetchList).toBeCalledWith(
      expect.any(Object),
      {
        params: { limit: MAX_LIMIT },
      },
      undefined,
    );
  });

  test('Has running enabled after trigger has running field', () => {
    const wrapper = factory({
      store,
      propsData: {
        form: {
          running: undefined,
        },
      },
    });

    const runningField = selectRunningField(wrapper);
    runningField.triggerCustomEvent('input', false);

    expect(wrapper).toEmit('input', { running: false });
  });

  test('With enabled after trigger with field', () => {
    const wrapper = factory({
      store,
      propsData: {
        form: {
          with: true,
        },
      },
    });

    const withField = selectWithField(wrapper);

    withField.triggerCustomEvent('input', false);

    expect(wrapper).toEmit('input', { with: false });
  });

  test('All disabled after trigger all field', () => {
    const wrapper = factory({
      store,
      propsData: {
        form: {
          all: true,
        },
      },
    });

    const allField = selectAllField(wrapper);

    allField.triggerCustomEvent('change', false);

    expect(wrapper).toEmit('input', { all: false });
  });

  test('All enabled after trigger all field', () => {
    const wrapper = factory({
      store,
      propsData: {
        form: {
          all: false,
        },
      },
    });

    const allField = selectAllField(wrapper);

    allField.triggerCustomEvent('change', true);

    expect(wrapper).toEmit('input', {
      all: true,
      manual: true,
      auto: true,
      instructions: [],
    });
  });

  test('Auto updated after trigger auto field', () => {
    const filter = {
      auto: true,
      instructions: [
        {
          type: REMEDIATION_INSTRUCTION_TYPES.auto,
        },
      ],
    };
    const wrapper = factory({
      store,
      propsData: {
        form: filter,
      },
    });

    const autoField = selectAutoField(wrapper);

    autoField.triggerCustomEvent('change', false);
    expect(wrapper).toEmit('input', {
      auto: false,
      instructions: filter.instructions,
    });
  });

  test('Manual updated after trigger auto field', () => {
    const filter = {
      manual: false,
      instructions: [
        {
          type: REMEDIATION_INSTRUCTION_TYPES.manual,
        },
      ],
    };
    const wrapper = factory({
      store,
      propsData: {
        form: filter,
      },
    });

    const manualField = selectManualField(wrapper);

    manualField.triggerCustomEvent('change', true);
    expect(wrapper).toEmit('input', {
      manual: true,
      instructions: [],
    });
  });

  test('Instructions updated after trigger auto field', () => {
    const wrapper = factory({
      store,
      propsData: {
        form: {
          manual: true,
        },
      },
    });

    const instructionsField = selectInstructionsField(wrapper);
    const instructions = [
      {
        _id: 'instruction',
        type: REMEDIATION_INSTRUCTION_TYPES.auto,
      },
    ];

    instructionsField.triggerCustomEvent('change', instructions);
    expect(wrapper).toEmit('input', {
      manual: true,
      instructions,
    });
  });

  test('Instructions updated after trigger auto field', () => {
    const autoInstruction = {
      _id: 'auto-instruction',
      type: REMEDIATION_INSTRUCTION_TYPES.auto,
    };
    const manualInstruction = {
      _id: 'manual-instruction',
      type: REMEDIATION_INSTRUCTION_TYPES.manual,
    };
    const customInstruction = {
      _id: 'custom-instruction',
      type: REMEDIATION_INSTRUCTION_TYPES.auto,
    };
    const availableInstructions = [
      autoInstruction,
      manualInstruction,
      customInstruction,
    ];
    const filters = [
      {
        manual: true,
      },
      {
        with: true,
        all: false,
      }, {
        with: true,
        manual: true,
      }, {
        with: true,
        manual: false,
        instructions: [autoInstruction],
      },
    ];
    const wrapper = factory({
      store: createMockedStoreModules([
        {
          ...remediationInstructionModule,
          getters: {
            ...remediationInstructionModule.getters,
            items: availableInstructions,
          },
        },
      ]),
      propsData: {
        form: {
          manual: true,
        },
        filters,
      },
    });

    const instructionsField = selectInstructionsField(wrapper);

    expect(instructionsField.vm.items).toEqual([
      {
        ...autoInstruction,
        disabled: true,
      },
      customInstruction,
    ]);
  });

  test('Renders `remediation-instructions-filter-form` with default props', () => {
    const wrapper = snapshotFactory({
      store,
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `remediation-instructions-filter-form` with custom props', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        form: {
          all: true,
          with: true,
          auto: true,
          manual: true,
          instructions: [],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
