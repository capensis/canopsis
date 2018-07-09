<template lang="pug">
  v-card
    v-card-title
      span.headline Create watcher
    v-form
      v-layout(wrap, justify-center)
        v-flex(xs10)
          v-text-field(label="Display name", v-model="form.display_name")
        v-flex(xs12)
          h3.text-xs-center Filter editor
          v-divider
          filter-editor
        v-btn(@click="submit") Submit
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import FilterEditor from '@/components/other/filter-editor/filter-editor.vue';

const { mapActions: watcherMapActions } = createNamespacedHelpers('watcher');
const { mapGetters: filterEditorMapGetters } = createNamespacedHelpers('mFilterEditor');

export default {
  components: {
    FilterEditor,
  },
  data() {
    return {
      form: {
        display_name: '',
        mfilter: '',
      },
    };
  },
  computed: {
    ...filterEditorMapGetters(['request']),
  },
  methods: {
    ...watcherMapActions(['create']),
    submit() {
      const formData = {
        ...this.form,
        _id: this.form.display_name,
        type: 'watcher',
        mfilter: JSON.stringify(this.request),
      };
      this.create(formData);
    },
  },
};
</script>
