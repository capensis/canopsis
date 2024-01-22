import { generateRenderer } from '@unit/utils/vue';
import { randomDurationValue } from '@unit/utils/duration';

import { TIME_UNITS } from '@/constants';

import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';

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
  const snapshotFactory = generateRenderer(StorageSettingsJunitForm, { stubs });

  test('Junit delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationValue();

    selectJunitDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', { ...form, delete_after: newValue });
  });

  test('Renders `storage-settings-junit-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: dataStorageSettingsToForm().junit,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `storage-settings-junit-form` with custom form and history', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
        history: 1611210000,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
