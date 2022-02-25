import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import Switcher from '@/components/sidebars/settings/fields/common/switcher.vue';

const localVue = createVueInstance();

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

const factory = (options = {}) => shallowMount(Switcher, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(Switcher, {
  localVue,
  stubs,

  parentComponent: {
    provide: {
      listClick: jest.fn(),
    },
  },

  ...options,
});

const selectSwitchField = wrapper => wrapper.find('input.v-switch');

describe('switcher', () => {
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

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `switcher` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Custom title',
        value: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
