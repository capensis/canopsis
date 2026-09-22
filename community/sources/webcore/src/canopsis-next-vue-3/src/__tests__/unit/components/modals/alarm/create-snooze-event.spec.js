import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { TIME_UNITS } from '@/constants';
import ClickOutside from '@/services/click-outside';

import CreateSnoozeEvent from '@/components/modals/alarm/create-snooze-event.vue';

const stubs = {
  'c-modal-wrapper': createModalWrapperStub('c-modal-wrapper'),
  'snooze-event-form': {
    props: { isNoteRequired: Boolean },
    template: '<div class="snooze-event-form-stub" />',
  },
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

describe('create-snooze-event', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(CreateSnoozeEvent, {
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

  test('submits default snooze values and closes the modal by id', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({ action });

    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith({
      duration: { unit: TIME_UNITS.minute, value: 1 },
      comment: '',
    });
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('forwards the note-required setting to the snooze form', () => {
    const wrapper = createWrapper({ isNoteRequired: true });

    expect(wrapper.vm.config.isNoteRequired).toBe(true);
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
