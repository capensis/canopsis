import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { ALARM_STATES } from '@/constants';
import ClickOutside from '@/services/click-outside';

import CreateChangeStateEvent from '@/components/modals/alarm/create-change-state-event.vue';

const stubs = {
  'c-modal-wrapper': createModalWrapperStub('c-modal-wrapper'),
  'c-change-state-field': true,
  'alarm-general-table': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

describe('create-change-state-event', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(CreateChangeStateEvent, {
    stubs,
    attachTo: document.body,
    parentComponent: { provide: { $clickOutside: new ClickOutside() } },
    mocks: { $modals, $popups },
  });

  const alarm = state => ({
    v: { state: { val: state } },
    entity: { type: 'component' },
  });

  const createWrapper = (config = {}) => factory({
    propsData: {
      modal: {
        id: 'test-modal-id',
        config,
      },
    },
  });

  test('initializes the state from a single alarm', () => {
    const wrapper = createWrapper({ items: [alarm(ALARM_STATES.critical)] });

    expect(wrapper.vm.values).toEqual({
      comment: '',
      state: ALARM_STATES.critical,
    });
  });

  test('uses the major state for multiple alarms', () => {
    const wrapper = createWrapper({
      items: [alarm(ALARM_STATES.critical), alarm(ALARM_STATES.minor)],
    });

    expect(wrapper.vm.values).toEqual({ comment: '', state: ALARM_STATES.major });
  });

  test('submits the form and closes the modal by id', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({ items: [alarm(ALARM_STATES.critical)], action });

    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith({ comment: '', state: ALARM_STATES.critical });
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('submits updated state-field values', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({ items: [alarm(ALARM_STATES.critical)], action });

    wrapper.vm.values = { comment: 'Comment', state: ALARM_STATES.major };
    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith({ comment: 'Comment', state: ALARM_STATES.major });
  });

  test('shows an error popup and keeps the modal open for unmapped action errors', async () => {
    const errors = { unavailableField: 'Error', anotherUnavailableField: 'Second error' };
    const action = vi.fn().mockRejectedValueOnce(errors);
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();
    const wrapper = createWrapper({ items: [alarm(ALARM_STATES.critical)], action });

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
