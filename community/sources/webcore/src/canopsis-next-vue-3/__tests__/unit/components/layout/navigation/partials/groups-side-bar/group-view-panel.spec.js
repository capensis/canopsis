import { generateShallowRenderer } from '@unit/utils/vue';

import GroupsPanel from '@/components/layout/navigation/partials/groups-side-bar/group-view-panel.vue';

describe('group-view-panel', () => {
  const factory = generateShallowRenderer(GroupsPanel);

  it('Change event emitted after trigger edit button', () => {
    const wrapper = factory({
      propsData: {
        view: {},
        editable: true,
      },
    });

    wrapper.vm.$emit('change');

    expect(wrapper).toHaveBeenEmit('change');
  });

  it('Duplicate event emitted after trigger duplicate button', () => {
    const wrapper = factory({
      propsData: {
        view: {},
        editable: true,
        duplicable: true,
      },
    });

    wrapper.vm.$emit('duplicate');

    expect(wrapper).toHaveBeenEmit('duplicate');
  });
});
