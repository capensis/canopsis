<template>
  <v-layout column>
    <v-layout align-center>
      <v-flex
        v-if="iconOriginalTitle"
        class="mr-2"
        xs2
      >
        <v-icon large>
          $vuetify.icons.{{ iconOriginalTitle }}
        </v-icon>
      </v-flex>
      <v-text-field
        v-field="form.title"
        v-validate="titleRules"
        :label="$t('common.title')"
        :error-messages="titleErrorMessages"
        name="title"
      />
    </v-layout>
    <icon-file-selector
      v-if="!iconOriginalTitle"
      v-field="form.file"
      :max-file-size="maxFileSize"
      class="mt-2"
    />
  </v-layout>
</template>

<script>
import { ICON_TITLE_REGEX, MAX_ICON_SIZE_IN_KB } from '@/constants';

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
  computed: {
    titleRules() {
      return {
        required: true,
        regex: ICON_TITLE_REGEX,
      };
    },

    titleErrorMessages() {
      return this.errors.collect('title', null, false)
        .map(({ rule, msg }) => (rule === 'regex' ? this.$t('icon.validation.iconTitleRegex') : msg));
    },

    maxFileSize() {
      return MAX_ICON_SIZE_IN_KB;
    },
  },
};
</script>
