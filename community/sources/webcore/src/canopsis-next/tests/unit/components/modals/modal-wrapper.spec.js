import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import ModalWrapper from '@/components/modals/modal-wrapper.vue';

const stubs = {
  'modal-title-buttons': true,
};

const selectModalTitleButtons = wrapper => wrapper.find('modal-title-buttons-stub');

describe('modal-wrapper', () => {
  const factory = generateShallowRenderer(ModalWrapper, {
    stubs,
    parentComponent: {
      provide: {
        $modal: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(ModalWrapper, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $modal: {},
      },
    },
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

    expect(close).toBeCalled();
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
            minimized: true,
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
