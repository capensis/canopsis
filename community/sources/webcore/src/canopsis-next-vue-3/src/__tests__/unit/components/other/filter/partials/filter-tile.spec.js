import { generateShallowRenderer } from '@unit/utils/vue';

import FilterTile from '@/components/other/filter/partials/filter-tile.vue';

const stubs = {
  'c-action-btn': true,
};

const selectEditButton = wrapper => wrapper.findAll('c-action-btn-stub').at(0);
const selectDeleteButton = wrapper => wrapper.findAll('c-action-btn-stub').at(1);

describe('filter-tile', () => {
  const factory = generateShallowRenderer(FilterTile, { stubs });

  it('Edit event emitted after trigger edit button', () => {
    const wrapper = factory({
      propsData: {
        editable: true,
      },
    });

    wrapper.vm.$emit('edit');

    expect(wrapper.emitted('edit')).toHaveLength(1);
  });

  it('Delete event emitted after trigger delete button', () => {
    const wrapper = factory({
      propsData: {
        editable: true,
      },
    });

    wrapper.vm.$emit('delete');

    expect(wrapper.emitted('delete')).toHaveLength(1);
  });
});
