<template>
  <c-layout column>
    <v-fade-transition>
      <v-progress-linear
        v-if="pending"
        color="primary"
        height="2"
        indeterminate
      />
    </v-fade-transition>
    <c-layout class="gap-3" column>
      <template v-if="comments.length">
        <entity-comment
          v-for="(comment, index) in comments"
          :key="comment._id"
          :comment="comment"
          :editable="index === 0 && editableFirst"
          @edit="edit"
        />
      </template>
      <span v-else class="font-italic text-center">
        <span v-if="pending" class="text-medium-emphasis">
          {{ $t('common.loadingItems') }}
        </span>
        <span v-else>
          {{ $t('entity.comments.emptyList') }}
        </span>
      </span>
    </c-layout>
  </c-layout>
</template>

<script>
import EntityComment from './entity-comment.vue';

export default {
  components: { EntityComment },

  props: {
    comments: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    editableFirst: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['edit'],

  setup(props, { emit }) {
    const edit = comment => emit('edit', comment);

    return {
      edit,
    };
  },
};
</script>
