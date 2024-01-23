import { generateRenderer } from '@unit/utils/vue';
import { createCheckboxInputStub } from '@unit/stubs/input';
import { randomDurationValue } from '@unit/utils/duration';

import { TIME_UNITS } from '@/constants';

import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import StorageSettingsWebhookForm from '@/components/other/storage-setting/form/storage-settings-webhook-form.vue';

const snapshotStubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
  'storage-settings-history-message': true,
  'c-enabled-field': true,
};
const stubs = {
  ...snapshotStubs,
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
  const snapshotFactory = generateRenderer(StorageSettingsWebhookForm, {
    stubs: snapshotStubs,
  });

  test('Webhook delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationValue();

    selectWebhookDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, delete_after: newValue });
  });

  test('Renders `storage-settings-webhook-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: dataStorageSettingsToForm().webhook,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `storage-settings-webhook-form` with custom form and history', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
        history: 1611250000,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
