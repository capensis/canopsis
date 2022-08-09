import flushPromises from 'flush-promises';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';

import { mockModals } from '@unit/utils/mock-hooks';
import Confirmation from '@/components/modals/common/confirmation.vue';

const localVue = createVueInstance();

const stubs = {
  'modal-wrapper': true,
};

const factory = (options = {}) => shallowMount(Confirmation, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(Confirmation, {
  localVue,
  stubs,

  ...options,
});

const selectSubmitButton = wrapper => wrapper.findAll('v-btn-stub').at(1);
const selectCancelButton = wrapper => wrapper.findAll('v-btn-stub').at(0);

describe('confirmation', () => {
  const $modals = mockModals();

  test('Submit action called after trigger submit button', async () => {
    const cancel = jest.fn();
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            cancel,
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.vm.$emit('click', new Event('click'));

    await flushPromises();

    expect(action).toBeCalledWith();
    expect($modals.hide).toBeCalledWith();

    wrapper.destroy();

    expect(cancel).not.toBeCalled();
  });

  test('Hidden action called after trigger submit button without action', async () => {
    const cancel = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            cancel,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.vm.$emit('click', new Event('click'));

    await flushPromises();

    expect($modals.hide).toBeCalledWith();

    wrapper.destroy();

    expect(cancel).not.toBeCalled();
  });

  test('Cancel action called after trigger cancel button', () => {
    const cancel = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            cancel,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const cancelButton = selectCancelButton(wrapper);

    cancelButton.vm.$emit('click', new Event('click'));

    expect($modals.hide).toBeCalledWith();

    wrapper.destroy();

    expect(cancel).toBeCalledWith(true);
  });

  test('Cancel action called after destroy', () => {
    const cancel = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            cancel,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    wrapper.destroy();

    expect(cancel).toBeCalledWith(false);
  });

  test('Renders `confirmation` with empty modal', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {},
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `confirmation` with text, title and actions', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            title: 'Confirmation title',
            text: 'Confirmation text',
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `confirmation` with hidden title', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            hideTitle: true,
            text: 'Confirmation text',
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
