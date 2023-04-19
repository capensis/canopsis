import Faker from 'faker';

import { createVueInstance, generateRenderer } from '@unit/utils/vue';
import { createCheckboxInputStub, createInputStub } from '@unit/stubs/input';
import { randomArrayItem } from '@unit/utils/array';
import { TIME_UNITS } from '@/constants';
import { dataStorageSettingsToForm } from '@/helpers/forms/data-storage';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import StorageSettingsForm from '@/components/other/storage-setting/form/storage-settings-form';

const localVue = createVueInstance();

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
  'c-enabled-field': true,
  'v-checkbox': createCheckboxInputStub('v-checkbox'),
  'v-radio-group': createInputStub('v-radio-group'),
};
const snapshotStubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
  'c-enabled-field': true,
};

const selectEnabledDurationFieldByIndex = (wrapper, index) => wrapper.findAll('c-enabled-duration-field-stub').at(index);
const selectAlarmArchiveAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 0);
const selectAlarmDeleteAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 1);
const selectRemediationDeleteAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 2);
const selectRemediationDeleteStatsAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 3);
const selectRemediationDeleteModStatsAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 4);
const selectPbehaviorDeleteAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 5);
const selectJunitDeleteAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 6);
const selectHealthCheckDeleteAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 7);
const selectWebhookDeleteAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 8);
const selectMetricsDeleteAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 9);
const selectPerfDataMetricsDeleteAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 10);

describe('storage-settings-form', () => {
  const form = dataStorageSettingsToForm({
    junit: {
      delete_after: {
        value: 1,
        unit: TIME_UNITS.month,
        enabled: false,
      },
    },
    remediation: {
      delete_after: {
        value: 1,
        unit: TIME_UNITS.month,
        enabled: false,
      },
      delete_stats_after: {
        value: 2,
        unit: TIME_UNITS.day,
        enabled: false,
      },
      delete_mod_stats_after: {
        value: 2,
        unit: TIME_UNITS.day,
        enabled: true,
      },
    },
    alarm: {
      archive_after: {
        value: 1,
        unit: TIME_UNITS.month,
        enabled: false,
      },
      delete_after: {
        value: 1,
        unit: TIME_UNITS.month,
        enabled: false,
      },
    },
    pbehavior: {
      delete_after: {
        value: 1,
        unit: TIME_UNITS.month,
        enabled: false,
      },
    },
    health_check: {
      delete_after: {
        value: 1,
        unit: TIME_UNITS.month,
        enabled: false,
      },
    },
    webhook: {
      log_credentials: true,
      delete_after: {
        value: 60,
        unit: TIME_UNITS.day,
        enabled: true,
      },
    },
    metrics: {
      delete_after: {
        value: 6,
        unit: TIME_UNITS.month,
        enabled: true,
      },
    },
    perf_data_metrics: {
      delete_after: {
        value: 6,
        unit: TIME_UNITS.month,
        enabled: false,
      },
    },
  });

  const randomDurationValue = () => ({
    unit: randomArrayItem(Object.values(TIME_UNITS)),
    value: Faker.datatype.number(),
    enabled: Faker.datatype.boolean(),
  });

  const factory = generateRenderer(StorageSettingsForm, {
    localVue,
    stubs,
  });
  const snapshotFactory = generateRenderer(StorageSettingsForm, {
    localVue,
    stubs: snapshotStubs,
  });

  test('Alarm archive after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectAlarmArchiveAfterField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      alarm: { ...form.alarm, archive_after: newValue },
    });
  });

  test('Alarm delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectAlarmDeleteAfterField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      alarm: { ...form.alarm, delete_after: newValue },
    });
  });

  test('Remediation delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectRemediationDeleteAfterField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      remediation: { ...form.remediation, delete_after: newValue },
    });
  });

  test('Remediation delete stats after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectRemediationDeleteStatsAfterField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      remediation: { ...form.remediation, delete_stats_after: newValue },
    });
  });

  test('Remediation delete mod stats after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectRemediationDeleteModStatsAfterField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      remediation: { ...form.remediation, delete_mod_stats_after: newValue },
    });
  });

  test('Pbehavior delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectPbehaviorDeleteAfterField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      pbehavior: { ...form.pbehavior, delete_after: newValue },
    });
  });

  test('Junit delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectJunitDeleteAfterField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      junit: { ...form.junit, delete_after: newValue },
    });
  });

  test('Health check delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectHealthCheckDeleteAfterField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      health_check: { ...form.health_check, delete_after: newValue },
    });
  });

  test('Webhook delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectWebhookDeleteAfterField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      webhook: { ...form.webhook, delete_after: newValue },
    });
  });

  test('Metrics delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectMetricsDeleteAfterField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      metrics: { ...form.metrics, delete_after: newValue },
    });
  });

  test('Perf data metrics delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectPerfDataMetricsDeleteAfterField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      perf_data_metrics: { ...form.perf_data_metrics, delete_after: newValue },
    });
  });

  test('Renders `storage-settings-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: dataStorageSettingsToForm(),
        history: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `storage-settings-form` with custom form and history', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
        history: {
          junit: 1611210000,
          remediation: 1611220000,
          pbehavior: 1611230000,
          health_check: 1611240000,
          webhook: 1611250000,
          alarm: {
            time: 1611260000,
            deleted: 1611270000,
            archived: 1611280000,
          },
          entity: {
            time: 1611290000,
            deleted: 1611300000,
            archived: 1611310000,
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
