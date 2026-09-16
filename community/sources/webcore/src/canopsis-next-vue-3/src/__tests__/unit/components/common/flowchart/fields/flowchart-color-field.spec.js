import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { MODALS } from '@/constants';

import FlowchartColorField from '@/components/common/flowchart/fields/flowchart-color-field.vue';

const selectCheckboxField = wrapper => wrapper.find('v-checkbox-stub');
const selectButton = wrapper => wrapper.find('v-btn-stub');

describe('flowchart-color-field', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(FlowchartColorField);

  test('Value changed to first palette color after trigger select field with true', () => {
    const firstColor = Faker.internet.color();
    const wrapper = factory({
      propsData: {
        palette: [firstColor],
      },
    });

    const checkboxField = selectCheckboxField(wrapper);

    checkboxField.triggerCustomEvent('change', true);

    expect(wrapper).toEmitInput(firstColor);
  });

  test('Value changed to transparent after trigger select field with false', () => {
    const wrapper = factory({
      propsData: {
        color: Faker.internet.color(),
      },
    });

    const checkboxField = selectCheckboxField(wrapper);

    checkboxField.triggerCustomEvent('change', false);

    expect(wrapper).toEmitInput('transparent');
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
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.colorPicker,
      config: {
        color,
        palette,
        action: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    const newColor = Faker.internet.color();

    modalArguments.config.action(newColor);

    expect(wrapper).toEmitInput(newColor);
  });
});
