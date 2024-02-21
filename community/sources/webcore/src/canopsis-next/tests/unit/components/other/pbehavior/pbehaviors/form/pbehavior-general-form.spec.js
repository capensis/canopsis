import Faker from 'faker';
import flushPromises from 'flush-promises';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createCheckboxInputStub } from '@unit/stubs/input';
import { createMockedStoreModules, createPbehaviorTypesModule } from '@unit/utils/store';
import { PBEHAVIOR_TYPE_TYPES, TIME_UNITS } from '@/constants';

import PbehaviorGeneralForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-general-form.vue';

const stubs = {
  'c-name-field': true,
  'c-enabled-field': true,
  'c-duration-field': true,
  'date-time-splitted-range-picker-field': true,
  'pbehavior-comments-field': true,
  'recurrence-rule-form': true,
  'pbehavior-recurrence-rule-exceptions-field': true,
  'c-enabled-color-picker-field': true,
  'c-pbehavior-reason-field': true,
  'c-pbehavior-type-field': true,
  'c-color-picker-field': true,
  'c-collapse-panel': true,
  'v-checkbox': createCheckboxInputStub('v-checkbox'),
};

const snapshotStubs = {
  'c-name-field': true,
  'c-enabled-field': true,
  'c-duration-field': true,
  'date-time-splitted-range-picker-field': true,
  'pbehavior-comments-field': true,
  'recurrence-rule-form': true,
  'pbehavior-recurrence-rule-exceptions-field': true,
  'c-pbehavior-reason-field': true,
  'c-pbehavior-type-field': true,
  'c-enabled-color-picker-field': true,
  'c-collapse-panel': true,
};

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectDurationField = wrapper => wrapper.find('c-duration-field-stub');
const selectDateTimePickerField = wrapper => wrapper.find('date-time-splitted-range-picker-field-stub');
const selectEnabledFields = wrapper => wrapper.findAll('c-enabled-field-stub');
const selectReasonField = wrapper => wrapper.find('c-pbehavior-reason-field-stub');
const selectTypeField = wrapper => wrapper.find('c-pbehavior-type-field-stub');
const selectEnabledField = wrapper => selectEnabledFields(wrapper)
  .at(0);
const selectStartOnTriggerField = wrapper => selectEnabledFields(wrapper)
  .at(1);
const selectCheckboxFields = wrapper => wrapper.findAll('.v-checkbox');
const selectFullDayCheckbox = wrapper => selectCheckboxFields(wrapper)
  .at(0);
const selectNoEndingCheckbox = wrapper => selectCheckboxFields(wrapper)
  .at(1);
const selectPbehaviorCommentsField = wrapper => wrapper.find('pbehavior-comments-field-stub');
const selectEnabledColorPickerField = wrapper => wrapper.find('c-enabled-color-picker-field-stub');

