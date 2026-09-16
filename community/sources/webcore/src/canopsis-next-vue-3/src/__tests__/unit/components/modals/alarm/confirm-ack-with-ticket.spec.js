import { createButtonStub } from '@unit/stubs/button';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { ACK_MODAL_ACTIONS_TYPES } from '@/constants';

import ClickOutside from '@/services/click-outside';

import ConfirmAckWithTicket from '@/components/modals/alarm/confirm-ack-with-ticket.vue';

const stubs = {
  'c-modal-wrapper': createModalWrapperStub('c-modal-wrapper'),
  'v-btn': createButtonStub('v-btn'),
  'v-form': {
    emits: ['submit'],
    template: '<form @submit.prevent="$emit(\'submit\', $event)"><slot /></form>',
  },
  'ack-event-form': true,
  'alarm-general-table': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);

describe('confirm-ack-with-ticket', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(ConfirmAckWithTicket, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  test('Ack action called after trigger submit button', async () => {
    const action = vi.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: { action },
        },
      },
      mocks: {
        $modals,
      },
    });

    selectSubmitButton(wrapper).trigger('click');
    wrapper.find('form').trigger('submit');

    await flushPromises();

    expect(action).toHaveBeenCalledWith(
      {
        ack_resources: false,
        comment: '',
      },
      {
        needAssociateTicket: false,
        needDeclareTicket: false,
      },
    );
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('Ack with ticket action called after trigger submit button', async () => {
    const action = vi.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: { action },
        },
      },
      mocks: {
        $modals,
      },
    });

    wrapper.vm.actionType = ACK_MODAL_ACTIONS_TYPES.ackAndAssociateTicket;

    selectSubmitButton(wrapper).trigger('click');
    wrapper.find('form').trigger('submit');

    await flushPromises();

    expect(action).toHaveBeenCalledWith(
      {
        ack_resources: false,
        comment: '',
      },
      {
        needAssociateTicket: true,
        needDeclareTicket: false,
      },
    );
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('Ack with declare ticket action called after trigger submit button', async () => {
    const action = vi.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
            items: [
              {
                assigned_declare_ticket_rules: ['rule'],
                entity: { type: 'component' },
              },
            ],
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    wrapper.vm.actionType = ACK_MODAL_ACTIONS_TYPES.ackAndDeclareTicket;

    selectSubmitButton(wrapper).trigger('click');
    wrapper.find('form').trigger('submit');

    await flushPromises();

    expect(action).toHaveBeenCalledWith(
      {
        ack_resources: false,
        comment: '',
      },
      {
        needAssociateTicket: false,
        needDeclareTicket: true,
      },
    );
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
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

    selectCancelButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });
});
