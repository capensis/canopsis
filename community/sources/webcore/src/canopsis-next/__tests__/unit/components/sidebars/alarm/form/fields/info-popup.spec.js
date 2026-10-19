import { createButtonStub } from '@unit/stubs/button';
import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { MODALS } from '@/constants';

import InfoPopup from '@/components/sidebars/alarm/form/fields/info-popup.vue';

const stubs = {
  'widget-settings-item': true,
  'v-btn': createButtonStub('v-btn'),
};

const selectCreateOrEditButton = wrapper => wrapper.find('button.v-btn');

describe('info-popup', () => {
  const $modals = mockModals();
  const popups = [
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
      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    },
    {
      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    },
  ];

  const factory = generateShallowRenderer(InfoPopup, { stubs });

  it('Info popup setting modal opened after trigger create button', () => {
    const wrapper = factory({
      propsData: {
        value: popups,
        columns,
      },
      mocks: {
        $modals,
      },
    });

    const createButton = selectCreateOrEditButton(wrapper);

    createButton.trigger('click');

    expect($modals.show).toHaveBeenCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.infoPopupSetting,
      config: {
        columns,
        infoPopups: popups,
        variables: [],
        action: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = [];

    modalArguments.config.action(actionValue);

    expect(wrapper).toEmitInput(actionValue);
  });
});
