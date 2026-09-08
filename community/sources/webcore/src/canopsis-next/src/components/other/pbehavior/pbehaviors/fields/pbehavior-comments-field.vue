<template>
  <c-form-block-array-field
    v-field="comments"
    :item-to-form="pbehaviorCommentItemToForm"
    :add-button-label="$t('modals.createPbehavior.steps.comments.buttons.addComment')"
    :label="$tc('common.comment', 2)"
  >
    <template #item="{ index, remove }">
      <pbehavior-comment-field
        v-field="comments[index]"
        @remove="remove"
      />
    </template>
  </c-form-block-array-field>
</template>

<script>
import { pbehaviorCommentItemToForm } from '@/helpers/entities/pbehavior/form';

import PbehaviorCommentField from './pbehavior-comment-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorCommentField },
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
  setup() {
    return {
      pbehaviorCommentItemToForm,
    };
  },
};
</script>
