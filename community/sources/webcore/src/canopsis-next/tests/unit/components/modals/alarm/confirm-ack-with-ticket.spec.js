import flushPromises from 'flush-promises';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createModalWrapperStub } from '@unit/stubs/modal';
import ClickOutside from '@/services/click-outside';

import ConfirmAckWithTicket from '@/components/modals/alarm/confirm-ack-with-ticket.vue';

const localVue = createVueInstance();

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
};

const factory = (options = {}) => shallowMount(ConfirmAckWithTicket, {
  localVue,
  stubs,
  attachTo: document.body,

  parentComponent: {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(ConfirmAckWithTicket, {
  localVue,
  stubs: snapshotStubs,

  parentComponent: {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  },

  ...options,
});

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectContinueWithTicketButton = wrapper => selectButtons(wrapper).at(2);
const selectContinueButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);

describe('confirm-ack-with-ticket', () => {
  const $modals = mockModals();

  test('Ack confirmed after trigger continue button', async () => {
    const continueAction = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: { continueAction },
        },
      },
      mocks: {
        $modals,
      },
    });

    const continueButton = selectContinueButton(wrapper);

    continueButton.trigger('click');

    await flushPromises();

    expect(continueAction).toBeCalledWith();
    expect($modals.hide).toBeCalledWith();
  });

  test('Ack confirmed after trigger continue with ticket button', async () => {
    const continueWithTicketAction = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: { continueWithTicketAction },
        },
      },
      mocks: {
        $modals,
      },
    });

    const continueWithTicketButton = selectContinueWithTicketButton(wrapper);

    continueWithTicketButton.trigger('click');

    await flushPromises();

    expect(continueWithTicketAction).toBeCalledWith();
    expect($modals.hide).toBeCalledWith();
  });

  test('Modal hidden after trigger continue button without action', async () => {
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

    const continueButton = selectContinueButton(wrapper);

    continueButton.trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalledWith();
  });

  test('Modal hidden after trigger continue with ticket button without action', async () => {
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

    const continueWithTicketButton = selectContinueWithTicketButton(wrapper);

    continueWithTicketButton.trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalledWith();
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

    await flushPromises();

    expect($modals.hide).toBeCalled();
  });

  test('Renders `confirm-ack-with-ticket` with empty modal', () => {
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
});
