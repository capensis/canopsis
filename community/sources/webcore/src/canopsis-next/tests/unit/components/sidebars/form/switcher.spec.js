import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import Switcher from '@/components/sidebars/form/fields/switcher.vue';

const stubs = {
  'v-switch': {
    props: ['inputValue'],
    template: `
      <input
        :checked="inputValue"
        type="checkbox"
        class="v-switch"
        @change="$listeners.change($event.target.checked)"
      />
    `,
  },
};

const selectSwitchField = wrapper => wrapper.find('input.v-switch');

describe('switcher', () => {
  const factory = generateShallowRenderer(Switcher, { stubs });
  const snapshotFactory = generateRenderer(Switcher, {
    stubs,
    parentComponent: {
      provide: {
        listClick: jest.fn(),
      },
    },
  });

  it('Value changed after trigger switch field', () => {
    const wrapper = factory({
      propsData: {
        title: '',
        value: false,
      },
    });

    const switchField = selectSwitchField(wrapper);

    switchField.setChecked(true);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(true);
  });

  it('Renders `switcher` with default and required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `switcher` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Custom title',
        value: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
