import { generateRenderer } from '@unit/utils/vue';
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import EntityColumnPbehaviorInfo from '@/components/widgets/context/columns-formatting/entity-column-pbehavior-info.vue';

describe('entity-column-pbehavior-info', () => {
  const snapshotFactory = generateRenderer(EntityColumnPbehaviorInfo, {
    attachTo: document.body,
  });

  it('Renders `entity-column-pbehavior-info` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it.each(Object.values(PBEHAVIOR_TYPE_TYPES))('Renders `entity-column-pbehavior-info` with type: %s', async (type) => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehaviorInfo: {
          canonical_type: type,
          name: `pbehavior-${type}-name`,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});
