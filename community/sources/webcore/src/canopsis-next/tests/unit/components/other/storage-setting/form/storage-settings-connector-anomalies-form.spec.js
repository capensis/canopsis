import { generateRenderer } from '@unit/utils/vue';
import { randomDurationEnabledValue } from '@unit/utils/duration';

import { TIME_UNITS } from '@/constants';

import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import StorageSettingsConnectorAnomaliesForm from '@/components/other/storage-setting/form/storage-settings-connector-anomalies-form.vue';

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
};

const selectEventAnomalyDeleteAfterField = wrapper => wrapper.find('c-enabled-duration-field-stub');

describe('storage-settings-connector-anomalies-form', () => {
  const form = {
    delete_after: {
      value: 1,
      unit: TIME_UNITS.month,
      enabled: true,
    },
  };

  const factory = generateRenderer(StorageSettingsConnectorAnomaliesForm, { stubs });
  const snapshotFactory = generateRenderer(StorageSettingsConnectorAnomaliesForm, { stubs });

  test('Event anomaly delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationEnabledValue();

    selectEventAnomalyDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, delete_after: newValue });
  });

  test('Renders `storage-settings-connector-anomalies-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: dataStorageSettingsToForm().connector_anomalies,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `storage-settings-connector-anomalies-form` with custom form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
