import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';

import { MODALS } from '@/constants';

import FieldTemplate from '@/components/sidebars/counter/form/fields/template.vue';

const stubs = {
  'widget-settings-item': true,
};

const selectButton = wrapper => wrapper.find('v-btn-stub');

describe('field-template', () => {
  const $modals = mockModals();
  const factory = generateShallowRenderer(FieldTemplate, {

    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(FieldTemplate, { stubs });

  test('Value changed after trigger color indicator field', () => {
    const title = Faker.datatype.string();
    const value = Faker.datatype.string();
    const variables = [{
      value: Faker.datatype.string(),
    }];

    const wrapper = factory({
      propsData: {
        title,
        value,
        variables,
      },
    });

    selectButton(wrapper).triggerCustomEvent('click');

    expect($modals.show).toBeCalledWith({
      name: MODALS.textEditor,
      config: {
        title,
        text: value,
        variables,
        action: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    const newTemplate = Faker.datatype.string();

    modalArguments.config.action(newTemplate);

    expect(wrapper).toEmit('input', newTemplate);
  });

  test('Renders `field-template` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: '<div>value-template</div>',
        title: 'Custom required title',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `field-template` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: '<div>template</div>',
        title: 'Custom title',
        variables: [{}],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
