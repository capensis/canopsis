import { generateShallowRenderer } from '@unit/utils/vue';

import GroupsPanel from '@/components/layout/navigation/partials/groups-side-bar/group-panel.vue';

describe('group-panel', () => {
  const factory = generateShallowRenderer(GroupsPanel);

  test('Change event emitted after trigger button', async () => {
    const wrapper = factory({
      propsData: {
        editable: true,
        group: {},
      },
    });

    wrapper.vm.handleChange();

    expect(wrapper).toHaveBeenEmit('change');
  });
});
