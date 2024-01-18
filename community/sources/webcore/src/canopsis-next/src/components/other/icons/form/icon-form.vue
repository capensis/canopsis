<template>
  <v-layout column>
    <v-layout align-center>
      <v-flex
        class="mr-2"
        v-if="iconOriginalTitle"
        xs2
      >
        <v-icon
          large
        >
          $vuetify.icons.{{ iconOriginalTitle }}
        </v-icon>
      </v-flex>
      <c-name-field
        v-field="form.title"
        :label="$t('common.title')"
        :error-messages="errors.collect('title')"
        name="title"
        required
      />
    </v-layout>
    <icon-file-selector
      v-if="!iconOriginalTitle"
      v-field="form.file"
    />
  </v-layout>
</template>

<script>
import IconFileSelector from './fields/icon-file-selector.vue';

export default {
  inject: ['$validator'],
  components: {
    IconFileSelector,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      iconOriginalTitle: this.form.title,
    };
  },
};
</script>
