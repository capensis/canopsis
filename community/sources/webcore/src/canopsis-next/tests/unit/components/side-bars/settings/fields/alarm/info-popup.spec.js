import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { MODALS } from '@/constants';

import InfoPopup from '@/components/sidebars/settings/fields/alarm/info-popup.vue';

const localVue = createVueInstance();

const stubs = {
  'v-btn': createButtonStub('v-btn'),
};

const factory = (options = {}) => shallowMount(InfoPopup, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(InfoPopup, {
  localVue,

  ...options,
});

const selectCreateOrEditButton = wrapper => wrapper.find('button.v-btn');

describe('info-popup', () => {
  const $modals = mockModals();
  const popups = [{
    column: Faker.datatype.string(),
    template: Faker.datatype.string(),
  }, {
    column: Faker.datatype.string(),
    template: Faker.datatype.string(),
  }];
  const columns = [{
    label: Faker.datatype.string(),
    value: Faker.datatype.string(),
  }, {
    label: Faker.datatype.string(),
    value: Faker.datatype.string(),
  }];

  it('Info popup setting modal opened after trigger create button', () => {
    const wrapper = factory({
      propsData: {
        popups,
        columns,
      },
      mocks: {
        $modals,
      },
    });

    const createButton = selectCreateOrEditButton(wrapper);

    createButton.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.infoPopupSetting,
        config: {
          columns,
          infoPopups: popups,
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = [];

    modalArguments.config.action(actionValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(actionValue);
  });

  it('Renders `info-popup` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `info-popup` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        popups: [
          {
            column: 'Column',
            template: 'Template',
          },
        ],
        columns: [
          {
            value: 'Value',
            label: 'Label',
          },
        ],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