describe('pbehavior-general-form', () => {
  const { pbehaviorTypesModule } = createPbehaviorTypesModule();
  const store = createMockedStoreModules([
    pbehaviorTypesModule,
  ]);

  const form = {
    name: Faker.datatype.string(),
    enabled: Faker.datatype.boolean(),
    start_on_trigger: Faker.datatype.boolean(),
    tstart: new Date(1610393200000),
    tstop: new Date(1757480799000),
    reason: {},
    type: {},
    color: Faker.internet.color(),
    comments: [],
  };

  const factory = generateShallowRenderer(PbehaviorGeneralForm, {
    stubs,
    store,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  const snapshotFactory = generateRenderer(PbehaviorGeneralForm, {
    stubs: snapshotStubs,
    store,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Name changed after trigger name field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newName = Faker.datatype.string();

    selectNameField(wrapper).vm.$emit('input', newName);

    expect(wrapper).toEmit('input', { ...form, name: newName });
  });

  test('Enabled changed after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newEnabled = !form.enabled;

    selectEnabledField(wrapper).vm.$emit('input', newEnabled);

    expect(wrapper).toEmit('input', { ...form, enabled: newEnabled });
  });

  test('Start on trigger enabled after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          start_on_trigger: false,
        },
        withStartOnTrigger: true,
      },
    });

    selectStartOnTriggerField(wrapper).vm.$emit('input', true);

    expect(wrapper).toEmit('input', {
      ...form,

      start_on_trigger: true,
      tstart: null,
      tstop: null,
    });
  });

  test('Start on trigger disabled after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          start_on_trigger: true,
        },
        withStartOnTrigger: true,
      },
    });

    selectStartOnTriggerField(wrapper).vm.$emit('input', false);

    expect(wrapper).toEmit('input', {
      ...form,

      start_on_trigger: false,
    });
  });

  test('Duration changed after trigger duration field', () => {
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          start_on_trigger: true,
        },
        withStartOnTrigger: true,
      },
    });

    const newDuration = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.week,
    };

    selectDurationField(wrapper).vm.$emit('input', newDuration);

    expect(wrapper).toEmit('input', {
      ...form,

      start_on_trigger: true,
      duration: newDuration,
    });
  });

  test('Start time changed after trigger date time picker field', () => {
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          start_on_trigger: false,
        },
      },
    });

    const newTimeStart = new Date(1233266565);

    selectDateTimePickerField(wrapper).vm.$emit('update:start', newTimeStart);

    expect(wrapper).toEmit('input', {
      ...form,

      start_on_trigger: false,
      tstart: newTimeStart,
    });
  });

  test('Stop time changed after trigger date time picker field', () => {
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          start_on_trigger: false,
        },
      },
    });

    const newTimeStart = new Date(1233266565);

    selectDateTimePickerField(wrapper).vm.$emit('update:end', newTimeStart);

    expect(wrapper).toEmit('input', {
      ...form,

      start_on_trigger: false,
      tstop: new Date(1233266565),
    });
  });

  test('Stop time didn\'t change after trigger date time picker field without value', () => {
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          start_on_trigger: false,
        },
      },
    });

    selectDateTimePickerField(wrapper).vm.$emit('update:end');

    expect(wrapper).toEmit('input', {
      ...form,

      start_on_trigger: false,
      tstop: undefined,
    });
  });

  test('Full day changed after trigger checkbox field', async () => {
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          start_on_trigger: false,
        },
      },
    });

    selectFullDayCheckbox(wrapper).vm.$emit('change', true);
    await flushPromises();

    expect(wrapper).toEmit('input', {
      ...form,

      start_on_trigger: false,
      tstart: new Date(1610319600000),
      tstop: new Date(1757541599000),
    });
  });

  test('Full day didn\'t change after trigger checkbox field without tstart', async () => {
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          tstart: null,
          start_on_trigger: false,
        },
      },
    });

    selectFullDayCheckbox(wrapper).vm.$emit('change', true);
    await flushPromises();

    expect(wrapper).not.toEmit('input');
  });

  test('No ending enabled after trigger checkbox field', async () => {
    const type = {
      type: PBEHAVIOR_TYPE_TYPES.pause,
    };
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          type,
          start_on_trigger: false,
        },
        noEnding: true,
      },
    });

    selectNoEndingCheckbox(wrapper).vm.$emit('change', true);
    await flushPromises();

    expect(wrapper).toEmit('input', {
      ...form,

      type,
      start_on_trigger: false,
      tstop: null,
    });
  });

  test('No ending disabled after trigger checkbox field', async () => {
    const type = {
      type: PBEHAVIOR_TYPE_TYPES.pause,
    };
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          tstop: null,
          type,
          start_on_trigger: false,
        },
        noEnding: true,
      },
    });

    selectNoEndingCheckbox(wrapper).vm.$emit('change', false);
    await flushPromises();

    expect(wrapper).toEmit('input', {
      ...form,

      type,
      start_on_trigger: false,
      tstop: new Date(1610396800000),
    });
  });

  test('No ending didn\'t changed after trigger checkbox field without tstart', async () => {
    const type = {
      type: PBEHAVIOR_TYPE_TYPES.pause,
    };
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          tstart: null,
          tstop: null,
          type,
          start_on_trigger: false,
        },
        noEnding: true,
      },
    });

    selectNoEndingCheckbox(wrapper).vm.$emit('change', false);
    await flushPromises();

    expect(wrapper).not.toEmit('input');
  });

  test('No ending changed after form type changed field', async () => {
    const type = {
      type: PBEHAVIOR_TYPE_TYPES.pause,
    };
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          tstop: null,
          type,
        },
      },
    });

    const newForm = {
      ...form,
      tstop: null,
      type: {
        type: PBEHAVIOR_TYPE_TYPES.active,
      },
    };

    await wrapper.setProps({
      form: newForm,
    });
    await flushPromises();

    expect(wrapper).toEmit('input', {
      ...newForm,
      tstop: new Date(1610396800000),
    });
  });

  test('No ending changed after form type changed field with full day', async () => {
    const type = {
      type: PBEHAVIOR_TYPE_TYPES.pause,
    };
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          tstart: new Date(1610319600000),
          tstop: null,
          type,
        },
      },
    });

    const newForm = {
      ...form,
      type: {
        type: PBEHAVIOR_TYPE_TYPES.active,
      },
    };

    await wrapper.setProps({
      form: newForm,
    });
    await flushPromises();

    expect(wrapper).toEmit('input', {
      ...newForm,
      tstop: new Date(1610492399999),
    });
  });

  test('No ending didn\'t change after form type changed field with paused type', async () => {
    const type = {
      type: PBEHAVIOR_TYPE_TYPES.active,
    };
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          tstop: null,
          type,
        },
      },
    });

    const newForm = {
      ...form,
      tstop: null,
      type: {
        type: PBEHAVIOR_TYPE_TYPES.pause,
      },
    };

    await wrapper.setProps({
      form: newForm,
    });
    await flushPromises();

    expect(wrapper).not.toEmit('input');
  });

  test('Reason changed after trigger reason field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newReason = {
      _id: Faker.datatype.string(),
    };

    selectReasonField(wrapper).vm.$emit('input', newReason);

    expect(wrapper).toEmit('input', {
      ...form,

      reason: newReason,
    });
  });

  test('Type changed after trigger type field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newType = {
      _id: Faker.datatype.string(),
    };

    selectTypeField(wrapper).vm.$emit('input', newType);

    expect(wrapper).toEmit('input', {
      ...form,

      type: newType,
    });
  });

  test('Comments updated after trigger pbehavior comments field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newComments = [
      {
        key: Faker.datatype.string(),
        message: Faker.datatype.string(),
      },
    ];

    selectPbehaviorCommentsField(wrapper).vm.$emit('input', newComments);

    expect(wrapper).toEmit('input', {
      ...form,
      comments: newComments,
    });
  });

  test('Color changed after trigger color field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newColor = Faker.internet.color();

    selectEnabledColorPickerField(wrapper).vm.$emit('input', newColor);

    expect(wrapper).toEmit('input', {
      ...form,

      color: newColor,
    });
  });

  test('Renders `pbehavior-general-form` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          name: 'pbehavior',
          enabled: false,
          tstart: new Date(1614861000000),
          tstop: new Date(1614861200000),
          reason: {},
          type: {},
          comments: [],
          color: '#123123',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pbehavior-general-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          name: 'pbehavior',
          enabled: false,
          start_on_trigger: true,
          duration: {
            unit: TIME_UNITS.hour,
            value: 2,
          },
          reason: {},
          type: {},
          comments: [],
          color: '#123123',
        },
        noEnabled: true,
        withStartOnTrigger: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
