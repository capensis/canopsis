import { randomDurationEnabledValue } from '@unit/utils/duration';
import { generateRenderer } from '@unit/utils/vue';

import { TIME_UNITS } from '@/constants';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import StorageSettingsMetricsForm from '@/components/other/storage-setting/form/storage-settings-metrics-form.vue';

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
};

const selectMetricsDeleteAfterField = wrapper => wrapper.find('c-enabled-duration-field-stub');

describe('storage-settings-metrics-form', () => {
  const form = {
    delete_after: {
      value: 6,
      unit: TIME_UNITS.month,
      enabled: true,
    },
  };

  const factory = generateRenderer(StorageSettingsMetricsForm, { stubs });

  test('Metrics delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationEnabledValue();

    selectMetricsDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, delete_after: newValue });
  });
});
