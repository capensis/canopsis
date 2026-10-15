import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { ACK_MODAL_ACTIONS_TYPES } from '@/constants';

import ClickOutside from '@/services/click-outside';

import CreateAckEvent from '@/components/modals/alarm/create-ack-event.vue';

const stubs = {
  'c-modal-wrapper': createModalWrapperStub('c-modal-wrapper'),
  'alarm-general-table': true,
  'ack-event-form': true,
  'c-radio-group': true,
  'v-radio': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

describe('create-ack-event', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(CreateAckEvent, {
    stubs,
    attachTo: document.body,
    parentComponent: { provide: { $clickOutside: new ClickOutside() } },
    mocks: { $modals, $popups },
  });

  const createWrapper = (config = {}) => factory({ propsData: { modal: { config } } });

  test('submits default form values and closes the modal by id', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({ action });

    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith(
      { ack_resources: false, comment: '' },
      { needAssociateTicket: false, needDeclareTicket: false },
    );
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('passes the associate-ticket action flag', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({ action });

    wrapper.vm.actionType = ACK_MODAL_ACTIONS_TYPES.ackAndAssociateTicket;
    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith(
      { ack_resources: false, comment: '' },
      { needAssociateTicket: true, needDeclareTicket: false },
    );
  });

  test('passes the declare-ticket action flag when declare rules exist', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({
      action,
      items: [{ assigned_declare_ticket_rules: ['rule'], entity: { type: 'component' } }],
    });

    wrapper.vm.actionType = ACK_MODAL_ACTIONS_TYPES.ackAndDeclareTicket;
    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith(
      { ack_resources: false, comment: '' },
      { needAssociateTicket: false, needDeclareTicket: true },
    );
  });

  test('shows an error popup and keeps the modal open for unmapped action errors', async () => {
    const errors = { unavailableField: 'Error', anotherUnavailableField: 'Second error' };
    const action = vi.fn().mockRejectedValueOnce(errors);
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();
    const wrapper = createWrapper({ action });

    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(consoleErrorSpy).toHaveBeenCalledWith(errors);
    expect($popups.error).toHaveBeenCalledWith({ text: 'Error\nSecond error' });
    expect($modals.hide).not.toHaveBeenCalled();
    consoleErrorSpy.mockRestore();
  });

  test('closes without an action when cancel is clicked', async () => {
    const wrapper = createWrapper();

    await wrapper.findAll('button.v-btn').at(0).trigger('click');
    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });
});
