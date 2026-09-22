import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import ClickOutside from '@/services/click-outside';

import CreateServicePauseEvent from '@/components/modals/service/create-service-pause-event.vue';

const stubs = {
  'c-modal-wrapper': createModalWrapperStub('c-modal-wrapper'),
  'service-pause-event-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

describe('create-service-pause-event', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(CreateServicePauseEvent, {
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

  test('initializes the default pause form', () => {
    const wrapper = createWrapper();

    expect(wrapper.vm.form).toEqual({ comment: '', reason: '' });
  });

  test('submits updated form values and closes by id', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({ action });
    const form = { comment: 'Comment', reason: 'Maintenance' };

    wrapper.vm.form = form;
    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith(form);
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('closes without an action', async () => {
    const wrapper = createWrapper();

    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
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

  test('closes when cancel is clicked', async () => {
    const wrapper = createWrapper();

    await wrapper.findAll('button.v-btn').at(0).trigger('click');
    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });
});
