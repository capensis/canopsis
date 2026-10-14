import { randomDurationEnabledValue } from '@unit/utils/duration';
import { generateRenderer } from '@unit/utils/vue';

import { TIME_UNITS } from '@/constants';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import StorageSettingsAlarmForm from '@/components/other/storage-setting/form/storage-settings-alarm-form.vue';

const stubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': true,
  'c-enabled-duration-field': true,
  'storage-settings-history-message': true,
};

const selectEnabledDurationFieldByIndex = (wrapper, index) =>
  wrapper.findAll('c-enabled-duration-field-stub').at(index);
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

  const factory = generateRenderer(StorageSettingsAlarmForm, {
    stubs,
  });

  test('Alarm archive after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationEnabledValue();

    selectAlarmArchiveAfterField(wrapper).triggerCustomEvent('update:duration', newValue);

    expect(wrapper).toEmit('update:form', { ...form, archive_after: newValue });
  });

  test('Alarm delete after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationEnabledValue();

    selectAlarmDeleteAfterField(wrapper).triggerCustomEvent('update:duration', newValue);

    expect(wrapper).toEmit('update:form', { ...form, delete_after: newValue });
  });
});
