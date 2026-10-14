import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockConsole, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import ClickOutside from '@/services/click-outside';

import { formToEntity, entityToForm } from '@/helpers/entities/entity/form';

import CreateEntity from '@/components/modals/entity/create-entity.vue';

const stubs = {
  'c-modal-wrapper': createModalWrapperStub('c-modal-wrapper'),
  'entity-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
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

  const render = generateShallowRenderer(CreateEntity, {
    stubs,
    attachTo: document.body,
    mocks: { $modals, $popups },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  const factory = options =>
    render({
      ...options,
      propsData: {
        ...options?.propsData,
        modal: {
          id: 'test-modal-id',
          ...options?.propsData?.modal,
        },
      },
    });

  test('Form submitted after trigger submit button', async () => {
    const action = vi.fn();
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

    expect(action).toHaveBeenCalledWith(defaultEntity);
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('submits the default form through the Vue 3 form contract', async () => {
    const action = vi.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
          },
        },
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledWith(defaultEntity);
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
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

    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      name: 'Name error',
      enabled: 'Category error',
      infos: 'Category error',
    };
    const action = vi.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
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

    expect(action).toHaveBeenCalledWith(defaultEntity);
    expect($modals.hide).not.toHaveBeenCalledWith({ id: 'test-modal-id' });
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = vi.fn().mockRejectedValue(errors);
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

    expect(consoleMock.error).toHaveBeenCalledWith(errors);
    expect($popups.error).toHaveBeenCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toHaveBeenCalledWith({
      ...defaultEntity,
      name: customEntity.name,
    });
    expect($modals.hide).not.toHaveBeenCalled();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const action = vi.fn();
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

    expect(action).toHaveBeenCalledWith({
      ...defaultEntity,
      name: newForm.name,
    });
    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
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

    expect($modals.hide).toHaveBeenCalledWith({ id: 'test-modal-id' });
  });
});
