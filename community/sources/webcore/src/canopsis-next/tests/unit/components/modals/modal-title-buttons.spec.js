import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';

import { mockModals } from '@unit/utils/mock-hooks';
import { createMockedStoreModules } from '@unit/utils/store';
import ModalTitleButtons from '@/components/modals/modal-title-buttons.vue';

const localVue = createVueInstance();

const stubs = {
  'modal-title-buttons': true,
};

const factory = (options = {}) => shallowMount(ModalTitleButtons, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(ModalTitleButtons, {
  localVue,
  stubs,

  ...options,
});

const selectButton = wrapper => wrapper.find('v-btn-stub');

describe('modal-title-buttons', () => {
  const $modals = mockModals();

  const hasMinimizedModal = jest.fn().mockReturnValue(true);
  const modalsModule = {
    name: 'modals',
    getters: {
      hasMinimizedModal,
    },
  };

  const store = createMockedStoreModules([
    modalsModule,
  ]);

  test('Modals hide handler called after close button', () => {
    const modal = {
      id: Faker.datatype.string(),
    };
    const wrapper = factory({
      store,
      propsData: {
        close: true,
      },
      parentComponent: {
        provide: {
          $modal: modal,
        },
      },
      mocks: {
        $modals,
      },
    });

    const closeButton = selectButton(wrapper);

    closeButton.vm.$emit('click');

    expect($modals.hide).toBeCalledWith({ id: modal.id });
  });

  test('Custom close handler called after close button', () => {
    const modal = {
      id: Faker.datatype.string(),
    };
    const close = jest.fn();
    const wrapper = factory({
      store,
      propsData: {
        close,
      },
      parentComponent: {
        provide: {
          $modal: modal,
        },
      },
      mocks: {
        $modals,
      },
    });

    const closeButton = selectButton(wrapper);

    closeButton.vm.$emit('click');

    expect(close).toBeCalled();
  });

  test('Modals minimize handler called after trigger minimize button', () => {
    const modal = {
      id: Faker.datatype.string(),
    };
    const wrapper = factory({
      store,
      propsData: {
        minimize: true,
      },
      parentComponent: {
        provide: {
          $modal: modal,
        },
      },
      mocks: {
        $modals,
      },
    });

    const minimizeButton = selectButton(wrapper);

    minimizeButton.vm.$emit('click');

    expect($modals.minimize).toBeCalledWith({ id: modal.id });
  });

  test('Modals maximize handler called after trigger maximize button', () => {
    const modal = {
      id: Faker.datatype.string(),
      minimized: true,
    };
    const wrapper = factory({
      store,
      propsData: {
        minimize: true,
      },
      parentComponent: {
        provide: {
          $modal: modal,
        },
      },
      mocks: {
        $modals,
      },
    });

    const maximizeButton = selectButton(wrapper);

    maximizeButton.vm.$emit('click');

    expect($modals.maximize).toBeCalledWith({ id: modal.id });
  });

  test('Renders `modal-title-buttons` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `modal-title-buttons` with custom props', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        minimize: true,
        close: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `modal-title-buttons` with minimized modal', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        minimize: true,
      },

      parentComponent: {
        provide: {
          $modal: {
            minimized: true,
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
