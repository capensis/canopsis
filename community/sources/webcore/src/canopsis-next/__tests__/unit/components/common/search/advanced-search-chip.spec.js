import { flushPromises, shallowMount } from '@vue/test-utils';

import AdvancedSearchChip from '@/components/common/search/partials/advanced-search-chip.vue';
import VariablesList from '@/components/common/text-editor/variables-list.vue';

describe('advanced-search-chip', () => {
  let wrapper;

  const factory = (props) => {
    wrapper = shallowMount(AdvancedSearchChip, {
      props,
      global: {
        stubs: {
          'c-tooltip': true,
          'c-icon': true,
          'v-menu': { template: '<div><slot /></div>' },
        },
      },
    });

    return wrapper.findComponent(VariablesList);
  };

  afterEach(() => wrapper?.unmount());

  it('updates the dropdown when aliases arrive after dictionary keys', async () => {
    const infos = {
      text: 'Infos',
      value: 'infos',
      items: [{ text: 'environment', value: 'environment' }],
    };
    const list = factory({ items: [infos] });

    await flushPromises();

    expect(list.props('items')).toEqual([infos]);

    const alias = { text: 'ENV', value: 'ENV', alias: true };

    await wrapper.setProps({ items: [infos, alias] });
    await flushPromises();

    expect(list.props('items')).toEqual([infos, alias]);

    list.vm.$emit('update:value', alias);
    await flushPromises();

    expect(wrapper.emitted('update:modelValue')).toEqual([['ENV']]);
  });

  it('replaces aliases in the dropdown when the number of items stays the same', async () => {
    const list = factory({ items: [{ text: 'ENV', value: 'ENV', alias: true }] });

    await flushPromises();

    const alias = { text: 'Environment', value: 'Environment', alias: true };

    await wrapper.setProps({ items: [alias] });
    await flushPromises();

    expect(list.props('items')).toEqual([alias]);
  });

  it('keeps cached suggestions when items are recreated with the same content', async () => {
    const getItems = () => [
      {
        text: 'Infos',
        value: 'infos',
        items: [{ text: 'environment', value: 'environment' }],
      },
      { text: 'ENV', value: 'ENV', alias: true },
    ];
    const list = factory({ items: getItems() });

    await flushPromises();

    const cachedItems = list.props('items');

    await wrapper.setProps({ items: getItems() });
    await flushPromises();

    expect(list.props('items')).toBe(cachedItems);
  });
});
