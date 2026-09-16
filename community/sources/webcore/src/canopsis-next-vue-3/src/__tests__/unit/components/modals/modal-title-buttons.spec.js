import { createActivatorElementStub } from '@unit/stubs/vuetify';
import { mockModals } from '@unit/utils/mock-hooks';
import { createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import ModalTitleButtons from '@/components/modals/modal-title-buttons.vue';

const stubs = {
  'modal-title-buttons': true,
  'v-tooltip': createActivatorElementStub('v-tooltip'),
};

const selectButton = wrapper => wrapper.find('v-btn-stub');

describe('modal-title-buttons', () => {
  const $modals = mockModals();

  const hasMinimizedModal = vi.fn().mockReturnValue(true);
  const modalsModule = {
    name: 'modals',
    getters: {
      hasMinimizedModal,
    },
  };

  const store = createMockedStoreModules([modalsModule]);

  const factory = generateShallowRenderer(ModalTitleButtons, { stubs });

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

    closeButton.triggerCustomEvent('click');

    expect($modals.hide).toHaveBeenCalledWith({ id: modal.id });
  });

  test('Custom close handler called after close button', () => {
    const modal = {
      id: Faker.datatype.string(),
    };
    const close = vi.fn();
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

    closeButton.triggerCustomEvent('click');

    expect(close).toHaveBeenCalled();
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

    selectButton(wrapper).triggerCustomEvent('click');

    expect($modals.minimize).toHaveBeenCalledWith({ id: modal.id });
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

    maximizeButton.triggerCustomEvent('click');

    expect($modals.maximize).toHaveBeenCalledWith({ id: modal.id });
  });
});
