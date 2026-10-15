import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { MODALS } from '@/constants';

import TextEditor from '@/components/sidebars/form/fields/text-editor.vue';

const stubs = {
  'settings-button-field': {
    template: `
      <div class="settings-button-field">
        <button class="edit" @click="$emit('edit', $event)" />
      </div>
    `,
  },
};

const selectEditButton = wrapper => wrapper.find('.settings-button-field .edit');

describe('field-template', () => {
  mockModals();

  const factory = generateShallowRenderer(TextEditor, { stubs });

  test('Value changed after trigger color indicator field', () => {
    const title = Faker.datatype.string();
    const value = Faker.datatype.string();
    const variables = [
      {
        value: Faker.datatype.string(),
      },
    ];

    const wrapper = factory({
      propsData: {
        value,
        title,
        variables,
      },
      mocks: {
        $modals: globalThis.__TEST_MODALS__,
      },
    });

    selectEditButton(wrapper).trigger('click');

    expect(globalThis.__TEST_MODALS__.show).toHaveBeenCalledWith({
      name: MODALS.textEditor,
      dialogProps: undefined,
      config: {
        text: value,
        variables,
        sanitizeOptions: undefined,
        action: expect.any(Function),
      },
    });

    const [modalArguments] = globalThis.__TEST_MODALS__.show.mock.calls[0];

    const newTemplate = Faker.datatype.string();

    modalArguments.config.action(newTemplate);

    expect(wrapper).toEmitInput(newTemplate);
  });
});
