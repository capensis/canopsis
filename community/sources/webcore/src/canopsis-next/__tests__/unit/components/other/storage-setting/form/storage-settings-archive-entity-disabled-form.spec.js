import { generateRenderer } from '@unit/utils/vue';

import StorageSettingsArchiveEntityDisabledForm from '@/components/other/storage-setting/form/storage-settings-archive-entity-disabled-form.vue';

const stubs = {
  'c-checkbox': true,
  'c-help-icon': true,
};

const selectWithDependenciesCheckbox = wrapper => wrapper.find('c-checkbox-stub');

describe('storage-settings-archive-entity-disabled-form', () => {
  const form = {
    with_dependencies: true,
  };

  const factory = generateRenderer(StorageSettingsArchiveEntityDisabledForm, { stubs });

  test('Form changed after trigger with dependencies checkbox', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = false;

    selectWithDependenciesCheckbox(wrapper).triggerCustomEvent('update:modelValue', newValue);

    expect(wrapper).toEmit('update:modelValue', { ...form, with_dependencies: newValue });
  });
});
