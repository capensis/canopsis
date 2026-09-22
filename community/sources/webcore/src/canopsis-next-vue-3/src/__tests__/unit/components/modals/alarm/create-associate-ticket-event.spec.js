import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import ClickOutside from '@/services/click-outside';

import CreateAssociateTicketEvent from '@/components/modals/declare-ticket/create-associate-ticket-event.vue';

const stubs = {
  'c-modal-wrapper': createModalWrapperStub('c-modal-wrapper'),
  'c-alert': true,
  'c-checkbox': true,
  'alarm-general-table': true,
  'associate-ticket-event-form': true,
  'c-description-field': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

describe('create-associate-ticket-event', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(CreateAssociateTicketEvent, {
    stubs,
    attachTo: document.body,
    parentComponent: { provide: { $clickOutside: new ClickOutside() } },
    mocks: { $modals, $popups },
  });

  const createWrapper = (config = {}) => factory({
    propsData: {
      modal: {
        id: 'test-modal-id',
        config,
      },
    },
  });

  test('submits the default form and closes the modal by id', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({ action });

    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith({
      comment: '',
      ticket: '',
      data: {},
      system_name: '',
      ticket_resources: false,
      url: '',
      url_title: '',
    });
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('maps updated form values before submitting', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({ action });

    wrapper.vm.updateForm({
      comment: 'Comment',
      ticket: 'INC-42',
      ticket_url: 'https://tickets.test/INC-42',
      ticket_url_title: 'Ticket',
      ticket_system_name: 'Service desk',
      ticket_resources: true,
      mapping: [{ key: 'pair-id', text: 'priority', value: 'high' }],
    });
    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith({
      comment: 'Comment',
      ticket: 'INC-42',
      data: { priority: 'high' },
      system_name: 'Service desk',
      ticket_resources: true,
      url: 'https://tickets.test/INC-42',
      url_title: 'Ticket',
    });
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
