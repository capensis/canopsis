import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import AddLocationBtn from '@/components/other/map/form/partials/add-location-btn.vue';

const localVue = createVueInstance();

const factory = (options = {}) => shallowMount(AddLocationBtn, {
  localVue,

  ...options,
});

const snapshotFactory = (options = {}) => mount(AddLocationBtn, {
  localVue,

  ...options,
});

const selectBtnToggleNode = wrapper => wrapper.vm.$children[0];

describe('add-location-btn', () => {
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

  test('Renders `add-location-btn` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: false,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
