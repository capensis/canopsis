import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import TextEditor from '@/components/common/text-editor/text-editor.vue';

vi.mock('@/plugins/vuetify/components/icons', () => ({}));

const selectVariablesButton = wrapper => wrapper.find('.jodit-ui-group__variables button');

describe('text-editor variables late loading', () => {
  const factory = generateRenderer(TextEditor, {
    stubs: {
      'c-label': true,
      'variables-menu': true,
    },
  });

  test('adds the variables button when variables arrive after mount', async () => {
    const wrapper = factory({
      props: {
        variables: [],
      },
    });

    expect(selectVariablesButton(wrapper).exists()).toBe(false);

    await wrapper.setProps({
      variables: [
        {
          text: 'env',
          variables: [{ text: 'Location', value: '{{ env.Location }}' }],
        },
      ],
    });
    await flushPromises();

    expect(selectVariablesButton(wrapper).exists()).toBe(true);
  });
});
