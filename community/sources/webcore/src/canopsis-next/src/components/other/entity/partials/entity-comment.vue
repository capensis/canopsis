<template>
  <v-layout column>
    <v-layout align-center>
      <v-layout class="gap-2" column>
        <v-layout column>
          <span class="comment-author font-weight-bold">{{ comment.author.display_name }}</span>
          <span class="text-subtitle-2 grey--text">{{ comment.t | date }}</span>
        </v-layout>
        <span class="pre-wrap comment-message">{{ comment.message }}</span>
      </v-layout>
      <v-flex shrink>
        <v-layout class="gap-1" column>
          <c-action-btn
            v-if="editable"
            type="edit"
            @click="edit"
          />
        </v-layout>
      </v-flex>
    </v-layout>
    <v-divider class="mt-3" />
  </v-layout>
</template>

<script>
export default {
  props: {
    comment: {
      type: Object,
      default: () => ({}),
    },
    editable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const edit = () => emit('edit', props.comment);

    return {
      edit,
    };
  },
};
</script>
