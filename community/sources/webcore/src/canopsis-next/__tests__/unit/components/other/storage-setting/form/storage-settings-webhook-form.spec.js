import { createCheckboxInputStub } from '@unit/stubs/input';
import { randomDurationEnabledValue } from '@unit/utils/duration';
import { generateRenderer } from '@unit/utils/vue';

import { TIME_UNITS } from '@/constants';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import StorageSettingsWebhookForm from '@/components/other/storage-setting/form/storage-settings-webhook-form.vue';

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
  'storage-settings-history-message': true,
  'c-enabled-field': true,
  'v-checkbox': createCheckboxInputStub('v-checkbox'),
};

const selectWebhookDeleteAfterField = wrapper => wrapper.find('c-enabled-duration-field-stub');

describe('storage-settings-webhook-form', () => {
  const form = {
    log_credentials: true,
    delete_after: {
      value: 60,
      unit: TIME_UNITS.day,
      enabled: true,
    },
  };

  const factory = generateRenderer(StorageSettingsWebhookForm, { stubs });

  test('Webhook delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationEnabledValue();

    selectWebhookDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, delete_after: newValue });
  });
});
