import { generateRenderer } from '@unit/utils/vue';
import { randomDurationValue } from '@unit/utils/duration';

import { TIME_UNITS } from '@/constants';

import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';

import StorageSettingsForm from '@/components/other/storage-setting/form/storage-settings-form.vue';

const stubs = {
  'storage-settings-entity-clean-form': true,
  'storage-settings-entity-unlinked-form': true,
  'storage-settings-entity-disabled-form': true,
  'storage-settings-event-filter-failure-form': true,
  'storage-settings-perf-data-metrics-form': true,
  'storage-settings-metrics-form': true,
  'storage-settings-webhook-form': true,
  'storage-settings-health-check-form': true,
  'storage-settings-junit-form': true,
  'storage-settings-pbehavior-form': true,
  'storage-settings-remediation-form': true,
  'storage-settings-alarm-form': true,
};

const selectStorageSettingsPerfDataMetricsForm = wrapper => wrapper.find('storage-settings-perf-data-metrics-form-stub');
const selectStorageSettingsMetricsForm = wrapper => wrapper.find('storage-settings-metrics-form-stub');
const selectStorageSettingsWebhookForm = wrapper => wrapper.find('storage-settings-webhook-form-stub');
const selectStorageSettingsHealthCheckForm = wrapper => wrapper.find('storage-settings-health-check-form-stub');
const selectStorageSettingsJunitForm = wrapper => wrapper.find('storage-settings-junit-form-stub');
const selectStorageSettingsPbehaviorForm = wrapper => wrapper.find('storage-settings-pbehavior-form-stub');
const selectStorageSettingsRemediationForm = wrapper => wrapper.find('storage-settings-remediation-form-stub');
const selectStorageSettingsAlarmForm = wrapper => wrapper.find('storage-settings-alarm-form-stub');

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

  const factory = generateRenderer(StorageSettingsForm, { stubs });
  const snapshotFactory = generateRenderer(StorageSettingsForm, { stubs });

  test('Alarm storage settings changed after trigger alarm settings', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = {
      archive_after: randomDurationValue(),
      delete_after: randomDurationValue(),
    };

    selectStorageSettingsAlarmForm(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      alarm: newValue,
    });
  });

  test('Remediation storage settings changed after trigger remediation settings', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = {
      delete_after: randomDurationValue(),
      delete_stats_after: randomDurationValue(),
      delete_mod_stats_after: randomDurationValue(),
    };

    selectStorageSettingsRemediationForm(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      remediation: newValue,
    });
  });

  test('Pbehavior storage settings changed after trigger pbehavior settings', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = {
      delete_after: randomDurationValue(),
    };

    selectStorageSettingsPbehaviorForm(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      pbehavior: newValue,
    });
  });

  test('Junit storage settings changed after trigger junit settings', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = {
      delete_after: randomDurationValue(),
    };

    selectStorageSettingsJunitForm(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      junit: newValue,
    });
  });

  test('Health check storage settings changed after trigger health check settings', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = {
      delete_after: randomDurationValue(),
    };

    selectStorageSettingsHealthCheckForm(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      health_check: newValue,
    });
  });

  test('Webhook storage settings changed after trigger webhook settings', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = {
      delete_after: randomDurationValue(),
    };

    selectStorageSettingsWebhookForm(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      webhook: newValue,
    });
  });

  test('Metrics storage settings changed after trigger metrics settings', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = {
      delete_after: randomDurationValue(),
    };

    selectStorageSettingsMetricsForm(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      metrics: newValue,
    });
  });

  test('Perf data metrics storage settings changed after trigger perf data metrics settings', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = {
      delete_after: randomDurationValue(),
    };

    selectStorageSettingsPerfDataMetricsForm(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', {
      ...form,
      perf_data_metrics: newValue,
    });
  });

  test('Renders `storage-settings-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: dataStorageSettingsToForm(),
        history: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
