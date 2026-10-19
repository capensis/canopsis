import { randomDurationValue } from '@unit/utils/duration';
import { generateRenderer } from '@unit/utils/vue';

import { TIME_UNITS } from '@/constants';

import StorageSettingsArchiveEntityUnlinkedForm from '@/components/other/storage-setting/form/storage-settings-archive-entity-unlinked-form.vue';

const stubs = {
  'c-duration-field': true,
};

const selectDurationField = wrapper => wrapper.find('c-duration-field-stub');

describe('storage-settings-archive-entity-unlinked-form', () => {
  const duration = {
    value: 60,
    unit: TIME_UNITS.day,
  };

  const factory = generateRenderer(StorageSettingsArchiveEntityUnlinkedForm, { stubs });

  test('Duration changed after trigger duration field', () => {
    const wrapper = factory({
      propsData: {
        duration,
      },
    });

    const newValue = randomDurationValue();

    selectDurationField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
