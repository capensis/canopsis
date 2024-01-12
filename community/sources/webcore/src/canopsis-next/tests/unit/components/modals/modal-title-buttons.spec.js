import Faker from 'faker';
import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';
import { createMockedStoreModules } from '@unit/utils/store';
import { createActivatorElementStub } from '@unit/stubs/vuetify';

import ModalTitleButtons from '@/components/modals/modal-title-buttons.vue';

const stubs = {
  'modal-title-buttons': true,
  'v-tooltip': createActivatorElementStub('v-tooltip'),
};

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

  const factory = generateShallowRenderer(ModalTitleButtons, { stubs });
  const snapshotFactory = generateRenderer(ModalTitleButtons, { stubs });

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

  test('Modals minimize handler called after trigger minimize button', async () => {
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

    await flushPromises();

    selectButton(wrapper).vm.$emit('click');

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

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `modal-title-buttons` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        minimize: true,
        close: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
