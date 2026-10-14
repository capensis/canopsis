import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import ClickOutside from '@/services/click-outside';

import CreateCommentEvent from '@/components/modals/common/create-comment-event.vue';

const stubs = {
  'c-modal-wrapper': createModalWrapperStub('c-modal-wrapper'),
  'alarm-general-table': true,
  'alarm-comment-template-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

describe('create-comment-event', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(CreateCommentEvent, {
    stubs,
    attachTo: document.body,
    parentComponent: { provide: { $clickOutside: new ClickOutside() } },
    mocks: { $modals, $popups },
  });

  const createWrapper = (config = {}) =>
    factory({
      propsData: {
        modal: {
          id: 'test-modal-id',
          config,
        },
      },
    });

  test('submits a plain comment and closes the modal by id', async () => {
    const action = vi.fn();
    const wrapper = createWrapper({ action });

    wrapper.vm.values = { template: null, comment: 'Comment' };
    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith({ comment: 'Comment' });
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('converts template fields to struct_comment', async () => {
    const action = vi.fn();
    const template = {
      _id: 'template-id',
      fields: [
        { name: 'summary', required: true },
        { name: 'details', required: false },
      ],
    };
    const wrapper = createWrapper({ action, templates: [template] });

    wrapper.vm.values = {
      template,
      summary: 'Summary',
      details: 'Details',
    };
    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith({
      struct_comment: [
        { field: 'summary', message: 'Summary' },
        { field: 'details', message: 'Details' },
      ],
    });
  });

  test('omits empty template fields from struct_comment', async () => {
    const action = vi.fn();
    const template = { _id: 'template-id', fields: [{ name: 'summary' }, { name: 'details' }] };
    const wrapper = createWrapper({ action, templates: [template] });

    wrapper.vm.values = { template, summary: 'Summary', details: '' };
    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith({
      struct_comment: [{ field: 'summary', message: 'Summary' }],
    });
  });

  test('passes configured templates to the view model', () => {
    const templates = [{ _id: '1', name: 'Template', fields: [] }];
    const wrapper = createWrapper({ templates });

    expect(wrapper.vm.templates).toEqual(templates);
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
    const wrapper = createWrapper({});

    await wrapper.findAll('button.v-btn').at(0).trigger('click');
    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });
});
