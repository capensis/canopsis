import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';
import { MODALS } from '@/constants';

import TextEditor from '@/components/sidebars/form/fields/text-editor.vue';

const stubs = {
  'settings-button-field': {
    template: `
      <div class="settings-button-field">
        <button class="create" @click="$listeners.create" />
        <button class="edit" @click="$listeners.edit" />
        <button class="delete" @click="$listeners.delete" />
      </div>
    `,
  },
};

const snapshotStubs = {
  'settings-button-field': true,
};

const selectSettingsCreateButton = wrapper => wrapper.find('.settings-button-field .create');
const selectSettingsEditButton = wrapper => wrapper.find('.settings-button-field .edit');
const selectSettingsDeleteButton = wrapper => wrapper.find('.settings-button-field .delete');

describe('text-editor', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(TextEditor, { stubs });
  const snapshotFactory = generateRenderer(TextEditor, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        listClick: jest.fn(),
      },
    },
  });

  it('Text editor modal opened after trigger create button', () => {
    const wrapper = factory({
      mocks: {
        $modals,
      },
    });

    const createButton = selectSettingsCreateButton(wrapper);

    createButton.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.textEditor,
        config: {
          action: expect.any(Function),
          text: '',
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = Faker.datatype.string();

    modalArguments.config.action(actionValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual(actionValue);
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
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.textEditor,
        config: {
          action: expect.any(Function),
          text: value,
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = Faker.datatype.string();

    modalArguments.config.action(actionValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual(actionValue);
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
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action(Faker.datatype.string());

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual('');
  });

  it('Renders `text-editor` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `text-editor` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Custom title',
        value: 'Custom value',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
