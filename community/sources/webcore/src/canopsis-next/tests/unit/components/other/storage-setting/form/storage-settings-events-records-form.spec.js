import { generateRenderer } from '@unit/utils/vue';
import { randomDurationValue } from '@unit/utils/duration';

import { TIME_UNITS } from '@/constants';

import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import StorageSettingsEventsRecordsForm from '@/components/other/storage-setting/form/storage-settings-events-records-form.vue';

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'storage-settings-duration-field': true,
  'storage-settings-history-message': true,
};

const selectDeleteAfterField = wrapper => wrapper.find('storage-settings-duration-field-stub');

describe('storage-settings-junit-form', () => {
  const form = {
    delete_after: {
      value: 2,
      unit: TIME_UNITS.month,
    },
  };

  const factory = generateRenderer(StorageSettingsEventsRecordsForm, { stubs });
  const snapshotFactory = generateRenderer(StorageSettingsEventsRecordsForm, { stubs });

  test('Junit delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = randomDurationValue();

    selectDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, delete_after: newValue });
  });

  test('Renders `storage-settings-events-records-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: dataStorageSettingsToForm().event_records,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `storage-settings-events-records-form` with custom form and history', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
        history: 1611210000,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
