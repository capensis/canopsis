import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import MermaidAddLocationBtn from '@/components/other/map/form/partials/mermaid-add-location-btn.vue';

const localVue = createVueInstance();

const stubs = {
  'code-editor': true,
};

const factory = (options = {}) => shallowMount(MermaidAddLocationBtn, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(MermaidAddLocationBtn, {
  localVue,
  stubs,

  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },

  ...options,
});

const selectBtnToggleNode = wrapper => wrapper.vm.$children[0];

describe('mermaid-add-location-btn', () => {
  test('Value updated after click button', () => {
    const wrapper = factory({
      propsData: {
        value: true,
      },
    });

    const btnToggleNode = selectBtnToggleNode(wrapper);

    btnToggleNode.$emit('change');

    expect(wrapper).toEmit('input', false);
  });

  test('Renders `mermaid-add-location-btn` with custom props ', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: false,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
