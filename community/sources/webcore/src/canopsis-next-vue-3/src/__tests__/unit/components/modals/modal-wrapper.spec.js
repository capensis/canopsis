import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';

import ModalWrapper from '@/components/modals/modal-wrapper.vue';

const MODAL_ID = 'modal-wrapper-test';

const stubs = {
  'modal-title-buttons': true,
  'modal-mass-actions-panel': true,
};

const defaultProvide = {
  $modal: {
    id: MODAL_ID,
  },
};

const selectModalTitleButtons = wrapper => wrapper.find('modal-title-buttons-stub');
const selectModalMassActionsPanel = wrapper => wrapper.find('modal-mass-actions-panel-stub');

describe('modal-wrapper', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(ModalWrapper, {
    stubs,
    mocks: { $modals },
    parentComponent: {
      provide: defaultProvide,
    },
  });

  test('Modal mass actions panel receives id from injected modal', () => {
    const modalId = 'alarms-list-modal';
    const wrapper = factory({
      slots: {
        title: 'Title',
      },
      parentComponent: {
        provide: {
          $modal: {
            id: modalId,
          },
        },
      },
    });

    const modalMassActionsPanel = selectModalMassActionsPanel(wrapper);

    expect(modalMassActionsPanel.attributes('id')).toBe(modalId);
  });

  test('Close handler called after trigger close in the title', () => {
    const close = vi.fn();
    const wrapper = factory({
      propsData: {
        close,
      },
      slots: {
        title: 'Title',
      },
    });

    const modalTitleButtons = selectModalTitleButtons(wrapper);

    modalTitleButtons.vm.close();

    expect(close).toHaveBeenCalled();
  });
});
