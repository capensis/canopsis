<template lang="pug">
  div
    pbehavior-comment-field(
      v-for="(comment, index) in comments",
      v-field="comments[index]",
      :key="comment.key",
      @remove="removeItemFromArray(index)"
    )
    v-layout(row)
      v-btn.ml-0.primary(
        type="button",
        @click="addComment"
      ) {{ $t('modals.createPbehavior.steps.comments.buttons.addComment') }}
</template>

<script>
import uid from '@/helpers/uid';

import authMixin from '@/mixins/auth';
import formArrayMixin from '@/mixins/form/array';

import PbehaviorCommentField from '../fields/pbehavior-comment-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorCommentField },
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
