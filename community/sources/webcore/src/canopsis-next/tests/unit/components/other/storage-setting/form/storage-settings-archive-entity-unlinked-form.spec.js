import { generateRenderer } from '@unit/utils/vue';
import { randomDurationValue } from '@unit/utils/duration';

import { TIME_UNITS } from '@/constants';

import { dataStorageEntityUnlinkedSettingsToForm } from '@/helpers/entities/data-storage/form';

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
  const snapshotFactory = generateRenderer(StorageSettingsArchiveEntityUnlinkedForm, { stubs });

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

  test('Renders `storage-settings-archive-entity-unlinked-form` with default duration', () => {
    const wrapper = snapshotFactory({
      propsData: {
        duration: dataStorageEntityUnlinkedSettingsToForm().archive_after,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `storage-settings-archive-entity-unlinked-form` with custom duration', () => {
    const wrapper = snapshotFactory({
      propsData: {
        duration,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
