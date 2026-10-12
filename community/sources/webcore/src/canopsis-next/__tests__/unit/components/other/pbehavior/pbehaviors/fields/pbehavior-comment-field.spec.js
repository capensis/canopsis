import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import PbehaviorCommentField from '@/components/other/pbehavior/pbehaviors/fields/pbehavior-comment-field.vue';

const stubs = {
  'c-textarea': {
    name: 'CTextarea',
    emits: ['update:modelValue'],
    template: '<div />',
  },
  'c-action-btn': {
    name: 'CActionBtn',
    emits: ['click'],
    template: '<button @click="$emit(\'click\', $event)" />',
  },
};

const selectMessageField = wrapper => wrapper.findComponent({ name: 'CTextarea' });
const selectRemoveButton = wrapper => wrapper.findComponent({ name: 'CActionBtn' });

describe('pbehavior-comment-field', () => {
  const factory = generateShallowRenderer(PbehaviorCommentField, { stubs });

  test('Message changed after trigger textarea', async () => {
    const comment = {
      key: Faker.datatype.string(),
      message: Faker.datatype.string(),
    };
    const wrapper = factory({
      props: { comment },
    });

    const newMessage = Faker.datatype.string();

    selectMessageField(wrapper).vm.$emit('update:modelValue', newMessage);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:comment')).toEqual([
      [
        {
          ...comment,
          message: newMessage,
        },
      ],
    ]);
  });

  test('Remove event emitted after trigger remove button', async () => {
    const wrapper = factory({
      props: { comment: {} },
    });

    selectRemoveButton(wrapper).vm.$emit('click');
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('remove')).toHaveLength(1);
  });
});
