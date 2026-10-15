import { generateRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CList from '@/components/common/list/c-list.vue';

const selectListItem = (wrapper, index) => wrapper.findAll('.v-list-item').at(index);

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
  ])(
    'Submenu opened when mouseover list item (childrenKey: %s)',
    async (childrenKey, internalKey) => {
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

      vi.spyOn(listItem.element, 'getBoundingClientRect').mockReturnValue({
        top: 101,
        left: 112,
        width: 88,
      });
      await listItem.trigger('mouseenter');

      expect(wrapper.vm.subItemsShown).toBe(true);
      expect(wrapper.vm.subItemsPosition).toEqual({ x: 200, y: 101 });
      expect(wrapper.vm.parentItem[internalKey]).toEqual(subItems);
      expect(wrapper.props('zIndex')).toEqual(zIndex);
    },
  );

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

    expect(wrapper.vm.subItemsShown).toBe(true);

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

    const value = Faker.datatype.string();

    wrapper.vm.selectSubItem(value);

    expect(wrapper).toEmitInput(`${parentValue}.${value}`);
  });
});
