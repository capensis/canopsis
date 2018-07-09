<template lang="pug">
  v-card
    v-card-title
      span.headline Create watcher
    v-form
      v-layout(wrap)
        v-flex(xs10)
          v-text-field(label="Display name", v-model="form.name")
        v-flex(xs12)
          h3.text-xs-center Filter editor
          v-divider
          filter-editor
        v-btn(@click="submit") Submit
</template>

<script>
import FilterEditor from '@/components/other/filter-editor/filter-editor.vue';

export default {
  components: {
    FilterEditor,
  },
  data() {
    return {
      form: {
        name: '',
        filter: '',
      },
    };
  },
  methods: {
    submit() {
      const formData = { ...this.form, _id: 'df', type: 'watcher' };
    },
  },
};
</script>
