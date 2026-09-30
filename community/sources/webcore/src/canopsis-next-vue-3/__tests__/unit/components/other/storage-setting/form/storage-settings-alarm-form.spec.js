import { randomDurationEnabledValue } from '@unit/utils/duration';
import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

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

  const validator = {
    validate: vi.fn(),
  };

  const factory = generateRenderer(StorageSettingsAlarmForm, {
    stubs,

    provide: {
      $validator: validator,
    },
  });

  test('Alarm archive after changed after trigger enabled duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationEnabledValue();

    selectAlarmArchiveAfterField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, archive_after: newValue });
  });

  test('Alarm delete after changed after trigger enabled duration field', async () => {
    const wrapper = factory({
      propsData: {
        form,
        history: {},
      },
    });

    const newValue = randomDurationEnabledValue();

    selectAlarmDeleteAfterField(wrapper).triggerCustomEvent('input', newValue);

    await flushPromises();

    expect(wrapper).toEmitInput({ ...form, delete_after: newValue });
    expect(validator.validate).toHaveBeenCalledWith('alarm.delete_after');
  });
});
