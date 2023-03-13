import { mount, shallowMount } from '@unit/utils/vue';

import CAlarmActionsChips from '@/components/common/chips/c-alarm-actions-chips.vue';

const stubs = {
  'c-alarm-action-chip': {
    template: '<span v-on="$listeners" class="c-alarm-action-chip"></span>',
  },
};

const factory = (options = {}) => shallowMount(CAlarmActionsChips, {
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CAlarmActionsChips, {
  stubs,

  ...options,
});

const selectChip = wrapper => wrapper.find('.c-alarm-action-chip');

describe('c-alarm-actions-chips', () => {
  const items = [
    { text: 'item1', color: 'color1' },
    { text: 'item2', color: 'color2' },
    { text: 'item3', color: 'color3' },
  ];
  const activeItem = items[2].text;

  test('Should emit `select` event on `click`', () => {
    const wrapper = factory({
      propsData: {
        items,
      },
    });

    const chip = selectChip(wrapper);
    chip.trigger('click');

    expect(wrapper).toEmit('select', items[0].text);
  });

  test('Should emit `select` event on `click` with custom itemValue', () => {
    const wrapper = factory({
      propsData: {
        items,
        itemValue: 'color',
      },
    });

    const chip = selectChip(wrapper);
    chip.trigger('click');

    expect(wrapper).toEmit('select', items[0].color);
  });

  test('Renders `c-alarm-tags-chips` without selected tag and dropdown', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items,
        inlineCount: 3,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-alarm-tags-chips` without selected tag and with dropdown', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items,
        inlineCount: 1,
      },
    });

    const dropdownContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  test('Renders `c-alarm-tags-chips` with selected tag and dropdown', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items,
        activeItem,
      },
    });

    const dropdownContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });
});
