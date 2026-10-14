import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import FieldTextEditorWithTemplate from '@/components/sidebars/form/fields/text-editor-with-template.vue';

describe('text-editor-with-template field', () => {
  const $modals = mockModals();
  const factory = generateShallowRenderer(FieldTextEditorWithTemplate, {
    mocks: { $modals },
    stubs: {
      'settings-button-field': true,
    },
  });

  test('emits the edited text together with the template update', () => {
    const wrapper = factory();
    const text = Faker.datatype.string();
    const template = Faker.datatype.uuid();

    wrapper.vm.updateText(text, template);

    expect(wrapper.emitted('update:modelValue')).toEqual([[text]]);
    expect(wrapper.emitted('update:template')).toEqual([[template, text]]);
  });
});
