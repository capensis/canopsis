import { generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import CAlarmTagsChips from '@/components/common/chips/c-alarm-tags-chips.vue';

const stubs = {
  'c-alarm-actions-chips': true,
};

describe('c-alarm-tags-chips', () => {
  const tags = [
    { value: 'tag1', color: 'color1' },
    { value: 'tag2', color: 'color2' },
    { value: 'tag3', color: 'color3' },
  ];
  const selectedTags = [tags[2].value];
  const alarm = {
    tags: [tags[0].value],
  };
  const alarmTagModule = {
    name: 'alarmTag',
    getters: {
      items: () => tags,
    },
  };

  const store = createMockedStoreModules([alarmTagModule]);

  const snapshotFactory = generateRenderer(CAlarmTagsChips, { stubs });

  it('Renders `c-alarm-tags-chips` with selectedTags and without dropdown', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm,
        selectedTags,
      },
      store,
    });

    expect(wrapper).toMatchSnapshot();
  });
});
