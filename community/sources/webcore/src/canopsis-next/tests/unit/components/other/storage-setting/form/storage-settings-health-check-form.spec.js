import { generateRenderer } from '@unit/utils/vue';
import { TIME_UNITS } from '@/constants';
import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';
import { randomDurationValue } from '@unit/utils/duration';
import CInformationBlock from '@/components/common/block/c-information-block.vue';

import StorageSettingsHealthCheckForm from '@/components/other/storage-setting/form/storage-settings-health-check-form.vue';

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
  'storage-settings-history-message': true,
};

const selectHealthCheckDeleteAfterField = wrapper => wrapper.find('c-enabled-duration-field-stub');

describe('storage-settings-health-check-form', () => {
  const form = {
    delete_after: {
      value: 1,
      unit: TIME_UNITS.month,
      enabled: false,
    },
  };

  const factory = generateRenderer(StorageSettingsHealthCheckForm, { stubs });
  const snapshotFactory = generateRenderer(StorageSettingsHealthCheckForm, { stubs });

  test('Health check delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationValue();

    selectHealthCheckDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', { ...form, delete_after: newValue });
  });

  test('Renders `storage-settings-health-check-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: dataStorageSettingsToForm().health_check,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `storage-settings-health-check-form` with custom form and history', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
        history: 1611240000,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
