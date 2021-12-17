import { mount, createVueInstance } from '@unit/utils/vue';

import GridOverviewWidget from '@/components/widgets/grid-overview-widget.vue';

const localVue = createVueInstance();

const stubs = {
  'grid-overview-item': true,
};

const snapshotFactory = (options = {}) => mount(GridOverviewWidget, {
  localVue,
  stubs,

  ...options,
});

describe('grid-overview-widget', () => {
  it('Renders `grid-overview-widget` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tab: {
          widgets: [],
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `grid-overview-widget` with widgets', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tab: {
          widgets: [{ _id: 'id1' }, { _id: 'id2' }],
        },
      },
      scopedSlots: {
        default(props) {
          return this.$createElement(
            'div',
            { attrs: { class: 'default-slot' } },
            JSON.stringify(props),
          );
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
