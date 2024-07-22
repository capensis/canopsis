import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockConsole, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';

import ClickOutside from '@/services/click-outside';

import { formToEntity, entityToForm } from '@/helpers/entities/entity/form';

import CreateEntity from '@/components/modals/entity/create-entity.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'entity-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'entity-form': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectEntityForm = wrapper => wrapper.find('entity-form-stub');

describe('create-entity', () => {
  const $modals = mockModals();
  const $popups = mockPopups();
  const consoleMock = mockConsole();

  const defaultEntityForm = entityToForm();
  const defaultEntity = formToEntity(defaultEntityForm);

  const factory = generateShallowRenderer(CreateEntity, {
    stubs,
    attachTo: document.body,
    mocks: { $modals, $popups },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });
  const snapshotFactory = generateRenderer(CreateEntity, {
    stubs: snapshotStubs,
    mocks: { $modals, $popups },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  test('Form submitted after trigger submit button', async () => {
    const action = jest.fn();
    const modal = {
      config: {
        action,
      },
    };
    const wrapper = factory({
      propsData: {
        modal,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith(defaultEntity);
    expect($modals.hide).toHaveBeenCalledWith(modal);
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
          },
        },
      },
    });

    const entityForm = selectEntityForm(wrapper);
    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      vm: entityForm.vm,
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).not.toHaveBeenCalled();
    expect($modals.hide).not.toHaveBeenCalled();
  });

  test('Form submitted after trigger submit button without action', async () => {
    const modal = {
      config: {},
    };
    const wrapper = factory({
      propsData: {
        modal,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalledWith(modal);
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      name: 'Name error',
      enabled: 'Category error',
      infos: 'Category error',
    };
    const action = jest.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
    const modal = {
      config: {
        action,
      },
    };
    const wrapper = factory({
      propsData: {
        modal,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toBeCalledWith(defaultEntity);
    expect($modals.hide).not.toHaveBeenCalled();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = jest.fn().mockRejectedValue(errors);
    const customEntity = {
      ...defaultEntity,
      name: 'Custom name',
    };
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            entity: customEntity,
            action,
          },
        },
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(consoleMock.error).toBeCalledWith(errors);
    expect($popups.error).toBeCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toBeCalledWith({
      ...defaultEntity,
      name: customEntity.name,
    });
    expect($modals.hide).not.toHaveBeenCalled();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
          },
        },
      },
    });

    const newForm = {
      ...defaultEntityForm,
      name: Faker.datatype.string(),
    };

    await selectEntityForm(wrapper).triggerCustomEvent('input', newForm);
    await selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith({
      ...defaultEntity,
      name: newForm.name,
    });
    expect($modals.hide).toBeCalled();
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {},
        },
      },
    });

    selectCancelButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalled();
  });

  test('Renders `create-entity` with empty modal', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {},
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `create-entity` with pbehavior', () => {
    const entity = {
      ...defaultEntity,
      name: 'Entity name',
    };
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            entity,
            title: 'Custom create entity title',
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
