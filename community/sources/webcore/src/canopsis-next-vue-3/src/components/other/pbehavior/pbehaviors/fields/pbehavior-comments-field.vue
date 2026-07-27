<template>
  <c-form-block-array-field
    :form="comments"
    :item-to-form="pbehaviorCommentItemToForm"
    :add-button-label="$t('modals.createPbehavior.steps.comments.buttons.addComment')"
    :label="$t('common.comment', 2)"
    @update:form="$emit('update:comments', $event)"
  >
    <template #item="{ index, remove }">
      <pbehavior-comment-field
        :comment="comments[index]"
        @update:comment="updateItemInArray(index, $event)"
        @remove="remove"
      />
    </template>
  </c-form-block-array-field>
</template>

<script>
import { pbehaviorCommentItemToForm } from '@/helpers/entities/pbehavior/form';
import { useArrayModelField } from '@/hooks/form/array-model-field';

import PbehaviorCommentField from './pbehavior-comment-field.vue';

export default {
  components: { PbehaviorCommentField },

  props: {
    comments: {
      type: Array,
      required: true,
    },
  },

  emits: ['update:comments'],

  setup(props, { emit }) {
    const { updateItemInArray } = useArrayModelField(props, emit);

    return {
      pbehaviorCommentItemToForm,
      updateItemInArray,
    };
  },
};
</script>
