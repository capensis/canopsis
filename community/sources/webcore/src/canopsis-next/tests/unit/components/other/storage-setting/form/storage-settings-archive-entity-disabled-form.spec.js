import { generateRenderer } from '@unit/utils/vue';

import StorageSettingsArchiveEntityDisabledForm from '@/components/other/storage-setting/form/storage-settings-archive-entity-disabled-form.vue';

const stubs = {
  'c-help-icon': true,
};

const selectWithDependenciesCheckbox = wrapper => wrapper.find('input[type="checkbox"]');

describe('storage-settings-archive-entity-disabled-form', () => {
  const form = {
    with_dependencies: true,
  };

  const factory = generateRenderer(StorageSettingsArchiveEntityDisabledForm, { stubs });
  const snapshotFactory = generateRenderer(StorageSettingsArchiveEntityDisabledForm, { stubs });

  test('Form changed after trigger with dependencies checkbox', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = false;

    const checkbox = selectWithDependenciesCheckbox(wrapper);
    checkbox.setChecked(newValue);

    expect(wrapper).toEmitInput({ ...form, with_dependencies: newValue });
  });

  test('Renders `storage-settings-archive-entity-disabled-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `storage-settings-archive-entity-disabled-form` with custom form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
