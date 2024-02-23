import { generateRenderer } from '@unit/utils/vue';
import { randomDurationValue } from '@unit/utils/duration';

import { TIME_UNITS } from '@/constants';

import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';

import StorageSettingsPbehaviorForm from '@/components/other/storage-setting/form/storage-settings-pbehavior-form.vue';
import CInformationBlock from '@/components/common/block/c-information-block.vue';

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
  const snapshotFactory = generateRenderer(StorageSettingsPbehaviorForm, { stubs });

  test('Pbehavior delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationValue();

    selectPbehaviorDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, delete_after: newValue });
  });

  test('Renders `storage-settings-pbehavior-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: dataStorageSettingsToForm().pbehavior,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `storage-settings-pbehavior-form` with custom form and history', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
        history: 1611230000,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
