import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import ClickOutside from '@/services/click-outside';
import { MODALS } from '@/constants';

import InfoPopupSetting from '@/components/modals/alarm/info-popup-setting/info-popup-setting.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'date-interval-selector': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'date-interval-selector': true,
};

const selectModalActionsButtons = wrapper => wrapper.findAll('.actions button.v-btn');
const selectSubmitButton = wrapper => selectModalActionsButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectModalActionsButtons(wrapper).at(0);
const selectAddPopupButton = wrapper => wrapper.find('.text button.v-btn');
const selectPopupCards = wrapper => wrapper.findAll('.text  v-card-stub');
const selectPopupCardByIndex = (wrapper, index) => selectPopupCards(wrapper).at(index);
const selectRemovePopupButtonByIndex = (wrapper, index) => selectPopupCardByIndex(wrapper, index)
  .findAll('button.v-btn').at(0);
const selectEditPopupButtonByIndex = (wrapper, index) => selectPopupCardByIndex(wrapper, index)
  .findAll('button.v-btn').at(1);

describe('info-popup-setting', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(InfoPopupSetting, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });
  const snapshotFactory = generateRenderer(InfoPopupSetting, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  beforeAll(() => jest.useFakeTimers());

  test('Form submitted after trigger submit button', async () => {
    const action = jest.fn();
    const infoPopups = [
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
    ];
    const config = { infoPopups, action };

    const wrapper = factory({
      propsData: {
        modal: {
          config,
        },
      },
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledWith(infoPopups);
    expect($modals.hide).toBeCalledWith();
  });

  test('Form submitted after trigger submit button without action', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {},
        },
      },
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    expect($modals.hide).toBeCalledWith();
  });

  test('Popup added after trigger add button', async () => {
    const action = jest.fn();
    const infoPopups = [
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
    ];
    const columns = [
      {
        value: Faker.datatype.string(),
      },
    ];
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            columns,
            infoPopups,
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const addPopupButton = selectAddPopupButton(wrapper);

    const newPopup = {
      column: Faker.datatype.string(),
      template: Faker.datatype.string(),
    };

    addPopupButton.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.addInfoPopup,
        config: {
          columns,
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action(newPopup);

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledWith([
      ...infoPopups,
      newPopup,
    ]);
    expect($modals.hide).toBeCalled();
  });

  test('Popup edited after trigger edit button', async () => {
    const action = jest.fn();
    const infoPopups = [
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
    ];
    const columns = [
      {
        value: Faker.datatype.string(),
      },
    ];
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            columns,
            infoPopups,
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const editPopupButton = selectEditPopupButtonByIndex(wrapper, 1);

    const newPopupData = {
      column: Faker.datatype.string(),
      template: Faker.datatype.string(),
    };

    editPopupButton.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.addInfoPopup,
        config: {
          columns,
          popup: infoPopups[1],
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action(newPopupData);

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledWith([
      infoPopups[0],
      newPopupData,
    ]);
    expect($modals.hide).toBeCalled();
  });

  test('Popup removed after trigger delete button', async () => {
    const action = jest.fn();
    const infoPopups = [
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
    ];
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            infoPopups,
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const removePopupButton = selectRemovePopupButtonByIndex(wrapper, 1);

    removePopupButton.trigger('click');

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledWith([infoPopups[0]]);
    expect($modals.hide).toBeCalled();
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {},
        },
      },
      mocks: {
        $modals,
      },
    });

    const cancelButton = selectCancelButton(wrapper);

    cancelButton.trigger('click');

    await flushPromises(true);

    expect($modals.hide).toBeCalled();
  });

  test('Renders `info-popup-setting` with empty modal', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {},
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `info-popup-setting` with interval', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            infoPopups: [
              {
                column: 'alarm.v.connector',
                template: 'Connector template',
              },
              {
                column: 'extra_details',
                template: 'Extra details template',
              },
            ],
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
