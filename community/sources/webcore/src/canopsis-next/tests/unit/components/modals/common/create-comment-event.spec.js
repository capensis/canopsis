import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createInputStub } from '@unit/stubs/input';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { fakeAlarms } from '@unit/data/alarm';

import ClickOutside from '@/services/click-outside';

import CreateCommentEvent from '@/components/modals/common/create-comment-event.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarm-general-table': true,
  'alarm-comment-template-form': createInputStub('alarm-comment-template-form'),
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarm-general-table': true,
  'alarm-comment-template-form': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectCommentForm = wrapper => wrapper.find('.alarm-comment-template-form');

describe('create-comment-event', () => {
  const timestamp = 1386435600000;

  mockDateNow(timestamp);
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(CreateCommentEvent, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });
  const snapshotFactory = generateRenderer(CreateCommentEvent, {
    stubs: snapshotStubs,
    propsData: {
      modal: {
        config: {},
      },
    },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  test('Form submitted after trigger submit button', async () => {
    const action = jest.fn();
    const config = { action };

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

    const comment = Faker.datatype.string();

    wrapper.vm.form.comment = comment;

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledTimes(1);
    expect(action).toHaveBeenCalledWith({ comment });
    expect($modals.hide).toHaveBeenCalled();
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const action = jest.fn();
    const formErrors = {
      comment: 'Comment error',
    };
    action.mockRejectedValueOnce({ ...formErrors, unavailableField: 'Error' });
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

    const comment = Faker.datatype.string();

    wrapper.vm.form.comment = comment;

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toHaveBeenCalledTimes(1);
    expect(action).toHaveBeenCalledWith({ comment });
    expect($modals.hide).not.toHaveBeenCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = jest.fn();
    const config = { action };
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

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => wrapper.vm,
      vm: wrapper.vm,
    });

    selectSubmitButton(wrapper).trigger('click');

    expect(action).not.toHaveBeenCalled();
    expect($modals.hide).not.toHaveBeenCalled();

    validator.detach('name');
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const action = jest.fn();
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    action.mockRejectedValueOnce(errors);

    const wrapper = factory({
      propsData: {
        modal: {
          config: { action },
        },
      },
      mocks: {
        $modals,
        $popups,
      },
    });

    const comment = Faker.datatype.string();

    wrapper.vm.form.comment = comment;

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(consoleErrorSpy).toHaveBeenCalledWith(errors);
    expect($popups.error).toHaveBeenCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toHaveBeenCalledTimes(1);
    expect(action).toHaveBeenCalledWith({ comment });
    expect($modals.hide).not.toHaveBeenCalledWith();

    consoleErrorSpy.mockClear();
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

    expect($modals.hide).toHaveBeenCalled();
  });

  test('Comment form bound to form via v-field directive', async () => {
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

    const commentForm = selectCommentForm(wrapper);

    expect(commentForm.vm.value).toEqual({
      template: null,
      comment: '',
    });

    const newComment = Faker.datatype.string();
    wrapper.vm.form.comment = newComment;

    await wrapper.vm.$nextTick();

    expect(commentForm.vm.value.comment).toBe(newComment);
  });

  test('Submit button disabled state reflects submitting state', async () => {
    const action = jest.fn().mockImplementation(() => new Promise(() => {}));
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

    const submitButton = selectSubmitButton(wrapper);

    expect(submitButton.attributes('disabled')).toBeUndefined();
    expect(wrapper.vm.submitting).toBe(false);

    const comment = Faker.datatype.string();
    wrapper.vm.form.comment = comment;

    submitButton.trigger('click');

    await wrapper.vm.$nextTick();

    expect(wrapper.vm.submitting).toBe(true);
    expect(wrapper.vm.isDisabled).toBe(true);
  });

  test('Templates from config are passed to component', () => {
    const templates = [
      { _id: '1', name: 'Template 1', fields: [] },
      { _id: '2', name: 'Template 2', fields: [{ name: 'field1', required: true }] },
    ];

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            templates,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper.vm.templates).toEqual(templates);
  });

  test('Form submitted with template sends struct_comment instead of comment', async () => {
    const action = jest.fn();
    const template = {
      _id: '1',
      name: 'Template 1',
      fields: [
        { name: 'field1', label: 'Field 1', required: true },
        { name: 'field2', label: 'Field 2', required: false },
      ],
    };

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
            templates: [template],
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    wrapper.vm.form.template = template;
    wrapper.vm.form.field1 = Faker.datatype.string();
    wrapper.vm.form.field2 = Faker.datatype.string();

    const submitButton = selectSubmitButton(wrapper);
    submitButton.trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledTimes(1);
    expect(action).toHaveBeenCalledWith({
      struct_comment: [
        { field: 'field1', message: wrapper.vm.form.field1 },
        { field: 'field2', message: wrapper.vm.form.field2 },
      ],
    });
    expect($modals.hide).toHaveBeenCalled();
  });

  test('Form submitted without template sends comment', async () => {
    const action = jest.fn();
    const comment = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
            templates: [],
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    wrapper.vm.form.template = null;
    wrapper.vm.form.comment = comment;

    const submitButton = selectSubmitButton(wrapper);
    submitButton.trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledTimes(1);
    expect(action).toHaveBeenCalledWith({
      comment,
    });
    expect(action).not.toHaveBeenCalledWith(expect.objectContaining({
      struct_comment: expect.any(Array),
    }));
    expect($modals.hide).toHaveBeenCalled();
  });

  test('Struct_comment includes all template fields', async () => {
    const action = jest.fn();
    const template = {
      _id: '2',
      name: 'Multi Field Template',
      fields: [
        { name: 'summary', label: 'Summary', required: true },
        { name: 'details', label: 'Details', required: false },
        { name: 'resolution', label: 'Resolution', required: false },
      ],
    };

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
            templates: [template],
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    wrapper.vm.form.template = template;
    wrapper.vm.form.summary = Faker.datatype.string();
    wrapper.vm.form.details = Faker.datatype.string();
    wrapper.vm.form.resolution = Faker.datatype.string();

    const submitButton = selectSubmitButton(wrapper);
    submitButton.trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledTimes(1);
    expect(action).toHaveBeenCalledWith({
      struct_comment: [
        { field: 'summary', message: wrapper.vm.form.summary },
        { field: 'details', message: wrapper.vm.form.details },
        { field: 'resolution', message: wrapper.vm.form.resolution },
      ],
    });
    expect($modals.hide).toHaveBeenCalled();
  });

  test('Struct_comment handles empty template fields array', async () => {
    const action = jest.fn();
    const template = {
      _id: '3',
      name: 'Empty Template',
      fields: [],
    };

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
            templates: [template],
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    wrapper.vm.form.template = template;

    const submitButton = selectSubmitButton(wrapper);
    submitButton.trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledTimes(1);
    expect(action).toHaveBeenCalledWith({
      struct_comment: [],
    });
    expect($modals.hide).toHaveBeenCalled();
  });

  test('Form switches from comment to struct_comment when template is selected', async () => {
    const action = jest.fn();
    const template = {
      _id: '4',
      name: 'Switch Template',
      fields: [
        { name: 'field1', label: 'Field 1', required: true },
      ],
    };

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
            templates: [template],
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    // First, set a comment without template
    wrapper.vm.form.comment = Faker.datatype.string();
    wrapper.vm.form.template = null;

    // Then select a template
    wrapper.vm.form.template = template;
    wrapper.vm.form.field1 = Faker.datatype.string();

    const submitButton = selectSubmitButton(wrapper);
    submitButton.trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledTimes(1);
    expect(action).toHaveBeenCalledWith({
      struct_comment: [
        { field: 'field1', message: wrapper.vm.form.field1 },
      ],
    });
    expect(action).not.toHaveBeenCalledWith(expect.objectContaining({
      comment: expect.any(String),
    }));
    expect($modals.hide).toHaveBeenCalled();
  });

  test('Renders `create-comment-event` with empty modal', () => {
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

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `create-comment-event` with items', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            items: fakeAlarms(10),
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
