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

import { formArrayMixin } from '@/mixins/form';

import PbehaviorCommentField from './pbehavior-comment-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorCommentField },
  mixins: [formArrayMixin],
  model: {
    prop: 'comments',
    event: 'input',
  },
  props: {
    comments: {
      type: Array,
      required: true,
    },
  },
  methods: {
    addComment() {
      this.addItemIntoArray({
        key: uid(),
        message: '',
      });
    },
  },
};
</script>
