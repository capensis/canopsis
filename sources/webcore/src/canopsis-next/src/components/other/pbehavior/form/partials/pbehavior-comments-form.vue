<template lang="pug">
  div
    v-layout(
      :data-test="`pbehaviorComment-${index + 1}`",
      v-for="(comment, index) in comments",
      :key="comment.key",
      row,
      wrap,
      align-center
    )
      v-flex(xs11)
        v-textarea(
          data-test="pbehaviorCommentField",
          :disabled="!!comment._id",
          :label="$t('modals.createPbehavior.steps.comments.fields.message')",
          :value="comment.message",
          @input="updateFieldInArrayItem(index, 'message', $event)"
        )
      v-flex(xs1)
        v-btn(
          data-test="pbehaviorCommentDeleteButton",
          color="error",
          icon,
          @click="removeItemFromArray(index)"
        )
          v-icon delete
    v-layout(row)
      v-btn.ml-0.primary(
        data-test="pbehaviorAddCommentButton",
        type="button",
        @click="addComment"
      ) {{ $t('modals.createPbehavior.steps.comments.buttons.addComment') }}
</template>

<script>
import uid from '@/helpers/uid';

import authMixin from '@/mixins/auth';
import formArrayMixin from '@/mixins/form/array';

export default {
  mixins: [authMixin, formArrayMixin],
  model: {
    prop: 'comments',
    event: 'input',
  },
  props: {
    comments: {
      type: Array,
      required: true,
    },
    author: {
      type: String,
    },
  },
  methods: {
    addComment() {
      this.addItemIntoArray({
        key: uid(),
        author: this.author || this.currentUser._id,
        message: '',
      });
    },
  },
};
</script>
