import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import CTreeview from '@/components/common/treeview/c-treeview.vue';
import BroadcastMessageViewsTreeviewField from '@/components/other/broadcast-message/form/fields/broadcast-message-views-treeview-field.vue';

describe('broadcast-message-views-treeview-field', () => {
  const snapshotFactory = generateRenderer(BroadcastMessageViewsTreeviewField, {
    global: {
      components: { CTreeview },
    },
  });

  test('Renders selectable leaf items', async () => {
    const wrapper = snapshotFactory({
      props: {
        modelValue: ['login', 'profile'],
        items: [
          { value: 'login', name: 'Login' },
          { value: 'profile', name: 'Profile' },
        ],
      },
    });

    await flushPromises();

    expect(wrapper.findAll('.v-treeview-item')).toHaveLength(2);
  });
});
