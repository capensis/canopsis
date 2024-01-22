import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { mockModals } from '@unit/utils/mock-hooks';
import FlowchartColorField from '@/components/common/flowchart/fields/flowchart-color-field.vue';
import { MODALS } from '@/constants';

const selectCheckboxField = wrapper => wrapper.find('v-checkbox-stub');
const selectButton = wrapper => wrapper.find('v-btn-stub');

describe('flowchart-color-field', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(FlowchartColorField);
  const snapshotFactory = generateRenderer(FlowchartColorField);

  test('Value changed to first palette color after trigger select field with true', () => {
    const firstColor = Faker.internet.color();
    const wrapper = factory({
      propsData: {
        palette: [firstColor],
      },
    });

    const checkboxField = selectCheckboxField(wrapper);

    checkboxField.triggerCustomEvent('change', true);

    expect(wrapper).toEmit('input', firstColor);
  });

  test('Value changed to transparent after trigger select field with false', () => {
    const wrapper = factory({
      propsData: {
        color: Faker.internet.color(),
      },
    });

    const checkboxField = selectCheckboxField(wrapper);

    checkboxField.triggerCustomEvent('change', false);

    expect(wrapper).toEmit('input', 'transparent');
  });

  test('Color picker modal showed after trigger button', () => {
    const color = Faker.internet.color();
    const palette = [Faker.internet.color()];
    const wrapper = factory({
      propsData: {
        value: color,
        palette,
      },
      mocks: {
        $modals,
      },
    });

    const button = selectButton(wrapper);

    button.triggerCustomEvent('click');
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.colorPicker,
        config: {
          color,
          palette,
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const newColor = Faker.internet.color();

    modalArguments.config.action(newColor);

    expect(wrapper).toEmit('input', newColor);
  });

  test('Renders `flowchart-color-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `flowchart-color-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: '#000000',
        palette: ['#010101'],
        label: 'Custom label',
        hideCheckbox: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
