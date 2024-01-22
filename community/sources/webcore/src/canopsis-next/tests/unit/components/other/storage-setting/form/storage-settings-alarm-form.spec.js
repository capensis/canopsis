import { generateRenderer } from '@unit/utils/vue';
import { TIME_UNITS } from '@/constants';
import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';
import { randomDurationValue } from '@unit/utils/duration';

import StorageSettingsAlarmForm from '@/components/other/storage-setting/form/storage-settings-alarm-form.vue';
import CInformationBlock from '@/components/common/block/c-information-block.vue';

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
  'storage-settings-history-message': true,
};

const selectEnabledDurationFieldByIndex = (wrapper, index) => wrapper.findAll('c-enabled-duration-field-stub').at(index);
const selectAlarmArchiveAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 0);
const selectAlarmDeleteAfterField = wrapper => selectEnabledDurationFieldByIndex(wrapper, 1);

describe('storage-settings-alarm-form', () => {
  const form = {
    archive_after: {
      value: 1,
      unit: TIME_UNITS.month,
      enabled: false,
    },
    delete_after: {
      value: 1,
      unit: TIME_UNITS.month,
      enabled: false,
    },
  };

  const factory = generateRenderer(StorageSettingsAlarmForm, { stubs });
  const snapshotFactory = generateRenderer(StorageSettingsAlarmForm, { stubs });

  test('Alarm archive after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectAlarmArchiveAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', { ...form, archive_after: newValue });
  });

  test('Alarm delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationValue();

    selectAlarmDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', { ...form, delete_after: newValue });
  });

  test('Renders `storage-settings-alarm-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: dataStorageSettingsToForm().alarm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `storage-settings-alarm-form` with custom form and history', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
        history: {
          time: 1611260000,
          deleted: 1611270000,
          archived: 1611280000,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
