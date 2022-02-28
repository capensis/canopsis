import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import FilterTile from '@/components/other/filter/form/partials/filter-tile.vue';

const localVue = createVueInstance();

const stubs = {
  'c-action-btn': true,
};

const factory = (options = {}) => shallowMount(FilterTile, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(FilterTile, {
  localVue,
  stubs,

  ...options,
});

const selectEditButton = wrapper => wrapper.findAll('c-action-btn-stub').at(0);
const selectDeleteButton = wrapper => wrapper.findAll('c-action-btn-stub').at(1);

describe('filter-tile', () => {
  it('Edit event emitted after trigger edit button', () => {
    const wrapper = factory({
      propsData: {
        editable: true,
      },
    });

    const editButton = selectEditButton(wrapper);

    editButton.vm.$emit('click');

    const editEvents = wrapper.emitted('edit');

    expect(editEvents).toHaveLength(1);
  });

  it('Delete event emitted after trigger delete button', () => {
    const wrapper = factory({
      propsData: {
        editable: true,
      },
    });

    const deleteButton = selectDeleteButton(wrapper);

    deleteButton.vm.$emit('click');

    const deleteEvents = wrapper.emitted('delete');

    expect(deleteEvents).toHaveLength(1);
  });

  it('Renders `filter-tile` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
