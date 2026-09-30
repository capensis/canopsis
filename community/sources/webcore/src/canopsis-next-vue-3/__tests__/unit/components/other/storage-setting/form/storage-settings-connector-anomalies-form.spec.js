import { randomDurationEnabledValue } from '@unit/utils/duration';
import { generateRenderer } from '@unit/utils/vue';

import { TIME_UNITS } from '@/constants';

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
});
