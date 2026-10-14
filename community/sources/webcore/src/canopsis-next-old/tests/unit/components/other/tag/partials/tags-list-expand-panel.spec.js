import { generateRenderer } from '@unit/utils/vue';

import TagsListExpandPanel from '@/components/other/tag/partials/tags-list-expand-panel.vue';

const stubs = {
  'tag-patterns-form': true,
};

describe('tags-list-expand-panel', () => {
  const snapshotFactory = generateRenderer(TagsListExpandPanel, { stubs });

  it('Renders `tags-list-expand-panel` with patterns', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tag: {
          alarm_pattern: [[]],
          entity_pattern: [[]],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
