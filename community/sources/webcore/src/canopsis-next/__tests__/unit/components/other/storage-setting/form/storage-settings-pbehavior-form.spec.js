import { randomDurationEnabledValue } from '@unit/utils/duration';
import { generateRenderer } from '@unit/utils/vue';

import { TIME_UNITS } from '@/constants';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import StorageSettingsPbehaviorForm from '@/components/other/storage-setting/form/storage-settings-pbehavior-form.vue';

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
  'storage-settings-history-message': true,
};

const selectPbehaviorDeleteAfterField = wrapper => wrapper.find('c-enabled-duration-field-stub');

describe('storage-settings-pbehavior-form', () => {
  const form = {
    delete_after: {
      value: 1,
      unit: TIME_UNITS.month,
      enabled: false,
    },
  };

  const factory = generateRenderer(StorageSettingsPbehaviorForm, { stubs });

  test('Pbehavior delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationEnabledValue();

    selectPbehaviorDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, delete_after: newValue });
  });
});
