import { randomDurationEnabledValue } from '@unit/utils/duration';
import { generateRenderer } from '@unit/utils/vue';

import { TIME_UNITS } from '@/constants';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import StorageSettingsJunitForm from '@/components/other/storage-setting/form/storage-settings-junit-form.vue';

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
  'storage-settings-history-message': true,
};

const selectJunitDeleteAfterField = wrapper => wrapper.find('c-enabled-duration-field-stub');

describe('storage-settings-junit-form', () => {
  const form = {
    junit: {
      delete_after: {
        value: 1,
        unit: TIME_UNITS.month,
        enabled: false,
      },
    },
  };

  const factory = generateRenderer(StorageSettingsJunitForm, { stubs });

  test('Junit delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationEnabledValue();

    selectJunitDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, delete_after: newValue });
  });
});
