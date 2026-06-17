import Faker from 'faker';

import { generateRenderer } from '@unit/utils/vue';

import CList from '@/components/common/list/c-list.vue';

const selectListItem = (wrapper, index) => wrapper.findAll('.v-list-item').at(index);
const selectMenu = wrapper => wrapper.find('.v-menu');
const selectNestedList = wrapper => wrapper.findAllComponents({ name: 'c-list' }).at(1);

describe('c-list', () => {
  const factory = generateRenderer(CList, {
    attachTo: document.body,
  });

  test('Item selected after click on list item', () => {
    const value = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        items: [{ value }],
      },
    });

    const listItem = selectListItem(wrapper, 0);

    listItem.trigger('click');

    expect(wrapper).toEmitInput(value);
  });

  test('Item selected after click on list item (returnObject is enabled)', () => {
    const value = Faker.datatype.string();
    const item = { value };

    const wrapper = factory({
      propsData: {
        items: [item],
        returnObject: true,
      },
    });

    const listItem = selectListItem(wrapper, 0);

    listItem.trigger('click');

    expect(wrapper).toEmitInput(item);
  });

  test.each([
    [undefined, 'items'],
    ['variables', 'variables'],
  ])('Submenu opened when mouseover list item (childrenKey: %s)', async (childrenKey, internalKey) => {
    const zIndex = Faker.datatype.number();
    const subItems = [
      {
        value: 'parent.child',
      },
    ];
    const wrapper = factory({
      propsData: {
        zIndex,
        childrenKey: childrenKey ?? 'items',
        value: 'parent.child',
        items: [
          {
            value: 'parent',
            [internalKey]: subItems,
          },
        ],
      },
    });

    const listItem = selectListItem(wrapper, 0);

    jest.spyOn(listItem.element, 'getBoundingClientRect').mockReturnValue({
      top: 101,
      left: 112,
      width: 88,
    });
    await listItem.trigger('mouseenter');

    const menu = selectMenu(wrapper);
    expect(menu.exists()).toBe(true);
    expect(menu.vm.positionX).toEqual(200);
    expect(menu.vm.positionY).toEqual(101);
    expect(menu.vm.zIndex).toEqual(zIndex);

    const nestedList = selectNestedList(wrapper);

    expect(nestedList.exists()).toBe(true);
    expect(nestedList.vm.items).toEqual(subItems);
    expect(nestedList.props('zIndex')).toBe(zIndex + 1);
  });

  test('Submenu closed after mouseover on other list item', async () => {
    const wrapper = factory({
      propsData: {
        items: [
          {
            value: 'first',
            items: [],
          },
          {
            value: 'second',
          },
        ],
      },
    });

    const firstListItem = selectListItem(wrapper, 0);

    await firstListItem.trigger('mouseenter');

    expect(selectMenu(wrapper).exists()).toBe(true);

    const secondListItem = selectListItem(wrapper, 1);
    await secondListItem.trigger('mouseenter');

    expect(wrapper.vm.subItemsShown).toBe(false);
  });

  test('Sub item selected after trigger input on nested list', async () => {
    const parentValue = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        items: [
          {
            value: parentValue,
            items: [],
          },
        ],
      },
    });

    const firstListItem = selectListItem(wrapper, 0);

    await firstListItem.trigger('mouseenter');

    const nestedList = selectNestedList(wrapper);

    const value = Faker.datatype.string();

    nestedList.triggerCustomEvent('input', value);

    expect(wrapper).toHaveBeenEmit('input');
    expect(wrapper).toEmitInput(`${parentValue}.${value}`);
  });
});
