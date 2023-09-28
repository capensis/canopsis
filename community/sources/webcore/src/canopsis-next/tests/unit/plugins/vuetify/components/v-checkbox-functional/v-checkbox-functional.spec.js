import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import VCheckboxFunctional from '@/plugins/vuetify/components/v-checkbox-functional/v-checkbox-functional.vue';

const selectControlField = wrapper => wrapper.find('.v-input--selection-controls__input');
const selectLabel = wrapper => wrapper.find('label');

describe('v-checkbox-functional', () => {
  const factory = generateShallowRenderer(VCheckboxFunctional);
  const snapshotFactory = generateRenderer(VCheckboxFunctional);

  it('Value changed after click on the control element', () => {
    const wrapper = factory({
      propsData: {
        inputValue: true,
      },
    });

    const control = selectControlField(wrapper);

    control.trigger('click');

    const changeEvents = wrapper.emitted('change');
    expect(changeEvents).toHaveLength(1);

    const [eventData] = changeEvents[0];

    expect(eventData).toBe(false);
  });

  it('Value changed after click on the label', () => {
    const wrapper = factory({
      propsData: {
        inputValue: false,
      },
    });

    const label = selectLabel(wrapper);

    label.trigger('click');

    const changeEvents = wrapper.emitted('change');
    expect(changeEvents).toHaveLength(1);

    const [eventData] = changeEvents[0];

    expect(eventData).toBe(true);
  });

  it('Renders `v-checkbox-functional` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `v-checkbox-functional` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        inputValue: true,
        hideDetails: true,
        label: 'Custom label',
        disabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
