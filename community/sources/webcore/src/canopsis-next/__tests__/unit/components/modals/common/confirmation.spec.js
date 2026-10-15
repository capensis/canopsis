import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import Confirmation from '@/components/modals/common/confirmation.vue';

const stubs = {
  'modal-wrapper': true,
};

const selectSubmitButton = wrapper => wrapper.findAll('v-btn-stub').at(1);
const selectCancelButton = wrapper => wrapper.findAll('v-btn-stub').at(0);

describe('confirmation', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(Confirmation, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Submit action called after trigger submit button', async () => {
    const cancel = vi.fn();
    const action = vi.fn();
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

    submitButton.triggerCustomEvent('click', new Event('click'));

    await flushPromises();

    expect(action).toHaveBeenCalledWith();
    expect($modals.hide).toHaveBeenCalled();

    wrapper.destroy();

    expect(cancel).not.toHaveBeenCalled();
  });

  test('Hidden action called after trigger submit button without action', async () => {
    const cancel = vi.fn();
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

    submitButton.triggerCustomEvent('click', new Event('click'));

    await flushPromises();

    expect($modals.hide).toHaveBeenCalled();

    wrapper.destroy();

    expect(cancel).not.toHaveBeenCalled();
  });

  test('Cancel action called after trigger cancel button', () => {
    const cancel = vi.fn();
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

    selectCancelButton(wrapper).triggerCustomEvent('click', new Event('click'));

    expect($modals.hide).toHaveBeenCalled();

    wrapper.destroy();

    expect(cancel).toHaveBeenCalledWith(true);
  });

  test('Cancel action called after destroy', () => {
    const cancel = vi.fn();
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

    expect(cancel).toHaveBeenCalledWith(false);
  });
});
