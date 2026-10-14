import { fakeTimestamp } from '@unit/data/date';
import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { MODALS } from '@/constants';

import LiveReporting from '@/components/sidebars/form/fields/live-reporting.vue';

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

describe('live-reporting', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(LiveReporting, { stubs });

  it('Live reporting modal opened after trigger create button', () => {
    const liveReporting = { tstart: fakeTimestamp(), tstop: fakeTimestamp() };
    const wrapper = factory({
      propsData: {
        value: liveReporting,
      },
    });

    const createButton = selectSettingsCreateButton(wrapper);

    createButton.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.editLiveReporting,
      config: {
        action: expect.any(Function),
        ...liveReporting,
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = Faker.datatype.string();

    modalArguments.config.action(actionValue);

    expect(wrapper).toEmitInput(actionValue);
  });

  it('Text editor modal opened after trigger edit button', () => {
    const liveReporting = { tstart: fakeTimestamp(), tstop: fakeTimestamp() };
    const wrapper = factory({
      propsData: {
        value: liveReporting,
      },
    });

    const editButton = selectSettingsEditButton(wrapper);

    editButton.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.editLiveReporting,
      config: {
        ...liveReporting,
        action: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = Faker.datatype.string();

    modalArguments.config.action(actionValue);

    expect(wrapper).toEmitInput(actionValue);
  });

  it('Value removed after trigger delete button', () => {
    const wrapper = factory();

    const deleteButton = selectSettingsDeleteButton(wrapper);

    deleteButton.trigger('click');

    expect(wrapper).toEmitInput({});
  });
});
