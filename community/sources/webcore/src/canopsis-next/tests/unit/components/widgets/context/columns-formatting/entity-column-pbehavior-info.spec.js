import { mount, createVueInstance } from '@unit/utils/vue';
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import EntityColumnPbehaviorInfo from '@/components/widgets/context/columns-formatting/entity-column-pbehavior-info.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(EntityColumnPbehaviorInfo, {
  localVue,
  attachTo: document.body,

  ...options,
});

describe('entity-column-pbehavior-info', () => {
  it('Renders `entity-column-pbehavior-info` with default props', () => {
    const wrapper = snapshotFactory();

    const tooltip = wrapper.findTooltip();

    expect(tooltip.element).toMatchSnapshot();
    expect(wrapper.element).toMatchSnapshot();
  });

  it.each(Object.values(PBEHAVIOR_TYPE_TYPES))('Renders `entity-column-pbehavior-info` with type: %s', (type) => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehaviorInfo: {
          canonical_type: type,
          name: `pbehavior-${type}-name`,
        },
      },
    });

    const tooltip = wrapper.findTooltip();

    expect(tooltip.element).toMatchSnapshot();
    expect(wrapper.element).toMatchSnapshot();
  });
});
