import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

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
  const factory = generateShallowRenderer(ModalWrapper, {
    stubs,
    parentComponent: {
      provide: defaultProvide,
    },
  });
  const snapshotFactory = generateRenderer(ModalWrapper, {
    stubs,
    attachTo: document.body,
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
    const close = jest.fn();
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

  test('Renders `modal-wrapper` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `modal-wrapper` with custom props and slots', () => {
    const wrapper = snapshotFactory({
      propsData: {
        fillHeight: true,
        minimize: true,
        close: true,
        titleColor: 'red',
      },

      slots: {
        title: '<div>Title</div>',
        text: '<div>Text</div>',
        actions: '<div>Actions</div>',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `modal-wrapper` with minimized modal', () => {
    const wrapper = snapshotFactory({
      slots: {
        title: '<div>Title</div>',
        text: '<div>Text</div>',
        actions: '<div>Actions</div>',
      },

      parentComponent: {
        provide: {
          $modal: {
            id: MODAL_ID,
            minimized: true,
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
