import { generateRenderer } from '@unit/utils/vue';
import { TIME_UNITS } from '@/constants';
import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';
import { randomDurationValue } from '@unit/utils/duration';
import CInformationBlock from '@/components/common/block/c-information-block.vue';

import StorageSettingsPerfDataMetricsForm from '@/components/other/storage-setting/form/storage-settings-perf-data-metrics-form.vue';

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
};

const selectPerfDataMetricsDeleteAfterField = wrapper => wrapper.find('c-enabled-duration-field-stub');

describe('storage-settings-perf-data-metrics-form', () => {
  const form = {
    delete_after: {
      value: 6,
      unit: TIME_UNITS.month,
      enabled: false,
    },
  };

  const factory = generateRenderer(StorageSettingsPerfDataMetricsForm, { stubs });
  const snapshotFactory = generateRenderer(StorageSettingsPerfDataMetricsForm, { stubs });

  test('Perf data metrics delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationValue();

    selectPerfDataMetricsDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', { ...form, delete_after: newValue });
  });

  test('Renders `storage-settings-perf-data-metrics-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: dataStorageSettingsToForm().perf_data_metrics,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `storage-settings-perf-data-metrics-form` with custom form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
