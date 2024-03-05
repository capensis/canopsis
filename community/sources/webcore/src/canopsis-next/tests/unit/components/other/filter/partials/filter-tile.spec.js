import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import FilterTile from '@/components/other/filter/partials/filter-tile.vue';

const stubs = {
  'c-action-btn': true,
};

const selectEditButton = wrapper => wrapper.findAll('c-action-btn-stub').at(0);
const selectDeleteButton = wrapper => wrapper.findAll('c-action-btn-stub').at(1);

describe('filter-tile', () => {
  const factory = generateShallowRenderer(FilterTile, { stubs });
  const snapshotFactory = generateRenderer(FilterTile, { stubs });

  it('Edit event emitted after trigger edit button', () => {
    const wrapper = factory({
      propsData: {
        editable: true,
      },
    });

    selectEditButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toHaveBeenEmit('edit');
  });

  it('Delete event emitted after trigger delete button', () => {
    const wrapper = factory({
      propsData: {
        editable: true,
      },
    });

    selectDeleteButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toHaveBeenEmit('delete');
  });

  it('Renders `filter-tile` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `filter-tile` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filter: {
          title: 'Filter title',
        },
        editable: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `filter-tile` with old pattern', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filter: {
          title: 'Filter title',
          old_mongo_query: {},
        },
        editable: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
