import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';
import { fakeTimestamp } from '@unit/data/date';

import { MODALS } from '@/constants';

import LiveReporting from '@/components/sidebars/alarm/form/fields/live-reporting.vue';

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

describe('live-reporting', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(LiveReporting, { stubs });
  const snapshotFactory = generateRenderer(LiveReporting, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        listClick: jest.fn(),
      },
    },
  });

  it('Live reporting modal opened after trigger create button', () => {
    const liveReporting = { tstart: fakeTimestamp(), tstop: fakeTimestamp() };
    const wrapper = factory({
      propsData: {
        value: liveReporting,
      },
      mocks: {
        $modals,
      },
    });

    const createButton = selectSettingsCreateButton(wrapper);

    createButton.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.editLiveReporting,
        config: {
          action: expect.any(Function),
          ...liveReporting,
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = Faker.datatype.string();

    modalArguments.config.action(actionValue);

    expect(wrapper).toEmit('input', actionValue);
  });

  it('Text editor modal opened after trigger edit button', () => {
    const liveReporting = { tstart: fakeTimestamp(), tstop: fakeTimestamp() };
    const wrapper = factory({
      propsData: {
        value: liveReporting,
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
        name: MODALS.editLiveReporting,
        config: {
          ...liveReporting,
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = Faker.datatype.string();

    modalArguments.config.action(actionValue);

    expect(wrapper).toEmit('input', actionValue);
  });

  it('Value removed after trigger delete button', () => {
    const wrapper = factory({
      mocks: {
        $modals,
      },
    });

    const deleteButton = selectSettingsDeleteButton(wrapper);

    deleteButton.trigger('click');

    expect(wrapper).toEmit('input', {});
  });

  it('Renders `live-reporting` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `live-reporting` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          tstart: 1,
          tstop: 2,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
