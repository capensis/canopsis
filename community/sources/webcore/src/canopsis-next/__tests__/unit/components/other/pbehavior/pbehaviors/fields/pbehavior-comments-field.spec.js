import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import PbehaviorCommentsField from '@/components/other/pbehavior/pbehaviors/fields/pbehavior-comments-field.vue';

const stubs = {
  'c-form-block-array-field': {
    name: 'CFormBlockArrayField',
    props: ['form', 'itemToForm'],
    emits: ['update:form'],
    template: `
      <div>
        <div v-for="(item, index) in form" :key="item.key">
          <slot name="item" :index="index" :remove="() => remove(index)" />
        </div>
        <button type="button" @click="add">add</button>
      </div>
    `,
    methods: {
      add() {
        this.$emit('update:form', [...this.form, this.itemToForm()]);
      },
      remove(index) {
        this.$emit(
          'update:form',
          this.form.filter((item, itemIndex) => itemIndex !== index),
        );
      },
    },
  },
  'pbehavior-comment-field': {
    name: 'PbehaviorCommentField',
    props: ['comment'],
    emits: ['update:comment', 'remove'],
    template: '<div />',
  },
};

const selectAddCommentButton = wrapper => wrapper.find('button');
const selectCommentFieldByIndex = (wrapper, index) =>
  wrapper.findAllComponents({ name: 'PbehaviorCommentField' }).at(index);

describe('pbehavior-comments-field', () => {
  const factory = generateShallowRenderer(PbehaviorCommentsField, {
    stubs,
  });

  test('Comment added after trigger create button', async () => {
    const comments = [{ key: Faker.datatype.string(), message: Faker.datatype.string() }];
    const wrapper = factory({
      props: {
        comments,
      },
    });

    await selectAddCommentButton(wrapper).trigger('click');

    expect(wrapper.emitted('update:comments')).toEqual([
      [
        [
          ...comments,
          {
            key: expect.any(String),
            message: '',
          },
        ],
      ],
    ]);
  });

  test('Comment changed after trigger comment field', async () => {
    const comments = [
      { key: 'comment-1', message: 'message-1' },
      { key: 'comment-2', message: 'message-2' },
    ];
    const wrapper = factory({
      props: {
        comments,
      },
    });

    const newComment = {
      key: Faker.datatype.string(),
      message: Faker.datatype.number(),
    };

    selectCommentFieldByIndex(wrapper, 1).vm.$emit('update:comment', newComment);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:comments')).toEqual([[[comments[0], newComment]]]);
  });

  test('Comment removed after trigger comment field', async () => {
    const comments = [
      { key: 'comment-1', message: 'message-1' },
      { key: 'comment-2', message: 'message-2' },
    ];
    const wrapper = factory({
      props: {
        comments,
      },
    });

    selectCommentFieldByIndex(wrapper, 1).vm.$emit('remove');
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:comments')).toEqual([[[comments[0]]]]);
  });
});
