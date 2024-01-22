import { generateRenderer } from '@unit/utils/vue';
import { randomDurationValue } from '@unit/utils/duration';

import { TIME_UNITS } from '@/constants';

import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import StorageSettingsRemediationForm from '@/components/other/storage-setting/form/storage-settings-remediation-form.vue';

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
  'storage-settings-history-message': true,
};

const selectEnabledDurationFieldByIndex = (wrapper, index) => wrapper.findAll('c-enabled-duration-field-stub').at(index);
const selectRemediationDeleteAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 0);
const selectRemediationDeleteStatsAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 1);
const selectRemediationDeleteModStatsAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 2);

describe('storage-settings-remediation-form', () => {
  const form = {
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
  };

  const factory = generateRenderer(StorageSettingsRemediationForm, { stubs });
  const snapshotFactory = generateRenderer(StorageSettingsRemediationForm, { stubs });

  test('Remediation delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationValue();

    selectRemediationDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', { ...form, delete_after: newValue });
  });

  test('Remediation delete stats after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationValue();

    selectRemediationDeleteStatsAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', { ...form, delete_stats_after: newValue });
  });

  test('Remediation delete mod stats after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationValue();

    selectRemediationDeleteModStatsAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', { ...form, delete_mod_stats_after: newValue });
  });

  test('Renders `storage-settings-remediation-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: dataStorageSettingsToForm().remediation,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `storage-settings-remediation-form` with custom form and history', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
        history: 1611220000,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
