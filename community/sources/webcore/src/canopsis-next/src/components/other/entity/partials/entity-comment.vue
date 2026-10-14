<template>
  <c-layout column>
    <c-layout align="center">
      <c-layout class="gap-2" column>
        <c-layout column>
          <span class="font-weight-bold">
            {{ comment.author.display_name }}
          </span>
          <span class="text-subtitle-2 text-medium-emphasis">
            {{ formattedDate }}
          </span>
        </c-layout>
        <span class="pre-wrap">
          {{ comment.message }}
        </span>
      </c-layout>
      <c-layout v-if="editable">
        <c-action-btn type="edit" @click="edit" />
      </c-layout>
    </c-layout>
    <v-divider class="mt-3" />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

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

  emits: ['edit'],

  setup(props, { emit }) {
    const formattedDate = computed(() => convertDateToString(props.comment.t));
    const edit = () => emit('edit', props.comment);

    return {
      formattedDate,
      edit,
    };
  },
};
</script>
