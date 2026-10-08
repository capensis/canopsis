import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';
import { merge } from 'lodash';

import { MODALS } from '@/constants';

import ClickOutside from '@/services/click-outside';

import InfoPopupSetting from '@/components/modals/alarm/info-popup-setting/info-popup-setting.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'date-interval-selector': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const selectModalActionsButtons = wrapper => wrapper.findAll('.actions button.v-btn');
const selectCancelButton = wrapper => selectModalActionsButtons(wrapper).at(0);
const modalId = wrapper => ({ id: wrapper.props().modal.id });

const withModalInject = (options = {}) => {
  const rawModal = options.propsData?.modal ?? { config: {} };
  const modal = rawModal.id ? rawModal : { id: 'test-modal-id', ...rawModal };

  return merge({}, options, {
    propsData: {
      ...options.propsData,
      modal,
    },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
        $modal: modal,
      },
    },
  });
};

describe('info-popup-setting', () => {
  const $modals = mockModals();

  const shallowInfoPopupSetting = generateShallowRenderer(InfoPopupSetting, {
    stubs,
    attachTo: document.body,
  });
  const factory = (options = {}) => shallowInfoPopupSetting(withModalInject(options));

  test('Form submitted after trigger submit button', async () => {
    const action = vi.fn();
    const infoPopups = [
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
    ];
    const wrapper = factory({
      propsData: {
        modal: {
          config: { infoPopups, action },
        },
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(infoPopups);
    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Form submitted after trigger submit button without action', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {},
        },
      },
    });

    await wrapper.vm.submit();

    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Popup added after trigger add button', async () => {
    const action = vi.fn();
    const infoPopups = [
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
    ];
    const columns = [
      {
        value: Faker.datatype.string(),
      },
    ];
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            columns,
            infoPopups,
            action,
          },
        },
      },
    });

    const newPopup = {
      column: Faker.datatype.string(),
      template: Faker.datatype.string(),
    };

    wrapper.vm.addPopup();

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.addInfoPopup,
      config: {
        columns,
        variables: undefined,
        action: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action(newPopup);

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith([...infoPopups, newPopup]);
    expect($modals.hide).toBeCalled();
  });

  test('Popup edited after trigger edit button', async () => {
    const action = vi.fn();
    const infoPopups = [
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
    ];
    const columns = [
      {
        value: Faker.datatype.string(),
      },
    ];
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            columns,
            infoPopups,
            action,
          },
        },
      },
    });

    const newPopupData = {
      column: Faker.datatype.string(),
      template: Faker.datatype.string(),
    };

    wrapper.vm.editPopup(1, infoPopups[1]);

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.addInfoPopup,
      config: {
        columns,
        variables: undefined,
        popup: infoPopups[1],
        action: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action(newPopupData);

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith([infoPopups[0], newPopupData]);
    expect($modals.hide).toBeCalled();
  });

  test('Popup removed after trigger delete button', async () => {
    const action = vi.fn();
    const infoPopups = [
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
    ];
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            infoPopups,
            action,
          },
        },
      },
    });

    wrapper.vm.deletePopup(1);

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith([infoPopups[0]]);
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

    await selectCancelButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalled();
  });
});
