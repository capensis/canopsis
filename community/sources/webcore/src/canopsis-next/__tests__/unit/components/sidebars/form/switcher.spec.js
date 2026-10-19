import { generateShallowRenderer } from '@unit/utils/vue';

import Switcher from '@/components/sidebars/form/fields/switcher.vue';

const stubs = {
  'v-switch': {
    props: ['inputValue'],
    template: `
      <input
        :checked="inputValue"
        type="checkbox"
        class="v-switch"
        @change="$emit('change', $event.target.checked)"
      />
    `,
  },
};

const selectSwitchField = wrapper => wrapper.find('input.v-switch');

describe('switcher', () => {
  const factory = generateShallowRenderer(Switcher, { stubs });

  it('Value changed after trigger switch field', () => {
    const wrapper = factory({
      propsData: {
        title: '',
        value: false,
      },
    });

    selectSwitchField(wrapper).setChecked(true);

    expect(wrapper).toEmitInput(true);
  });
});
