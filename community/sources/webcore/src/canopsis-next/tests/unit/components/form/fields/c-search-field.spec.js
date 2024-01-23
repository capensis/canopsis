import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import CSearchField from '@/components/forms/fields/c-search-field.vue';

const mockData = {
  search: Faker.lorem.words(),
  newSearch: Faker.lorem.words(),
};

const stubs = {
  'v-text-field': {
    props: ['value'],
    template: `
      <input
        class="v-text-field"
        :value="value"
        @input="$emit('input', $event.target.value)"
        @keydown="$listeners.keydown"
      />
    `,
  },
  'v-combobox': createSelectInputStub('v-combobox'),
  'c-action-btn': true,
};

const snapshotStubs = {
  'c-action-btn': true,
};

const selectTextInput = wrapper => wrapper.find('input.v-text-field');
const selectCombobox = wrapper => wrapper.find('select.v-combobox');
const selectSubmitButton = wrapper => wrapper.findAll('c-action-btn-stub').at(0);
const selectClearButton = wrapper => wrapper.findAll('c-action-btn-stub').at(1);

describe('c-search-field', () => {
  const factory = generateShallowRenderer(CSearchField, { stubs });
  const snapshotFactory = generateRenderer(CSearchField, { stubs: snapshotStubs });

  it('Not empty value was pass into props and it was applied to input field', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const input = selectTextInput(wrapper);

    expect(input.element.value).toBe(search);
  });

  it('Input value was update, after prop change and it was applied to input field', async () => {
    const { search, newSearch } = mockData;

    const wrapper = factory({ propsData: { value: search } });

    await wrapper.setProps({ value: newSearch });

    const input = selectTextInput(wrapper);

    expect(input.element.value).toBe(newSearch);
  });

  it('Set value into input element', () => {
    const { search } = mockData;

    const wrapper = factory();
    const input = selectTextInput(wrapper);

    input.setValue(search);

    expect(wrapper.vm.localValue).toEqual(search);
  });

  it('Keyup without enter key on input element', () => {
    const { search, newSearch } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const input = selectTextInput(wrapper);

    input.setValue(newSearch);
    input.trigger('keyup');

    expect(wrapper).not.toHaveBeenEmit('submit');
  });

  it('Keyup with enter key on input element', async () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const input = selectTextInput(wrapper);

    await input.trigger('keydown.enter');

    expect(wrapper).toEmit('submit', search);
  });

  it('Submit search button is the first button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const submitButton = selectSubmitButton(wrapper);

    expect(submitButton.attributes('icon')).toBe('search');
  });

  it('Clear search button is the second button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const clearButton = selectClearButton(wrapper);

    expect(clearButton.attributes('icon')).toBe('clear');
  });

  it('Click on submit button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const submitButton = selectSubmitButton(wrapper);

    submitButton.triggerCustomEvent('click');

    expect(wrapper).toEmit('submit', search);
  });

  it('Click on clear button', () => {
    const { search } = mockData;

    const wrapper = factory({ propsData: { value: search } });
    const clearButton = selectClearButton(wrapper);

    clearButton.triggerCustomEvent('click');

    expect(wrapper).toHaveBeenEmit('clear');
  });

  it('Set value into combobox element', () => {
    const items = [
      { search: Faker.lorem.words(), pinned: false },
      { search: Faker.lorem.words(), pinned: false },
    ];

    const wrapper = factory({
      propsData: {
        combobox: true,
        items,
      },
    });

    const combobox = selectCombobox(wrapper);

    combobox.setValue(items[0].search);

    expect(wrapper).toEmit('submit', items[0].search);
  });

  it('Renders `c-search-field` correctly', () => {
    const wrapper = snapshotFactory({
      propsData: { value: 'c-search-field' },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-search-field` correctly with combobox and items', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'c-search-field',
        combobox: true,
        items: [
          { search: 'item 1', pinned: true },
          { search: 'item 2', pinned: false },
        ],
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
