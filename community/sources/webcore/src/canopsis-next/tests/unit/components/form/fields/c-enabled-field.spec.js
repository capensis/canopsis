import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CEnabledField from '@/components/forms/fields/c-enabled-field.vue';

const stubs = {
  'v-switch': {
    props: ['value'],
    template: `
      <input :checked="value" type="checkbox" class="v-switch" @change="$listeners.input($event.target.checked)" />
    `,
  },
};

describe('c-enabled-field', () => {
  const factory = generateShallowRenderer(CEnabledField, { stubs });
  const snapshotFactory = generateRenderer(CEnabledField);

  it('Value set in the input', () => {
    const wrapper = factory({ propsData: { value: false } });
    const input = wrapper.find('input.v-switch');

    expect(input.element.checked).toBe(false);
  });

  it('Value changed in the input', () => {
    const wrapper = factory({ propsData: { value: false } });
    const input = wrapper.find('input.v-switch');

    input.setChecked(true);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);
    expect(inputEvents[0]).toEqual([true]);
  });

  it('Renders `c-enabled-field` with default props correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-enabled-field` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: false,
        label: 'Custom label',
        color: 'customcolor',
        disabled: true,
        hideDetails: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
