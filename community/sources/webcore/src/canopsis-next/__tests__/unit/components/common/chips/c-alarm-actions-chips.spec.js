import { generateShallowRenderer } from '@unit/utils/vue';

import CAlarmActionsChips from '@/components/common/chips/c-alarm-actions-chips.vue';

const stubs = {
  'c-chip': {
    template: '<span class="c-chip" @click="$emit(\'click\')"></span>',
    emits: ['click'],
  },
};

const selectChip = wrapper => wrapper.find('.c-chip');

describe('c-alarm-actions-chips', () => {
  const factory = generateShallowRenderer(CAlarmActionsChips, { stubs });

  const items = [
    { text: 'item1', color: 'color1' },
    { text: 'item2', color: 'color2' },
    { text: 'item3', color: 'color3' },
  ];

  test('Should emit `select` event on `click`', () => {
    const wrapper = factory({
      propsData: {
        items,
      },
    });

    const chip = selectChip(wrapper);
    chip.trigger('click');

    expect(wrapper).toHaveBeenEmitWith('select', items[0].text);
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

    expect(wrapper).toHaveBeenEmitWith('select', items[0].color);
  });

  test('Should emit `select` event on `click` with returnObject', () => {
    const wrapper = factory({
      propsData: {
        items,
        returnObject: true,
      },
    });

    const chip = selectChip(wrapper);
    chip.trigger('click');

    expect(wrapper).toHaveBeenEmitWith('select', items[0]);
  });

  test('Should not emit `select` when clicking an already active item', () => {
    const wrapper = factory({
      propsData: {
        items,
        activeItems: [items[0].text],
      },
    });

    selectChip(wrapper).trigger('click');

    expect(wrapper).not.toHaveBeenEmit('select');
  });
});
