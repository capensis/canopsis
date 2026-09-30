import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { MODALS } from '@/constants';

import TextEditor from '@/components/sidebars/form/fields/text-editor.vue';

const stubs = {
  'settings-button-field': {
    template: `
      <div class="settings-button-field">
        <button class="create" @click="$emit('create', $event)" />
        <button class="edit" @click="$emit('edit', $event)" />
        <button class="delete" @click="$emit('delete', $event)" />
      </div>
    `,
  },
};

const selectSettingsCreateButton = wrapper => wrapper.find('.settings-button-field .create');
const selectSettingsEditButton = wrapper => wrapper.find('.settings-button-field .edit');
const selectSettingsDeleteButton = wrapper => wrapper.find('.settings-button-field .delete');

describe('text-editor', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(TextEditor, { stubs });

  it('Text editor modal opened after trigger create button', () => {
    const wrapper = factory({
      mocks: {
        $modals,
      },
    });

    const createButton = selectSettingsCreateButton(wrapper);

    createButton.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.textEditor,
      config: {
        action: expect.any(Function),
        text: '',
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = Faker.datatype.string();

    modalArguments.config.action(actionValue);

    expect(wrapper).toEmitInput(actionValue);
  });

  it('Text editor modal opened after trigger edit button', () => {
    const value = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        value,
      },
      mocks: {
        $modals,
      },
    });

    const editButton = selectSettingsEditButton(wrapper);

    editButton.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.textEditor,
      config: {
        action: expect.any(Function),
        text: value,
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = Faker.datatype.string();

    modalArguments.config.action(actionValue);

    expect(wrapper).toEmitInput(actionValue);
  });

  it('Confirmation modal opened after trigger delete button', () => {
    const wrapper = factory({
      mocks: {
        $modals,
      },
    });

    const deleteButton = selectSettingsDeleteButton(wrapper);

    deleteButton.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action(Faker.datatype.string());

    expect(wrapper).toEmitInput('');
  });
});
