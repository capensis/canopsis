import { generateShallowRenderer } from '@unit/utils/vue';
import { shallowMount } from '@vue/test-utils';

import { PAGINATION_LIMIT, PAGINATION_PER_PAGE_VALUES } from '@/config';

import CItemsPerPageField from '@/components/forms/fields/c-items-per-page-field.vue';

const fullscreenState = vi.hoisted(() => ({
  isFullscreen: { value: false, __v_isRef: true },
}));

vi.mock('@/hooks/fullscreen', () => ({
  useFullscreen: () => fullscreenState,
}));

const stubs = {
  'c-select': {
    name: 'CItemsPerPageSelectStub',
    inheritAttrs: false,
    props: {
      modelValue: [Number, String],
      items: Array,
      menuProps: Object,
    },
    emits: ['update:modelValue'],
    template: `
      <select class="c-select" :value="modelValue">
        <option v-for="item in items" :key="item" :value="item">{{ item }}</option>
      </select>
    `,
  },
};

const selectCSelect = wrapper => wrapper.findComponent({ name: 'CItemsPerPageSelectStub' });
const selectOptionValues = wrapper =>
  selectCSelect(wrapper)
    .findAll('option')
    .map(option => Number(option.element.value));

describe('c-items-per-page-field', () => {
  const factory = generateShallowRenderer(CItemsPerPageField, { stubs });

  afterEach(() => {
    fullscreenState.isFullscreen.value = false;
    Object.defineProperty(document, 'fullscreenElement', {
      configurable: true,
      value: null,
    });
  });

  it('Default items is equal to PAGINATION_PER_PAGE_VALUES', () => {
    const wrapper = factory();

    expect(wrapper.props('items')).toEqual(PAGINATION_PER_PAGE_VALUES);
  });

  it('Default items property applied to select element', () => {
    const wrapper = factory();

    expect(selectOptionValues(wrapper)).toEqual(PAGINATION_PER_PAGE_VALUES);
  });

  it('Attaches the menu to the fullscreen root', () => {
    const fullscreenElement = document.createElement('div');

    fullscreenState.isFullscreen.value = true;
    Object.defineProperty(document, 'fullscreenElement', {
      configurable: true,
      value: fullscreenElement,
    });

    const wrapper = shallowMount(CItemsPerPageField, {
      global: { stubs },
    });

    expect(selectCSelect(wrapper).props('menuProps')).toEqual({ attach: fullscreenElement });

    wrapper.unmount();
  });

  it('Custom items property applied to select element', () => {
    const items = [1, 2, 3];
    const wrapper = factory({ propsData: { items } });

    expect(selectOptionValues(wrapper)).toEqual(items);
  });

  it('Default value property is equal to PAGINATION_LIMIT', () => {
    const wrapper = factory();

    expect(wrapper.props('modelValue')).toBe(PAGINATION_LIMIT);
  });

  it('Default value property applied to select element', () => {
    const wrapper = factory();

    expect(Number(selectCSelect(wrapper).element.value)).toBe(PAGINATION_LIMIT);
  });

  it('Custom value property applied to select element', () => {
    const value = PAGINATION_PER_PAGE_VALUES[0];
    const wrapper = factory({ propsData: { modelValue: value } });

    expect(Number(selectCSelect(wrapper).element.value)).toBe(value);
  });

  it('Custom items and value properties applied to select element', () => {
    const items = [1, 2, 3, 4, 5];
    const value = items[0];
    const wrapper = factory({ propsData: { modelValue: value, items } });

    expect(Number(selectCSelect(wrapper).element.value)).toBe(value);
    expect(selectOptionValues(wrapper)).toEqual(items);
  });

  it('Set value into select element', () => {
    const value = PAGINATION_PER_PAGE_VALUES[0];
    const wrapper = factory();

    selectCSelect(wrapper).triggerCustomEvent('update:modelValue', value);

    expect(wrapper).toEmitInput(value);
  });
});
